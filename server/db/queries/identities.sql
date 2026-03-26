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

-- name: GetIdentityOptions :many
SELECT DISTINCT ON (i.id)
    i.id, 
    k.id AS kid_id,
    k.name AS kid_name,
    u.id AS user_id,
    u.name AS user_name,
    ifi.storage_key
FROM 
    identities AS i
    LEFT JOIN kids AS k ON i.id = k.identity_id
    LEFT JOIN users AS u ON i.id = u.identity_id
    LEFT JOIN identity_face_imgs AS ifi ON i.id = ifi.identity_id
    JOIN media_items AS mi ON ifi.media_item_id = mi.id
WHERE 
    i.family_id = $1
ORDER BY 
    i.id,
    mi.taken_at DESC;