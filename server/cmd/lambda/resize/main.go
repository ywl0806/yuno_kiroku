package main

import (
	"context"
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

func handler(ctx context.Context, event events.S3Event) error {
	for _, record := range event.Records {
		key, err := url.QueryUnescape(record.S3.Object.Key)
		if err != nil {
			key = record.S3.Object.Key
		}
		log.Printf("S3 리사이즈 처리: %s", key)
		if err := resizeSvc.ProcessResize(ctx, key); err != nil {
			log.Printf("리사이즈 실패 [%s]: %v", key, err)
			return err
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
