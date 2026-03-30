-- name: CreateMediaItem :one
INSERT INTO
    media_items (
        family_id,
        album_id,
        upload_batch_id,
        taken_at,
        file_name,
        taken_location_latitude,
        taken_location_longitude
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    id,
    family_id,
    album_id,
    upload_batch_id,
    upload_status,
    taken_location_latitude,
    taken_location_longitude,
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
    agp.group_id = sqlc.arg(group_id)::int
    AND agp.permission = 'R'
    AND mi.taken_at >= sqlc.arg(taken_at_from)::timestamp
    AND mi.taken_at <= sqlc.arg(taken_at_to)::timestamp
    AND mi.upload_status = '03' -- 03: completed
    AND (sqlc.narg(album_id)::int IS NULL OR mi.album_id = sqlc.narg(album_id)::int)
    AND (
        cardinality(sqlc.arg(identity_ids)::int[]) = 0
        OR EXISTS (
            SELECT 1 FROM face_detections fd
            WHERE fd.media_item_id = mi.id
              AND fd.identity_id = ANY(sqlc.arg(identity_ids)::int[])
        )
    )
ORDER BY
    mi.taken_at DESC;

-- name: GetMediaItemsByTakenAtHome :many
SELECT
    mi.*,
    COALESCE(mf_orig.storage_key, '')  AS original_storage_key,
    COALESCE(mf_thumb.storage_key, '') AS thumbnail_storage_key,
    COALESCE(mf_view.storage_key, '')  AS view_storage_key,
    COALESCE(mf_orig.width, 0)   AS original_width,
    COALESCE(mf_orig.height, 0)  AS original_height,
    COALESCE(mf_thumb.width, 0)  AS thumbnail_width,
    COALESCE(mf_thumb.height, 0) AS thumbnail_height,
    COALESCE(mf_view.width, 0)   AS view_width,
    COALESCE(mf_view.height, 0)  AS view_height
FROM
    media_items AS mi
    INNER JOIN albums AS a ON mi.album_id = a.id
    INNER JOIN album_groups_permissions AS agp ON a.id = agp.album_id
    LEFT JOIN media_files AS mf_orig  ON mf_orig.media_item_id  = mi.id AND mf_orig.role  = '01'
    LEFT JOIN media_files AS mf_thumb ON mf_thumb.media_item_id = mi.id AND mf_thumb.role = '02'
    LEFT JOIN media_files AS mf_view  ON mf_view.media_item_id  = mi.id AND mf_view.role  = '03'
WHERE
    agp.group_id = sqlc.arg(group_id)::int
    AND agp.permission = 'R'
    AND mi.taken_at >= sqlc.arg(taken_at_from)::timestamp
    AND mi.taken_at <= sqlc.arg(taken_at_to)::timestamp
    AND mi.upload_status = '03'
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

-- name: UpdateMediaItemTakenAt :exec
UPDATE media_items
SET taken_at = $2,
    taken_location_latitude = $3,
    taken_location_longitude = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetMediaItemByID :one
SELECT * FROM media_items WHERE id = $1 LIMIT 1;

-- name: SearchMediaItems :many
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
    agp.group_id = sqlc.arg(group_id)::int
    AND agp.permission = 'R'
    AND mi.upload_status = '03'
    AND (sqlc.narg(taken_at_from)::timestamp IS NULL OR mi.taken_at >= sqlc.narg(taken_at_from)::timestamp)
    AND (sqlc.narg(taken_at_to)::timestamp IS NULL OR mi.taken_at <= sqlc.narg(taken_at_to)::timestamp)
    AND (sqlc.narg(album_id)::int IS NULL OR mi.album_id = sqlc.narg(album_id)::int)
    AND (
        cardinality(sqlc.arg(identity_ids)::int[]) = 0
        OR EXISTS (
            SELECT 1 FROM face_detections fd
            WHERE fd.media_item_id = mi.id
              AND fd.identity_id = ANY(sqlc.arg(identity_ids)::int[])
        )
    )
ORDER BY
    mi.taken_at DESC
LIMIT sqlc.arg(page_size)::int
OFFSET sqlc.arg(page_offset)::int;

-- name: GetUploadStatuses :many
SELECT
    mi.id,
    mi.upload_status
FROM media_items AS mi
WHERE mi.upload_batch_id = $1
ORDER BY mi.id ASC;
