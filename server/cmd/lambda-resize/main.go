package main

import (
	"context"
	"encoding/json"
	"log"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ywl0806/yuno_kiroku/internal/worker"
	workerServices "github.com/ywl0806/yuno_kiroku/internal/worker/services"
	"github.com/ywl0806/yuno_kiroku/pkg/setting"
)

var resizeSvc *workerServices.ResizeService

func init() {
	setting.SettingEnv()
	resizeSvc = worker.InitResize()
}

// SQS가 S3 이벤트 JSON을 메시지 body에 감싸서 전달하므로 언래핑 후 처리
func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, sqsRecord := range sqsEvent.Records {
		var s3Event events.S3Event
		if err := json.Unmarshal([]byte(sqsRecord.Body), &s3Event); err != nil {
			log.Printf("SQS 메시지 파싱 실패: %v", err)
			return err
		}
		for _, record := range s3Event.Records {
			key, err := url.QueryUnescape(record.S3.Object.Key)
			if err != nil {
				key = record.S3.Object.Key
			}
			log.Printf("S3 리사이즈 처리: %s", key)
			if err := resizeSvc.ProcessResize(ctx, key); err != nil {
				log.Printf("리사이즈 실패 [%s]: %v", key, err)
				// 에러 반환 시 SQS 메시지가 visibility timeout 후 자동 재시도됨
				return err
			}
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
