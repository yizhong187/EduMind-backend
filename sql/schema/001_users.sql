-- +goose Up

CREATE TABLE users (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL
  valid BOOLEAN DEFAULT FALSE
);

-- +goose Down
DROP TABLE users;