INSERT INTO
    families (id, name)
VALUES
    (1, 'family1');

SELECT setval('families_id_seq', (SELECT MAX(id) FROM families));

INSERT INTO
    groups (id, family_id, name, is_admin)
VALUES
    (1, 1, 'admin', true),
    (2, 1, '이가', false),
    (3, 1, '村岡家', false);

SELECT setval('groups_id_seq', (SELECT MAX(id) FROM groups));

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
    (1, 'admin', 'admin', 'password', 1, 1),
    (2, 'yongwoo', 'lee', 'password', 1, 2),
    (3, 'yukina', 'muraoka', 'password', 1, 3);

SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));

INSERT INTO
    albums (id, family_id, name, is_common)
VALUES
    (1, 1, '공용', TRUE),
    (2, 1, 'lee', FALSE),
    (3, 1, 'muraoka', FALSE),
    (4, 1, 'admin', FALSE);

SELECT setval('albums_id_seq', (SELECT MAX(id) FROM albums));

INSERT INTO
    album_groups_permissions (album_id, group_id, permission)
VALUES
    (1, 1, 'R'),
    (1, 1, 'W'),
    (1, 2, 'R'),
    (1, 3, 'R'),
    (2, 1, 'R'),
    (2, 1, 'W'),
    (2, 2, 'R'),
    (2, 2, 'W'),
    (3, 1, 'R'),
    (3, 1, 'W'),
    (3, 3, 'R'),
    (3, 3, 'W'),
    (4, 1, 'R'),
    (4, 1, 'W');

INSERT INTO tags (name, is_preset) VALUES
    ('誕生日',      true),
    ('旅行',      true),
    ('卒業',      true),
    ('運動会',    true),
    ('クリスマス', true),
    ('日常',      true),
    ('初めて',      true),
    ('公園',  true);
