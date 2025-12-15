CREATE VIEW attendees_count AS
SELECT events.event_id, COUNT(attendance.user_id) 
FROM events LEFT JOIN attendance ON attendance.event_id = events.event_id 
GROUP BY events.event_id;

CREATE VIEW find_event_with_tags AS
SELECT e.event_id, e.name, e.description, e.date, e.latitude, e.longitude, e.fee, e.organizer_id, COALESCE(ARRAY_AGG(DISTINCT t.name ORDER BY t.name), ARRAY[]::text[]) AS tags 
FROM events e
LEFT JOIN event_tag et ON et.event_id = e.event_id
LEFT JOIN tags t ON t.tag_id = et.tag_id
GROUP BY
  e.event_id,
  e.name,
  e.description,
  e.date,
  e.latitude,
  e.longitude,
  e.fee,
  e.organizer_id;
