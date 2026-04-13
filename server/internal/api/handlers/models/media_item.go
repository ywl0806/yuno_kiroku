package models

import (
	"time"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/enums"
	"github.com/ywl0806/yuno_kiroku/internal/services"
	internalutils "github.com/ywl0806/yuno_kiroku/internal/utils"
)

// PresignedUploadRequest Presigned URL 발급 요청
type PresignedUploadRequest struct {
	FileName    string `json:"file_name" validate:"required"`
	ContentType string `json:"content_type" validate:"required"`
}

// PresignedUploadResponse Presigned URL 발급 응답
type PresignedUploadResponse struct {
	MediaItemID  int32  `json:"media_item_id"`
	PresignedURL string `json:"presigned_url"`
	StorageKey   string `json:"storage_key"`
	ExpiresIn    int    `json:"expires_in"`
}

// CreateUploadBatchResponse 업로드 배치 생성 응답
type CreateUploadBatchResponse struct {
	ID        int32     `json:"id"`
	AlbumID   int32     `json:"album_id"`
	CreatedAt time.Time `json:"created_at"`
}

// UploadImageResponse 이미지 업로드 응답
type UploadImageResponse struct {
	MediaItemID int32  `json:"media_item_id"`
	Status      string `json:"status"`
}

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
	FamilyID        int32     `json:"family_id"`
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

	return &MediaItemResponse{
		ID:              mediaItem.ID,
		FamilyID:        mediaItem.FamilyID,
		AlbumID:         mediaItem.AlbumID,
		TakenAt:         mediaItem.TakenAt,
		FileName:        mediaItem.FileName.String,
		CreatedAt:       mediaItem.CreatedAt,
		UpdatedAt:       mediaItem.UpdatedAt,
		OriginalUrl:     internalutils.ParseStoragePath(mediaItem.OriginalStorageKey),
		OriginalWidth:   mediaItem.OriginalWidth,
		OriginalHeight:  mediaItem.OriginalHeight,
		ThumbnailUrl:    internalutils.ParseStoragePath(mediaItem.ThumbnailStorageKey),
		ThumbnailWidth:  mediaItem.ThumbnailWidth,
		ThumbnailHeight: mediaItem.ThumbnailHeight,
		ViewUrl:         internalutils.ParseStoragePath(mediaItem.ViewStorageKey),
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

func NewMediaItemResponseFromHome(mediaItem *db.GetMediaItemsByTakenAtHomeRow) *MediaItemResponse {
	return &MediaItemResponse{
		ID:              mediaItem.ID,
		FamilyID:        mediaItem.FamilyID,
		AlbumID:         mediaItem.AlbumID,
		TakenAt:         mediaItem.TakenAt,
		FileName:        mediaItem.FileName.String,
		CreatedAt:       mediaItem.CreatedAt,
		UpdatedAt:       mediaItem.UpdatedAt,
		OriginalUrl:     internalutils.ParseStoragePath(mediaItem.OriginalStorageKey),
		OriginalWidth:   mediaItem.OriginalWidth,
		OriginalHeight:  mediaItem.OriginalHeight,
		ThumbnailUrl:    internalutils.ParseStoragePath(mediaItem.ThumbnailStorageKey),
		ThumbnailWidth:  mediaItem.ThumbnailWidth,
		ThumbnailHeight: mediaItem.ThumbnailHeight,
		ViewUrl:         internalutils.ParseStoragePath(mediaItem.ViewStorageKey),
		ViewWidth:       mediaItem.ViewWidth,
		ViewHeight:      mediaItem.ViewHeight,
	}
}

func NewMediaItemsResponseFromHome(mediaItems []db.GetMediaItemsByTakenAtHomeRow) *[]MediaItemResponse {
	mediaItemResponses := make([]MediaItemResponse, len(mediaItems))
	for i, mediaItem := range mediaItems {
		mediaItemResponses[i] = *NewMediaItemResponseFromHome(&mediaItem)
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

type BatchThumbnailResponse struct {
	ID              int32  `json:"id"`
	ThumbnailUrl    string `json:"thumbnail_url"`
	ThumbnailWidth  int32  `json:"thumbnail_width"`
	ThumbnailHeight int32  `json:"thumbnail_height"`
}

type UploadBatchWithThumbnailsResponse struct {
	ID         int32                    `json:"id"`
	AlbumID    int32                    `json:"album_id"`
	UploadAt   time.Time                `json:"upload_at"`
	Count      int64                    `json:"count"`
	Thumbnails []BatchThumbnailResponse `json:"thumbnails"`
}

type GetUploadBatchesResponse struct {
	Items   []UploadBatchWithThumbnailsResponse `json:"items"`
	HasNext bool                                `json:"has_next"`
	Page    int                                 `json:"page"`
}

func NewUploadBatchWithThumbnailsResponse(batch services.UploadBatchWithThumbnails) UploadBatchWithThumbnailsResponse {
	thumbnails := make([]BatchThumbnailResponse, len(batch.Thumbnails))
	for i, t := range batch.Thumbnails {
		thumbnails[i] = BatchThumbnailResponse{
			ID:              t.ID,
			ThumbnailUrl:    t.ThumbnailUrl,
			ThumbnailWidth:  t.ThumbnailWidth,
			ThumbnailHeight: t.ThumbnailHeight,
		}
	}
	return UploadBatchWithThumbnailsResponse{
		ID:         batch.ID,
		AlbumID:    batch.AlbumID,
		UploadAt:   batch.UploadAt,
		Count:      batch.Count,
		Thumbnails: thumbnails,
	}
}

func NewUploadBatchItemResponse(item *db.GetMediaItemsByUploadBatchIdRow) *MediaItemResponse {
	return &MediaItemResponse{
		ID:              item.ID,
		FamilyID:        item.FamilyID,
		AlbumID:         item.AlbumID,
		TakenAt:         item.TakenAt,
		FileName:        item.FileName.String,
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
		OriginalUrl:     internalutils.ParseStoragePath(item.OriginalStorageKey),
		OriginalWidth:   item.OriginalWidth,
		OriginalHeight:  item.OriginalHeight,
		ThumbnailUrl:    internalutils.ParseStoragePath(item.ThumbnailStorageKey),
		ThumbnailWidth:  item.ThumbnailWidth,
		ThumbnailHeight: item.ThumbnailHeight,
		ViewUrl:         internalutils.ParseStoragePath(item.ViewStorageKey),
		ViewWidth:       item.ViewWidth,
		ViewHeight:      item.ViewHeight,
	}
}

type GetMediaItemsRequest struct {
	From        *time.Time `query:"from" validate:"required"`
	To          *time.Time `query:"to" validate:"required"`
	IdentityIDs []int32    `query:"identity_ids"`
	AlbumID     *int32     `query:"album_id"`
}

type SearchMediaItemsRequest struct {
	From        *time.Time `query:"from"`
	To          *time.Time `query:"to"`
	IdentityIDs []int32    `query:"identity_ids"`
	AlbumID     *int32     `query:"album_id"`
	Page        int        `query:"page"`
}

type SearchMediaItemsResponse struct {
	Items   []MediaItemResponse `json:"items"`
	HasNext bool                `json:"has_next"`
	Page    int                 `json:"page"`
}

func NewSearchMediaItemResponse(item *db.SearchMediaItemsRow) *MediaItemResponse {
	return &MediaItemResponse{
		ID:              item.ID,
		FamilyID:        item.FamilyID,
		AlbumID:         item.AlbumID,
		TakenAt:         item.TakenAt,
		FileName:        item.FileName.String,
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
		OriginalUrl:     internalutils.ParseStoragePath(item.OriginalStorageKey),
		OriginalWidth:   item.OriginalWidth,
		OriginalHeight:  item.OriginalHeight,
		ThumbnailUrl:    internalutils.ParseStoragePath(item.ThumbnailStorageKey),
		ThumbnailWidth:  item.ThumbnailWidth,
		ThumbnailHeight: item.ThumbnailHeight,
		ViewUrl:         internalutils.ParseStoragePath(item.ViewStorageKey),
		ViewWidth:       item.ViewWidth,
		ViewHeight:      item.ViewHeight,
	}
}
