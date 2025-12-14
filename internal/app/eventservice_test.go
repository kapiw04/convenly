package app

import (
	"errors"
	"testing"
	"time"

	"github.com/kapiw04/convenly/internal/domain/event"
	mock_event "github.com/kapiw04/convenly/internal/domain/event/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestEventService_CreateEvent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	testEvent := &event.Event{
		EventID:     "event-1",
		Name:        "Test Event",
		Description: "Description",
		Date:        time.Now().Add(24 * time.Hour),
		Latitude:    52.0,
		Longitude:   21.0,
		Fee:         10.0,
		OrganizerID: "organizer-1",
	}

	eventRepo.EXPECT().Save(testEvent).Return(nil)

	svc := NewEventService(eventRepo)
	err := svc.CreateEvent(testEvent)

	require.NoError(t, err)
}

func TestEventService_CreateEvent_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	testEvent := &event.Event{
		EventID: "event-1",
		Name:    "Test Event",
	}

	eventRepo.EXPECT().Save(testEvent).Return(errors.New("database error"))

	svc := NewEventService(eventRepo)
	err := svc.CreateEvent(testEvent)

	require.Error(t, err)
}

func TestEventService_GetEventByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	expected := &event.Event{
		EventID: "event-1",
		Name:    "Test Event",
	}

	eventRepo.EXPECT().FindByID("event-1").Return(expected, nil)

	svc := NewEventService(eventRepo)
	result, err := svc.GetEventByID("event-1")

	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestEventService_GetEventByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	eventRepo.EXPECT().FindByID("nonexistent").Return(nil, errors.New("not found"))

	svc := NewEventService(eventRepo)
	_, err := svc.GetEventByID("nonexistent")

	require.Error(t, err)
}

func TestEventService_GetAllEvents_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	expected := []*event.Event{
		{EventID: "event-1", Name: "Event 1"},
		{EventID: "event-2", Name: "Event 2"},
	}

	eventRepo.EXPECT().FindAll().Return(expected, nil)

	svc := NewEventService(eventRepo)
	result, err := svc.GetAllEvents()

	require.NoError(t, err)
	require.Len(t, result, 2)
}

func TestEventService_GetEventsWithFilters_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	minFee := float32(5.0)
	filter := &event.EventFilter{
		MinFee: &minFee,
	}

	expected := []*event.Event{
		{EventID: "event-1", Name: "Event 1", Fee: 10.0},
	}

	eventRepo.EXPECT().FindAllWithFilters(filter).Return(expected, nil)

	svc := NewEventService(eventRepo)
	result, err := svc.GetEventsWithFilters(filter)

	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestEventService_GetEventByTag_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	expected := []*event.Event{
		{EventID: "event-1", Name: "Event 1", Tags: []string{"music"}},
	}

	eventRepo.EXPECT().FindAllByTags([]string{"music"}).Return(expected, nil)

	svc := NewEventService(eventRepo)
	result, err := svc.GetEventByTag([]string{"music"})

	require.NoError(t, err)
	require.Len(t, result, 1)
}

func TestEventService_RegisterAttendance_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	eventRepo.EXPECT().RegisterAttendance("user-1", "event-1").Return(nil)

	svc := NewEventService(eventRepo)
	err := svc.RegisterAttendance("user-1", "event-1")

	require.NoError(t, err)
}

func TestEventService_RegisterAttendance_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	eventRepo.EXPECT().RegisterAttendance("user-1", "event-1").Return(errors.New("already registered"))

	svc := NewEventService(eventRepo)
	err := svc.RegisterAttendance("user-1", "event-1")

	require.Error(t, err)
}

func TestEventService_GetAttendees_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	expected := []string{"user-1", "user-2"}

	eventRepo.EXPECT().GetAttendees("event-1").Return(expected, nil)

	svc := NewEventService(eventRepo)
	result, err := svc.GetAttendees("event-1")

	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestEventService_RemoveAttendance_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	eventRepo.EXPECT().RemoveAttendance("user-1", "event-1").Return(nil)

	svc := NewEventService(eventRepo)
	err := svc.RemoveAttendance("user-1", "event-1")

	require.NoError(t, err)
}

func TestEventService_RemoveAttendance_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	eventRepo.EXPECT().RemoveAttendance("user-1", "event-1").Return(errors.New("not registered"))

	svc := NewEventService(eventRepo)
	err := svc.RemoveAttendance("user-1", "event-1")

	require.Error(t, err)
}

func TestEventService_GetHostingEvents_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	expected := []*event.Event{
		{EventID: "event-1", Name: "Hosted Event 1", OrganizerID: "user-1"},
		{EventID: "event-2", Name: "Hosted Event 2", OrganizerID: "user-1"},
	}

	eventRepo.EXPECT().FindByOrganizer("user-1", (*event.Pagination)(nil)).Return(expected, nil)

	svc := NewEventService(eventRepo)
	result, err := svc.GetHostingEvents("user-1", nil)

	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, expected, result)
}

func TestEventService_GetHostingEvents_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	eventRepo.EXPECT().FindByOrganizer("user-1", (*event.Pagination)(nil)).Return(nil, errors.New("database error"))

	svc := NewEventService(eventRepo)
	_, err := svc.GetHostingEvents("user-1", nil)

	require.Error(t, err)
}

func TestEventService_GetHostingEvents_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	eventRepo.EXPECT().FindByOrganizer("user-1", (*event.Pagination)(nil)).Return([]*event.Event{}, nil)

	svc := NewEventService(eventRepo)
	result, err := svc.GetHostingEvents("user-1", nil)

	require.NoError(t, err)
	require.Len(t, result, 0)
}

func TestEventService_GetAttendingEvents_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	expected := []*event.Event{
		{EventID: "event-1", Name: "Attending Event 1"},
		{EventID: "event-2", Name: "Attending Event 2"},
	}

	eventRepo.EXPECT().FindAttendingEvents("user-1", (*event.Pagination)(nil)).Return(expected, nil)

	svc := NewEventService(eventRepo)
	result, err := svc.GetAttendingEvents("user-1", nil)

	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, expected, result)
}

func TestEventService_GetAttendingEvents_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	eventRepo.EXPECT().FindAttendingEvents("user-1", (*event.Pagination)(nil)).Return(nil, errors.New("database error"))

	svc := NewEventService(eventRepo)
	_, err := svc.GetAttendingEvents("user-1", nil)

	require.Error(t, err)
}

func TestEventService_GetAttendingEvents_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	eventRepo.EXPECT().FindAttendingEvents("user-1", (*event.Pagination)(nil)).Return([]*event.Event{}, nil)

	svc := NewEventService(eventRepo)
	result, err := svc.GetAttendingEvents("user-1", nil)

	require.NoError(t, err)
	require.Len(t, result, 0)
}

func TestEventService_DeleteEvent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	eventRepo.EXPECT().Delete("event-1").Return(nil)

	svc := NewEventService(eventRepo)
	err := svc.DeleteEvent("event-1")

	require.NoError(t, err)
}

func TestEventService_DeleteEvent_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventRepo := mock_event.NewMockEventRepo(ctrl)

	eventRepo.EXPECT().Delete("event-1").Return(errors.New("database error"))

	svc := NewEventService(eventRepo)
	err := svc.DeleteEvent("event-1")

	require.Error(t, err)
}
