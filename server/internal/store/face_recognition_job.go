package store

import (
	"context"
	"database/sql"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// FaceRecognitionJobStore 얼굴 인식 배치 작업 데이터 접근 인터페이스
type FaceRecognitionJobStore interface {
	CreateFaceRecognitionJob(ctx context.Context, arg db.CreateFaceRecognitionJobParams) (db.FaceRecognitionJob, error)
	FetchPendingFaceRecognitionJobs(ctx context.Context, limit int32) ([]db.FaceRecognitionJob, error)
	GetFaceRecognitionJobByID(ctx context.Context, id int32) (db.FaceRecognitionJob, error)
	CompleteFaceRecognitionJob(ctx context.Context, id int32) error
	FailFaceRecognitionJob(ctx context.Context, id int32, lastError string) error
}

type faceRecognitionJobStore struct {
	queries *db.Queries
}

func NewFaceRecognitionJobStore(queries *db.Queries) FaceRecognitionJobStore {
	return &faceRecognitionJobStore{queries: queries}
}

func (s *faceRecognitionJobStore) CreateFaceRecognitionJob(ctx context.Context, arg db.CreateFaceRecognitionJobParams) (db.FaceRecognitionJob, error) {
	return wrapErr(s.queries.CreateFaceRecognitionJob(ctx, arg))
}

func (s *faceRecognitionJobStore) FetchPendingFaceRecognitionJobs(ctx context.Context, limit int32) ([]db.FaceRecognitionJob, error) {
	return wrapErr(s.queries.FetchPendingFaceRecognitionJobs(ctx, limit))
}

func (s *faceRecognitionJobStore) GetFaceRecognitionJobByID(ctx context.Context, id int32) (db.FaceRecognitionJob, error) {
	return wrapErr(s.queries.GetFaceRecognitionJobByID(ctx, id))
}

func (s *faceRecognitionJobStore) CompleteFaceRecognitionJob(ctx context.Context, id int32) error {
	return s.queries.CompleteFaceRecognitionJob(ctx, id)
}

func (s *faceRecognitionJobStore) FailFaceRecognitionJob(ctx context.Context, id int32, lastError string) error {
	return s.queries.FailFaceRecognitionJob(ctx, db.FailFaceRecognitionJobParams{
		ID:        id,
		LastError: sql.NullString{String: lastError, Valid: lastError != ""},
	})
}
