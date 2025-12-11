package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/kapiw04/convenly/internal/domain/event"
	"github.com/lib/pq"
)

type PostgresEventRepo struct {
	DB      *sql.DB
	TagRepo event.TagRepo
}

func NewPostgresEventRepo(db *sql.DB, tr event.TagRepo) *PostgresEventRepo {
	return &PostgresEventRepo{DB: db, TagRepo: tr}
}

func (p *PostgresEventRepo) FindByID(eventId string) (*event.Event, error) {
	e, err := findEvent(p, eventId)
	if err != nil {
		return nil, err
	}
	tags, err := findTagNames(p, eventId)
	if err != nil {
		return nil, err
	}
	e.Tags = tags
	return e, nil
}

func findTagNames(p *PostgresEventRepo, eventId string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	query := "SELECT t.name FROM event_tag et INNER JOIN tags t ON et.tag_id = t.tag_id WHERE et.event_id = $1"

	eid, err := uuid.Parse(eventId)
	if err != nil {
		return nil, err
	}

	rows, err := p.DB.QueryContext(ctx, query, eid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tagNames []string

	for rows.Next() {
		var tagName string
		if err := rows.Scan(&tagName); err != nil {
			return nil, err
		}

		tagNames = append(tagNames, tagName)
	}
	return tagNames, nil
}

func findEvent(p *PostgresEventRepo, eventId string) (*event.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	query := "SELECT event_id, name, description, date, latitude, longitude, fee, organizer_id FROM events WHERE event_id = $1"
	rows, err := p.DB.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, err
	}

	rows.Next()
	var e event.Event
	if err := rows.Scan(&e.EventID, &e.Name, &e.Description, &e.Date, &e.Latitude, &e.Longitude, &e.Fee, &e.OrganizerID); err != nil {
		return nil, err
	}

	return &e, nil
}

func (p *PostgresEventRepo) Save(e *event.Event) error {
	err := saveEvent(e, p)
	if err != nil {
		return err
	}
	return saveEventTag(e, p)
}

func saveEventTag(e *event.Event, p *PostgresEventRepo) error {
	eventId, err := uuid.Parse(e.EventID)
	if err != nil {
		return err
	}
	for _, tag := range e.Tags {
		t, err := p.TagRepo.FindByName(tag)
		if err != nil {
			return err
		}
		if t == nil {
			return sql.ErrNoRows
		}

		query := "INSERT INTO event_tag (event_id, tag_id) VALUES ($1, $2)"
		_, err = p.DB.Exec(query, eventId, t.TagID)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveEvent(e *event.Event, p *PostgresEventRepo) error {
	query := "INSERT INTO events" +
		"(event_id, name, description, date, latitude, longitude, fee, organizer_id)" +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	eventId, err := uuid.Parse(e.EventID)
	if err != nil {
		return err
	}

	organizerId, err := uuid.Parse(e.OrganizerID)
	if err != nil {
		return err
	}
	_, err = p.DB.Exec(
		query,
		eventId,
		e.Name,
		e.Description,
		e.Date,
		e.Latitude,
		e.Longitude,
		e.Fee,
		organizerId,
	)
	return err
}

func (p *PostgresEventRepo) FindAll() ([]*event.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	query := "SELECT event_id, name, description, date, latitude, longitude, fee, organizer_id FROM events"
	rows, err := p.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*event.Event
	for rows.Next() {
		var e event.Event
		if err := rows.Scan(&e.EventID, &e.Name, &e.Description, &e.Date, &e.Latitude, &e.Longitude, &e.Fee, &e.OrganizerID); err != nil {
			return nil, err
		}
		events = append(events, &e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, e := range events {
		tags, err := findTagNames(p, e.EventID)
		if err != nil {
			return nil, err
		}
		e.Tags = tags
	}

	return events, nil
}

func (p *PostgresEventRepo) RegisterAttendance(userID, eventID string) error {
	query := "INSERT INTO attendance (user_id, event_id) VALUES ($1, $2)"

	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	eid, err := uuid.Parse(eventID)
	if err != nil {
		return err
	}

	_, err = p.DB.Exec(query, uid, eid)
	return err
}

func (p *PostgresEventRepo) GetAttendees(eventID string) ([]string, error) {
	query := "SELECT user_id FROM attendance WHERE event_id = $1"

	eid, err := uuid.Parse(eventID)
	if err != nil {
		return nil, err
	}

	rows, err := p.DB.Query(query, eid)
	attendees := []string{}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var uid uuid.UUID
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		attendees = append(attendees, uid.String())
	}
	return attendees, nil
}

func (p *PostgresEventRepo) RemoveAttendance(userID, eventID string) error {
	query := "DELETE FROM attendance WHERE user_id = $1 AND event_id = $2"

	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	eid, err := uuid.Parse(eventID)
	if err != nil {
		return err
	}

	_, err = p.DB.Exec(query, uid, eid)
	return err
}

func (p *PostgresEventRepo) FindAllByTags(tagNames []string) ([]*event.Event, error) {
	if len(tagNames) == 0 {
		return []*event.Event{}, nil
	}

	query := `
		SELECT DISTINCT e.event_id, e.name, e.description, e.date, e.latitude, e.longitude, e.fee, e.organizer_id
		FROM events e 
		INNER JOIN event_tag et ON et.event_id = e.event_id
		INNER JOIN tags t ON t.tag_id = et.tag_id
		WHERE t.name = ANY($1)
	`
	rows, err := p.DB.Query(query, pq.Array(tagNames))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*event.Event
	for rows.Next() {
		var e event.Event
		if err := rows.Scan(&e.EventID, &e.Name, &e.Description, &e.Date, &e.Latitude, &e.Longitude, &e.Fee, &e.OrganizerID); err != nil {
			return nil, err
		}
		events = append(events, &e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

var _ event.EventRepo = &PostgresEventRepo{}
