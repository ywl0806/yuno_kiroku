package services

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/spf13/viper"
)

type VideoJobDispatcher interface {
	Dispatch(ctx context.Context, params VideoJobParams) error
}

type VideoJobParams struct {
	MediaItemID        string `json:"media_item_id"`
	FamilyID           string `json:"family_id"`
	OriginalStorageKey string `json:"original_storage_key"`
	FileName           string `json:"file_name"`
	MimeType           string `json:"mime_type"`
}

type SQSVideoJobDispatcher struct {
	log       *slog.Logger
	sqsClient *sqs.Client
	queueUrl  string
}

func NewSQSVideoJobDispatcher(sqsClient *sqs.Client, queueUrl string) *SQSVideoJobDispatcher {
	if sqsClient == nil {
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			slog.Error("AWS 설정 로드 실패", "error", err)
			os.Exit(1)
		}
		if ep := viper.GetString("SQS_ENDPOINT_URL"); ep != "" {
			cfg.BaseEndpoint = aws.String(ep)
		}
		sqsClient = sqs.NewFromConfig(cfg)
	}

	if queueUrl == "" {
		queueUrl = viper.GetString("VIDEO_SQS_QUEUE_URL")
	}

	return &SQSVideoJobDispatcher{
		log:       slog.Default().With("layer", "service", "component", "video_dispatcher"),
		sqsClient: sqsClient,
		queueUrl:  queueUrl,
	}
}

func (d *SQSVideoJobDispatcher) Dispatch(ctx context.Context, params VideoJobParams) error {
	body, err := json.Marshal(params)
	if err != nil {
		d.log.ErrorContext(ctx, "video job marshal failed", "error", err)
		return err
	}
	_, err = d.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(d.queueUrl),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		d.log.ErrorContext(ctx, "video job SQS send failed", "media_item_id", params.MediaItemID, "error", err)
		return err
	}
	d.log.InfoContext(ctx, "video job dispatched", "media_item_id", params.MediaItemID)
	return nil
}
