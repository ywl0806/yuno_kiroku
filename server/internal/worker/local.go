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

// InitLocal 로컬 개발용 — resize 라우트를 Echo에 등록
func InitLocal(e *echo.Echo) {
	conn, err := sql.Open("pgx", viper.GetString("DATABASE_URL"))
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}
	st := store.New(conn, db.New(conn))

	storageService := providers.NewStorageProvider().StorageService()
	imageUploader := services.NewImageUploader(storageService)
	faceDispatcher := services.NewSQSFaceRecognitionDispatcher(nil, "")

	resizeService := workerServices.NewResizeService(st.MediaItem(), imageUploader, faceDispatcher)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/resize", handlers.NewResizeHandler(resizeService).HandleMinioEvent)
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})
}
