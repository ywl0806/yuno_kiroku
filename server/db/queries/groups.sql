-- name: FindGroupByID :one
SELECT
    id,
    name
FROM
    groups
WHERE
    id = $1;