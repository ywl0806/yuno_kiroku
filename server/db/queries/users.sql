-- name: FindUsers :many
SELECT
    id,
    name,
    username,
    group_id,
    clan_group_id
FROM
    users;

-- name: FindUserByUsername :one
SELECT
    *
FROM
    users
WHERE
    username = $1;

-- name: CreateUser :one
INSERT INTO
    users (name, username, password, group_id, clan_group_id)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING
    *;
