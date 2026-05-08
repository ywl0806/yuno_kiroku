package store

import (
	"context"
	"database/sql"

	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// Transactor 트랜잭션 범위 관리 인터페이스
type Transactor interface {
	Transact(ctx context.Context, fn func(tx *Store) error) error
	TransactWithAdvisoryLock(ctx context.Context, key int64, fn func(tx *Store) error) error
}

// Store 모든 도메인 Store를 묶은 컨테이너
type Store struct {
	db      *sql.DB     // BeginTx 전용
	queries *db.Queries // WithTx 전용

	User                 UserStore
	Family               FamilyStore
	Group                GroupStore
	Album                AlbumStore
	AlbumGroupPermission AlbumGroupPermissionStore
	Identity             IdentityStore
	Face                 FaceStore
	MediaItem            MediaItemStore
	InviteToken          InviteTokenStore
	Kid                  KidStore
	IdentityFaceImg      IdentityFaceImgStore
	Like                 LikeStore
	Tag                  TagStore
}

// New Store 컨테이너 생성 (각 도메인 Store 구현체 주입)
func New(sqlDB *sql.DB, queries *db.Queries) *Store {
	return &Store{
		db:      sqlDB,
		queries: queries,

		User:                 NewUserStore(queries),
		Family:               NewFamilyStore(queries),
		Group:                NewGroupStore(queries),
		Album:                NewAlbumStore(queries),
		AlbumGroupPermission: NewAlbumGroupPermissionStore(queries),
		Identity:             NewIdentityStore(queries),
		Face:                 NewFaceStore(queries),
		MediaItem:            NewMediaItemStore(queries),
		InviteToken:          NewInviteTokenStore(queries),
		Kid:                  NewKidStore(queries),
		IdentityFaceImg:      NewIdentityFaceImgStore(queries),
		Like:                 NewLikeStore(queries),
		Tag:                  NewTagStore(queries),
	}
}

// AcquireAdvisoryXactLock 트랜잭션 범위 advisory lock을 획득
// 트랜잭션 종료(커밋/롤백) 시 자동 해제. Transact 내부에서만 사용해야 함.
func AcquireAdvisoryXactLock(ctx context.Context, key int64, tx db.DBTX) error {
	_, err := tx.ExecContext(ctx, "SELECT pg_advisory_xact_lock($1)", key)
	return err
}

// ReleaseAdvisoryXactLock 트랜잭션 범위 advisory lock을 해제
// 트랜잭션 종료(커밋/롤백) 시 자동 해제. Transact 내부에서만 사용해야 함.
func ReleaseAdvisoryXactLock(ctx context.Context, key int64, tx db.DBTX) error {
	_, err := tx.ExecContext(ctx, "SELECT pg_advisory_xact_unlock($1)", key)
	return err
}

// Transact 트랜잭션 범위 내에서 fn을 실행한다. fn이 에러를 반환하면 롤백, 성공하면 커밋한다.
func (s *Store) Transact(ctx context.Context, fn func(tx *Store) error) error {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	// panic 발생 시 트랜잭션을 롤백한 뒤 패닉을 던짐
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	txStore := New(s.db, s.queries.WithTx(tx))
	if err := fn(txStore); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// TransactWithAdvisoryLock 트랜잭션 범위 내에서 advisory lock을 획득한 후 fn을 실행한다.
func (s *Store) TransactWithAdvisoryLock(ctx context.Context, key int64, fn func(tx *Store) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// advisory lock 획득
	if err = AcquireAdvisoryXactLock(ctx, key, tx); err != nil {
		tx.Rollback()
		return err
	}
	defer ReleaseAdvisoryXactLock(ctx, key, tx) // advisory lock 해제

	// panic 발생 시 롤백한 뒤 패닉을 다시 던짐
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	// fn 실행
	txStore := New(s.db, s.queries.WithTx(tx))
	if err := fn(txStore); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
