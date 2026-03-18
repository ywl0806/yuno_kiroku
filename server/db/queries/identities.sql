-- name: CreateIdentity :one
INSERT INTO
    identities (family_id)
VALUES
    ($1)
RETURNING
    *;

-- name: FindIdentityByIdAndFamilyId :one
SELECT
    *
FROM
    identities
WHERE
    id = sqlc.arg (id)::int
    AND family_id = sqlc.arg (family_id)::int;

-- name: FindIdentitiesByFamilyId :many
SELECT
    *
FROM
    identities
WHERE
    family_id = sqlc.arg (family_id)::int;
