-- +goose Up
CREATE TABLE students (
  student_id UUID PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  valid BOOLEAN DEFAULT FALSE NOT NULL,
  hashed_password VARCHAR(64) NOT NULL
);

-- +goose Down
DROP TABLE students;