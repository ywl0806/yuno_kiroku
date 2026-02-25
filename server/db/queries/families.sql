-- name: FindFamilyByID :one
SELECT
    id,
    name
FROM
    families
WHERE
    id = $1;

-- name: CreateFamily :one
INSERT INTO
    families (name)
VALUES
    ($1)
RETURNING
    *;
