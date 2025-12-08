INSERT INTO
    groups (id, name)
VALUES
    (1, 'group1');

INSERT INTO
    clan_groups (id, group_id, name, is_admin)
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
        group_id,
        clan_group_id
    )
VALUES
    (1, 'admin', 'admin', 'password', 1, 1),
    (2, 'lee', 'lee', 'password', 1, 2),
    (3, 'muraoka', 'muraoka', 'password', 1, 3);

INSERT INTO
    albums (id, group_id, name)
VALUES
    (1, 1, 'album 공용'),
    (2, 1, 'album lee'),
    (3, 1, 'album muraoka'),
    (4, 1, 'album admin');

INSERT INTO
    album_clan_groups_permissions (album_id, clan_group_id, permission)
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