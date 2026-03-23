package services

import (
	"context"
	"database/sql"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ywl0806/yuno_kiroku/internal/consts"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/enums"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	imagepkg "github.com/ywl0806/yuno_kiroku/pkg/image"
)

type MediaItemService struct {
	mediaItemStore       store.MediaItemStore
	imageUploader        *ImageUploader
	faceService          *FaceService
	identityFaceImgStore store.IdentityFaceImgStore
}

func NewMediaItemService(
	mediaItemStore store.MediaItemStore,
	imageUploader *ImageUploader,
	faceService *FaceService,
	identityFaceImgStore store.IdentityFaceImgStore,
) *MediaItemService {
	return &MediaItemService{
		mediaItemStore:       mediaItemStore,
		imageUploader:        imageUploader,
		faceService:          faceService,
		identityFaceImgStore: identityFaceImgStore,
	}
}

// UploadImageResult는 이미지 업로드 시작 후 반환하는 응답입니다.
type UploadImageResult struct {
	MediaItemID int32
	Status      string
}

// bgParams는 백그라운드 이미지 처리에 필요한 파라미터입니다.
type bgParams struct {
	data        []byte
	meta        imagepkg.Meta
	uploadPath  string
	mediaItemID int32
	familyId    int32
	albumId     int32
	retry       bool
}

// multipart 파일을 읽어 원본을 업로드하고 DB 레코드를 생성한 후,
// 리사이즈 이미지와 얼굴 인식을 백그라운드에서 처리
func (s *MediaItemService) UploadImage(
	ctx context.Context,
	file *multipart.FileHeader,
	familyId, albumId, uploadBatchID int32,
	retry bool,
) (*UploadImageResult, error) {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Filename), "."))

	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	originalData, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	meta, err := imagepkg.Parse(originalData, ext)
	if err != nil {
		log.Println("이미지 파싱 실패:", err)
		return nil, err
	}

	uploadPath := s.imageUploader.BuildUploadPath(familyId, albumId)

	// 원본 이미지 업로드
	originalKey, err := s.imageUploader.UploadOriginal(originalData, uploadPath, meta)
	if err != nil {
		return nil, err
	}

	// media_item과 원본 media_file 레코드 생성
	mediaItem, err := s.createMediaItemWithOriginal(ctx, meta, originalKey, file.Filename, familyId, albumId, uploadBatchID)
	if err != nil {
		return nil, err
	}

	// 리사이즈, 얼굴 인식, 업로드를 백그라운드에서 처리
	go func() {
		s.processImageInBackground(context.Background(), &bgParams{
			data:        originalData,
			meta:        meta,
			uploadPath:  uploadPath,
			mediaItemID: mediaItem.ID,
			familyId:    familyId,
			albumId:     albumId,
			retry:       retry,
		})
	}()

	return &UploadImageResult{MediaItemID: mediaItem.ID, Status: "uploaded"}, nil
}

// 리사이즈, 얼굴 인식, 업로드 백그라운드에서 처리
func (s *MediaItemService) processImageInBackground(ctx context.Context, p *bgParams) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("백그라운드 처리 중 패닉 발생:", r)
			s.UpdateMediaItemUploadStatus(ctx, p.mediaItemID, enums.UploadStatusFailed)
		}
	}()

	// 1. 리사이즈 (원본 데이터가 필요한 모든 작업을 먼저 수행)
	viewImg, err := imagepkg.Resize(p.data, consts.VIEW_MAX_LENGTH)
	if err != nil {
		log.Println("뷰 이미지 생성 실패:", err)
		s.UpdateMediaItemUploadStatus(ctx, p.mediaItemID, enums.UploadStatusFailed)
		return
	}

	thumbImg, err := imagepkg.Resize(p.data, consts.THUMBNAIL_MAX_LENGTH)
	if err != nil {
		log.Println("썸네일 이미지 생성 실패:", err)
		s.UpdateMediaItemUploadStatus(ctx, p.mediaItemID, enums.UploadStatusFailed)
		return
	}

	p.data = nil // 원본 bytes 즉시 해제 (GC 대상)

	// 2. 얼굴 인식
	faceDetections, err := s.faceService.GetFaceDetection(viewImg.Data)
	if err != nil {
		log.Println("얼굴 인식 실패:", err)
	}

	// 3. 중복 검사
	if !p.retry && faceDetections != nil {
		if err = s.faceService.CheckImageDuplicateByFaceDetection(ctx, p.familyId, p.albumId, *faceDetections); err != nil {
			log.Println("중복 사진 발견:", err)
			s.UpdateMediaItemUploadStatus(ctx, p.mediaItemID, enums.UploadStatusDuplicate)
			return
		}
	}

	// 4. 얼굴 크롭 + 저장 (view bytes가 유효한 동안 처리)
	if faceDetections != nil {
		createdDetections, err := s.faceService.SearchAndSaveFaceDetections(ctx, p.familyId, p.mediaItemID, *faceDetections)
		if err != nil {
			log.Println("얼굴 인식 결과 저장 실패:", err)
		} else {
			s.saveFaceCropImages(ctx, viewImg.Data, viewImg.Width, viewImg.Height, createdDetections, p.mediaItemID)
		}
	}

	// 5. view 업로드 후 bytes 해제
	viewResult, err := s.imageUploader.UploadResized(viewImg, p.uploadPath, consts.VIEW_STORAGE_PREFIX, p.meta.TakenAt)
	viewImg.Data = nil
	if err != nil {
		log.Println("뷰 이미지 업로드 실패:", err)
		s.UpdateMediaItemUploadStatus(ctx, p.mediaItemID, enums.UploadStatusFailed)
		return
	}

	// 6. thumbnail 업로드 후 bytes 해제
	thumbResult, err := s.imageUploader.UploadResized(thumbImg, p.uploadPath, consts.THUMBNAIL_STORAGE_PREFIX, p.meta.TakenAt)
	thumbImg.Data = nil
	if err != nil {
		log.Println("썸네일 이미지 업로드 실패:", err)
		s.UpdateMediaItemUploadStatus(ctx, p.mediaItemID, enums.UploadStatusFailed)
		return
	}

	// 7. DB에 리사이즈 파일 레코드 생성
	if err = s.saveResizedMediaFiles(ctx, p.mediaItemID, thumbResult, viewResult); err != nil {
		log.Println("미디어 파일 저장 실패:", err)
		s.UpdateMediaItemUploadStatus(ctx, p.mediaItemID, enums.UploadStatusFailed)
		return
	}

	s.UpdateMediaItemUploadStatus(ctx, p.mediaItemID, enums.UploadStatusCompleted)
}

// media_item과 원본 media_file 레코드를 생성
func (s *MediaItemService) createMediaItemWithOriginal(
	ctx context.Context,
	meta imagepkg.Meta,
	originalKey, filename string,
	familyId, albumId, uploadBatchID int32,
) (*db.MediaItem, error) {
	var nullLat, nullLon sql.NullFloat64
	if meta.Lat != nil {
		nullLat = sql.NullFloat64{Float64: *meta.Lat, Valid: true}
	}
	if meta.Lon != nil {
		nullLon = sql.NullFloat64{Float64: *meta.Lon, Valid: true}
	}

	mediaItem, err := s.mediaItemStore.CreateMediaItem(ctx, db.CreateMediaItemParams{
		FamilyID:               familyId,
		AlbumID:                albumId,
		UploadBatchID:          uploadBatchID,
		TakenAt:                meta.TakenAt,
		FileName:               sql.NullString{String: filename, Valid: true},
		TakenLocationLatitude:  nullLat,
		TakenLocationLongitude: nullLon,
	})
	if err != nil {
		return nil, err
	}

	_, err = s.mediaItemStore.CreateMediaFile(ctx, db.CreateMediaFileParams{
		MediaItemID: mediaItem.ID,
		Role:        string(enums.MediaItemRoleOriginal),
		StorageKey:  originalKey,
		Width:       sql.NullInt32{Int32: int32(meta.Width), Valid: true},
		Height:      sql.NullInt32{Int32: int32(meta.Height), Valid: true},
		FileSize:    sql.NullInt64{Int64: meta.FileSize, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &mediaItem, nil
}

// 썸네일과 뷰 media_file 레코드를 생성
func (s *MediaItemService) saveResizedMediaFiles(
	ctx context.Context,
	mediaItemID int32,
	thumb, view UploadResult,
) error {
	_, err := s.mediaItemStore.CreateMediaFile(ctx, db.CreateMediaFileParams{
		MediaItemID: mediaItemID,
		Role:        string(enums.MediaItemRoleThumbnail),
		StorageKey:  thumb.StorageKey,
		Width:       sql.NullInt32{Int32: int32(thumb.Width), Valid: true},
		Height:      sql.NullInt32{Int32: int32(thumb.Height), Valid: true},
		FileSize:    sql.NullInt64{Int64: thumb.FileSize, Valid: true},
	})
	if err != nil {
		return err
	}

	_, err = s.mediaItemStore.CreateMediaFile(ctx, db.CreateMediaFileParams{
		MediaItemID: mediaItemID,
		Role:        string(enums.MediaItemRoleView),
		StorageKey:  view.StorageKey,
		Width:       sql.NullInt32{Int32: int32(view.Width), Valid: true},
		Height:      sql.NullInt32{Int32: int32(view.Height), Valid: true},
		FileSize:    sql.NullInt64{Int64: view.FileSize, Valid: true},
	})
	return err
}

// 감지된 각 얼굴을 뷰 이미지에서 크롭하여 저장
// 이미 얼굴 이미지가 있는 identity는 건너뜀
func (s *MediaItemService) saveFaceCropImages(
	ctx context.Context,
	viewData []byte, imgWidth, imgHeight int,
	detections []db.FaceDetection,
	mediaItemID int32,
) {
	for _, fd := range detections {
		hasImg, err := s.identityFaceImgStore.HasIdentityFaceImg(ctx, fd.IdentityID)
		if err != nil || hasImg {
			continue
		}

		storageKey, err := s.imageUploader.CropFaceAndUpload(
			viewData, imgWidth, imgHeight,
			fd.LocationTop, fd.LocationRight, fd.LocationBottom, fd.LocationLeft,
			0.3, fd.IdentityID, mediaItemID,
		)
		if err != nil {
			log.Println("얼굴 크롭 이미지 업로드 실패:", err)
			continue
		}

		_, err = s.identityFaceImgStore.CreateIdentityFaceImg(ctx, db.CreateIdentityFaceImgParams{
			IdentityID:  fd.IdentityID,
			MediaItemID: mediaItemID,
			StorageKey:  storageKey,
		})
		if err != nil {
			log.Println("identity_face_imgs 저장 실패:", err)
		}
	}
}

// 미디어 아이템의 업로드 상태를 업데이트
func (s *MediaItemService) UpdateMediaItemUploadStatus(ctx context.Context, mediaItemID int32, uploadStatus enums.UploadStatus) error {
	_, err := s.mediaItemStore.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           mediaItemID,
		UploadStatus: string(uploadStatus),
	})
	return err
}

// GetMediaItemsByTakenAt는 촬영 시간 범위 내의 미디어 아이템을 반환합니다.
func (s *MediaItemService) GetMediaItemsByTakenAt(ctx context.Context, params *db.GetMediaItemsByTakenAtParams) ([]db.GetMediaItemsByTakenAtRow, error) {
	mediaItems, err := s.mediaItemStore.GetMediaItemsByTakenAt(ctx, *params)
	if err != nil {
		log.Println("get media items by taken at error:", err)
		return nil, err
	}
	return mediaItems, nil
}

// MediaItemRange는 미디어 아이템이 존재하는 연/월 조합입니다.
type MediaItemRange struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

// GetMediaItemRange는 미디어 아이템이 있는 연/월 목록을 반환합니다.
func (s *MediaItemService) GetMediaItemRange(ctx context.Context, groupId int32) ([]MediaItemRange, error) {
	ranges, err := s.mediaItemStore.GetMediaItemRange(ctx, groupId)
	if err != nil {
		return nil, err
	}

	mediaItemRanges := make([]MediaItemRange, len(ranges))
	for i, r := range ranges {
		year, err := strconv.Atoi(r.Year)
		if err != nil {
			return nil, err
		}
		month, err := strconv.Atoi(r.Month)
		if err != nil {
			return nil, err
		}
		mediaItemRanges[i] = MediaItemRange{Year: year, Month: month}
	}

	return mediaItemRanges, nil
}

// CreateUploadBatch는 앨범에 대한 새 업로드 배치를 생성합니다.
func (s *MediaItemService) CreateUploadBatch(ctx context.Context, familyId, albumId int32) (*db.UploadBatch, error) {
	uploadBatch, err := s.mediaItemStore.CreateUploadBatch(ctx, albumId)
	if err != nil {
		return nil, err
	}
	return &uploadBatch, nil
}

// GetUploadStatuses는 배치 내 모든 아이템의 업로드 상태를 반환합니다.
func (s *MediaItemService) GetUploadStatuses(ctx context.Context, uploadBatchID int32) ([]db.GetUploadStatusesRow, error) {
	return s.mediaItemStore.GetUploadStatuses(ctx, uploadBatchID)
}
