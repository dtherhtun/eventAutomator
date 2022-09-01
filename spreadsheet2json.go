package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type data struct {
	App string `json:"app"`
	Cpu string `json:"cpu"`
	Mem string `json:"mem"`
	Min int    `json:"min"`
	Max int    `json:"max"`
}

type scale struct {
	Data   []data `json:"data"`
	Nats   bool   `json:"nats"`
	Commit string `json:"commit"`
}

func spreadSheet2Data(sheetdata [][]interface{}) []data {
	var allData []data

	for i := 0; i < len(sheetdata); i++ {
		temp := data{
			App: fmt.Sprintf("%v", sheetdata[i][0]),
			Cpu: fmt.Sprintf("%v", sheetdata[i][1]),
			Mem: fmt.Sprintf("%v", sheetdata[i][2]),
			Min: string2Int(fmt.Sprintf("%v", sheetdata[i][3])),
			Max: string2Int(fmt.Sprintf("%v", sheetdata[i][4])),
		}
		allData = append(allData, temp)
	}
	return allData
}

func string2Int(s string) int {
	number, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("Unable to parse string to int: %v", err)
	}
	return number
}

func toJson(data []data, nats bool, commit string) ([]byte, error) {
	scaledata := scale{
		Data:   data,
		Nats:   nats,
		Commit: commit,
	}
	u, err := json.Marshal(scaledata)
	if err != nil {
		return nil, err
	}
	return u, nil
}
