-- +goose Up
CREATE TABLE chats (
  chat_id SERIAL PRIMARY KEY,
  student_id UUID NOT NULL,
  tutor_id UUID,
  created_at TIMESTAMP NOT NULL,
  subject VARCHAR(50) NOT NULL CHECK (subject IN ('chemistry', 'physics', 'math')),
  topic TEXT,
  header TEXT NOT NULL,
  completed BOOLEAN DEFAULT FALSE NOT NULL,

  FOREIGN KEY (student_id) REFERENCES students(student_id),
  FOREIGN KEY (tutor_id) REFERENCES tutors(tutor_id)
);

-- +goose Down
DROP TABLE chats;