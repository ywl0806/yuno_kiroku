-- name: GetAlbumGroupPermissions :many
SELECT
    agp.*
FROM
    album_groups_permissions as agp
    INNER JOIN albums as a on agp.album_id = a.id AND a.family_id = $2
WHERE
    agp.album_id = $1;

-- name: InsertAlbumGroupPermission :exec
INSERT INTO
    album_groups_permissions (album_id, group_id, permission)
VALUES
    ($1, $2, $3) ON CONFLICT (album_id, group_id, permission) DO NOTHING;

-- name: DeleteAlbumGroupPermissionsByAlbumID :exec
DELETE FROM album_groups_permissions
WHERE
    album_id = $1;

-- name: GetWritableAlbumIDsByGroupID :many
SELECT album_id FROM album_groups_permissions
WHERE group_id = $1 AND permission = 'W';

-- name: CheckUserHasPermissionForAlbum :one
SELECT EXISTS(
    SELECT 1 FROM album_groups_permissions as agp
    INNER JOIN groups as g on agp.group_id = g.id
    INNER JOIN users as u on g.id = u.group_id AND u.id = sqlc.arg(user_id)::uuid
    WHERE agp.album_id = sqlc.arg(album_id)::uuid AND agp.permission = sqlc.arg(permission)::varchar
) AS has_permission;