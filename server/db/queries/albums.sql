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