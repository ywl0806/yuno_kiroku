package services

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/ywl0806/yuno_kiroku/internal/api/appErrors"
	"github.com/ywl0806/yuno_kiroku/internal/api/handlers/models"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
	"github.com/ywl0806/yuno_kiroku/internal/db"
)

const (
	// 얼굴 유사도 임계값
	FACE_SEARCH_THRESHOLD = 0.5
)

type FaceService struct {
	queries *db.Queries
}

func NewFaceService(queries *db.Queries) *FaceService {
	return &FaceService{queries: queries}
}

// 얼굴 인식 결과를 검색 후 저장
//
// 1. 얼굴 임베딩을 데이터베이스에서 검색
// 2. 얼굴 임베딩이 있으면 FaceDetection을 생성
// 3. 얼굴 임베딩이 없으면 Identity를 생성 후 FaceDetection을 생성
func (s *FaceService) SearchAndSaveFaceDetections(ctx context.Context, groupId int32, mediaItemId int32, faceDetections []models.FaceDetection) ([]db.FaceDetection, error) {

	// 얼굴인식 결과를 저장할 파라미터 리스트
	var createFaceDetectionParams []db.CreateFaceDetectionParams
	// 얼굴인식 결과리스트
	var createdFaceDetections []db.FaceDetection

	for _, faceDetection := range faceDetections {
		// 얼굴인식 결과를 저장할 파라미터 생성
		faceDetectionParams := db.CreateFaceDetectionParams{
			MediaItemID:    mediaItemId,
			LocationTop:    int32(faceDetection.FaceLocation.Top),
			LocationRight:  int32(faceDetection.FaceLocation.Right),
			LocationBottom: int32(faceDetection.FaceLocation.Bottom),
			LocationLeft:   int32(faceDetection.FaceLocation.Left),
			Embedding:      utils.Float64SliceToVectorString(faceDetection.Embedding),
		}
		// 얼굴인식 결과를 검색
		similarFace, err := s.FindMostSimilarFace(ctx, groupId, faceDetection.Embedding)

		// 검색 에러 시 로깅, 다음 얼굴인식 결과 처리
		if err != nil {
			log.Println("FindMostSimilarFace error: ", err)
			continue
		}

		// 유사 얼굴 검색 결과가 있으면 해당 사람 ID 설정
		if similarFace != nil {
			if similarFace.Distance == 0 {
				log.Println("Distance is 0")
				continue
			}

			faceDetectionParams.IdentityID = similarFace.IdentityID
		} else {
			// 유사 얼굴 검색 결과가 없으면 새로운 사람 생성
			newIdentity, err := s.queries.CreateIdentity(ctx, db.CreateIdentityParams{
				Name:    sql.NullString{String: "", Valid: false},
				GroupID: groupId,
			})
			if err != nil {
				log.Println("CreateIdentity error: ", err)
				continue
			}
			faceDetectionParams.IdentityID = newIdentity.ID
		}

		// 얼굴인식 결과를 저장할 파라미터 리스트에 추가
		createFaceDetectionParams = append(createFaceDetectionParams, faceDetectionParams)
	}

	// 얼굴인식 결과를 저장
	for _, faceDetectionParams := range createFaceDetectionParams {
		createdFaceDetection, err := s.queries.CreateFaceDetection(ctx, faceDetectionParams)
		if err != nil {
			log.Println("CreateFaceDetection error: ", err)
			continue
		}
		createdFaceDetections = append(createdFaceDetections, createdFaceDetection)
	}
	return createdFaceDetections, nil
}

/*
*

	얼굴 임베딩을 데이터베이스에서 검색
	1. 얼굴 임베딩이 없으면 nil 반환
*/
func (s *FaceService) FindMostSimilarFace(ctx context.Context, groupId int32, embedding []float64) (*db.FindMostSimilarFaceRow, error) {
	faceDetection, err := s.queries.FindMostSimilarFace(ctx, db.FindMostSimilarFaceParams{
		GroupID:             groupId,
		Embedding:           utils.Float64SliceToVectorString(embedding),
		SimilarityThreshold: FACE_SEARCH_THRESHOLD,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &faceDetection, nil
}

/*
*

	얼굴 인식
	1. 이미지를 얼굴 인식 API에 전송합니다.
	2. 얼굴 인식 API에서 얼굴 인식 결과를 반환합니다.

	Args:
		image: 이미지 파일
	Returns:
		*[]FaceDetection: 얼굴 인식 결과
		error: 얼굴 인식 실패시 에러
*/
func (s *FaceService) GetFaceDetection(image []byte) (*[]models.FaceDetection, error) {

	// API URL 없으면 에러
	aiApiUrl := os.Getenv("AI_API_URL")
	if aiApiUrl == "" {
		return nil, errors.New("AI_API_URL is not set")
	}

	// multipart/form-data 생성
	var bodyBuffer bytes.Buffer
	writer := multipart.NewWriter(&bodyBuffer)
	// 파일 파트 생성
	part, err := writer.CreateFormFile("file", "image.jpeg")
	if err != nil {
		return nil, err
	}
	// 파일 내용 복사
	copiedBytes, err := part.Write(image)
	if err != nil {
		return nil, err
	}
	log.Println("copied bytes to multipart: ", copiedBytes)
	writer.Close()
	log.Println("bodyBuffer size after close: ", bodyBuffer.Len())

	// 요청 생성
	request, err := http.NewRequest("POST", aiApiUrl+"/face-detection/file", &bodyBuffer)
	if err != nil {
		return nil, err
	}
	// 요청 헤더 설정
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// 클라이언트 생성
	client := &http.Client{}
	// 요청 전송
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// 응답 바디 읽기
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// 응답 바디 파싱
	var faceDetectionResponse models.FaceDetectionResponse

	// 응답 바디가 없으면 빈 얼굴 인식 결과 반환
	if response.StatusCode == http.StatusNoContent {
		return &[]models.FaceDetection{}, nil
	}

	err = json.Unmarshal(responseBody, &faceDetectionResponse)
	if err != nil {
		return nil, err
	}

	return &faceDetectionResponse.Faces, nil
}

//	얼굴 인식 결과로 사진이 중복되는지 확인합니다.
//	1. 얼굴 인식 결과를 벡터로 변환합니다.
//	2. 데이터베이스에서 얼굴 인식 결과와 일치하는 사진을 조회합니다.
//	3. 조회된 사진이 있으면 중복된 사진이 있다는 에러를 반환합니다.
//
// */
func (s *FaceService) CheckImageDuplicateByFaceDetection(ctx context.Context, groupId int32, albumId int32, faceDetections []models.FaceDetection) error {

	embeddings := make([]interface{}, len(faceDetections))
	for i, faceDetection := range faceDetections {
		embeddings[i] = utils.Float64SliceToVectorString(faceDetection.Embedding)
	}
	mediaItem, err := s.queries.GetMediaItemByFaceDetection(ctx, db.GetMediaItemByFaceDetectionParams{
		GroupID:    groupId,
		Embeddings: embeddings,
	})
	if err == sql.ErrNoRows {
		return nil
	}
	if mediaItem.ID != 0 {
		return appErrors.NewDuplicateError("")
	}

	return nil
}
