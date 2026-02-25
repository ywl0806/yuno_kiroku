-- name: FindGroupByID :one
SELECT
    id,
    name
FROM
    groups
WHERE
    id = $1;

-- name: CreateGroup :one
INSERT INTO
    groups (name)
VALUES
    ($1)
RETURNING
    *;