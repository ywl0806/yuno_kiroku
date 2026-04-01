package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/url"

	"github.com/labstack/echo/v4"
	workerServices "github.com/ywl0806/yuno_kiroku/internal/worker/services"
)

type minioEvent struct {
	EventName string `json:"EventName"`
	Key       string `json:"Key"`
	Records   []struct {
		S3 struct {
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

type ResizeHandler struct {
	resizeService *workerServices.ResizeService
}

func NewResizeHandler(resizeService *workerServices.ResizeService) *ResizeHandler {
	return &ResizeHandler{resizeService: resizeService}
}

// HandleMinioEvent POST /resize — MinIO bucket notification webhook 수신
func (h *ResizeHandler) HandleMinioEvent(c echo.Context) error {
	var event minioEvent
	if err := json.NewDecoder(c.Request().Body).Decode(&event); err != nil {
		log.Printf("이벤트 파싱 실패: %v", err)
		return c.JSON(400, map[string]string{"error": "invalid event"})
	}
	if len(event.Records) == 0 {
		return c.JSON(200, map[string]string{"status": "no records"})
	}

	rawKey := event.Records[0].S3.Object.Key
	originalKey, err := url.QueryUnescape(rawKey)
	if err != nil {
		originalKey = rawKey
	}

	log.Printf("리사이즈 처리 시작: %s", originalKey)
	go func() {
		if err := h.resizeService.ProcessResize(context.Background(), originalKey); err != nil {
			log.Printf("리사이즈 처리 실패 [%s]: %v", originalKey, err)
		}
	}()

	return c.JSON(200, map[string]string{"status": "processing"})
}
