-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS clan_groups (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups (id),
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

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

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users (id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS photos (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups (id),
    clan_group_id INTEGER REFERENCES clan_groups (id),
    thumbnail_url VARCHAR(255) NOT NULL,
    original_url VARCHAR(255),
    live_url VARCHAR(255),
    original_live_url VARCHAR(255),
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    orientation INTEGER NOT NULL,
    photo_created_at TIMESTAMP NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS people (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    group_id INTEGER NOT NULL REFERENCES groups (id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS face_detections (
    id SERIAL PRIMARY KEY,
    photo_id INTEGER NOT NULL REFERENCES photos (id),
    person_id INTEGER NOT NULL REFERENCES people (id),
    location_top INTEGER NOT NULL,
    location_right INTEGER NOT NULL,
    location_bottom INTEGER NOT NULL,
    location_left INTEGER NOT NULL,
    embedding vector (512),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_face_detections_embedding ON face_detections USING ivfflat (embedding vector_cosine_ops);

CREATE MATERIALIZED VIEW IF NOT EXISTS average_face_embeddings AS
SELECT
    person_id,
    AVG(embedding)::vector (512) AS embedding
FROM
    face_detections
GROUP BY
    person_id;

CREATE INDEX IF NOT EXISTS idx_average_face_embeddings_embedding ON average_face_embeddings USING ivfflat (embedding vector_cosine_ops);

-- MATERIALIZED VIEW 자동 갱신을 위한 트리거 함수
-- 참고: CONCURRENTLY는 트랜잭션 내에서 실행할 수 없으므로 일반 REFRESH 사용
CREATE OR REPLACE FUNCTION refresh_average_face_embeddings () RETURNS TRIGGER AS $$
BEGIN
    -- 비동기적으로 갱신 (트랜잭션 완료 후 실행)
    -- CONCURRENTLY는 트랜잭션 외부에서만 실행 가능하므로 일반 REFRESH 사용
    REFRESH MATERIALIZED VIEW average_face_embeddings;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- face_detections 테이블 변경 시 자동 갱신 트리거
-- FOR EACH STATEMENT: 각 문장마다 실행 (성능 최적화)
CREATE TRIGGER trigger_refresh_average_face_embeddings
AFTER INSERT
OR
UPDATE
OR DELETE ON face_detections FOR EACH STATEMENT
EXECUTE FUNCTION refresh_average_face_embeddings ();