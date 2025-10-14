package google

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// CalendarService provides methods for interacting with the Google Calendar API.
type CalendarService struct {
	service *calendar.Service
}

// NewCalendarService creates a new CalendarService.
// The token source is used to authenticate with the Google Calendar API.
func NewCalendarService(ctx context.Context, ts oauth2.TokenSource) (*CalendarService, error) {
	client := oauth2.NewClient(ctx, ts)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Calendar client: %v", err)
	}
	return &CalendarService{service: srv}, nil
}

// CreateEvent creates a new event on the specified calendar.
// calendarID is the ID of the calendar, typically "primary".
func (s *CalendarService) CreateEvent(calendarID string, event *calendar.Event) (*calendar.Event, error) {
	return s.service.Events.Insert(calendarID, event).Do()
}

// GetEvent retrieves an event from the specified calendar.
func (s *CalendarService) GetEvent(calendarID, eventID string) (*calendar.Event, error) {
	return s.service.Events.Get(calendarID, eventID).Do()
}

// UpdateEvent updates an existing event.
func (s *CalendarService) UpdateEvent(calendarID, eventID string, event *calendar.Event) (*calendar.Event, error) {
	return s.service.Events.Update(calendarID, eventID, event).Do()
}

// DeleteEvent deletes an event.
func (s *CalendarService) DeleteEvent(calendarID, eventID string) error {
	return s.service.Events.Delete(calendarID, eventID).Do()
}

// ListUpcomingEvents lists the upcoming events on a calendar.
// maxResults is the maximum number of events to return.
func (s *CalendarService) ListUpcomingEvents(calendarID string, maxResults int64) (*calendar.Events, error) {
	t := time.Now().Format(time.RFC3339)
	return s.service.Events.List(calendarID).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(t).
		MaxResults(maxResults).
		OrderBy("startTime").
		Do()
}

// SearchEvents searches for events on a calendar.
// The query is a free-text search query.
func (s *CalendarService) SearchEvents(calendarID, query string) (*calendar.Events, error) {
	return s.service.Events.List(calendarID).Q(query).Do()
}
