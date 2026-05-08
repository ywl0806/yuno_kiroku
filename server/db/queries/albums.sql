-- name: FindAlbumsForWrite :many
SELECT
    a.*
FROM
    albums as a
    INNER JOIN album_groups_permissions as agp on a.id = agp.album_id
    INNER JOIN groups as g on agp.group_id = g.id
WHERE
    a.family_id = sqlc.arg(family_id)::int
    AND g.id = sqlc.arg(group_id)::int
    AND agp.permission = 'W'
UNION
SELECT
    a.*
FROM
    albums as a
WHERE
    a.family_id = sqlc.arg(family_id)::int
    AND a.is_common = TRUE
ORDER BY
    id;

-- name: GetAlbumsOptions :many
SELECT
    a.id,
    a.name,
    a.is_common
FROM
    albums as a
    INNER JOIN album_groups_permissions as agp on a.id = agp.album_id
    INNER JOIN groups as g on agp.group_id = g.id
WHERE
    a.family_id = sqlc.arg(family_id)::int
    AND g.id = sqlc.arg(group_id)::int
    AND agp.permission = 'R'
UNION
SELECT
    a.id,
    a.name,
    a.is_common
FROM
    albums as a
WHERE
    a.family_id = sqlc.arg(family_id)::int
    AND a.is_common = TRUE
ORDER BY
    id;

-- name: FindAlbumsByFamilyID :many
SELECT
    *
FROM
    albums
WHERE
    family_id = $1
ORDER BY
    id;

-- name: FindAlbumByIDAndFamilyID :one
SELECT
    *
FROM
    albums
WHERE
    id = $1
    AND family_id = $2;

-- name: CreateAlbum :one
INSERT INTO
    albums (family_id, name, is_common)
VALUES
    ($1, $2, $3)
RETURNING
    *;

-- name: UpdateAlbum :one
UPDATE albums
SET
    name = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $2
RETURNING
    *;

-- name: DeleteAlbum :exec
DELETE FROM albums
WHERE
    id = $1;