INSERT INTO
    families (id, name)
VALUES
    (1, 'lee');

INSERT INTO
    groups (id, family_id, name, is_admin)
VALUES
    (1, 1, '', true),

INSERT INTO
    users (
        id,
        name,
        username,
        password,
        family_id,
        group_id
    )
VALUES
    (1, '용우', 'ywl0806', 'password', 1, 1)

INSERT INTO
    albums (id, family_id, name, is_common)
VALUES
    (1, 1, '', TRUE),

INSERT INTO
    album_groups_permissions (album_id, group_id, permission)
VALUES
    (1, 1, 'R'),
    (1, 1, 'W')

INSERT INTO tags (name, is_preset) VALUES
    ('誕生日',      true),
    ('旅行',      true),
    ('卒業',      true),
    ('運動会',    true),
    ('クリスマス', true),
    ('日常',      true),
    ('初めて',      true),
    ('公園',  true);
