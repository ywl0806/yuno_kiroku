package models

import (
	"time"

	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/enums"
	"github.com/ywl0806/yuno_kiroku/internal/utils"
)

type MediaFileResponse struct {
	ID          int32  `json:"id"`
	MediaItemID int32  `json:"media_item_id"`
	Role        string `json:"role"`
	Url         string `json:"url"`
	MimeType    string `json:"mime_type"`
	Width       int32  `json:"width"`
	Height      int32  `json:"height"`
}

type MediaItemResponse struct {
	ID              int32     `json:"id"`
	GroupID         int32     `json:"group_id"`
	AlbumID         int32     `json:"album_id"`
	TakenAt         time.Time `json:"taken_at"`
	FileName        string    `json:"file_name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	OriginalUrl     string    `json:"original_url"`
	OriginalWidth   int32     `json:"original_width"`
	OriginalHeight  int32     `json:"original_height"`
	ThumbnailUrl    string    `json:"thumbnail_url"`
	ThumbnailWidth  int32     `json:"thumbnail_width"`
	ThumbnailHeight int32     `json:"thumbnail_height"`
	ViewUrl         string    `json:"view_url"`
	ViewWidth       int32     `json:"view_width"`
	ViewHeight      int32     `json:"view_height"`
}

func NewMediaItemResponse(mediaItem *db.GetMediaItemsByTakenAtRow) *MediaItemResponse {
	storageUrl := viper.GetString("STORAGE_URL")

	return &MediaItemResponse{
		ID:              mediaItem.ID,
		GroupID:         mediaItem.GroupID,
		AlbumID:         mediaItem.AlbumID,
		TakenAt:         mediaItem.TakenAt,
		FileName:        mediaItem.FileName.String,
		CreatedAt:       mediaItem.CreatedAt,
		UpdatedAt:       mediaItem.UpdatedAt,
		OriginalUrl:     utils.ParsePath(storageUrl, mediaItem.OriginalStorageKey),
		OriginalWidth:   mediaItem.OriginalWidth,
		OriginalHeight:  mediaItem.OriginalHeight,
		ThumbnailUrl:    utils.ParsePath(storageUrl, mediaItem.ThumbnailStorageKey),
		ThumbnailWidth:  mediaItem.ThumbnailWidth,
		ThumbnailHeight: mediaItem.ThumbnailHeight,
		ViewUrl:         utils.ParsePath(storageUrl, mediaItem.ViewStorageKey),
		ViewWidth:       mediaItem.ViewWidth,
		ViewHeight:      mediaItem.ViewHeight,
	}
}

func NewMediaItemsResponse(mediaItems []db.GetMediaItemsByTakenAtRow) *[]MediaItemResponse {
	mediaItemResponses := make([]MediaItemResponse, len(mediaItems))
	for i, mediaItem := range mediaItems {
		mediaItemResponses[i] = *NewMediaItemResponse(&mediaItem)
	}
	return &mediaItemResponses
}

type UploadBatchStatus struct {
	ID           int32  `json:"id"`
	UploadStatus string `json:"upload_status"`
}
type UploadBatchStatusResponse struct {
	Statuses    []UploadBatchStatus `json:"statuses"`
	IsCompleted bool                `json:"is_completed"`
}

func NewUploadBatchStatusResponse(uploadStatuses []db.GetUploadStatusesRow) *UploadBatchStatusResponse {
	statuses := make([]UploadBatchStatus, len(uploadStatuses))
	isCompleted := true
	for i, uploadStatus := range uploadStatuses {
		if uploadStatus.UploadStatus == string(enums.UploadStatusProcessing) || uploadStatus.UploadStatus == string(enums.UploadStatusPending) {
			isCompleted = false
		}
		statuses[i] = UploadBatchStatus{
			ID:           uploadStatus.ID,
			UploadStatus: uploadStatus.UploadStatus,
		}
	}
	return &UploadBatchStatusResponse{
		Statuses:    statuses,
		IsCompleted: isCompleted,
	}
}
