package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type data struct {
	App interface{} `json:"app"`
	Cpu interface{} `json:"cpu"`
	Mem interface{} `json:"mem"`
}

type scale struct {
	Data   []data `json:"data"`
	Commit string `json:"commit"`
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	shsrv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	spreadsheetId := "1jheKrlqGnfQ6kx5lfB2hm6aYo4atM8GCrWn5D1jt9u0"
	readRange := "c1!A2:C"

	resp, err := shsrv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	name := []string{"app", "cpu", "mem"}
	var datas []data
	m := make(map[string]interface{})
	var datamap []map[string]interface{}

	for i := 0; i < len(resp.Values); i++ {
		for j := 0; j < len(resp.Values[i]); j++ {
			m[name[j]] = resp.Values[i][j]
		}
		datamap = append(datamap, m)
	}
	for _, v := range datamap {
		temp := data{
			App: v["app"],
			Cpu: v["cpu"],
			Mem: v["mem"],
		}
		datas = append(datas, temp)
	}
	aeiou := scale{
		Data:   datas,
		Commit: "scale up",
	}
	u, _ := json.Marshal(aeiou)
	fmt.Println(string(u))

	url := "http://localhost:8080"

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(u))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	c := &http.Client{}
	response, err := c.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	t := time.Now().Format(time.RFC3339)
	eventsListCall := srv.Events.List("c_ul2f5s0g93ib5efh11r23c7ink@group.calendar.google.com").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime")

	fmt.Println("Upcoming events:")

	for {
		events, err := eventsListCall.Do()
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
				fmt.Println(item.Summary, " Event is Starting...........")
			}
		}

	}
}
