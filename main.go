package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

func main() {
	var nats bool = false
	p := bluemonday.StrictPolicy()
	p.AddSpaceWhenStrippingTag(false)

	ctx := context.Background()

	cfgPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	clientGoogle := oauthClient(ctx, cfg.Credentail)

	clientTarget := NewClient(cfg.Target.URL, cfg.Target.ApiKey)

	chkevent := eventList(ctx, clientGoogle, cfg.Calendar.Id)

	for {
		events, err := chkevent.Do()
		if err != nil {
			log.Printf("Unable to retrieve the events: %v", err)
		}
		for _, item := range events.Items {

			start, _ := time.Parse(time.RFC3339, item.Start.DateTime)
			now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

			if item.Summary == "[ASC] TESTING" && start.Equal(now) {
				time.Sleep(5 * time.Second)
				log.Println(item.Summary, " Event is Starting...........")
				description := p.Sanitize(item.Description)
				description = strings.TrimSpace(description)
				commit := fmt.Sprintf("%s-%s", item.Summary, description)
				p.Sanitize(item.Description)

				switch description {
				case "c1Up", "c1Down", "c2Up", "c2Down":
					if description == "c2up" {
						nats = true
					}
					log.Println("nats->", nats)
					scaleData := cfg.getScaleData(ctx, clientGoogle, nats, commit, description)
					req := postRequest(clientTarget.BaseURL, scaleData)

					if err := clientTarget.sendRequest(req); err != nil {
						log.Printf("Unable to send request: %v", err)
					}
				case "multiClusterUp":
					scaleData := cfg.getScaleData(ctx, clientGoogle, nats, commit, "c1Up")
					req := postRequest(clientTarget.BaseURL, scaleData)

					if err := clientTarget.sendRequest(req); err != nil {
						log.Printf("Unable to send request: %v", err)
					}
					time.Sleep(10 * time.Second)
					scaleData2 := cfg.getScaleData(ctx, clientGoogle, true, commit, "c2Up")
					req2 := postRequest(clientTarget.BaseURL, scaleData2)

					if err := clientTarget.sendRequest(req2); err != nil {
						log.Printf("Unable to send request: %v", err)
					}
				case "multiClusterDown", "multiClusterDownWithoutNats":
					if description == "multiClusterDownWithoutNats" {
						nats = true
					}
					scaleData := cfg.getScaleData(ctx, clientGoogle, nats, commit, "c1Down")
					req := postRequest(clientTarget.BaseURL, scaleData)

					if err := clientTarget.sendRequest(req); err != nil {
						log.Printf("Unable to send request: %v", err)
					}
					time.Sleep(5 * time.Second)
					scaleData2 := cfg.getScaleData(ctx, clientGoogle, nats, commit, "c2Down")
					req2 := postRequest(clientTarget.BaseURL, scaleData2)

					if err := clientTarget.sendRequest(req2); err != nil {
						log.Printf("Unable to send request: %v", err)
					}
				default:
					log.Printf("Wrong Event instruction: %s", description)
				}
			}
		}
	}
}
