-- +goose Up
CREATE TABLE messages (
  message_id UUID PRIMARY KEY,
  chat_id INT NOT NULL,
  user_id UUID UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted BOOLEAN DEFAULT FALSE NOT NULL,
  content TEXT NOT NULL,

  FOREIGN KEY (chat_id) REFERENCES chats(chat_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- +goose Down
DROP TABLE messages;