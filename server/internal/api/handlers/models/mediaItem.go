package models

import (
	"time"

	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/internal/api/utils"
	"github.com/ywl0806/yuno_kiroku/internal/db"
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
}

func NewMediaItemResponse(mediaItem *db.GetMediaItemsByTakenAtRow) *MediaItemResponse {
	originalStorageUrl := viper.GetString("ORIGINAL_STORAGE_URL")
	thumbnailStorageUrl := viper.GetString("THUMBNAIL_STORAGE_URL")
	return &MediaItemResponse{
		ID:              mediaItem.ID,
		GroupID:         mediaItem.GroupID,
		AlbumID:         mediaItem.AlbumID,
		TakenAt:         mediaItem.TakenAt,
		FileName:        mediaItem.FileName.String,
		CreatedAt:       mediaItem.CreatedAt,
		UpdatedAt:       mediaItem.UpdatedAt,
		OriginalUrl:     utils.ParsePath(originalStorageUrl, mediaItem.OriginalStorageKey),
		OriginalWidth:   mediaItem.OriginalWidth.Int32,
		OriginalHeight:  mediaItem.OriginalHeight.Int32,
		ThumbnailUrl:    utils.ParsePath(thumbnailStorageUrl, mediaItem.ThumbnailStorageKey),
		ThumbnailWidth:  mediaItem.ThumbnailWidth.Int32,
		ThumbnailHeight: mediaItem.ThumbnailHeight.Int32,
	}
}

func NewMediaItemsResponse(mediaItems []db.GetMediaItemsByTakenAtRow) *[]MediaItemResponse {
	mediaItemResponses := make([]MediaItemResponse, len(mediaItems))
	for i, mediaItem := range mediaItems {
		mediaItemResponses[i] = *NewMediaItemResponse(&mediaItem)
	}
	return &mediaItemResponses
}
