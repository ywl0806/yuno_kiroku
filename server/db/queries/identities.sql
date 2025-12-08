-- name: CreateIdentity :one
INSERT INTO
    identities (name, group_id)
VALUES
    ($1, $2)
RETURNING
    *;

-- name: FindIdentityByIdAndGroupId :one
SELECT
    *
FROM
    identities
WHERE
    id = sqlc.arg (id)::int
    AND group_id = sqlc.arg (group_id)::int;

-- name: FindIdentitiesByGroupId :many
SELECT
    *
FROM
    identities
WHERE
    group_id = sqlc.arg (group_id)::int;


-- name: UpdateIdentityByIdAndGroupId :one
UPDATE
    identities
SET
    name = sqlc.arg (name)
WHERE
    id = sqlc.arg (id)::int
    AND group_id = sqlc.arg (group_id)::int
RETURNING
    *;