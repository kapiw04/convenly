CREATE TABLE users (
  user_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  email TEXT UNIQUE NOT NULL
    CONSTRAINT users_email_format CHECK (
      email ~ '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
    ),
  password_hash TEXT NOT NULL,
  name TEXT NOT NULL
    CONSTRAINT users_name_len CHECK (char_length(name) > 4),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
