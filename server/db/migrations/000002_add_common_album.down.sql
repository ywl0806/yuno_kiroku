DROP INDEX IF EXISTS idx_albums_family_id_is_common;
ALTER TABLE albums DROP COLUMN IF EXISTS is_common;
