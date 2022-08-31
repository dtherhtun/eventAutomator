package main

import (
	"context"
	"fmt"
	"log"
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

	c := NewClient(cfg.Target.URL, cfg.Target.ApiKey)

	scaleData := cfg.getScaleData(ctx, client, false, "scale up", "c1up")

	req := postRequest(c.BaseURL, scaleData)

	if err := c.sendRequest(req); err != nil {
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
