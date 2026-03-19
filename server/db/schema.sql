-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- 가족 테이블 (전체 가족 단위)
CREATE TABLE IF NOT EXISTS families (
    id SERIAL PRIMARY KEY, -- 가족 ID
    name VARCHAR(255) NOT NULL, -- 가족 이름
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 생성일시
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 수정일시
);

-- 그룹 테이블 (가족 내 쪽 단위: 아빠 쪽, 엄마 쪽 등)
CREATE TABLE IF NOT EXISTS groups (
    -- 소속 가족 ID
    id SERIAL PRIMARY KEY, -- 그룹 ID
    family_id INTEGER NOT NULL REFERENCES families (id), -- 소속 가족 ID
    is_admin BOOLEAN NOT NULL DEFAULT FALSE, -- 관리자 여부
    name VARCHAR(255) NOT NULL, -- 그룹 이름
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 생성일시
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 수정일시
);

-- 신원 테이블
CREATE TABLE IF NOT EXISTS identities (
    id SERIAL PRIMARY KEY, -- 신원 ID
    family_id INTEGER NOT NULL REFERENCES families (id), -- 소속 가족 ID
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 생성일시
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 수정일시
);

-- 사용자 테이블
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY, -- 사용자 ID
    name VARCHAR(255), -- 표시 이름
    username VARCHAR(255) NOT NULL UNIQUE, -- 로그인 아이디
    password VARCHAR(255) NOT NULL, -- 비밀번호(해시)
    family_id INTEGER NOT NULL REFERENCES families (id), -- 소속 가족 ID
    group_id INTEGER NOT NULL REFERENCES groups (id), -- 소속 그룹 ID
    identity_id INTEGER REFERENCES identities (id), -- 연결된 신원 ID
    provider VARCHAR(50) DEFAULT NULL, -- OAuth 제공자 (line/kakao)
    provider_user_id VARCHAR(255) DEFAULT NULL, -- OAuth 제공자 측 사용자 ID
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 생성일시
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 수정일시
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_provider_user ON users (provider, provider_user_id) WHERE provider IS NOT NULL AND provider_user_id IS NOT NULL;

-- 아이 테이블
CREATE TABLE IF NOT EXISTS kids (
    id SERIAL PRIMARY KEY, -- 아이 ID
    name VARCHAR(255), -- 아이 이름
    birth_date DATE, -- 생년월일
    identity_id INTEGER REFERENCES identities (id), -- 연결된 신원 ID
    family_id INTEGER NOT NULL REFERENCES families (id), -- 소속 가족 ID
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 생성일시
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 수정일시
);

-- 초대 토큰 테이블 (초대 링크로 가족/그룹 가입용)
CREATE TABLE IF NOT EXISTS invite_tokens (
    id SERIAL PRIMARY KEY, -- 초대 토큰 ID
    token VARCHAR(64) NOT NULL UNIQUE, -- 초대 토큰 값
    family_id INTEGER NOT NULL REFERENCES families (id), -- 초대 대상 가족 ID
    group_id INTEGER NOT NULL REFERENCES groups (id), -- 초대 대상 그룹 ID
    created_by_user_id INTEGER NOT NULL REFERENCES users (id), -- 초대 생성 사용자 ID
    expires_at TIMESTAMP NOT NULL, -- 만료일시
    used_at TIMESTAMP DEFAULT NULL, -- 사용일시
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 생성일시
);
CREATE INDEX IF NOT EXISTS idx_invite_tokens_token ON invite_tokens (token);
CREATE INDEX IF NOT EXISTS idx_invite_tokens_expires_at ON invite_tokens (expires_at);

-- 리프레시 토큰 테이블
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY, -- 리프레시 토큰 ID
    token VARCHAR(255) NOT NULL, -- 토큰 값
    expires_at TIMESTAMP NOT NULL, -- 만료일시
    user_id INTEGER NOT NULL REFERENCES users (id), -- 소유 사용자 ID
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 생성일시
);

-- 앨범 테이블
CREATE TABLE IF NOT EXISTS albums (
    id SERIAL PRIMARY KEY, -- 앨범 ID
    family_id INTEGER NOT NULL REFERENCES families (id), -- 소속 가족 ID
    name VARCHAR(255) NOT NULL, -- 앨범 이름
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 생성일시
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 수정일시
);

-- 앨범 그룹 권한 테이블
CREATE TABLE IF NOT EXISTS album_groups_permissions (
    album_id INTEGER NOT NULL REFERENCES albums (id), -- 앨범 ID
    group_id INTEGER NOT NULL REFERENCES groups (id), -- 그룹 ID
    permission VARCHAR(1) NOT NULL, -- 권한 (R: 읽기, W: 쓰기)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 생성일시
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 수정일시
);

-- 업로드 배치 테이블
CREATE TABLE IF NOT EXISTS upload_batches (
    id SERIAL PRIMARY KEY, -- 업로드 배치 ID
    album_id INTEGER NOT NULL REFERENCES albums (id), -- 대상 앨범 ID
    upload_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 업로드 일시

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 생성일시
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 수정일시
);

-- 미디어 아이템 테이블
CREATE TABLE IF NOT EXISTS media_items (
    id SERIAL PRIMARY KEY, -- 미디어 아이템 ID
    family_id INTEGER NOT NULL REFERENCES families (id), -- 소속 가족 ID
    album_id INTEGER NOT NULL REFERENCES albums (id), -- 소속 앨범 ID
    upload_batch_id INTEGER NOT NULL REFERENCES upload_batches (id), -- 업로드 배치 ID
    upload_status VARCHAR(2) NOT NULL DEFAULT '01', -- 업로드 상태 (01: pending | 02: processing | 03: completed | 04: failed | 05: duplicate)
    taken_location_latitude DOUBLE PRECISION, -- 촬영 위치 위도
    taken_location_longitude DOUBLE PRECISION, -- 촬영 위치 경도
    taken_at TIMESTAMP NOT NULL, -- 촬영 일시
    file_name VARCHAR(255), -- 원본 파일명
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 생성일시
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 수정일시
);


-- 미디어 파일 테이블
CREATE TABLE IF NOT EXISTS media_files (
    id SERIAL PRIMARY KEY, -- 미디어 파일 ID
    media_item_id INTEGER NOT NULL REFERENCES media_items (id) ON DELETE CASCADE, -- 소속 미디어 아이템 ID

    role VARCHAR(2) NOT NULL, -- 파일 역할 (01: original | 02: thumbnail | 03: view | 04: live)
    storage_key VARCHAR(512) NOT NULL, -- 스토리지 저장 경로
    mime_type VARCHAR(50), -- MIME 타입 (image/jpeg, image/png, video/mp4, ...)

    width INTEGER, -- 가로 크기(px)
    height INTEGER, -- 세로 크기(px)
    file_size BIGINT, -- 파일 크기(bytes)

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 생성일시
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 수정일시
);

-- 신원 얼굴 이미지 테이블
CREATE TABLE IF NOT EXISTS identity_face_imgs (
    id SERIAL PRIMARY KEY, -- 얼굴 이미지 ID
    identity_id INTEGER NOT NULL REFERENCES identities (id), -- 연결된 신원 ID
    media_item_id INTEGER NOT NULL REFERENCES media_items (id), -- 소속 미디어 아이템 ID
    storage_key VARCHAR(255) NOT NULL, -- 스토리지 저장 경로
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 생성일시
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 수정일시
);

-- 얼굴 감지 테이블
CREATE TABLE IF NOT EXISTS face_detections (
    id SERIAL PRIMARY KEY, -- 얼굴 감지 ID
    media_item_id INTEGER NOT NULL REFERENCES media_items (id), -- 소속 미디어 아이템 ID
    identity_id INTEGER NOT NULL REFERENCES identities (id), -- 감지된 신원 ID
    location_top INTEGER NOT NULL, -- 얼굴 위치 상단(px)
    location_right INTEGER NOT NULL, -- 얼굴 위치 우측(px)
    location_bottom INTEGER NOT NULL, -- 얼굴 위치 하단(px)
    location_left INTEGER NOT NULL, -- 얼굴 위치 좌측(px)
    embedding vector (512), -- 얼굴 특징 벡터(512차원)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 생성일시
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 수정일시
);

CREATE INDEX IF NOT EXISTS idx_face_detections_embedding ON face_detections USING ivfflat (embedding vector_cosine_ops);

CREATE INDEX IF NOT EXISTS idx_face_detections_media_item_id ON face_detections (media_item_id);

CREATE INDEX IF NOT EXISTS idx_face_detections_identity_id ON face_detections (identity_id);

CREATE INDEX IF NOT EXISTS idx_media_items_family_album_id_year_month ON media_items (
    family_id,
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


CREATE INDEX IF NOT EXISTS idx_media_items_family_year_month ON media_items (
    family_id,
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

CREATE UNIQUE INDEX IF NOT EXISTS idx_album_groups_permissions_album_id_group_id ON album_groups_permissions (album_id, group_id, permission);

CREATE INDEX IF NOT EXISTS idx_groups_family_id ON groups (family_id);

CREATE INDEX IF NOT EXISTS idx_users_family_id ON users (family_id);

CREATE INDEX IF NOT EXISTS idx_users_group_id ON users (group_id);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens (user_id);

CREATE INDEX IF NOT EXISTS idx_identities_family_id ON identities (family_id);

CREATE INDEX IF NOT EXISTS idx_media_items_family_id ON media_items (family_id);

CREATE INDEX IF NOT EXISTS idx_media_items_family_album_id ON media_items (family_id, album_id);

CREATE INDEX IF NOT EXISTS idx_media_items_family_created_at ON media_items (family_id, created_at);

CREATE INDEX IF NOT EXISTS idx_media_items_album_id ON media_items (album_id);

CREATE INDEX IF NOT EXISTS idx_media_items_taken_at ON media_items (taken_at);

CREATE INDEX IF NOT EXISTS idx_media_items_upload_batch_id ON media_items (upload_batch_id);

CREATE INDEX IF NOT EXISTS idx_media_items_family_album_id_created_at ON media_items (family_id, album_id, created_at);

CREATE INDEX IF NOT EXISTS idx_media_items_family_created_at ON media_items (family_id, created_at);

CREATE INDEX IF NOT EXISTS idx_media_items_upload_status ON media_items (upload_status);

CREATE INDEX IF NOT EXISTS idx_albums_family_id ON albums (family_id);

CREATE INDEX IF NOT EXISTS idx_album_groups_permissions_album_id ON album_groups_permissions (album_id);

CREATE INDEX IF NOT EXISTS idx_album_groups_permissions_group_id ON album_groups_permissions (group_id);

CREATE INDEX IF NOT EXISTS idx_identity_face_imgs_media_item_id ON identity_face_imgs (media_item_id);

CREATE INDEX IF NOT EXISTS idx_identity_face_imgs_identity_id ON identity_face_imgs (identity_id);

CREATE INDEX IF NOT EXISTS idx_upload_batches_upload_at ON upload_batches (upload_at);

CREATE INDEX IF NOT EXISTS idx_upload_batches_album_id ON upload_batches (album_id);
