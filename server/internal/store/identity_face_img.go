package store

import (
	"context"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// IdentityFaceImgStore identity 얼굴 이미지 데이터 접근 인터페이스
type IdentityFaceImgStore interface {
	CreateIdentityFaceImg(ctx context.Context, arg db.CreateIdentityFaceImgParams) (db.IdentityFaceImg, error)
	HasIdentityFaceImg(ctx context.Context, identityID int32) (bool, error)
}

type identityFaceImgStore struct {
	queries *db.Queries
}

// NewIdentityFaceImgStore IdentityFaceImgStore 구현체 생성
func NewIdentityFaceImgStore(queries *db.Queries) IdentityFaceImgStore {
	return &identityFaceImgStore{queries: queries}
}

func (s *identityFaceImgStore) CreateIdentityFaceImg(ctx context.Context, arg db.CreateIdentityFaceImgParams) (db.IdentityFaceImg, error) {
	return wrapErr(s.queries.CreateIdentityFaceImg(ctx, arg))
}

func (s *identityFaceImgStore) HasIdentityFaceImg(ctx context.Context, identityID int32) (bool, error) {
	return s.queries.HasIdentityFaceImg(ctx, identityID)
}
