-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- 그룹 테이블
CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 클랜 그룹 테이블
CREATE TABLE IF NOT EXISTS clan_groups (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups (id),
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 사용자 테이블
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    group_id INTEGER NOT NULL REFERENCES groups (id),
    clan_group_id INTEGER NOT NULL REFERENCES clan_groups (id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 리프레시 토큰 테이블
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users (id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 앨범 테이블
CREATE TABLE IF NOT EXISTS albums (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups (id),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 앨범 클랜 그룹 권한 테이블
CREATE TABLE IF NOT EXISTS album_clan_groups_permissions (
    album_id INTEGER NOT NULL REFERENCES albums (id),
    clan_group_id INTEGER NOT NULL REFERENCES clan_groups (id),
    permission VARCHAR(1) NOT NULL, -- R: 읽기, W: 쓰기
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 업로드 배치 테이블
CREATE TABLE IF NOT EXISTS upload_batches (
    id SERIAL PRIMARY KEY,
    album_id INTEGER NOT NULL REFERENCES albums (id),
    upload_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 미디어 아이템 테이블
CREATE TABLE IF NOT EXISTS media_items (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups (id),
    album_id INTEGER NOT NULL REFERENCES albums (id),
    upload_batch_id INTEGER NOT NULL REFERENCES upload_batches (id),
    upload_status VARCHAR(2) NOT NULL DEFAULT '01', -- 01: pending | 02: processing | 03: completed | 04: failed | 05: duplicate

    taken_at TIMESTAMP NOT NULL,
    file_name VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- 미디어 파일 테이블
CREATE TABLE IF NOT EXISTS media_files (
    id SERIAL PRIMARY KEY,
    media_item_id INTEGER NOT NULL REFERENCES media_items (id) ON DELETE CASCADE,

    role VARCHAR(2) NOT NULL, -- 01: original | 02: thumbnail | 03: view | 04: live
    storage_key VARCHAR(512) NOT NULL,
    mime_type VARCHAR(50), -- image/jpeg, image/png, video/mp4, video/quicktime, video/mov, video/avi, video/wmv, video/flv, video/webm, video/mkv

    width INTEGER,
    height INTEGER,
    file_size BIGINT,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- 신원 테이블
CREATE TABLE IF NOT EXISTS identities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    group_id INTEGER NOT NULL REFERENCES groups (id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 신원 얼굴 이미지 테이블
CREATE TABLE IF NOT EXISTS identity_face_imgs (
    id SERIAL PRIMARY KEY,
    identity_id INTEGER NOT NULL REFERENCES identities (id),
    media_item_id INTEGER NOT NULL REFERENCES media_items (id),
    img_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 얼굴 감지 테이블
CREATE TABLE IF NOT EXISTS face_detections (
    id SERIAL PRIMARY KEY,
    media_item_id INTEGER NOT NULL REFERENCES media_items (id),
    identity_id INTEGER NOT NULL REFERENCES identities (id),
    location_top INTEGER NOT NULL,
    location_right INTEGER NOT NULL,
    location_bottom INTEGER NOT NULL,
    location_left INTEGER NOT NULL,
    embedding vector (512),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_face_detections_embedding ON face_detections USING ivfflat (embedding vector_cosine_ops);

CREATE INDEX IF NOT EXISTS idx_face_detections_media_item_id ON face_detections (media_item_id);

CREATE INDEX IF NOT EXISTS idx_face_detections_identity_id ON face_detections (identity_id);

CREATE INDEX IF NOT EXISTS idx_media_items_group_album_id_year_month ON media_items (
    group_id,
    album_id,
    EXTRACT(
        YEAR
        FROM
            taken_at
    ),
    EXTRACT(
        MONTH
        FROM
            taken_at
    )
);


CREATE INDEX IF NOT EXISTS idx_media_items_group_year_month ON media_items (
    group_id,
    EXTRACT(
        YEAR
        FROM
            taken_at
    ),
    EXTRACT(
        MONTH
        FROM
            taken_at
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_album_clan_groups_permissions_album_id_clan_group_id ON album_clan_groups_permissions (album_id, clan_group_id, permission);

CREATE INDEX IF NOT EXISTS idx_clan_groups_group_id ON clan_groups (group_id);

CREATE INDEX IF NOT EXISTS idx_users_group_id ON users (group_id);

CREATE INDEX IF NOT EXISTS idx_users_clan_group_id ON users (clan_group_id);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens (user_id);

CREATE INDEX IF NOT EXISTS idx_identities_group_id ON identities (group_id);

CREATE INDEX IF NOT EXISTS idx_media_items_group_id ON media_items (group_id);

CREATE INDEX IF NOT EXISTS idx_media_items_group_album_id ON media_items (group_id, album_id);

CREATE INDEX IF NOT EXISTS idx_media_items_group_created_at ON media_items (group_id, created_at);

CREATE INDEX IF NOT EXISTS idx_media_items_album_id ON media_items (album_id);

CREATE INDEX IF NOT EXISTS idx_media_items_taken_at ON media_items (taken_at);

CREATE INDEX IF NOT EXISTS idx_media_items_upload_batch_id ON media_items (upload_batch_id);

CREATE INDEX IF NOT EXISTS idx_media_items_group_album_id_created_at ON media_items (group_id, album_id, created_at);

CREATE INDEX IF NOT EXISTS idx_media_items_group_created_at ON media_items (group_id, created_at);

CREATE INDEX IF NOT EXISTS idx_media_items_upload_status ON media_items (upload_status);

CREATE INDEX IF NOT EXISTS idx_albums_group_id ON albums (group_id);

CREATE INDEX IF NOT EXISTS idx_album_clan_groups_permissions_album_id ON album_clan_groups_permissions (album_id);

CREATE INDEX IF NOT EXISTS idx_album_clan_groups_permissions_clan_group_id ON album_clan_groups_permissions (clan_group_id);

CREATE INDEX IF NOT EXISTS idx_identity_face_imgs_media_item_id ON identity_face_imgs (media_item_id);

CREATE INDEX IF NOT EXISTS idx_identity_face_imgs_identity_id ON identity_face_imgs (identity_id);

CREATE INDEX IF NOT EXISTS idx_upload_batches_upload_at ON upload_batches (upload_at);

CREATE INDEX IF NOT EXISTS idx_upload_batches_album_id ON upload_batches (album_id);
