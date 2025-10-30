-- name: CreatePhoto :one
INSERT INTO
    photos (
        group_id,
        clan_group_id,
        thumbnail_url,
        original_url,
        live_url,
        original_live_url,
        width,
        height,
        orientation,
        photo_created_at,
        file_name
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING
    id,
    group_id,
    clan_group_id,
    thumbnail_url,
    original_url,
    live_url,
    original_live_url,
    width,
    height,
    orientation,
    photo_created_at,
    file_name,
    created_at,
    updated_at;