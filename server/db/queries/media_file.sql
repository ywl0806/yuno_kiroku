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
