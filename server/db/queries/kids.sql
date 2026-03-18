-- name: FindKidsByFamilyID :many
SELECT * FROM kids WHERE family_id = $1 ORDER BY id;
