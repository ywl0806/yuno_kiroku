-- name: FindMembersByFamilyID :many
SELECT
    *
FROM
    users
WHERE
    family_id = $1
ORDER BY
    id;

-- name: FindUsers :many
SELECT
    id,
    name,
    username,
    family_id,
    group_id
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
    users (name, username, password, family_id, group_id)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING
    *;

-- name: CreateUserOAuth :one
INSERT INTO
    users (name, username, password, family_id, group_id, provider, provider_user_id, family_title, custom_family_title)
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING
    *;

-- name: FindUserByID :one
SELECT
    *
FROM
    users
WHERE
    id = $1;

-- name: FindUserByIDAndFamilyID :one
SELECT
    *
FROM
    users
WHERE
    id = $1
    AND family_id = $2;
-- name: UpdateUserName :one
UPDATE users
SET
    name = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $2
RETURNING
    *;

-- name: UpdateMember :one
UPDATE users
SET
    group_id = $1,
    family_title = $2,
    custom_family_title = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $4
    AND family_id = $5
RETURNING
    *;
