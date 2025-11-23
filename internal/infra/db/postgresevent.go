package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/kapiw04/convenly/internal/domain/event"
)

type PostgresEventRepo struct {
	DB *sql.DB
}

func NewPostgresEventRepo(db *sql.DB) *PostgresEventRepo {
	return &PostgresEventRepo{DB: db}
}

func (p *PostgresEventRepo) FindByID(eventId string) (*event.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	query := "SELECT event_id, name, description, date, latitude, longitude, fee, organizer_id FROM events WHERE event_id = $1"
	rows, err := p.DB.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Next()
	var e event.Event
	if err := rows.Scan(&e.EventID, &e.Name, &e.Description, &e.Date, &e.Latitude, &e.Longitude, &e.Fee, &e.OrganizerID); err != nil {
		return nil, err
	}
	return &e, nil
}

func (p *PostgresEventRepo) Save(e *event.Event) error {
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

var _ event.EventRepo = &PostgresEventRepo{}
