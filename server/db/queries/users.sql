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

-- name: FindUserByProvider :one
SELECT
    *
FROM
    users
WHERE
    provider = $1
    AND provider_user_id = $2;

-- name: CreateUser :one
INSERT INTO
    users (name, username, password, group_id, clan_group_id)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: CreateUserOAuth :one
INSERT INTO
    users (name, username, password, group_id, clan_group_id, provider, provider_user_id)
VALUES
    ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    *;
