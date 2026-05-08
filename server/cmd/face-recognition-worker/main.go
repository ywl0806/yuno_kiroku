package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/providers"
	"github.com/ywl0806/yuno_kiroku/internal/services"
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

	var input cliInput
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		log.Fatalf("입력 파싱 실패: %v", err)
	}

	conn, err := sql.Open("pgx", viper.GetString("DATABASE_URL"))
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}
	defer conn.Close()

	st := store.New(conn, db.New(conn))
	storageService := providers.NewStorageProvider().StorageService()
	imageUploader := services.NewImageUploader(storageService)

	svc := workerServices.NewFaceRecognitionService(
		st,
		st.Face,
		st.Identity,
		st.IdentityFaceImg,
		imageUploader,
	)

	ctx := context.Background()
	if err := svc.ProcessFaces(ctx, input.ProcessFacesParams, input.Faces); err != nil {
		log.Fatalf("얼굴 인식 처리 실패: %v", err)
	}
}
