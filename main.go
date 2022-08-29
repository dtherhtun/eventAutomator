package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	filename := "credentials.json"

	client := oauthClient(ctx, filename)

	spreadsheetId := "1jheKrlqGnfQ6kx5lfB2hm6aYo4atM8GCrWn5D1jt9u0"
	readRange := "c1!A2:C"

	calendarId := "c_ul2f5s0g93ib5efh11r23c7ink@group.calendar.google.com"

	//url := "http://localhost:8000"
	url := "http://rundeck-staging.upstra-next.ekomedia.technology:8080/api/41/webhook/KixguWigHzred0sCWf5SdVDCSi3pdI21#ParseJSon"

	c := NewClient(url)

	fmt.Println("Upcoming events:")
	chkevent := eventList(ctx, client, calendarId)

	for {
		events, err := chkevent.Do()
		if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
		}
		for _, item := range events.Items {
			//fmt.Println(item.Summary)
			start, _ := time.Parse(time.RFC3339, item.Start.DateTime)
			now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			//fmt.Println(start)
			//fmt.Println(now)
			if item.Summary == "[ASC] TESTING" && start.Equal(now) {
				time.Sleep(5 * time.Second)
				fmt.Println(item.Summary, " Event is Starting...........")
				sheetData := fetchData(ctx, client, spreadsheetId, readRange)
				req, err := http.NewRequest(http.MethodPost, c.BaseURL, bytes.NewBuffer(spreadsheet2json(sheetData)))

				if err != nil {
					fmt.Println("newreq err", err)
				}
				abc := req.WithContext(ctx)

				if err := c.sendRequest(abc); err != nil {
					fmt.Println("sendreq err", err)
				}
			}
		}
	}
}
