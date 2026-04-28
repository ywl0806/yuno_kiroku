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

// ECSFaceRecognitionDispatcher 프로덕션용: job 삽입 후 ECS task 트리거
type SQSFaceRecognitionDispatcher struct {
	sqsClient *sqs.Client
	queueUrl  string
}

func NewSQSFaceRecognitionDispatcher(
	sqsClient *sqs.Client,
	queueUrl string,
) *SQSFaceRecognitionDispatcher {
	if sqsClient == nil {
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			log.Fatalf("AWS 설정 로드 실패: %v", err)
		}
		sqsClient = sqs.NewFromConfig(cfg)
	}

	if queueUrl == "" {
		queueUrl = viper.GetString("SQS_QUEUE_URL")
	}

	return &SQSFaceRecognitionDispatcher{
		sqsClient: sqsClient,
		queueUrl:  queueUrl,
	}
}

type FaceRecognitionMessage struct {
	MediaItemID    int32  `json:"media_item_id"`
	FamilyID       int32  `json:"family_id"`
	ViewStorageKey string `json:"view_storage_key"`
}

func (d *SQSFaceRecognitionDispatcher) Dispatch(ctx context.Context, params FaceRecognitionJobParams) error {
	body, err := json.Marshal(params)
	if err != nil {
		log.Printf("JSON Marshal 실패: %v", err)
		return err
	}
	_, err = d.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(d.queueUrl),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		log.Printf("SQS SendMessage 실패: %v", err)
		return err
	}

	return nil
}
