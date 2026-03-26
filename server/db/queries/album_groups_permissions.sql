-- name: GetAlbumGroupPermissions :many
SELECT
    *
FROM
    album_groups_permissions
WHERE
    album_id = $1;

-- name: InsertAlbumGroupPermission :exec
INSERT INTO
    album_groups_permissions (album_id, group_id, permission)
VALUES
    ($1, $2, $3) ON CONFLICT (album_id, group_id, permission) DO NOTHING;

-- name: DeleteAlbumGroupPermissionsByAlbumID :exec
DELETE FROM album_groups_permissions
WHERE
    album_id = $1;
