CREATE TABLE events (
    event_id UUID PRIMARY KEY,
    name TEXT UNIQUE,
    description TEXT,
    date DATE,
    latitude DECIMAL,
    longitude DECIMAL,
    fee DECIMAL,
    organizer_id UUID REFERENCES users(user_id) ON DELETE CASCADE
);
