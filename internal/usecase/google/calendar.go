package usecase

import (
	"fmt"

	gcalendar "google.golang.org/api/calendar/v3"
)

// CalendarIntegration defines the interface for calendar operations.
// This allows for decoupling the use case from a specific calendar provider.
type CalendarIntegration interface {
	CreateEvent(calendarID string, event *gcalendar.Event) (*gcalendar.Event, error)
	GetEvent(calendarID, eventID string) (*gcalendar.Event, error)
	UpdateEvent(calendarID, eventID string, event *gcalendar.Event) (*gcalendar.Event, error)
	DeleteEvent(calendarID, eventID string) error
	ListUpcomingEvents(calendarID string, maxResults int64) (*gcalendar.Events, error)
	SearchEvents(calendarID, query string) (*gcalendar.Events, error)
}

// CalendarUseCase handles the business logic for calendar operations.
type CalendarUseCase struct {
	calendarIntegration CalendarIntegration
}

// NewCalendarUseCase creates a new CalendarUseCase.
func NewCalendarUseCase(ci CalendarIntegration) *CalendarUseCase {
	return &CalendarUseCase{calendarIntegration: ci}
}

// CreateEvent creates a new event using the calendar integration.
// It checks for existing events with the same summary and start time to avoid duplicates.
func (uc *CalendarUseCase) CreateEvent(calendarID string, event *gcalendar.Event) (*gcalendar.Event, error) {
	if event.Summary != "" && event.Start != nil && event.Start.DateTime != "" {
		query := fmt.Sprintf("summary='%s'", event.Summary)
		existingEvents, err := uc.calendarIntegration.SearchEvents(calendarID, query)
		if err != nil {
			return nil, fmt.Errorf("failed to search for existing events: %w", err)
		}

		for _, item := range existingEvents.Items {
			if item.Start != nil && item.Start.DateTime == event.Start.DateTime {
				// Event with same summary and start time already exists.
				return item, nil
			}
		}
	}

	return uc.calendarIntegration.CreateEvent(calendarID, event)
}

// GetEvent retrieves a specific event.
func (uc *CalendarUseCase) GetEvent(calendarID, eventID string) (*gcalendar.Event, error) {
	return uc.calendarIntegration.GetEvent(calendarID, eventID)
}

// UpdateEvent updates an existing event.
func (uc *CalendarUseCase) UpdateEvent(calendarID, eventID string, event *gcalendar.Event) (*gcalendar.Event, error) {
	return uc.calendarIntegration.UpdateEvent(calendarID, eventID, event)
}

// DeleteEvent deletes an event.
func (uc *CalendarUseCase) DeleteEvent(calendarID, eventID string) error {
	return uc.calendarIntegration.DeleteEvent(calendarID, eventID)
}

// ListUpcomingEvents lists upcoming events.
func (uc *CalendarUseCase) ListUpcomingEvents(calendarID string, maxResults int64) (*gcalendar.Events, error) {
	return uc.calendarIntegration.ListUpcomingEvents(calendarID, maxResults)
}

// SearchEvents searches for events.
func (uc *CalendarUseCase) SearchEvents(calendarID, query string) (*gcalendar.Events, error) {
	return uc.calendarIntegration.SearchEvents(calendarID, query)
}
