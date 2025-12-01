package photo

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

	"github.com/ywl0806/yuno_kiroku/internal/api/services/photo/models"
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

/*
*

	얼굴 인식 결과를 검색 후 저장
	1. 얼굴 임베딩을 데이터베이스에서 검색
	2. 얼굴 임베딩이 있으면 FaceDetection을 생성
	3. 얼굴 임베딩이 없으면 Person을 생성 후 FaceDetection을 생성
*/
func (s *FaceService) SearchAndSaveFaceDetections(ctx context.Context, groupId int32, photoId int32, faceDetections []models.FaceDetection) ([]db.FaceDetection, error) {

	// 얼굴인식 결과를 저장할 파라미터 리스트
	var createFaceDetectionParams []db.CreateFaceDetectionParams
	// 얼굴인식 결과리스트
	var createdFaceDetections []db.FaceDetection

	for _, faceDetection := range faceDetections {
		// 얼굴인식 결과를 저장할 파라미터 생성
		faceDetectionParams := db.CreateFaceDetectionParams{
			PhotoID:        photoId,
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
			faceDetectionParams.PersonID = similarFace.PersonID
		} else {
			// 유사 얼굴 검색 결과가 없으면 새로운 사람 생성
			newPerson, err := s.queries.CreatePerson(ctx, db.CreatePersonParams{
				Name:    sql.NullString{String: "", Valid: false},
				GroupID: groupId,
			})
			if err != nil {
				log.Println("CreatePerson error: ", err)
				continue
			}
			faceDetectionParams.PersonID = newPerson.ID
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
	log.Println("faceDetection distance: ", faceDetection.Distance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &faceDetection, nil
}

func (s *FaceService) GetFaceDetections(ctx context.Context, photoId int32) ([]db.GetFaceDetectionsByPhotoIdRow, error) {
	faceDetections, err := s.queries.GetFaceDetectionsByPhotoId(ctx, photoId)
	if err != nil {
		return nil, err
	}
	return faceDetections, nil
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
func GetFaceDetection(image io.Reader) (*[]models.FaceDetection, error) {

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
	copiedBytes, err := io.Copy(part, image)
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
