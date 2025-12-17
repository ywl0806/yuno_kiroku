package models

import (
	"time"

	"github.com/guregu/null/v6"
	"github.com/ywl0806/yuno_kiroku/internal/db"
)

type faceDetectionResponse struct {
	ID             int32  `json:"id"`
	PhotoID        int32  `json:"photo_id"`
	IdentityID     int32  `json:"identity_id"`
	Name           string `json:"name"`
	LocationTop    int32  `json:"location_top"`
	LocationRight  int32  `json:"location_right"`
	LocationBottom int32  `json:"location_bottom"`
	LocationLeft   int32  `json:"location_left"`
}

type UploadPhotoResponse struct {
	ID              int32                   `json:"id"`
	GroupID         int32                   `json:"group_id"`
	AlbumID         int32                   `json:"album_id"`
	ThumbnailUrl    string                  `json:"thumbnail_url"`
	OriginalUrl     null.String             `json:"original_url"`
	LiveUrl         null.String             `json:"live_url"`
	OriginalLiveUrl null.String             `json:"original_live_url"`
	FileName        string                  `json:"file_name"`
	OriginalWidth   null.Int32              `json:"original_width"`
	OriginalHeight  null.Int32              `json:"original_height"`
	ThumbnailWidth  int32                   `json:"thumbnail_width"`
	ThumbnailHeight int32                   `json:"thumbnail_height"`
	PhotoCreatedAt  time.Time               `json:"photo_created_at"`
	FaceDetections  []faceDetectionResponse `json:"face_detections"`
}

func NewUploadPhotoResponse(photo *db.Photo, faceDetections []db.GetFaceDetectionsByPhotoIdRow) *UploadPhotoResponse {
	var faceDetectionResponses []faceDetectionResponse
	for _, faceDetection := range faceDetections {
		faceDetectionResponses = append(faceDetectionResponses, faceDetectionResponse{
			ID:             faceDetection.ID,
			PhotoID:        faceDetection.PhotoID,
			IdentityID:     faceDetection.IdentityID,
			Name:           faceDetection.Name.String,
			LocationTop:    faceDetection.LocationTop,
			LocationRight:  faceDetection.LocationRight,
			LocationBottom: faceDetection.LocationBottom,
			LocationLeft:   faceDetection.LocationLeft,
		})
	}

	return &UploadPhotoResponse{
		ID:              photo.ID,
		GroupID:         photo.GroupID,
		AlbumID:         photo.AlbumID,
		ThumbnailUrl:    photo.ThumbnailUrl,
		OriginalUrl:     null.NewString(photo.OriginalUrl.String, photo.OriginalUrl.Valid),
		LiveUrl:         null.NewString(photo.LiveUrl.String, photo.LiveUrl.Valid),
		OriginalLiveUrl: null.NewString(photo.OriginalLiveUrl.String, photo.OriginalLiveUrl.Valid),
		FileName:        photo.FileName,
		OriginalWidth:   null.NewInt32(photo.OriginalWidth.Int32, photo.OriginalWidth.Valid),
		OriginalHeight:  null.NewInt32(photo.OriginalHeight.Int32, photo.OriginalHeight.Valid),
		ThumbnailWidth:  photo.ThumbnailWidth,
		ThumbnailHeight: photo.ThumbnailHeight,
		PhotoCreatedAt:  photo.PhotoCreatedAt,
		FaceDetections:  faceDetectionResponses,
	}
}

type PhotoResponse struct {
	ID              int32       `json:"id"`
	GroupID         int32       `json:"group_id"`
	AlbumID         int32       `json:"album_id"`
	ThumbnailUrl    string      `json:"thumbnail_url"`
	OriginalUrl     null.String `json:"original_url"`
	LiveUrl         null.String `json:"live_url"`
	OriginalLiveUrl null.String `json:"original_live_url"`
	FileName        string      `json:"file_name"`
	OriginalWidth   null.Int32  `json:"original_width"`
	OriginalHeight  null.Int32  `json:"original_height"`
	ThumbnailWidth  int32       `json:"thumbnail_width"`
	ThumbnailHeight int32       `json:"thumbnail_height"`
	PhotoCreatedAt  time.Time   `json:"photo_created_at"`
}

func NewPhotoResponse(photo *db.Photo) *PhotoResponse {
	return &PhotoResponse{
		ID:              photo.ID,
		GroupID:         photo.GroupID,
		AlbumID:         photo.AlbumID,
		ThumbnailUrl:    photo.ThumbnailUrl,
		OriginalUrl:     null.NewString(photo.OriginalUrl.String, photo.OriginalUrl.Valid),
		LiveUrl:         null.NewString(photo.LiveUrl.String, photo.LiveUrl.Valid),
		OriginalLiveUrl: null.NewString(photo.OriginalLiveUrl.String, photo.OriginalLiveUrl.Valid),
		FileName:        photo.FileName,
		OriginalWidth:   null.NewInt32(photo.OriginalWidth.Int32, photo.OriginalWidth.Valid),
		OriginalHeight:  null.NewInt32(photo.OriginalHeight.Int32, photo.OriginalHeight.Valid),
		ThumbnailWidth:  photo.ThumbnailWidth,
		ThumbnailHeight: photo.ThumbnailHeight,
		PhotoCreatedAt:  photo.PhotoCreatedAt,
	}
}

func NewPhotosResponse(photos []db.Photo) *[]PhotoResponse {
	photoResponses := make([]PhotoResponse, len(photos))
	for i, photo := range photos {
		photoResponses[i] = *NewPhotoResponse(&photo)
	}
	return &photoResponses
}

type IdentityRandomPhotoResponse struct {
	ID              int32     `json:"id"`
	GroupID         int32     `json:"group_id"`
	AlbumID         int32     `json:"album_id"`
	ThumbnailUrl    string    `json:"thumbnail_url"`
	FileName        string    `json:"file_name"`
	OriginalWidth   int32     `json:"original_width"`
	OriginalHeight  int32     `json:"original_height"`
	ThumbnailWidth  int32     `json:"thumbnail_width"`
	ThumbnailHeight int32     `json:"thumbnail_height"`
	PhotoCreatedAt  time.Time `json:"photo_created_at"`
	LocationTop     int32     `json:"location_top"`
	LocationRight   int32     `json:"location_right"`
	LocationBottom  int32     `json:"location_bottom"`
	LocationLeft    int32     `json:"location_left"`
}

func NewIdentityRandomPhotoResponse(photo *db.GetIdentityRandomPhotoRow) *IdentityRandomPhotoResponse {
	return &IdentityRandomPhotoResponse{
		ID:              photo.ID,
		GroupID:         photo.GroupID,
		AlbumID:         photo.AlbumID,
		ThumbnailUrl:    photo.ThumbnailUrl,
		FileName:        photo.FileName,
		OriginalWidth:   photo.OriginalWidth.Int32,
		OriginalHeight:  photo.OriginalHeight.Int32,
		ThumbnailWidth:  photo.ThumbnailWidth,
		ThumbnailHeight: photo.ThumbnailHeight,
		PhotoCreatedAt:  photo.PhotoCreatedAt,
		LocationTop:     photo.LocationTop,
		LocationRight:   photo.LocationRight,
		LocationBottom:  photo.LocationBottom,
		LocationLeft:    photo.LocationLeft,
	}
}
