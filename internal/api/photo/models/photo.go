package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Photo struct {
	ID              primitive.ObjectID `json:"_id" dynamodbav:"id"`
	GroupId         primitive.ObjectID `json:"group_id" dynamodbav:"group_id"`
	AlbumId         primitive.ObjectID `json:"album_id" dynamodbav:"album_id"`
	ThumbnailUrl    string             `json:"thumbnail_url" dynamodbav:"thumbnail_url"`
	OriginalUrl     string             `json:"original_url" dynamodbav:"original_url"`
	LiveUrl         string             `json:"live_url" dynamodbav:"live_url"`
	OriginalLiveUrl string             `json:"original_live_url" dynamodbav:"original_live_url"`
	Width           int                `json:"width" dynamodbav:"width"`
	Height          int                `json:"height" dynamodbav:"height"`
	Orientation     int                `json:"orientation" dynamodbav:"orientation"`
	FileName        string             `json:"file_name" dynamodbav:"file_name"`
	PhotoCreatedAt  time.Time          `json:"photo_created_at" dynamodbav:"photo_created_at"`
	CreatedAt       time.Time          `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" dynamodbav:"updated_at"`
	CreatedBy       string             `json:"created_by" dynamodbav:"created_by"`
	UpdatedBy       string             `json:"updated_by" dynamodbav:"updated_by"`
}

type PhotoGroupId struct {
	Year  int `json:"year" bson:"year"`
	Month int `json:"month" bson:"month"`
}
type PhotoGroup struct {
	PhotoGroupId `json:"_id" bson:"_id"`
	Photos       []Photo `json:"photos" bson:"photos"`
}

type PhotoRange struct {
	Year  int `json:"year" bson:"year"`
	Month int `json:"month" bson:"month"`
}
type PhotoRanges struct {
	PhotoRange []PhotoRange `json:"photo_range" bson:"photo_range"`
}
