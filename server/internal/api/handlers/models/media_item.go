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
	AlbumID       string `json:"album_id" validate:"required"`
	UploadBatchID int32  `json:"upload_batch_id" validate:"required"`
	FileName      string `json:"file_name" validate:"required"`
	ContentType   string `json:"content_type" validate:"required"`
}

// BatchPresignedUploadRequestItem 배치 Presigned URL 발급 요청 항목
type BatchPresignedUploadRequestItem struct {
	FileName    string `json:"file_name" validate:"required"`
	ContentType string `json:"content_type" validate:"required"`
}

// BatchPresignedUploadRequest 배치 Presigned URL 발급 요청
type BatchPresignedUploadRequest struct {
	AlbumID       string                            `json:"album_id" validate:"required"`
	UploadBatchID int32                             `json:"upload_batch_id" validate:"required"`
	Files         []BatchPresignedUploadRequestItem `json:"files" validate:"required,min=1,dive"`
}

// BatchPresignedUploadSuccessItem 배치 Presigned URL 발급 성공 항목
type BatchPresignedUploadSuccessItem struct {
	MediaItemID   string `json:"media_item_id"`
	PresignedURL  string `json:"presigned_url"`
	StorageKey    string `json:"storage_key"`
	ExpiresIn     int    `json:"expires_in"`
	OriginalIndex int    `json:"index"`
}

// BatchPresignedUploadFailedItem 배치 Presigned URL 발급 실패 항목
type BatchPresignedUploadFailedItem struct {
	FileName string `json:"file_name"`
	Index    int    `json:"index"`
	Reason   string `json:"reason,omitempty"`
}

// BatchPresignedUploadResponse 배치 Presigned URL 발급 응답 (부분 실패 지원)
type BatchPresignedUploadResponse struct {
	Success []BatchPresignedUploadSuccessItem `json:"success"`
	Failed  []BatchPresignedUploadFailedItem  `json:"failed"`
}

// VideoUploadCompleteRequest 비디오 업로드 완료 요청
type VideoUploadCompleteRequest struct {
	MediaItemID string `json:"media_item_id" validate:"required"`
	FileName    string `json:"file_name" validate:"required"`
	MimeType    string `json:"mime_type" validate:"required"`
}

// PresignedUploadResponse Presigned URL 발급 응답
type PresignedUploadResponse struct {
	MediaItemID  string `json:"media_item_id"`
	PresignedURL string `json:"presigned_url"`
	StorageKey   string `json:"storage_key"`
	ExpiresIn    int    `json:"expires_in"`
}

// CreateUploadBatchResponse 업로드 배치 생성 응답
type CreateUploadBatchResponse struct {
	ID        int32     `json:"id"`
	AlbumID   string    `json:"album_id"`
	CreatedAt time.Time `json:"created_at"`
}

// UploadImageResponse 이미지 업로드 응답
type UploadImageResponse struct {
	MediaItemID string `json:"media_item_id"`
	Status      string `json:"status"`
}

type MediaFileResponse struct {
	ID          int32  `json:"id"`
	MediaItemID string `json:"media_item_id"`
	Role        string `json:"role"`
	Url         string `json:"url"`
	Width       int32  `json:"width"`
	Height      int32  `json:"height"`
}

type TagResponse struct {
	ID       int32   `json:"id"`
	FamilyID *string `json:"family_id"`
	Name     string  `json:"name"`
	IsPreset bool    `json:"is_preset"`
}

func NewTagResponse(tag db.Tag) TagResponse {
	var familyID *string
	if tag.FamilyID.Valid {
		v := tag.FamilyID.UUID.String()
		familyID = &v
	}
	return TagResponse{
		ID:       tag.ID,
		FamilyID: familyID,
		Name:     tag.Name,
		IsPreset: tag.IsPreset,
	}
}

func NewTagResponses(tags []db.Tag) []TagResponse {
	if tags == nil {
		return []TagResponse{}
	}
	result := make([]TagResponse, len(tags))
	for i, t := range tags {
		result[i] = NewTagResponse(t)
	}
	return result
}

type MediaItemResponse struct {
	ID              string        `json:"id"`
	FamilyID        string        `json:"family_id"`
	AlbumID         string        `json:"album_id"`
	TakenAt         time.Time     `json:"taken_at"`
	FileName        string        `json:"file_name"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	OriginalUrl     string        `json:"original_url"`
	OriginalWidth   int32         `json:"original_width"`
	OriginalHeight  int32         `json:"original_height"`
	ThumbnailUrl    string        `json:"thumbnail_url"`
	ThumbnailWidth  int32         `json:"thumbnail_width"`
	ThumbnailHeight int32         `json:"thumbnail_height"`
	ViewUrl         string        `json:"view_url"`
	ViewWidth       int32         `json:"view_width"`
	ViewHeight      int32         `json:"view_height"`
	LiveUrl         string        `json:"live_url"`
	VideoUrl        string        `json:"video_url"`
	IsLiked         bool          `json:"is_liked"`
	Tags            []TagResponse `json:"tags"`
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
		OriginalUrl:     internalutils.ParseStoragePath(mediaItem.OriginalStorageKey.String),
		OriginalWidth:   mediaItem.OriginalWidth.Int32,
		OriginalHeight:  mediaItem.OriginalHeight.Int32,
		ThumbnailUrl:    internalutils.ParseStoragePath(mediaItem.ThumbnailStorageKey.String),
		ThumbnailWidth:  mediaItem.ThumbnailWidth.Int32,
		ThumbnailHeight: mediaItem.ThumbnailHeight.Int32,
		ViewUrl:         internalutils.ParseStoragePath(mediaItem.ViewStorageKey.String),
		ViewWidth:       mediaItem.ViewWidth.Int32,
		ViewHeight:      mediaItem.ViewHeight.Int32,
		VideoUrl:        internalutils.ParseStoragePath(mediaItem.VideoStorageKey.String),
		IsLiked:         mediaItem.IsLiked.Int32 > 0,
		Tags:            []TagResponse{},
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
	ID            string `json:"id"`
	UploadStatus  string `json:"upload_status"`
	FailureReason string `json:"failure_reason,omitempty"`
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
			ID:            uploadStatus.ID,
			UploadStatus:  uploadStatus.UploadStatus,
			FailureReason: uploadStatus.FailureReason.String,
		}
	}
	return &UploadBatchStatusResponse{
		Statuses:    statuses,
		IsCompleted: isCompleted,
	}
}

type BatchThumbnailResponse struct {
	ID              string `json:"id"`
	ThumbnailUrl    string `json:"thumbnail_url"`
	ThumbnailWidth  int32  `json:"thumbnail_width"`
	ThumbnailHeight int32  `json:"thumbnail_height"`
}

type UploadBatchWithThumbnailsResponse struct {
	ID         int32                    `json:"id"`
	AlbumID    string                   `json:"album_id"`
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
		OriginalUrl:     internalutils.ParseStoragePath(item.OriginalStorageKey.String),
		OriginalWidth:   item.OriginalWidth.Int32,
		OriginalHeight:  item.OriginalHeight.Int32,
		ThumbnailUrl:    internalutils.ParseStoragePath(item.ThumbnailStorageKey.String),
		ThumbnailWidth:  item.ThumbnailWidth.Int32,
		ThumbnailHeight: item.ThumbnailHeight.Int32,
		ViewUrl:         internalutils.ParseStoragePath(item.ViewStorageKey.String),
		ViewWidth:       item.ViewWidth.Int32,
		ViewHeight:      item.ViewHeight.Int32,
		VideoUrl:        internalutils.ParseStoragePath(item.VideoStorageKey.String),
		IsLiked:         item.IsLiked.Int32 > 0,
		Tags:            []TagResponse{},
	}
}

type UpdateMediaItemAlbumRequest struct {
	AlbumID string `json:"album_id" validate:"required"`
}

type GetMediaItemsRequest struct {
	From *time.Time `query:"from" validate:"required"`
	To   *time.Time `query:"to" validate:"required"`
}

type SearchMediaItemsRequest struct {
	From        *time.Time `query:"from"`
	To          *time.Time `query:"to"`
	IdentityIDs []int32    `query:"identity_ids"`
	AlbumID     *string    `query:"album_id"`
	Liked       bool       `query:"liked"`
	TagIDs      []int32    `query:"tag_ids"`
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
		OriginalUrl:     internalutils.ParseStoragePath(item.OriginalStorageKey.String),
		OriginalWidth:   item.OriginalWidth.Int32,
		OriginalHeight:  item.OriginalHeight.Int32,
		ThumbnailUrl:    internalutils.ParseStoragePath(item.ThumbnailStorageKey.String),
		ThumbnailWidth:  item.ThumbnailWidth.Int32,
		ThumbnailHeight: item.ThumbnailHeight.Int32,
		ViewUrl:         internalutils.ParseStoragePath(item.ViewStorageKey.String),
		ViewWidth:       item.ViewWidth.Int32,
		ViewHeight:      item.ViewHeight.Int32,
		VideoUrl:        internalutils.ParseStoragePath(item.VideoStorageKey.String),
		IsLiked:         item.IsLiked.Int32 > 0,
		Tags:            []TagResponse{},
	}
}
