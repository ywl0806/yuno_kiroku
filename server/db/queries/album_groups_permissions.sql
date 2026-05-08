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
