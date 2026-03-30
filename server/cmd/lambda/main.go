package main

import (
	"github.com/labstack/echo/v4"

	_ "github.com/ywl0806/yuno_kiroku/docs"
	"github.com/ywl0806/yuno_kiroku/internal/api"

	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"

	"github.com/aws/aws-lambda-go/lambda"
)

var echoLambda *echoadapter.EchoLambda

// @BasePath /api
func main() {
	// config
	// setting.SettingEnv()
	// mode := viper.GetString("APP_MODE")

	// fmt.Println("mode: ", mode)
	e := echo.New()

	// /api/* 경로는 API 라우팅
	api.Init(e)

	echoLambda = echoadapter.New(e)
	lambda.Start(echoLambda.ProxyWithContext)

}
