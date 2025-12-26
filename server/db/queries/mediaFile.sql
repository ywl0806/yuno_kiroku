-- name: CreateMediaFile :one
INSERT INTO
    media_files (
        media_item_id,
        role,
        storage_key,
        mime_type,
        width,
        height,
        file_size
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    id,
    media_item_id,
    role,
    storage_key,
    mime_type,
    width,
    height,
    file_size,
    created_at,
    updated_at;