INSERT INTO
    groups (name)
VALUES
    ('group1');

INSERT INTO
    clan_groups (group_id, name, is_admin)
VALUES
    (1, 'admin', true),
    (1, 'lee', false),
    (1, 'muraoka', false);

INSERT INTO
    users (name, username, password, group_id, clan_group_id)
VALUES
    ('admin', 'admin', 'password', 1, 1),
    ('lee', 'lee', 'password', 1, 2),
    ('muraoka', 'muraoka', 'password', 1, 3);