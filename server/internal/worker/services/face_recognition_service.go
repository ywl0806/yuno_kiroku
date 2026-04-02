package services

import (
	"context"
	"fmt"
	"log"

	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/internal/consts"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/services"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	imagepkg "github.com/ywl0806/yuno_kiroku/pkg/image"
	"github.com/ywl0806/yuno_kiroku/pkg/utils"
)

// advisoryLockNamespace family 단위 advisory lock에 사용하는 네임스페이스
// 다른 용도의 advisory lock과 충돌 방지를 위해 상위 32비트에 고정값을 둠
const advisoryLockNamespace = int64(0x59554E4F) // "YUNO"

// FaceLocation Python batch에서 전달받는 얼굴 위치 좌표
type FaceLocation struct {
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
	Left   int `json:"left"`
}

// FaceResult Python batch에서 전달받는 단일 얼굴 결과
type FaceResult struct {
	FaceLocation FaceLocation `json:"face_location"`
	Embedding    []float64    `json:"embedding"`
}

// FaceRecognitionService Python batch로부터 임베딩 결과를 받아
// identity 매칭, face_detections 저장, 얼굴 크롭을 처리합니다.
type FaceRecognitionService struct {
	transactor      store.Transactor
	faceJobStore    store.FaceRecognitionJobStore
	mediaItemStore  store.MediaItemStore
	faceStore       store.FaceStore
	identityStore   store.IdentityStore
	identityFaceImg store.IdentityFaceImgStore
	imageUploader   *services.ImageUploader
}

func NewFaceRecognitionService(
	transactor store.Transactor,
	faceJobStore store.FaceRecognitionJobStore,
	mediaItemStore store.MediaItemStore,
	faceStore store.FaceStore,
	identityStore store.IdentityStore,
	identityFaceImg store.IdentityFaceImgStore,
	imageUploader *services.ImageUploader,

) *FaceRecognitionService {
	return &FaceRecognitionService{
		transactor:      transactor,
		faceJobStore:    faceJobStore,
		mediaItemStore:  mediaItemStore,
		faceStore:       faceStore,
		identityStore:   identityStore,
		identityFaceImg: identityFaceImg,
		imageUploader:   imageUploader,
	}
}

// Fail Python batch가 처리 실패 시 job 상태를 되돌립니다.
func (s *FaceRecognitionService) Fail(ctx context.Context, jobID int32, errMsg string) error {
	return s.faceJobStore.FailFaceRecognitionJob(ctx, jobID, errMsg)
}

// Complete Python batch가 임베딩 계산 완료 후 호출합니다.
// identity 매칭 → face_detections INSERT → 얼굴 크롭 → job 완료 처리
func (s *FaceRecognitionService) Complete(ctx context.Context, jobID int32, faces []FaceResult) error {
	job, err := s.faceJobStore.GetFaceRecognitionJobByID(ctx, jobID)
	if err != nil {
		return err
	}

	if err := s.processFaces(ctx, job, faces); err != nil {
		s.faceJobStore.FailFaceRecognitionJob(ctx, jobID, err.Error())
		return err
	}

	return s.faceJobStore.CompleteFaceRecognitionJob(ctx, jobID)
}

func (s *FaceRecognitionService) processFaces(ctx context.Context, job db.FaceRecognitionJob, faces []FaceResult) error {
	// view 이미지 다운로드 (얼굴 크롭에 사용)
	var viewData []byte
	var viewWidth, viewHeight int

	for _, face := range faces {
		// advisory lock + identity 매칭/생성 + face_detection INSERT를 하나의 트랜잭션으로 처리
		// → 동일 family의 동시 요청이 와도 중복 identity가 생성되지 않음
		var identityID int32
		err := s.transactor.TransactWithAdvisoryLock(ctx, advisoryLockNamespace<<32|int64(job.FamilyID), func(tx *store.Store) error {
			id, err := matchOrCreateIdentity(ctx, tx, job.FamilyID, face.Embedding)
			if err != nil {
				return err
			}
			identityID = id

			_, err = tx.Face.CreateFaceDetection(ctx, db.CreateFaceDetectionParams{
				MediaItemID:    job.MediaItemID,
				IdentityID:     identityID,
				LocationTop:    int32(face.FaceLocation.Top),
				LocationRight:  int32(face.FaceLocation.Right),
				LocationBottom: int32(face.FaceLocation.Bottom),
				LocationLeft:   int32(face.FaceLocation.Left),
				Embedding:      utils.Float64SliceToVectorString(face.Embedding),
			})
			return err
		})
		if err != nil {
			log.Printf("face 처리 실패 (무시): %v", err)
			continue
		}

		// 최초 1회만 view 이미지 다운로드
		if viewData == nil {
			viewData, err = s.imageUploader.GetFile(ctx, job.ViewStorageKey)
			if err != nil {
				log.Printf("view 이미지 다운로드 실패 (크롭 생략): %v", err)
				continue
			}
			viewWidth, viewHeight, err = imageDimensions(viewData)
			if err != nil {
				log.Printf("이미지 크기 파싱 실패 (크롭 생략): %v", err)
				viewData = nil
				continue
			}
		}

		storageKey, err := s.CropFaceAndUpload(
			ctx,
			viewData, viewWidth, viewHeight,
			int32(face.FaceLocation.Top), int32(face.FaceLocation.Right),
			int32(face.FaceLocation.Bottom), int32(face.FaceLocation.Left),
			0.3, identityID, job.MediaItemID,
		)
		if err != nil {
			log.Printf("얼굴 크롭 실패 (무시): %v", err)
			continue
		}

		if _, err = s.identityFaceImg.CreateIdentityFaceImg(ctx, db.CreateIdentityFaceImgParams{
			IdentityID:  identityID,
			MediaItemID: job.MediaItemID,
			StorageKey:  storageKey,
		}); err != nil {
			log.Printf("identity_face_img 생성 실패 (무시): %v", err)
		}
	}

	return nil
}

// matchOrCreateIdentity 얼굴 임베딩을 데이터베이스에서 검색 후 매칭된 identity가 없으면 새로운 identity를 생성합니다.
// 반드시 advisory lock이 획득된 트랜잭션 Store(tx) 위에서 호출해야 합니다.
func matchOrCreateIdentity(ctx context.Context, tx *store.Store, familyID int32, embedding []float64) (int32, error) {
	similar, err := tx.Face.FindMostSimilarFace(ctx, db.FindMostSimilarFaceParams{
		FamilyID:            familyID,
		Embedding:           utils.Float64SliceToVectorString(embedding),
		SimilarityThreshold: consts.FACE_SEARCH_THRESHOLD,
	})
	log.Printf("similar: %+v", similar)
	log.Printf("err: %+v", err)
	if err != nil {
		if apperr.IsAppError(err, apperr.NotFound) {
			newIdentity, err := tx.Identity.CreateIdentity(ctx, familyID)
			if err != nil {
				return 0, err
			}
			return newIdentity.ID, nil
		}
		return 0, err
	}

	return similar.IdentityID, nil

}

// imageDimensions view 이미지 바이트에서 width/height를 파싱합니다.
func imageDimensions(data []byte) (width, height int, err error) {
	parsed, err := imagepkg.Parse(data, "")
	if err != nil {
		return 0, 0, err
	}
	return parsed.Width, parsed.Height, nil
}

// 얼굴 bbox에 패딩을 적용한 크롭 영역을 계산하고,
// cropper로 크롭·리사이즈한 뒤 스토리지에 업로드합니다.
func (s *FaceRecognitionService) CropFaceAndUpload(
	ctx context.Context,
	viewData []byte, imgWidth, imgHeight int,
	top, right, bottom, left int32, padding float64,
	identityID, mediaItemID int32,
) (string, error) {
	// bbox 패딩 계산
	faceW := int(right - left)
	faceH := int(bottom - top)

	padX := int(float64(faceW) * padding)
	padY := int(float64(faceH) * padding)

	if faceW > faceH {
		padY += int(float64(faceW-faceH)*padding) + (faceW-faceH)/2
	} else {
		padX += int(float64(faceH-faceW)*padding) + (faceH-faceW)/2
	}

	cropLeft := clampMin(int(left)-padX, 0)
	cropTop := clampMin(int(top)-padY, 0)
	cropRight := clampMax(int(right)+padX, imgWidth)
	cropBottom := clampMax(int(bottom)+padY, imgHeight)

	cropW := cropRight - cropLeft
	cropH := cropBottom - cropTop

	croppedBytes, err := imagepkg.CropAndResize(viewData, cropLeft, cropTop, cropW, cropH, 512)
	if err != nil {
		return "", fmt.Errorf("얼굴 크롭 실패: %w", err)
	}
	defer func() { croppedBytes = nil }()

	storageKey, err := s.imageUploader.SaveFile(
		ctx,
		croppedBytes,
		fmt.Sprintf("identities/%d", identityID),
		fmt.Sprintf("face_%d.webp", mediaItemID),
	)
	if err != nil {
		return "", fmt.Errorf("얼굴 크롭 이미지 업로드 실패: %w", err)
	}

	return storageKey, nil
}

func clampMin(v, min int) int {
	if v < min {
		return min
	}
	return v
}

func clampMax(v, max int) int {
	if v > max {
		return max
	}
	return v
}
