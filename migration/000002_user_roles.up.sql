-- goose Up
INSERT INTO user_roles(id, type) VALUES (0, 'commom user');
INSERT INTO user_roles(id, type) VALUES (1, 'company user');

-- goose Down
-- DELETE FROM user_roles WHERE id = 0;
-- DELETE FROM user_roles WHERE id = 1;
