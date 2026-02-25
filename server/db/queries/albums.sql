-- name: FindAlbumsForWrite :many
SELECT
    a.*
FROM
    albums as a
    inner join album_groups_permissions as agp on a.id = agp.album_id
    inner join groups as g on agp.group_id = g.id
WHERE
    a.family_id = sqlc.arg(family_id)::int
    AND g.id = sqlc.arg(group_id)::int
    AND agp.permission = 'W';
