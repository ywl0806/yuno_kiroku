-- name: FindGroupByID :one
SELECT
    id,
    family_id,
    is_admin,
    name
FROM
    groups
WHERE
    id = $1;

-- name: FindGroupsByFamilyID :many
SELECT
    *
FROM
    groups
WHERE
    family_id = $1
ORDER BY
    id;

-- name: CreateGroup :one
INSERT INTO
    groups (family_id, is_admin, name)
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: UpdateGroup :one
UPDATE groups
SET
    name = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $2
RETURNING
    *;

-- name: DeleteGroup :exec
DELETE FROM groups
WHERE
    id = $1;
