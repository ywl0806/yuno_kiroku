-- name: LikeMediaItem :exec
INSERT INTO media_item_likes (media_item_id, user_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: UnlikeMediaItem :exec
DELETE FROM media_item_likes
WHERE media_item_id = $1 AND user_id = $2;

-- name: IsMediaItemLiked :one
SELECT EXISTS(
    SELECT 1 FROM media_item_likes
    WHERE media_item_id = $1 AND user_id = $2
) AS is_liked;
