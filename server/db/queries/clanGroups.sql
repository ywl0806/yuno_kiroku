-- name: FindClanGroupByID :one
SELECT
    id,
    group_id,
    is_admin,
    name
FROM
    clan_groups
WHERE
    id = $1;