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


-- name FindNewestIdentityFaceImgByFamilyId :one
SELECT
    identity_face_imgs.id,
    identity_face_imgs.identity_id,
    identity_face_imgs.media_item_id,
    identity_face_imgs.storage_key,
    i.id AS identity_id,
    i.name AS identity_name,
    k.id AS kid_id,
    k.name AS kid_name,
    u.id AS user_id,
    u.name AS user_name

FROM
    identity_face_imgs AS ifi
    INNER JOIN media_items AS mi ON ifi.media_item_id = mi.id
    INNER JOIN identities AS i ON ifi.identity_id = i.id
    OUTER JOIN kids AS k ON i.id = k.identity_id
    OUTER JOIN users AS u ON i.id = u.identity_id
WHERE
    i.family_id = $1
ORDER BY
    mi.taken_at DESC
GROUP BY
    i.id
LIMIT 1;