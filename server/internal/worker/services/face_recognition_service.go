package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/internal/consts"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	imagepkg "github.com/ywl0806/yuno_kiroku/pkg/image"
	"github.com/ywl0806/yuno_kiroku/pkg/storage"
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

// FaceCropInput 얼굴 크롭에 필요한 이미지 정보와 bbox 좌표
type FaceCropInput struct {
	ImgWidth  int
	ImgHeight int
	Top       int32
	Right     int32
	Bottom    int32
	Left      int32
	Padding   float64
}

// FaceResult Python batch에서 전달받는 단일 얼굴 결과
type FaceResult struct {
	FaceLocation FaceLocation `json:"face_location"`
	Embedding    []float64    `json:"embedding"`
}

// FaceRecognitionService Go CLI에서 호출되어 identity 매칭, face_detections 저장, 얼굴 크롭을 처리합니다.
type FaceRecognitionService struct {
	transactor      store.Transactor
	faceStore       store.FaceStore
	identityStore   store.IdentityStore
	identityFaceImg store.IdentityFaceImgStore
	storage         storage.StorageService
}

func NewFaceRecognitionService(
	transactor store.Transactor,
	faceStore store.FaceStore,
	identityStore store.IdentityStore,
	identityFaceImg store.IdentityFaceImgStore,
	storageService storage.StorageService,
) *FaceRecognitionService {
	return &FaceRecognitionService{
		transactor:      transactor,
		faceStore:       faceStore,
		identityStore:   identityStore,
		identityFaceImg: identityFaceImg,
		storage:         storageService,
	}
}

// matchOrCreateIdentity 얼굴 임베딩을 데이터베이스에서 검색 후 매칭된 identity가 없으면 새로운 identity를 생성합니다.
// 반드시 advisory lock이 획득된 트랜잭션 Store(tx) 위에서 호출해야 합니다.
func matchOrCreateIdentity(ctx context.Context, tx *store.Store, familyID string, embedding []float64) (int32, error) {
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

// ProcessFacesParams Go CLI에서 전달받는 얼굴 인식 처리 파라미터
type ProcessFacesParams struct {
	MediaItemID   string `json:"media_item_id"`
	FamilyID      string `json:"family_id"`
	ViewImagePath string `json:"view_image_path"` // Python이 다운로드한 로컬 임시 파일 경로
}

// ProcessFaces Go CLI 진입점에서 호출. view 이미지를 로컬 파일에서 읽어 얼굴 처리를 수행합니다.
func (s *FaceRecognitionService) ProcessFaces(ctx context.Context, params ProcessFacesParams, faces []FaceResult) error {
	var viewData []byte
	var viewWidth, viewHeight int

	for _, face := range faces {
		var identityID int32
		// advisory lock 획득 타임아웃 240초
		lockCtx, cancel := context.WithTimeout(ctx, 240*time.Second)
		err := s.transactor.TransactWithAdvisoryLock(lockCtx, advisoryLockNamespace<<32|utils.StringToHash(params.FamilyID), func(tx *store.Store) error {
			id, err := matchOrCreateIdentity(lockCtx, tx, params.FamilyID, face.Embedding)
			if err != nil {
				return err
			}
			identityID = id
			_, err = tx.Face.CreateFaceDetection(lockCtx, db.CreateFaceDetectionParams{
				MediaItemID:    params.MediaItemID,
				IdentityID:     identityID,
				LocationTop:    int32(face.FaceLocation.Top),
				LocationRight:  int32(face.FaceLocation.Right),
				LocationBottom: int32(face.FaceLocation.Bottom),
				LocationLeft:   int32(face.FaceLocation.Left),
				Embedding:      utils.Float64SliceToVectorString(face.Embedding),
			})
			return err
		})
		cancel()
		if err != nil {

			log.Printf("[ERROR] face 처리 실패 (SQS 재시도 유도): %v", err)
			return fmt.Errorf("face 처리 실패: %w", err)
		}

		if viewData == nil {
			viewData, err = os.ReadFile(params.ViewImagePath)
			if err != nil {
				log.Printf("view 이미지 파일 읽기 실패 (크롭 생략): %v", err)
				continue
			}
			viewWidth, viewHeight, err = imageDimensions(viewData)
			if err != nil {
				log.Printf("이미지 크기 파싱 실패 (크롭 생략): %v", err)
				viewData = nil
				continue
			}
		}

		cropInput := FaceCropInput{
			ImgWidth:  viewWidth,
			ImgHeight: viewHeight,
			Top:       int32(face.FaceLocation.Top),
			Right:     int32(face.FaceLocation.Right),
			Bottom:    int32(face.FaceLocation.Bottom),
			Left:      int32(face.FaceLocation.Left),
			Padding:   0.3,
		}
		croppedBytes, err := cropFace(viewData, cropInput)
		if err != nil {
			log.Printf("[ERROR] 얼굴 크롭 실패 (SQS 재시도 유도): %v", err)
			return fmt.Errorf("얼굴 크롭 실패: %w", err)
		}
		storageKey, err := s.uploadFaceImage(ctx, croppedBytes, params.FamilyID, identityID, params.MediaItemID)
		if err != nil {
			log.Printf("[ERROR] 얼굴 크롭 이미지 업로드 실패 (SQS 재시도 유도): %v", err)
			return fmt.Errorf("얼굴 크롭 이미지 업로드 실패: %w", err)
		}

		if _, err = s.identityFaceImg.CreateIdentityFaceImg(ctx, db.CreateIdentityFaceImgParams{
			IdentityID:  identityID,
			MediaItemID: params.MediaItemID,
			StorageKey:  storageKey,
		}); err != nil {
			log.Printf("[ERROR] identity_face_img DB 저장 실패, S3 파일 롤백: %v", err)
			if delErr := s.storage.DeleteFile(ctx, storageKey); delErr != nil {
				log.Printf("[ERROR] S3 롤백 실패 (고아 파일 발생 가능): key=%s, err=%v", storageKey, delErr)
			}
		}
	}

	return nil
}

// imageDimensions view 이미지 바이트에서 width/height를 파싱합니다.
func imageDimensions(data []byte) (width, height int, err error) {
	parsed, err := imagepkg.Parse(data, "")
	if err != nil {
		return 0, 0, err
	}
	return parsed.Width, parsed.Height, nil
}

// cropFace bbox에 패딩을 적용한 크롭 영역을 계산하고 크롭·리사이즈한 이미지 바이트를 반환합니다.
func cropFace(imageData []byte, in FaceCropInput) ([]byte, error) {
	faceW := int(in.Right - in.Left)
	faceH := int(in.Bottom - in.Top)

	padX := int(float64(faceW) * in.Padding)
	padY := int(float64(faceH) * in.Padding)

	if faceW > faceH {
		padY += int(float64(faceW-faceH)*in.Padding) + (faceW-faceH)/2
	} else {
		padX += int(float64(faceH-faceW)*in.Padding) + (faceH-faceW)/2
	}

	cropLeft := clampMin(int(in.Left)-padX, 0)
	cropTop := clampMin(int(in.Top)-padY, 0)
	cropRight := clampMax(int(in.Right)+padX, in.ImgWidth)
	cropBottom := clampMax(int(in.Bottom)+padY, in.ImgHeight)

	croppedBytes, err := imagepkg.CropAndResize(imageData, cropLeft, cropTop, cropRight-cropLeft, cropBottom-cropTop, 512)
	if err != nil {
		return nil, fmt.Errorf("얼굴 크롭 실패: %w", err)
	}
	return croppedBytes, nil
}

// uploadFaceImage 크롭된 얼굴 이미지를 스토리지에 업로드하고 storage key를 반환합니다.
func (s *FaceRecognitionService) uploadFaceImage(ctx context.Context, data []byte, familyID string, identityID int32, mediaItemID string) (string, error) {
	key := fmt.Sprintf("%s/%s/%d/face_%s.webp", familyID, consts.IDENTITY_STORAGE_PREFIX, identityID, mediaItemID)
	storageKey, err := s.storage.SaveFile(ctx, key, data)
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
