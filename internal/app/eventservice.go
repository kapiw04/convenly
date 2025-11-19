package app

import "github.com/kapiw04/convenly/internal/domain/event"

type EventService struct {
	eventRepo event.EventRepo
}

func NewEventService(repo event.EventRepo) *EventService {
	return &EventService{eventRepo: repo}
}

func (s *EventService) CreateEvent(e *event.Event) error {
	return s.eventRepo.Save(e)
}

func (s *EventService) GetEventByID(eventID string) (*event.Event, error) {
	return s.eventRepo.FindByID(eventID)
}

func (s *EventService) GetAllEvents() ([]*event.Event, error) {
	return s.eventRepo.FindAll()
}
