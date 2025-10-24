CREATE TABLE atttendances (
    event_id UUID REFERENCES events(event_id),
    user_id UUID REFERENCES users(user_id),
    CONSTRAINT event_user_pkey PRIMARY KEY (event_id, user_id)
);