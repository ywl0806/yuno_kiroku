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
