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

// FaceRecognitionDispatcher face_recognition_job을 생성하고 처리 트리거를 발행합니다.
type FaceRecognitionDispatcher interface {
	Dispatch(ctx context.Context, params FaceRecognitionJobParams) error
}

type FaceRecognitionJobParams struct {
	MediaItemID    string `json:"media_item_id"`
	FamilyID       string `json:"family_id"`
	ViewStorageKey string `json:"view_storage_key"`
}

// ECSFaceRecognitionDispatcher 프로덕션용: job 삽입 후 ECS task 트리거
type SQSFaceRecognitionDispatcher struct {
	log       *slog.Logger
	sqsClient *sqs.Client
	queueUrl  string
}

func NewSQSFaceRecognitionDispatcher(
	sqsClient *sqs.Client,
	queueUrl string,
) *SQSFaceRecognitionDispatcher {
	if sqsClient == nil {
		cfg, err := config.LoadDefaultConfig(context.Background())
		sqsEndpointURL := viper.GetString("SQS_ENDPOINT_URL")
		if sqsEndpointURL != "" {
			cfg.BaseEndpoint = aws.String(sqsEndpointURL)
		}
		if err != nil {
			slog.Error("AWS 설정 로드 실패", "error", err)
			os.Exit(1)
		}
		sqsClient = sqs.NewFromConfig(cfg)
	}

	if queueUrl == "" {
		queueUrl = viper.GetString("SQS_QUEUE_URL")
	}

	return &SQSFaceRecognitionDispatcher{
		log:       slog.Default().With("layer", "service", "component", "face_dispatcher"),
		sqsClient: sqsClient,
		queueUrl:  queueUrl,
	}
}

func (d *SQSFaceRecognitionDispatcher) Dispatch(ctx context.Context, params FaceRecognitionJobParams) error {
	body, err := json.Marshal(params)
	if err != nil {
		d.log.ErrorContext(ctx, "face recognition job marshal failed", "error", err)
		return err
	}
	_, err = d.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(d.queueUrl),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		d.log.ErrorContext(ctx, "face recognition SQS send failed", "media_item_id", params.MediaItemID, "error", err)
		return err
	}

	d.log.InfoContext(ctx, "face recognition job dispatched", "media_item_id", params.MediaItemID)
	return nil
}
