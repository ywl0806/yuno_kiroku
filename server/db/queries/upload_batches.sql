-- name: CreateUploadBatch :one
INSERT INTO upload_batches (
    album_id
)
VALUES ($1)
RETURNING id, album_id, upload_at, created_at, updated_at;

-- name: GetUploadBatchesAndMediaItemCounts :many
SELECT
    ub.id,
    ub.album_id,
    ub.upload_at,
    item_counts.count
FROM upload_batches AS ub
INNER JOIN albums AS a ON a.id = ub.album_id
INNER JOIN album_groups_permissions AS agp ON a.id = agp.album_id
INNER JOIN (
    SELECT upload_batch_id, COUNT(*) AS count
    FROM media_items
    WHERE upload_status = '03'
    GROUP BY upload_batch_id
) AS item_counts ON item_counts.upload_batch_id = ub.id
WHERE
    agp.group_id = sqlc.arg(group_id)::int
    AND agp.permission = 'R'
ORDER BY ub.upload_at DESC
LIMIT sqlc.arg(page_size)::int
OFFSET sqlc.arg(page_offset)::int;
