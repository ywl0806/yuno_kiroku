-- name: CreateInviteToken :one
INSERT INTO
    invite_tokens (token, family_id, group_id, created_by_user_id, expires_at, family_title, custom_family_title)
VALUES
    ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    *;

-- name: GetInviteTokenByToken :one
SELECT
    *
FROM
    invite_tokens
WHERE
    token = $1
    AND (used_at IS NULL)
    AND (expires_at > NOW());

-- name: MarkInviteTokenUsed :exec
UPDATE
    invite_tokens
SET
    used_at = NOW()
WHERE
    token = $1;
