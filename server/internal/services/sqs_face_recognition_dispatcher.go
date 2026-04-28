package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

// ECSFaceRecognitionDispatcher 프로덕션용: job 삽입 후 ECS task 트리거
type SQSFaceRecognitionDispatcher struct {
	jobStore  store.FaceRecognitionJobStore
	sqsClient *sqs.Client
	queueUrl  string
}

func NewSQSFaceRecognitionDispatcher(
	jobStore store.FaceRecognitionJobStore,
	sqsClient *sqs.Client,
	queueUrl string,
) *SQSFaceRecognitionDispatcher {
	return &SQSFaceRecognitionDispatcher{
		jobStore:  jobStore,
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
