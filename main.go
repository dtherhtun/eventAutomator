package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	cfgPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	client := oauthClient(ctx, cfg.Credentail)

	/*spreadsheetId := "1jheKrlqGnfQ6kx5lfB2hm6aYo4atM8GCrWn5D1jt9u0"
	readRange := "c1up!A2:E"

	calendarId := "c_ul2f5s0g93ib5efh11r23c7ink@group.calendar.google.com"

	url := "http://localhost:8000"*/
	//url := "http://rundeck-staging.upstra-next.ekomedia.technology:8080/api/41/webhook/KixguWigHzred0sCWf5SdVDCSi3pdI21#ParseJSon"

	c := NewClient(cfg.Target)

	rr := fmt.Sprintf("%s!%s", "c1up", cfg.SpreadSheet.ReadRange)

	fmt.Println("Upcoming events:")
	sheetData := fetchData(ctx, client, cfg.SpreadSheet.Id, rr)

	data := spreadSheet2Data(sheetData)

	scaleData, _ := toJson(data, true, "scaleup")
	req, _ := http.NewRequest(http.MethodPost, c.BaseURL, bytes.NewBuffer(scaleData))
	abc := req.WithContext(ctx)

	if err := c.sendRequest(abc); err != nil {
		log.Fatalf("Unable to send request: %v", err)
	}

	chkevent := eventList(ctx, client, cfg.Calendar.Id)

	events, _ := chkevent.Do()

	for _, item := range events.Items {
		if item.Description == "c1" {
			fmt.Println(item.Description)
		}
	}

	/*for {
		events, err := chkevent.Do()
		if err != nil {
			log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
		}
		for _, item := range events.Items {

			start, _ := time.Parse(time.RFC3339, item.Start.DateTime)
			now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

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
	}*/
}
