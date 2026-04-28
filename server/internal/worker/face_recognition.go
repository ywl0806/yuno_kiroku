package worker

import (
	"database/sql"
	"log"

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

func InitFaceRecognition(e *echo.Echo) {
	conn, err := sql.Open("pgx", viper.GetString("DATABASE_URL"))
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}
	st := store.New(conn, db.New(conn))

	storageService := providers.NewStorageProvider().StorageService()
	imageUploader := services.NewImageUploader(storageService)

	faceRecognitionService := workerServices.NewFaceRecognitionService(
		st,
		st.FaceRecognitionJob,
		st.MediaItem,
		st.Face,
		st.Identity,
		st.IdentityFaceImg,
		imageUploader,
	)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	faceRecognitionHandler := handlers.NewFaceRecognitionHandler(faceRecognitionService)
	e.POST("/face-recognition/complete", faceRecognitionHandler.Complete)
	e.POST("/face-recognition/fail", faceRecognitionHandler.Fail)
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})
}
