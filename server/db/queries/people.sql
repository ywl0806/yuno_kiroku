-- name: CreatePerson :one
INSERT INTO
    people (name, group_id)
VALUES
    ($1, $2)
RETURNING
    *;

-- name: FindPersonByID :one
SELECT
    id,
    name,
    group_id
FROM
    people
WHERE
    id = sqlc.arg (id)::int;