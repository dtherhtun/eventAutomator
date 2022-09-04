package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func fetchData(ctx context.Context, client *http.Client, spreadsheetId, readRange string) [][]interface{} {
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Unable to retrieve Sheets client: %v", err)
	}

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Printf("Unable to retrieve data from sheet: %v", err)
	}

	return resp.Values
}

func (cfg *Config) getScaleData(ctx context.Context, client *http.Client, nats bool, commit, cluster string) []byte {
	var allSheetData [][]interface{}
	var sheets []string
	clusters := strings.Split(cluster, ",")

	for _, sheet := range clusters {
		rr := fmt.Sprintf("%s!%s", sheet, cfg.SpreadSheet.ReadRange)
		sheets = append(sheets, rr)
	}

	for _, rr := range sheets {
		sheetData := fetchData(ctx, client, cfg.SpreadSheet.Id, rr)
		allSheetData = append(allSheetData, sheetData...)
	}

	data := spreadSheet2Data(allSheetData)

	scaleData, err := toJson(data, nats, commit)
	if err != nil {
		log.Printf("Unable to get Json data: %v", err)
	}

	return scaleData
}
