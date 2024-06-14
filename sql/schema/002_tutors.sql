-- +goose Up
CREATE TABLE tutors (
  tutor_id UUID PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  valid BOOLEAN DEFAULT FALSE NOT NULL,
  hashed_password VARCHAR(64) NOT NULL,
  yoe INT NOT NULL,
  subject VARCHAR(50) NOT NULL CHECK (subject IN ('chemistry', 'physics', 'math')),
  verified BOOLEAN DEFAULT FALSE NOT NULL,
  rating FLOAT CHECK (rating >= 1.0 AND rating <= 5.0),
  rating_count INT NOT NULL
);

-- +goose Down
DROP TABLE tutors;

