package event

import "time"

type Event struct {
	EventID     string    `json:"event_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Fee         float32   `json:"fee"`
	OrganizerID string    `json:"organizer_id"`
}

type EventRepo interface {
	Save(*Event) error
	FindByID(string) (*Event, error)
	FindAll() ([]*Event, error)
	RegisterAttendance(userID, eventID string) error
	GetAttendees(eventID string) ([]string, error)
	RemoveAttendance(userID, eventID string) error
}
