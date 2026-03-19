package store

import (
	"github.com/ywl0806/yuno_kiroku/internal/db"
)

// Store 모든 도메인 Store를 묶은 컨테이너
type Store struct {
	User            UserStore
	Family          FamilyStore
	Group           GroupStore
	Album           AlbumStore
	Identity        IdentityStore
	Face            FaceStore
	MediaItem       MediaItemStore
	InviteToken     InviteTokenStore
	Kid             KidStore
	IdentityFaceImg IdentityFaceImgStore
}

// New Store 컨테이너 생성 (각 도메인 Store 구현체 주입)
func New(queries *db.Queries) *Store {
	return &Store{
		User:            NewUserStore(queries),
		Family:          NewFamilyStore(queries),
		Group:           NewGroupStore(queries),
		Album:           NewAlbumStore(queries),
		Identity:        NewIdentityStore(queries),
		Face:            NewFaceStore(queries),
		MediaItem:       NewMediaItemStore(queries),
		InviteToken:     NewInviteTokenStore(queries),
		Kid:             NewKidStore(queries),
		IdentityFaceImg: NewIdentityFaceImgStore(queries),
	}
}
