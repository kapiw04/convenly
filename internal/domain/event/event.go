package event

import "time"

type Event struct {
	EventID     string
	Name        string
	Description string
	Date        time.Time
	Latitude    float64
	Longitude   float64
	Fee         float32
	OrganizerID string
}

type EventRepo interface {
	Save(*Event) error
	FindByID(string) (*Event, error)
	FindAll() ([]*Event, error)
}
