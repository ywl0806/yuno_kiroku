package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	workerServices "github.com/ywl0806/yuno_kiroku/internal/worker/services"
)

type faceRecognitionCompleteRequest struct {
	JobID int32                       `json:"job_id"`
	Faces []workerServices.FaceResult `json:"faces"`
}

type faceRecognitionFailRequest struct {
	JobID int32  `json:"job_id"`
	Error string `json:"error"`
}

type FaceRecognitionHandler struct {
	faceRecognitionService *workerServices.FaceRecognitionService
}

func NewFaceRecognitionHandler(faceRecognitionService *workerServices.FaceRecognitionService) *FaceRecognitionHandler {
	return &FaceRecognitionHandler{faceRecognitionService: faceRecognitionService}
}

// Fail POST /face-recognition/fail — Python batch가 처리 실패 시 job 상태 되돌림
func (h *FaceRecognitionHandler) Fail(c echo.Context) error {
	var req faceRecognitionFailRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.faceRecognitionService.Fail(c.Request().Context(), req.JobID, req.Error); err != nil {
		log.Printf("job 실패 처리 오류 [job_id=%d]: %v", req.JobID, err)
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "failed"})
}

// Complete POST /face-recognition/complete — Python batch가 임베딩 결과 전달
func (h *FaceRecognitionHandler) Complete(c echo.Context) error {
	var req faceRecognitionCompleteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	go func() {
		backgroundContext := context.Background()
		if err := h.faceRecognitionService.Complete(backgroundContext, req.JobID, req.Faces); err != nil {
			log.Printf("얼굴 인식 처리 실패 [job_id=%d]: %v", req.JobID, err)
		} else {
			log.Printf("얼굴 인식 처리 완료 [job_id=%d]", req.JobID)
		}
		backgroundContext.Done()
	}()
	return c.JSON(http.StatusOK, map[string]string{"status": "processing"})
}
