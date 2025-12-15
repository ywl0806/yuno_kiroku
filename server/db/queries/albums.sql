-- name: FindAlbumsForWrite :many
SELECT
    a.*
FROM
    albums as a
    inner join album_clan_groups_permissions as acgp on a.id = acgp.album_id
    inner join clan_groups as cg on acgp.clan_group_id = cg.id
WHERE
    a.group_id = sqlc.arg(group_id)::int
    AND cg.id = sqlc.arg(clan_group_id)::int
    AND acgp.permission = 'W';