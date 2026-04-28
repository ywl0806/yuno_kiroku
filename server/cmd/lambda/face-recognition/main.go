package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/worker"
	"github.com/ywl0806/yuno_kiroku/pkg/setting"
)

var echoLambda *echoadapter.EchoLambdaV2

func init() {
	setting.SettingEnv()
	e := echo.New()
	worker.InitFaceRecognition(e)
	echoLambda = echoadapter.NewV2(e)
}

func main() {
	lambda.Start(echoLambda.ProxyWithContext)
}
