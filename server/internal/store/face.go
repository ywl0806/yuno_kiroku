package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// FaceStore 얼굴 감지 데이터 접근 인터페이스
type FaceStore interface {
	CreateFaceDetection(ctx context.Context, arg db.CreateFaceDetectionParams) (db.FaceDetection, error)
	FindMostSimilarFace(ctx context.Context, arg db.FindMostSimilarFaceParams) (db.FindMostSimilarFaceRow, error)
}

type faceStore struct {
	queries *db.Queries
}

// NewFaceStore FaceStore 구현체 생성
func NewFaceStore(queries *db.Queries) FaceStore {
	return &faceStore{queries: queries}
}

func (s *faceStore) CreateFaceDetection(ctx context.Context, arg db.CreateFaceDetectionParams) (db.FaceDetection, error) {
	return wrapErr(s.queries.CreateFaceDetection(ctx, arg))
}

func (s *faceStore) FindMostSimilarFace(ctx context.Context, arg db.FindMostSimilarFaceParams) (db.FindMostSimilarFaceRow, error) {
	return wrapErr(s.queries.FindMostSimilarFace(ctx, arg))
}
