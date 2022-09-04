package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	apiKey     string
	HTTPClient *http.Client
}

type errorResponse struct {
	Err string `json:"err"`
}

type RundeckResponse struct {
	JobId       string `json:"jobId"`
	ExecutionId string `json:"executionId"`
}

func NewClient(BaseURLV1, apiKey string) *Client {
	return &Client{
		BaseURL: BaseURLV1,
		apiKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func (c *Client) sendRequest(req *http.Request) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", c.apiKey)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Err)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	body, _ := ioutil.ReadAll(res.Body)
	var response RundeckResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return err
	}
	log.Println("JobID -> ", response.JobId, " executionId -> ", response.ExecutionId)

	return nil
}

func (c *Config) postRequest(data []byte) *http.Request {
	reqBody := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, c.Target.URL, reqBody)
	if err != nil {
		log.Printf("Unable to make request: %v", err)
	}

	return req
}
