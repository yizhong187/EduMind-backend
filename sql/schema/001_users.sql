-- +goose Up

CREATE TABLE users (
  id UUID PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  valid BOOLEAN DEFAULT FALSE NOT NULL
);

-- +goose Down
DROP TABLE users;