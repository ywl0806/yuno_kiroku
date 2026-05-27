package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/logger"
	"github.com/ywl0806/yuno_kiroku/internal/providers"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	workerServices "github.com/ywl0806/yuno_kiroku/internal/worker/services"
	"github.com/ywl0806/yuno_kiroku/pkg/setting"
)

type cliInput struct {
	workerServices.ProcessFacesParams
	Faces []workerServices.FaceResult `json:"faces"`
}

func main() {
	setting.SettingEnv()
	logger.Init(viper.GetString("APP_ENV"), "yuno-face-worker")

	var input cliInput
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		slog.Error("입력 파싱 실패", "error", err)
		os.Exit(1)
	}

	conn, err := sql.Open("pgx", viper.GetString("DATABASE_URL"))
	if err != nil {
		slog.Error("DB 연결 실패", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	st := store.New(conn, db.New(conn))
	storageService := providers.NewStorageProvider().StorageService()

	svc := workerServices.NewFaceRecognitionService(
		st,
		st.Face,
		st.Identity,
		st.IdentityFaceImg,
		storageService,
	)

	ctx := context.Background()
	if err := svc.ProcessFaces(ctx, input.ProcessFacesParams, input.Faces); err != nil {
		slog.Error("얼굴 인식 처리 실패", "error", err)
		os.Exit(1)
	}
}
