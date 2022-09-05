package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func eventList(ctx context.Context, logger *log.Logger, client *http.Client, calendarId string) *calendar.EventsListCall {
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		logger.Printf("Unable to retrieve Calendar client: %v", err)
	}

	t := time.Now().Format(time.RFC3339)
	eventsListCall := srv.Events.List(calendarId).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime")

	return eventsListCall
}
