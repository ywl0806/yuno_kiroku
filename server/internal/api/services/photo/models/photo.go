package models

import (
	"time"

	"github.com/guregu/null/v6"
	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// type Photo struct {
// 	ID              primitive.ObjectID  `json:"_id" bson:"_id,omitempty"`
// 	GroupId         primitive.ObjectID  `json:"group_id" bson:"group_id"`
// 	ClanGroupId     *primitive.ObjectID `json:"clan_group_id" bson:"clan_group_id"`
// 	AlbumId         primitive.ObjectID  `json:"album_id" bson:"album_id"`
// 	ThumbnailUrl    string              `json:"thumbnail_url" bson:"thumbnail_url"`
// 	OriginalUrl     string              `json:"original_url" bson:"original_url"`
// 	LiveUrl         string              `json:"live_url" bson:"live_url"`
// 	OriginalLiveUrl string              `json:"original_live_url" bson:"original_live_url"`
// 	Width           int                 `json:"width" bson:"width"`
// 	Height          int                 `json:"height" bson:"height"`
// 	Orientation     int                 `json:"orientation" bson:"orientation"`
// 	FileName        string              `json:"file_name" bson:"file_name"`
// 	PhotoCreatedAt  time.Time           `json:"photo_created_at" bson:"photo_created_at"`
// 	CreatedAt       time.Time           `json:"created_at" bson:"created_at"`
// 	UpdatedAt       time.Time           `json:"updated_at" bson:"updated_at"`
// 	CreatedBy       string              `json:"created_by" bson:"created_by"`
// 	UpdatedBy       string              `json:"updated_by" bson:"updated_by"`
// }

// type PhotoGroupId struct {
// 	Year  int `json:"year" bson:"year"`
// 	Month int `json:"month" bson:"month"`
// }
// type PhotoGroup struct {
// 	PhotoGroupId `json:"_id" bson:"_id"`
// 	Photos       []Photo `json:"photos" bson:"photos"`
// }

// type PhotoRange struct {
// 	Year  int `json:"year" bson:"year"`
// 	Month int `json:"month" bson:"month"`
// }
// type PhotoRanges struct {
// 	PhotoRange []PhotoRange `json:"photo_range" bson:"photo_range"`
// }

// type PhotoResponse struct {
// 	ID              int32     `json:"id"`
// 	GroupID         int32     `json:"group_id"`
// 	ClanGroupID     *int32    `json:"clan_group_id"`
// 	ThumbnailUrl    string    `json:"thumbnail_url"`
// 	OriginalUrl     *string   `json:"original_url"`
// 	LiveUrl         *string   `json:"live_url"`
// 	OriginalLiveUrl *string   `json:"original_live_url"`
// 	Width           int32     `json:"width"`
// 	Height          int32     `json:"height"`
// 	Orientation     int32     `json:"orientation"`
// 	PhotoCreatedAt  time.Time `json:"photo_created_at"`
// 	CreatedAt       time.Time `json:"created_at"`
// }

// type PhotoGroupResponse struct {
// 	Year   int             `json:"year" `
// 	Month  int             `json:"month" `
// 	Photos []PhotoResponse `json:"photos"`
// }

type faceDetectionResponse struct {
	ID             int32  `json:"id"`
	PhotoID        int32  `json:"photo_id"`
	PersonID       int32  `json:"person_id"`
	Name           string `json:"name"`
	LocationTop    int32  `json:"location_top"`
	LocationRight  int32  `json:"location_right"`
	LocationBottom int32  `json:"location_bottom"`
	LocationLeft   int32  `json:"location_left"`
}

type UploadPhotoResponse struct {
	ID              int32                   `json:"id"`
	GroupID         int32                   `json:"group_id"`
	ClanGroupID     null.Int32              `json:"clan_group_id"`
	ThumbnailUrl    string                  `json:"thumbnail_url"`
	OriginalUrl     null.String             `json:"original_url"`
	LiveUrl         null.String             `json:"live_url"`
	OriginalLiveUrl null.String             `json:"original_live_url"`
	FileName        string                  `json:"file_name"`
	Width           int32                   `json:"width"`
	Height          int32                   `json:"height"`
	Orientation     int32                   `json:"orientation"`
	PhotoCreatedAt  time.Time               `json:"photo_created_at"`
	FaceDetections  []faceDetectionResponse `json:"face_detections"`
}

func NewUploadPhotoResponse(photo *db.Photo, faceDetections []db.GetFaceDetectionsByPhotoIdRow) *UploadPhotoResponse {
	var faceDetectionResponses []faceDetectionResponse
	for _, faceDetection := range faceDetections {
		faceDetectionResponses = append(faceDetectionResponses, faceDetectionResponse{
			ID:             faceDetection.ID,
			PhotoID:        faceDetection.PhotoID,
			PersonID:       faceDetection.PersonID,
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
		ClanGroupID:     null.NewInt32(photo.ClanGroupID.Int32, photo.ClanGroupID.Valid),
		ThumbnailUrl:    photo.ThumbnailUrl,
		OriginalUrl:     null.NewString(photo.OriginalUrl.String, photo.OriginalUrl.Valid),
		LiveUrl:         null.NewString(photo.LiveUrl.String, photo.LiveUrl.Valid),
		OriginalLiveUrl: null.NewString(photo.OriginalLiveUrl.String, photo.OriginalLiveUrl.Valid),
		FileName:        photo.FileName,
		Width:           photo.Width,
		Height:          photo.Height,
		Orientation:     photo.Orientation,
		PhotoCreatedAt:  photo.PhotoCreatedAt,
		FaceDetections:  faceDetectionResponses,
	}
}
