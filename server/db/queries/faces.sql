-- name: CreateFaceDetection :one
INSERT INTO
    face_detections (
        media_item_id,
        identity_id,
        location_top,
        location_right,
        location_bottom,
        location_left,
        embedding
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    *;

-- name: FindMostSimilarFace :one
SELECT
    fd.identity_id,
    i.name,
    i.group_id,
    fd.embedding <=> sqlc.arg (embedding)::vector AS distance
FROM
    face_detections AS fd
    INNER JOIN identities AS i ON fd.identity_id = i.id
WHERE
    i.group_id = sqlc.arg (group_id)::int
    AND fd.embedding <=> sqlc.arg (embedding)::vector < sqlc.arg (similarity_threshold)::float
ORDER BY
    fd.embedding <=> sqlc.arg (embedding)::vector ASC
LIMIT
    1;

-- name: GetFaceDetectionsByEmbeddings :one
SELECT
    fd.embedding
FROM
    face_detections AS fd
    JOIN media_items AS mi ON fd.media_item_id = mi.id
WHERE
    mi.group_id = sqlc.arg(group_id)::int
    AND mi.album_id = sqlc.arg(album_id)::int
    AND fd.embedding = ANY(sqlc.arg(embeddings)::vector[]);