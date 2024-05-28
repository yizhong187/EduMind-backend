-- +goose Up
CREATE TABLE users (
  user_id UUID PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  user_type VARCHAR(50) CHECK (user_type IN ('tutor', 'student', 'admin')) NOT NULL
);

-- +goose Down
DROP TABLE users;