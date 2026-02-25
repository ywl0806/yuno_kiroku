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

-- name: CreateClanGroup :one
INSERT INTO
    clan_groups (group_id, is_admin, name)
VALUES
    ($1, $2, $3)
RETURNING
    *;