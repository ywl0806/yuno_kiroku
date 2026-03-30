-- name: CreateFaceRecognitionJob :one
INSERT INTO face_recognition_jobs (media_item_id, family_id, view_storage_key)
VALUES ($1, $2, $3)
ON CONFLICT (media_item_id) DO NOTHING
RETURNING *;

-- name: FetchPendingFaceRecognitionJobs :many
UPDATE face_recognition_jobs
SET status = '02',
    attempt_count = attempt_count + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE id IN (
    SELECT id FROM face_recognition_jobs
    WHERE status = '01'
    ORDER BY created_at ASC
    LIMIT $1
    FOR UPDATE SKIP LOCKED
)
RETURNING *;

-- name: CompleteFaceRecognitionJob :exec
UPDATE face_recognition_jobs
SET status = '03',
    completed_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: FailFaceRecognitionJob :exec
UPDATE face_recognition_jobs
SET status = CASE WHEN attempt_count >= 3 THEN '09' ELSE '01' END,
    last_error = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetFaceRecognitionJobByID :one
SELECT * FROM face_recognition_jobs WHERE id = $1 LIMIT 1;

-- name: GetMediaItemIDByStorageKey :one
SELECT media_item_id FROM media_files WHERE storage_key = $1 LIMIT 1;
