package worker

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	awsecs "github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/providers"
	"github.com/ywl0806/yuno_kiroku/internal/services"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	"github.com/ywl0806/yuno_kiroku/internal/worker/handlers"
	workerServices "github.com/ywl0806/yuno_kiroku/internal/worker/services"
)

func Init(e *echo.Echo) {
	conn, err := sql.Open("pgx", viper.GetString("DATABASE_URL"))
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}

	queries := db.New(conn)
	st := store.New(queries)

	storageProvider := providers.NewStorageProvider()
	storageService := storageProvider.StorageService()

	imageUploader := services.NewImageUploader(storageService)

	var faceDispatcher services.FaceRecognitionDispatcher
	if viper.GetString("APP_MODE") == "local_dev" {
		faceDispatcher = services.NewLocalFaceRecognitionDispatcher(st.FaceRecognitionJob)
	} else {
		cfg, cfgErr := config.LoadDefaultConfig(context.Background())
		if cfgErr != nil {
			log.Fatalf("AWS 설정 로드 실패: %v", cfgErr)
		}
		faceDispatcher = services.NewECSFaceRecognitionDispatcher(
			st.FaceRecognitionJob,
			awsecs.NewFromConfig(cfg),
			viper.GetString("ECS_CLUSTER_ARN"),
			viper.GetString("ECS_TASK_DEF_ARN"),
			strings.Split(viper.GetString("ECS_SUBNETS"), ","),
			strings.Split(viper.GetString("ECS_SECURITY_GROUPS"), ","),
		)
	}

	resizeService := workerServices.NewResizeService(st.MediaItem, imageUploader, faceDispatcher)
	faceRecognitionService := workerServices.NewFaceRecognitionService(
		st.FaceRecognitionJob,
		st.MediaItem,
		st.Face,
		st.Identity,
		st.IdentityFaceImg,
		imageUploader,
	)

	resizeHandler := handlers.NewResizeHandler(resizeService)
	faceRecognitionHandler := handlers.NewFaceRecognitionHandler(faceRecognitionService)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/resize", resizeHandler.HandleMinioEvent)
	e.POST("/face-recognition/complete", faceRecognitionHandler.Complete)
	e.POST("/face-recognition/fail", faceRecognitionHandler.Fail)
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})
}
