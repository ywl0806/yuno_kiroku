-- name: CreateFaceDetection :one
INSERT INTO
    face_detections (
        photo_id,
        person_id,
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
    ave.person_id,
    p.name,
    p.group_id,
    ave.embedding <=> sqlc.arg (embedding)::vector AS distance
FROM
    average_face_embeddings AS ave
    INNER JOIN people AS p ON ave.person_id = p.id
WHERE
    p.group_id = sqlc.arg (group_id)::int
    AND ave.embedding <=> sqlc.arg (embedding)::vector < sqlc.arg (similarity_threshold)::float
ORDER BY
    ave.embedding <=> sqlc.arg (embedding)::vector ASC
LIMIT
    1;

-- name: GetFaceDetectionsByPhotoId :many
SELECT
    fd.id,
    fd.photo_id,
    fd.person_id,
    p.name,
    fd.location_top,
    fd.location_right,
    fd.location_bottom,
    fd.location_left,
    fd.embedding
FROM
    face_detections AS fd
    JOIN people AS p ON fd.person_id = p.id
WHERE
    fd.photo_id = $1;