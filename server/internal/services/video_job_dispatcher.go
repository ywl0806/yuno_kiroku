package services

import (
	"context"
	"encoding/json"
	"log"

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
	sqsClient *sqs.Client
	queueUrl  string
}

func NewSQSVideoJobDispatcher(sqsClient *sqs.Client, queueUrl string) *SQSVideoJobDispatcher {
	if sqsClient == nil {
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			log.Fatalf("AWS 설정 로드 실패: %v", err)
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
		sqsClient: sqsClient,
		queueUrl:  queueUrl,
	}
}

func (d *SQSVideoJobDispatcher) Dispatch(ctx context.Context, params VideoJobParams) error {
	body, err := json.Marshal(params)
	if err != nil {
		log.Printf("VideoJob JSON Marshal 실패: %v", err)
		return err
	}
	_, err = d.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(d.queueUrl),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		log.Printf("VideoJob SQS SendMessage 실패: %v", err)
		return err
	}
	return nil
}
