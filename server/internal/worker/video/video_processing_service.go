package video

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/ywl0806/yuno_kiroku/internal/consts"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/enums"
	"github.com/ywl0806/yuno_kiroku/internal/services"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	"github.com/ywl0806/yuno_kiroku/pkg/storage"
)

type VideoProcessingService struct {
	mediaItemStore store.MediaItemStore
	storage        storage.StorageService
}

func NewVideoProcessingService(
	mediaItemStore store.MediaItemStore,
	storageService storage.StorageService,
) *VideoProcessingService {
	return &VideoProcessingService{
		mediaItemStore: mediaItemStore,
		storage:        storageService,
	}
}

// ProcessVideo SQS job 파라미터를 받아 비디오 처리 전체 파이프라인을 실행합니다.
func (s *VideoProcessingService) ProcessVideo(ctx context.Context, params services.VideoJobParams) error {
	tmpDir, err := os.MkdirTemp("", "yuno-video-"+params.MediaItemID[:8]+"-*")
	if err != nil {
		return fmt.Errorf("temp dir 생성 실패: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Printf("[video-worker] temp dir cleanup 실패 (무시): %v", err)
		}
	}()

	inputPath := filepath.Join(tmpDir, "input.mp4")
	thumbPath := filepath.Join(tmpDir, "thumb.webp")
	outputPath := filepath.Join(tmpDir, "output.mp4")

	// 1. processing 상태로 변경
	if _, err = s.mediaItemStore.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           params.MediaItemID,
		UploadStatus: string(enums.UploadStatusProcessing),
	}); err != nil {
		return fmt.Errorf("status processing 업데이트 실패: %w", err)
	}

	// 2. S3 스트리밍 다운로드
	log.Printf("[video-worker] 다운로드 시작: %s", params.OriginalStorageKey)
	if err = s.storage.DownloadToFile(ctx, params.OriginalStorageKey, inputPath); err != nil {
		s.setFailed(ctx, params.MediaItemID, err)
		return fmt.Errorf("원본 비디오 다운로드 실패: %w", err)
	}

	// 3. 비디오 리사이즈 (최대 1280px, H.264, faststart)
	log.Printf("[video-worker] 비디오 리사이즈 중: %s", params.MediaItemID)
	if err = resizeVideo(ctx, inputPath, outputPath); err != nil {
		s.setFailed(ctx, params.MediaItemID, err)
		return fmt.Errorf("비디오 리사이즈 실패: %w", err)
	}
	videoW, videoH, err := probeSize(outputPath)
	if err != nil {
		log.Printf("[video-worker] 비디오 치수 조회 실패 (무시): %v", err)
	}

	// 4. 썸네일 추출 (리사이즈된 영상의 1초 지점, WebP)
	log.Printf("[video-worker] 썸네일 추출 중: %s", params.MediaItemID)
	if err = extractThumbnail(ctx, outputPath, thumbPath); err != nil {
		s.setFailed(ctx, params.MediaItemID, err)
		return fmt.Errorf("썸네일 추출 실패: %w", err)
	}
	thumbW, thumbH := videoW, videoH

	// 5. 썸네일 S3 업로드
	thumbKey := consts.THUMBNAIL_STORAGE_PREFIX + "/" + params.FamilyID + "/" + params.MediaItemID + ".webp"
	if err = s.storage.UploadFromFile(ctx, thumbKey, "image/webp", thumbPath); err != nil {
		s.setFailed(ctx, params.MediaItemID, err)
		return fmt.Errorf("썸네일 업로드 실패: %w", err)
	}

	// 6. 처리된 비디오 S3 업로드
	videoKey := "video/" + params.FamilyID + "/" + params.MediaItemID + ".mp4"
	if err = s.storage.UploadFromFile(ctx, videoKey, "video/mp4", outputPath); err != nil {
		s.setFailed(ctx, params.MediaItemID, err)
		return fmt.Errorf("비디오 업로드 실패: %w", err)
	}

	// 7. media_files DB 저장 (thumbnail + video)
	if _, err = s.mediaItemStore.CreateMediaFile(ctx, db.CreateMediaFileParams{
		MediaItemID: params.MediaItemID,
		Role:        string(enums.MediaItemRoleThumbnail),
		StorageKey:  thumbKey,
		MimeType:    sql.NullString{String: "image/webp", Valid: true},
		Width:       sql.NullInt32{Int32: thumbW, Valid: thumbW > 0},
		Height:      sql.NullInt32{Int32: thumbH, Valid: thumbH > 0},
	}); err != nil {
		s.setFailed(ctx, params.MediaItemID, err)
		return fmt.Errorf("thumbnail media_file 생성 실패: %w", err)
	}
	if _, err = s.mediaItemStore.CreateMediaFile(ctx, db.CreateMediaFileParams{
		MediaItemID: params.MediaItemID,
		Role:        string(enums.MediaItemRoleVideo),
		StorageKey:  videoKey,
		MimeType:    sql.NullString{String: "video/mp4", Valid: true},
		Width:       sql.NullInt32{Int32: videoW, Valid: videoW > 0},
		Height:      sql.NullInt32{Int32: videoH, Valid: videoH > 0},
	}); err != nil {
		s.setFailed(ctx, params.MediaItemID, err)
		return fmt.Errorf("video media_file 생성 실패: %w", err)
	}

	// 8. upload_status = completed
	if _, err = s.mediaItemStore.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           params.MediaItemID,
		UploadStatus: string(enums.UploadStatusCompleted),
	}); err != nil {
		log.Printf("[video-worker] status completed 업데이트 실패 (무시): %v", err)
	}

	log.Printf("[video-worker] 처리 완료: media_item_id=%s", params.MediaItemID)
	return nil
}

func (s *VideoProcessingService) setFailed(ctx context.Context, mediaItemID string, cause error) {
	if _, err := s.mediaItemStore.UpdateMediaItemUploadStatus(ctx, db.UpdateMediaItemUploadStatusParams{
		ID:           mediaItemID,
		UploadStatus: string(enums.UploadStatusFailed),
	}); err != nil {
		log.Printf("[video-worker] status failed 업데이트 실패: %v", err)
	}
	log.Printf("[video-worker] 처리 실패: media_item_id=%s cause=%v", mediaItemID, cause)
}

// extractThumbnail 1초 지점에서 WebP 썸네일 1장을 추출합니다.
func extractThumbnail(ctx context.Context, inputPath, outputPath string) error {
	tctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	s := ffmpeg.Input(inputPath, ffmpeg.KwArgs{"ss": "1"})
	s.Context = tctx
	return s.Output(outputPath, ffmpeg.KwArgs{
		"vframes": 1,
		"vcodec":  "libwebp",
		"q:v":     "75",
	}).
		OverWriteOutput().
		ErrorToStdOut().
		Run()
}

// probeSize ffprobe로 미디어 파일의 첫 번째 비디오 스트림 치수를 반환합니다.
func probeSize(path string) (width, height int32, err error) {
	data, err := ffmpeg.Probe(path)
	if err != nil {
		return 0, 0, fmt.Errorf("ffprobe 실패: %w", err)
	}
	var result struct {
		Streams []struct {
			Width  int32 `json:"width"`
			Height int32 `json:"height"`
		} `json:"streams"`
	}
	if err = json.Unmarshal([]byte(data), &result); err != nil {
		return 0, 0, fmt.Errorf("ffprobe JSON 파싱 실패: %w", err)
	}
	for _, s := range result.Streams {
		if s.Width > 0 && s.Height > 0 {
			return s.Width, s.Height, nil
		}
	}
	return 0, 0, fmt.Errorf("스트림에서 치수를 찾을 수 없음")
}

// resizeVideo 가로 최대 1280px(짝수 맞춤), H.264 CRF23, faststart로 재인코딩합니다.
func resizeVideo(ctx context.Context, inputPath, outputPath string) error {
	tctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	s := ffmpeg.Input(inputPath)
	s.Context = tctx
	return s.Output(outputPath, ffmpeg.KwArgs{
		"vf":       "scale='min(1280,iw)':'-2',format=yuv420p",
		"vcodec":   "libx264",
		"crf":      "23",
		"preset":   "ultrafast",
		"acodec":   "aac",
		"movflags": "+faststart",
	}).
		OverWriteOutput().
		ErrorToStdOut().
		Run()
}
