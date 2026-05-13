package worker

import (
	"database/sql"
	"log"

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

	faceDispatcher := services.NewSQSFaceRecognitionDispatcher(nil, "")
	videoDispatcher := services.NewSQSVideoJobDispatcher(nil, "")

	return workerServices.NewResizeService(st.MediaItem, imageUploader, faceDispatcher, videoDispatcher)
}
