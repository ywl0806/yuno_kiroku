package services

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
)

// FaceRecognitionDispatcher face_recognition_job을 생성하고 처리 트리거를 발행합니다.
type FaceRecognitionDispatcher interface {
	Dispatch(ctx context.Context, params db.CreateFaceRecognitionJobParams) error
}

// LocalFaceRecognitionDispatcher 로컬 개발용: job만 DB에 삽입 (Docker 컨테이너가 폴링)
type LocalFaceRecognitionDispatcher struct {
	jobStore store.FaceRecognitionJobStore
}

func NewLocalFaceRecognitionDispatcher(jobStore store.FaceRecognitionJobStore) *LocalFaceRecognitionDispatcher {
	return &LocalFaceRecognitionDispatcher{jobStore: jobStore}
}

func (d *LocalFaceRecognitionDispatcher) Dispatch(ctx context.Context, params db.CreateFaceRecognitionJobParams) error {
	_, err := d.jobStore.CreateFaceRecognitionJob(ctx, params)
	return err
}
