CREATE TABLE event_tag (
    event_id UUID REFERENCES events(event_id),
    tag_id UUID REFERENCES tags(tag_id),
    PRIMARY KEY (event_id, tag_id)
)