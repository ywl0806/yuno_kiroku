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
    fd.id,
    fd.photo_id,
    fd.person_id,
    p.name,
    ph.group_id,
FROM
    face_detections AS fd
    INNER JOIN people AS p ON fd.person_id = p.id
    INNER JOIN photos AS ph ON fd.photo_id = ph.id
WHERE
    ph.group_id = sqlc.arg (group_id)::int
    AND fd.embedding <=> sqlc.arg (embedding)::vector < sqlc.arg (similarity_threshold)::float
ORDER BY
    fd.embedding <=> sqlc.arg (embedding)::vector ASC
LIMIT
    1;