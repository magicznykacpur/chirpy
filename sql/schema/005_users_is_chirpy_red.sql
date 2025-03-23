-- +goose Up
ALTER TABLE users ADD COLUMN is_chirpy_red boolean;
UPDATE users SET is_chirpy_red = FALSE;
ALTER TABLE users ALTER COLUMN is_chirpy_red SET DEFAULT FALSE;

-- +goose Down
ALTER TABLE users DROP COLUMN is_chirpy_red;

