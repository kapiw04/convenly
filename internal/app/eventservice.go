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

func (s *EventService) GetEventByTag(tagNames []string) ([]*event.Event, error) {
	return s.eventRepo.FindAllByTags(tagNames)
}

func (s *EventService) GetAllEvents() ([]*event.Event, error) {
	return s.eventRepo.FindAll()
}

func (s *EventService) GetEventsWithFilters(filter *event.EventFilter) ([]*event.Event, error) {
	return s.eventRepo.FindAllWithFilters(filter)
}

func (s *EventService) RegisterAttendance(userID, eventID string) error {
	return s.eventRepo.RegisterAttendance(userID, eventID)
}

func (s *EventService) GetAttendees(eventID string) ([]string, error) {
	return s.eventRepo.GetAttendees(eventID)
}

func (s *EventService) RemoveAttendance(userID, eventID string) error {
	return s.eventRepo.RemoveAttendance(userID, eventID)
}

func (s *EventService) GetHostingEvents(userID string) ([]*event.Event, error) {
	return s.eventRepo.FindByOrganizer(userID)
}

func (s *EventService) GetAttendingEvents(userID string) ([]*event.Event, error) {
	return s.eventRepo.FindAttendingEvents(userID)
}

func (s *EventService) DeleteEvent(eventID string) error {
	return s.eventRepo.Delete(eventID)
}
