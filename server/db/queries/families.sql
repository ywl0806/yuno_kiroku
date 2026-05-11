-- name: FindFamilyByID :one
SELECT
    id
FROM
    families
WHERE
    id = $1;

-- name: GetFamilyByID :one
SELECT
    *
FROM
    families
WHERE
    id = $1;

-- name: CreateFamily :one
INSERT INTO
    families
DEFAULT VALUES
RETURNING
    *;
