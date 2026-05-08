ALTER TABLE albums ADD COLUMN is_common BOOLEAN NOT NULL DEFAULT FALSE;

CREATE UNIQUE INDEX idx_albums_family_id_is_common
    ON albums (family_id) WHERE is_common = TRUE;

UPDATE albums SET is_common = TRUE WHERE id = 1;
