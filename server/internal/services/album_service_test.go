package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ywl0806/yuno_kiroku/internal/apperr"
	"github.com/ywl0806/yuno_kiroku/internal/db"
	"github.com/ywl0806/yuno_kiroku/internal/store"
	"github.com/ywl0806/yuno_kiroku/internal/store/mocks"
)

// newTxTransactor Transact 콜백에 주어진 TxStore 목을 그대로 넘겨 즉시 실행하는 Transactor 목을 만든다.
func newTxTransactor(t *testing.T, tx store.TxStore) *mocks.MockTransactor {
	transactor := mocks.NewMockTransactor(t)
	transactor.EXPECT().Transact(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(store.TxStore) error) error {
			return fn(tx)
		})
	return transactor
}

// CreateAlbum은 트랜잭션 안에서 앨범을 만들고 권한을 건별로 삽입해야 한다.
func TestCreateAlbum_InsertsAlbumAndPermissions(t *testing.T) {
	created := db.Album{ID: 10, FamilyID: 1, Name: "여행"}

	albumStore := mocks.NewMockAlbumStore(t)
	albumStore.EXPECT().CreateAlbum(mock.Anything, db.CreateAlbumParams{FamilyID: 1, Name: "여행"}).
		Return(created, nil)

	permStore := mocks.NewMockAlbumGroupPermissionStore(t)
	permStore.EXPECT().InsertAlbumGroupPermission(mock.Anything, db.InsertAlbumGroupPermissionParams{
		AlbumID: 10, GroupID: 2, Permission: "write",
	}).Return(nil)

	tx := mocks.NewMockTxStore(t)
	tx.EXPECT().Album().Return(albumStore)
	tx.EXPECT().AlbumGroupPermission().Return(permStore)

	svc := NewAlbumService(mocks.NewMockAlbumStore(t), mocks.NewMockAlbumGroupPermissionStore(t), newTxTransactor(t, tx))

	got, err := svc.CreateAlbum(context.Background(), 1, "여행", []GroupPermission{{GroupID: 2, Permission: "write"}})

	assert.NoError(t, err)
	assert.Equal(t, created, got)
}

// 권한 삽입이 실패하면 에러가 그대로 전파되고 결과는 zero value여야 한다.
func TestCreateAlbum_PermissionInsertFails(t *testing.T) {
	wantErr := errors.New("insert failed")

	albumStore := mocks.NewMockAlbumStore(t)
	albumStore.EXPECT().CreateAlbum(mock.Anything, mock.Anything).Return(db.Album{ID: 10}, nil)

	permStore := mocks.NewMockAlbumGroupPermissionStore(t)
	permStore.EXPECT().InsertAlbumGroupPermission(mock.Anything, mock.Anything).Return(wantErr)

	tx := mocks.NewMockTxStore(t)
	tx.EXPECT().Album().Return(albumStore)
	tx.EXPECT().AlbumGroupPermission().Return(permStore)

	svc := NewAlbumService(mocks.NewMockAlbumStore(t), mocks.NewMockAlbumGroupPermissionStore(t), newTxTransactor(t, tx))

	got, err := svc.CreateAlbum(context.Background(), 1, "여행", []GroupPermission{{GroupID: 2, Permission: "write"}})

	assert.ErrorIs(t, err, wantErr)
	assert.Equal(t, db.Album{}, got)
}

// 다른 가족의 앨범을 지우려 하면 Forbidden이고 트랜잭션은 시작되지 않아야 한다.
func TestDeleteAlbum_ForbiddenWhenFamilyMismatch(t *testing.T) {
	albumStore := mocks.NewMockAlbumStore(t)
	albumStore.EXPECT().FindAlbumByID(mock.Anything, int32(10)).Return(db.Album{ID: 10, FamilyID: 999}, nil)

	transactor := mocks.NewMockTransactor(t)

	svc := NewAlbumService(albumStore, mocks.NewMockAlbumGroupPermissionStore(t), transactor)

	err := svc.DeleteAlbum(context.Background(), 10, 1)

	assert.True(t, apperr.IsAppError(err, apperr.Forbidden))
	transactor.AssertNotCalled(t, "Transact", mock.Anything, mock.Anything)
}

// DeleteAlbum은 권한을 먼저 지우고 앨범을 지워야 한다.
func TestDeleteAlbum_DeletesPermissionsThenAlbum(t *testing.T) {
	outerAlbumStore := mocks.NewMockAlbumStore(t)
	outerAlbumStore.EXPECT().FindAlbumByID(mock.Anything, int32(10)).Return(db.Album{ID: 10, FamilyID: 1}, nil)

	txAlbumStore := mocks.NewMockAlbumStore(t)
	txAlbumStore.EXPECT().DeleteAlbum(mock.Anything, int32(10)).Return(nil)

	permStore := mocks.NewMockAlbumGroupPermissionStore(t)
	permStore.EXPECT().DeleteAlbumGroupPermissionsByAlbumID(mock.Anything, int32(10)).Return(nil)

	tx := mocks.NewMockTxStore(t)
	tx.EXPECT().Album().Return(txAlbumStore)
	tx.EXPECT().AlbumGroupPermission().Return(permStore)

	svc := NewAlbumService(outerAlbumStore, mocks.NewMockAlbumGroupPermissionStore(t), newTxTransactor(t, tx))

	assert.NoError(t, svc.DeleteAlbum(context.Background(), 10, 1))
}
