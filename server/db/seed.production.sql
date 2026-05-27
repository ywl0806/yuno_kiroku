DO $$
DECLARE
    v_family_id uuid;
    v_group_id  integer;
    v_album_id  uuid;
BEGIN
    INSERT INTO families
    DEFAULT VALUES
    RETURNING id INTO v_family_id;

    INSERT INTO groups (family_id, name, is_admin)
    VALUES (v_family_id, '', true)
    RETURNING id INTO v_group_id;

    INSERT INTO users (name, username, password, family_id, group_id)
    VALUES ('용우', 'ywl0806', 'Y@ml9723623', v_family_id, v_group_id);

    INSERT INTO albums (family_id, name, is_common)
    VALUES (v_family_id, '', true)
    RETURNING id INTO v_album_id;

    INSERT INTO album_groups_permissions (album_id, group_id, permission)
    VALUES
        (v_album_id, v_group_id, 'R'),
        (v_album_id, v_group_id, 'W');
END $$;

INSERT INTO tags (name, is_preset) VALUES
    ('誕生日',      true),
    ('旅行',        true),
    ('卒業',        true),
    ('運動会',      true),
    ('クリスマス',  true),
    ('日常',        true),
    ('初めて',      true),
    ('公園',        true);
