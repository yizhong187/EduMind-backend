CREATE TABLE students (
  student_id UUID PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  valid BOOLEAN DEFAULT FALSE NOT NULL,
  hashed_password VARCHAR(64) NOT NULL
);

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

CREATE TABLE users (
  user_id UUID PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  user_type VARCHAR(50) CHECK (user_type IN ('tutor', 'student', 'admin')) NOT NULL
);

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

CREATE TABLE messages (
  message_id UUID PRIMARY KEY,
  chat_id INT NOT NULL,
  user_id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted BOOLEAN DEFAULT FALSE NOT NULL,
  content TEXT NOT NULL,

  FOREIGN KEY (chat_id) REFERENCES chats(chat_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id)
);
