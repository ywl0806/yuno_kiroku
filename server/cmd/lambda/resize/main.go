package main

import (
	"context"
	"encoding/json"
	"log"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/worker"
	workerServices "github.com/ywl0806/yuno_kiroku/internal/worker/services"
	"github.com/ywl0806/yuno_kiroku/pkg/setting"
)

var echoLambda *echoadapter.EchoLambdaV2
var resizeSvc *workerServices.ResizeService

func init() {
	setting.SettingEnv()
	e := echo.New()
	resizeSvc = worker.Init(e)
	echoLambda = echoadapter.NewV2(e)
}

func handler(ctx context.Context, event json.RawMessage) (interface{}, error) {
	var probe struct {
		Records []struct {
			EventSource string `json:"eventSource"`
		} `json:"Records"`
	}
	json.Unmarshal(event, &probe)

	if len(probe.Records) > 0 && probe.Records[0].EventSource == "aws:s3" {
		var s3Event events.S3Event
		if err := json.Unmarshal(event, &s3Event); err != nil {
			return nil, err
		}
		for _, record := range s3Event.Records {
			key, err := url.QueryUnescape(record.S3.Object.Key)
			if err != nil {
				key = record.S3.Object.Key
			}
			log.Printf("S3 리사이즈 처리: %s", key)
			if err := resizeSvc.ProcessResize(ctx, key); err != nil {
				log.Printf("리사이즈 실패 [%s]: %v", key, err)
				return nil, err
			}
		}
		return nil, nil
	}

	// Function URL HTTP 이벤트 (face-recognition 콜백)
	var req events.APIGatewayV2HTTPRequest
	if err := json.Unmarshal(event, &req); err != nil {
		return nil, err
	}
	return echoLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(handler)
}
