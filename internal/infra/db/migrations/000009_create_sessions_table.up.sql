CREATE TABLE sessions (
  user_id UUID REFERENCES users(user_id),
  session_id TEXT
); 
