package event

//go:generate mockgen -destination=./mocks/mock_eventrepo.go -package mock_event . EventRepo

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
	Tags        []string  `json:"tag,omitempty"`
}

type Pagination struct {
	Page     int
	PageSize int
}

func (p *Pagination) Offset() int {
	if p == nil || p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) Limit() int {
	if p == nil || p.PageSize <= 0 {
		return 0
	}
	return p.PageSize
}

type EventFilter struct {
	DateFrom   *time.Time
	DateTo     *time.Time
	MinFee     *float32
	MaxFee     *float32
	Tags       []string
	Pagination *Pagination
}

type EventRepo interface {
	Save(*Event) error
	FindByID(string) (*Event, error)
	FindAll() ([]*Event, error)
	FindAllWithFilters(filter *EventFilter) ([]*Event, error)
	FindAllByTags(tagNames []string) ([]*Event, error)
	RegisterAttendance(userID, eventID string) error
	GetAttendees(eventID string) ([]string, error)
	RemoveAttendance(userID, eventID string) error
	FindByOrganizer(userID string, pagination *Pagination) ([]*Event, error)
	FindAttendingEvents(userID string, pagination *Pagination) ([]*Event, error)
	Delete(eventID string) error
}
