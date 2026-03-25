-- name: CreateIdentityFaceImg :one
INSERT INTO
    identity_face_imgs (identity_id, media_item_id, storage_key)
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: HasIdentityFaceImg :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            identity_face_imgs
        WHERE
            identity_id = $1
    ) AS exists;


-- name: FindNewestIdentityFaceImgByFamilyId :many
SELECT DISTINCT ON (ifi.identity_id)
    ifi.id,
    ifi.identity_id,
    ifi.media_item_id,
    ifi.storage_key,
    k.id AS kid_id,
    k.name AS kid_name,
    u.id AS user_id,
    u.name AS user_name
FROM
    identity_face_imgs AS ifi
    JOIN media_items AS mi ON ifi.media_item_id = mi.id
    JOIN identities AS i ON ifi.identity_id = i.id
    LEFT JOIN kids AS k ON i.id = k.identity_id
    LEFT JOIN users AS u ON i.id = u.identity_id
WHERE
    i.family_id = $1
    AND (
        CASE WHEN sqlc.arg (only_not_linked)::boolean THEN
            (k.id IS NULL AND u.id IS NULL)
        END
        OR CASE WHEN sqlc.arg (with_kid_ids)::int[] IS NOT NULL THEN
            (k.id = ALL (sqlc.arg (with_kid_ids)::int[]))
        END
        OR CASE WHEN sqlc.arg (with_user_ids)::int[] IS NOT NULL THEN
            (u.id = ALL (sqlc.arg (with_user_ids)::int[]))
        END
    )
ORDER BY
    ifi.identity_id,
    mi.taken_at DESC;
