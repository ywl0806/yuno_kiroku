-- name: CreateMediaFile :one
INSERT INTO
    media_files (
        media_item_id,
        role,
        storage_key,
        width,
        height
    )
VALUES
    ($1, $2, $3, $4, $5)
RETURNING
    id,
    media_item_id,
    role,
    storage_key,
    width,
    height,
    created_at,
    updated_at;

-- name: DeleteMediaFileByItemAndRole :exec
-- S2-05: S3 업로드 실패 시 DB 레코드 롤백용
DELETE FROM media_files WHERE media_item_id = $1 AND role = $2;

-- name: UpsertMediaFile :one
-- S3-04: SQS 재시도 시 중복 INSERT 방지 (ON CONFLICT upsert)
INSERT INTO
    media_files (
        media_item_id,
        role,
        storage_key,
        width,
        height
    )
VALUES
    ($1, $2, $3, $4, $5)
ON CONFLICT (media_item_id, role) DO UPDATE
    SET storage_key = EXCLUDED.storage_key,
        width       = EXCLUDED.width,
        height      = EXCLUDED.height,
        updated_at  = CURRENT_TIMESTAMP
RETURNING
    id,
    media_item_id,
    role,
    storage_key,
    width,
    height,
    created_at,
    updated_at;
