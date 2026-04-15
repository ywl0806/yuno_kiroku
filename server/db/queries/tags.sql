-- name: GetTagsByFamilyID :many
SELECT * FROM tags
WHERE family_id = $1 OR is_preset = true
ORDER BY is_preset DESC, id ASC;

-- name: CreateTag :one
INSERT INTO tags (family_id, name, is_preset, created_by)
VALUES ($1, $2, false, $3)
RETURNING *;

-- name: DeleteTag :exec
DELETE FROM tags
WHERE id = $1 AND family_id = $2;

-- name: AddTagToMediaItem :exec
INSERT INTO media_item_tags (media_item_id, tag_id, tagged_by)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;

-- name: RemoveTagFromMediaItem :exec
DELETE FROM media_item_tags
WHERE media_item_id = $1 AND tag_id = $2;

-- name: GetTagsForMediaItem :many
SELECT t.* FROM tags t
JOIN media_item_tags mt ON t.id = mt.tag_id
WHERE mt.media_item_id = $1
ORDER BY t.id ASC;
