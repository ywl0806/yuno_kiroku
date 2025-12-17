-- name: CreatePhoto :one
INSERT INTO
    photos (
        group_id,
        album_id,
        thumbnail_url,
        original_url,
        live_url,
        original_live_url,
        original_width,
        original_height,
        thumbnail_width,
        thumbnail_height,
        photo_created_at,
        file_name
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING
    id,
    group_id,
    album_id,
    thumbnail_url,
    original_url,
    live_url,
    original_live_url,
    original_width,
    original_height,
    thumbnail_width,
    thumbnail_height,
    photo_created_at,
    file_name,
    created_at,
    updated_at;

-- name: FindPhotosByPhotoCreatedAt :many
SELECT
    p.*
FROM
    photos AS p
    INNER JOIN albums AS a ON p.album_id = a.id
    INNER JOIN album_clan_groups_permissions AS acgp ON a.id = acgp.album_id
WHERE
    acgp.clan_group_id = sqlc.arg (clan_group_id)::int
    AND acgp.permission = 'R'
    AND p.photo_created_at >= sqlc.arg (photo_created_at_from)::timestamp
    AND p.photo_created_at <= sqlc.arg (photo_created_at_to)::timestamp
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
    photos AS p
    INNER JOIN albums AS a ON p.album_id = a.id
    INNER JOIN album_clan_groups_permissions AS acgp ON a.id = acgp.album_id
WHERE
    acgp.clan_group_id = sqlc.arg (clan_group_id)::int
    AND acgp.permission = 'R'
GROUP BY
    year,
    month
ORDER BY
    year DESC,
    month DESC;

-- name: GetIdentityRandomPhoto :one
SELECT
    p.*,
    fd.location_top,
    fd.location_right,
    fd.location_bottom,
    fd.location_left
FROM
    (
        SELECT photo_id, location_top, location_right, location_bottom, location_left
        FROM face_detections 
        WHERE identity_id = sqlc.arg(identity_id)::int
        ORDER BY RANDOM()
        LIMIT 1
    ) AS fd
    INNER JOIN photos AS p ON fd.photo_id = p.id
    INNER JOIN albums AS a ON p.album_id = a.id
    INNER JOIN album_clan_groups_permissions AS acgp ON a.id = acgp.album_id
WHERE
    acgp.clan_group_id = sqlc.arg (clan_group_id)::int
    AND acgp.permission = 'R';

-- name: GetPhotoByFaceDetection :one
SELECT
    p.*
FROM
    photos AS p
    INNER JOIN face_detections AS fd ON p.id = fd.photo_id
WHERE
    p.group_id = sqlc.arg(group_id)::int
    AND fd.embedding = ANY(sqlc.arg(embeddings)::vector[])
GROUP BY
    p.id
LIMIT 1;