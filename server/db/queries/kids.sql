-- name: FindKidsByFamilyID :many
SELECT * FROM kids WHERE family_id = $1 ORDER BY id;

-- name: CreateKid :one
INSERT INTO kids (name, birth_date, identity_id, family_id) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateKid :one
UPDATE kids SET name = $2, birth_date = $3, identity_id = $4 WHERE id = $1 RETURNING *;

-- name: DeleteKid :exec
DELETE FROM kids WHERE id = $1;

-- name: GetKidByID :one
SELECT * FROM kids WHERE id = $1;

-- name: GetKidsWithRandomFaceImg :many
SELECT DISTINCT ON (k.id)
    k.id,
    k.name,
    k.birth_date,
    k.identity_id,
    k.family_id,
    ifi.media_item_id,
    ifi.storage_key,
    mi.taken_at
FROM kids AS k 
JOIN identity_face_imgs AS ifi ON k.identity_id = ifi.identity_id
JOIN media_items AS mi ON ifi.media_item_id = mi.id
WHERE k.family_id = $1
    AND CASE WHEN sqlc.narg(taken_at_to)::timestamp IS NOT NULL THEN
        mi.taken_at <= sqlc.narg(taken_at_to)::timestamp
        ELSE TRUE
    END
    AND CASE WHEN sqlc.narg(taken_at_from)::timestamp IS NOT NULL THEN
        mi.taken_at >= sqlc.narg(taken_at_from)::timestamp
        ELSE TRUE
    END
ORDER BY k.id, RANDOM();