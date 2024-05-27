-- +goose Up
ALTER TABLE users ADD COLUMN hashed_password VARCHAR(64) NOT NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN hashed_password;