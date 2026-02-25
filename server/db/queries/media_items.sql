-- name: CreateMediaItem :one
INSERT INTO
    media_items (
        family_id,
        album_id,
        upload_batch_id,
        taken_at,
        file_name
    )
VALUES
    ($1, $2, $3, $4, $5)
RETURNING
    id,
    family_id,
    album_id,
    upload_batch_id,
    upload_status,
    taken_at,
    file_name,
    created_at,
    updated_at;

-- name: GetMediaItemsByTakenAt :many
SELECT
    mi.*,
    COALESCE(original_media_file.storage_key, '') AS original_storage_key,
    COALESCE(thumbnail_media_file.storage_key, '') AS thumbnail_storage_key,
    COALESCE(view_media_file.storage_key, '') AS view_storage_key,
    COALESCE(original_media_file.width, 0) AS original_width,
    COALESCE(original_media_file.height, 0) AS original_height,
    COALESCE(thumbnail_media_file.width, 0) AS thumbnail_width,
    COALESCE(thumbnail_media_file.height, 0) AS thumbnail_height,
    COALESCE(view_media_file.width, 0) AS view_width,
    COALESCE(view_media_file.height, 0) AS view_height
FROM
    media_items AS mi
    INNER JOIN albums AS a ON mi.album_id = a.id
    INNER JOIN album_groups_permissions AS agp ON a.id = agp.album_id
    LEFT JOIN LATERAL (
        SELECT storage_key, width, height, media_item_id
        FROM media_files 
        WHERE role = '01' 
    ) AS original_media_file ON original_media_file.media_item_id = mi.id
    LEFT JOIN LATERAL (
        SELECT storage_key, width, height, media_item_id
        FROM media_files 
        WHERE role = '02' 
    ) AS thumbnail_media_file ON thumbnail_media_file.media_item_id = mi.id
    LEFT JOIN LATERAL (
        SELECT storage_key, width, height, media_item_id
        FROM media_files 
        WHERE role = '03' 
    ) AS view_media_file ON view_media_file.media_item_id = mi.id
WHERE
    agp.group_id = sqlc.arg (group_id)::int
    AND agp.permission = 'R'
    AND mi.taken_at >= sqlc.arg (taken_at_from)::timestamp
    AND mi.taken_at <= sqlc.arg (taken_at_to)::timestamp
    AND mi.upload_status = '03' -- 03: completed
ORDER BY
    mi.taken_at DESC;

-- name: GetMediaItemRange :many
SELECT
    EXTRACT(
        YEAR
        FROM
            taken_at
    ) AS year,
    EXTRACT(
        MONTH
        FROM
            taken_at
    ) AS month
FROM
    media_items AS mi
    INNER JOIN albums AS a ON mi.album_id = a.id
    INNER JOIN album_groups_permissions AS agp ON a.id = agp.album_id
WHERE
    agp.group_id = sqlc.arg (group_id)::int
    AND agp.permission = 'R'
GROUP BY
    year,
    month
ORDER BY
    year DESC,
    month DESC;

-- name: GetMediaItemByFaceDetection :one
SELECT
    mi.*
FROM
    media_items AS mi
    INNER JOIN face_detections AS fd ON mi.id = fd.media_item_id
WHERE
    mi.family_id = sqlc.arg(family_id)::int
    AND fd.embedding = ANY(sqlc.arg(embeddings)::vector[])
GROUP BY
    mi.id
LIMIT 1;

-- name: UpdateMediaItemUploadStatus :one
UPDATE media_items
SET upload_status = $2
WHERE id = $1
RETURNING id, upload_status, created_at, updated_at;

-- name: GetUploadStatuses :many
SELECT
    mi.id,
    mi.upload_status
FROM media_items AS mi
WHERE mi.upload_batch_id = $1
ORDER BY mi.id ASC;
