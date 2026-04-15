INSERT INTO
    families (id, name)
VALUES
    (1, 'family1');

INSERT INTO
    groups (id, family_id, name, is_admin)
VALUES
    (1, 1, 'admin', true),
    (2, 1, 'lee', false),
    (3, 1, 'muraoka', false);

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

INSERT INTO
    albums (id, family_id, name)
VALUES
    (1, 1, '공용'),
    (2, 1, 'lee'),
    (3, 1, 'muraoka'),
    (4, 1, 'admin');

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
    ('생일',      true),
    ('여행',      true),
    ('졸업',      true),
    ('운동회',    true),
    ('크리스마스', true),
    ('일상',      true),
    ('처음',      true),
    ('가족모임',  true),
    ('학교',      true),
    ('방학',      true),
    ('명절',      true),
    ('외식',      true),
    ('입학',      true),
    ('졸업식',    true),
    ('생일파티',  true);
