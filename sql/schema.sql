CREATE TABLE students (
  student_id UUID PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  valid BOOLEAN DEFAULT FALSE NOT NULL,
  hashed_password VARCHAR(64) NOT NULL,
  photo_url TEXT
);

CREATE TABLE tutors (
  tutor_id UUID PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  valid BOOLEAN DEFAULT FALSE NOT NULL,
  hashed_password VARCHAR(64) NOT NULL,
  verified BOOLEAN DEFAULT FALSE NOT NULL,
  rating FLOAT CHECK (rating >= 1.0 AND rating <= 5.0),
  rating_count INT NOT NULL,
  photo_url TEXT
);

CREATE TABLE subjects (
    subject_id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE tutor_subjects (
    tutor_id UUID REFERENCES tutors(tutor_id),
    subject_id INTEGER REFERENCES subjects(subject_id),
    yoe INT NOT NULL,
    PRIMARY KEY (tutor_id, subject_id)
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
  subject_id INTEGER NOT NULL,
  topic TEXT,
  header TEXT NOT NULL,
  photo_url TEXT,
  completed BOOLEAN DEFAULT FALSE NOT NULL,

  FOREIGN KEY (student_id) REFERENCES students(student_id),
  FOREIGN KEY (tutor_id) REFERENCES tutors(tutor_id),
  FOREIGN KEY (subject_id) REFERENCES subjects(subject_id)
);

CREATE TABLE topics (
    subject_id INTEGER NOT NULL,
    topic_id SERIAL,
    name TEXT NOT NULL,
    PRIMARY KEY (subject_id, topic_id),
    FOREIGN KEY (subject_id) REFERENCES subjects(subject_id),
    CONSTRAINT topics_topic_id_unique UNIQUE (topic_id)
);

CREATE TABLE chat_topics (
    chat_id INTEGER NOT NULL REFERENCES chats(chat_id),
    topic_id INTEGER NOT NULL REFERENCES topics(topic_id),
    PRIMARY KEY (chat_id, topic_id)
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
