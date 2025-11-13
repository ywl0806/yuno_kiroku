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

-- name: FindPhotosByPhotoCreatedAt :many
SELECT
    *
FROM
    photos
WHERE
    group_id = sqlc.arg (group_id)::int
    AND photo_created_at >= sqlc.arg (photo_created_at_from)::time
    AND photo_created_at <= sqlc.arg (photo_created_at_to)::time
    AND (
        CASE
            WHEN sqlc.narg (clan_group_id)::int IS NOT NULL THEN clan_group_id IS NULL
            OR clan_group_id = sqlc.narg (clan_group_id)::int
        END
    )
ORDER BY
    photo_created_at DESC;

-- name: GetPhotoRange :many
SELECT
    EXTRACT(
        YEAR
        FROM
            photo_created_at
    ) AS year,
    EXTRACT(
        MONTH
        FROM
            photo_created_at
    ) AS month
FROM
    photos
WHERE
    group_id = sqlc.arg (group_id)::int
    AND (
        CASE
            WHEN sqlc.narg (clan_group_id)::int IS NOT NULL THEN clan_group_id IS NULL
            OR clan_group_id = sqlc.narg (clan_group_id)::int
        END
    )
GROUP BY
    year,
    month
ORDER BY
    year DESC,
    month DESC;