package photo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/ywl0806/yuno_kiroku/internal/api/services/photo/models"
	"github.com/ywl0806/yuno_kiroku/internal/db"
)

const (
	// 얼굴 유사도 임계값
	FACE_SEARCH_THRESHOLD = 0.6
)

type FaceService struct {
	queries *db.Queries
}

func NewFaceService(queries *db.Queries) *FaceService {
	return &FaceService{queries: queries}
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
	err = json.Unmarshal(responseBody, &faceDetectionResponse)
	if err != nil {
		return nil, err
	}

	return &faceDetectionResponse.Faces, nil
}
