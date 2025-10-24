CREATE TABLE events (
    event_id UUID PRIMARY KEY,
    name TEXT UNIQUE,
    description TEXT,
    date DATE,
    geolocation POINT,
    fee DECIMAL,
    organiser_id UUID REFERENCES users(user_id) ON DELETE CASCADE
);
