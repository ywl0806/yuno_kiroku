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
    AND agp.permission = 'W';

-- name: GetAlbumsOptions :many
SELECT
    a.id,
    a.name
FROM
    albums as a
    INNER JOIN album_groups_permissions as agp on a.id = agp.album_id
    INNER JOIN groups as g on agp.group_id = g.id
WHERE
    a.family_id = sqlc.arg(family_id)::int
    AND g.id = sqlc.arg(group_id)::int
    AND agp.permission = 'R';

-- name: FindAlbumsByFamilyID :many
SELECT
    *
FROM
    albums
WHERE
    family_id = $1
ORDER BY
    id;

-- name: FindAlbumByID :one
SELECT
    *
FROM
    albums
WHERE
    id = $1;

-- name: CreateAlbum :one
INSERT INTO
    albums (family_id, name)
VALUES
    ($1, $2)
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