package worker

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	awsecs "github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/spf13/viper"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/providers"
	"github.com/ywl0806/yuno_kiroku/internal/services"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	workerServices "github.com/ywl0806/yuno_kiroku/internal/worker/services"
)

func InitResize() *workerServices.ResizeService {
	conn, err := sql.Open("pgx", viper.GetString("DATABASE_URL"))
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}
	st := store.New(conn, db.New(conn))

	storageService := providers.NewStorageProvider().StorageService()
	imageUploader := services.NewImageUploader(storageService)

	var faceDispatcher services.FaceRecognitionDispatcher
	if viper.GetString("APP_ENV") == "local" {
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

	return workerServices.NewResizeService(st.MediaItem, imageUploader, faceDispatcher)
}
