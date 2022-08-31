package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func fetchData(ctx context.Context, client *http.Client, spreadsheetId, readRange string) [][]interface{} {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	fmt.Println(resp.Values)

	return resp.Values
}

func (cfg *Config) getScaleData(ctx context.Context, client *http.Client, nats bool, commit, cluster string) []byte {
	rr := fmt.Sprintf("%s!%s", cluster, cfg.SpreadSheet.ReadRange)

	sheetData := fetchData(ctx, client, cfg.SpreadSheet.Id, rr)

	data := spreadSheet2Data(sheetData)

	scaleData, err := toJson(data, nats, commit)
	if err != nil {
		log.Fatalf("Unable to get Json data: %v", err)
	}

	return scaleData
}
