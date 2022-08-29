package main

import (
	"encoding/json"
	"fmt"
)

type data struct {
	App string `json:"app"`
	Cpu string `json:"cpu"`
	Mem string `json:"mem"`
}

type scale struct {
	Data   []data `json:"data"`
	Commit string `json:"commit"`
}

func spreadsheet2json(sheetdata [][]interface{}) []byte {
	var allData []data

	for i := 0; i < len(sheetdata); i++ {
		temp := data{
			App: fmt.Sprintf("%v", sheetdata[i][0]),
			Cpu: fmt.Sprintf("%v", sheetdata[i][1]),
			Mem: fmt.Sprintf("%v", sheetdata[i][2]),
		}
		allData = append(allData, temp)
	}

	aeiou := scale{
		Data:   allData,
		Commit: "scale up",
	}
	u, _ := json.Marshal(aeiou)
	return u
}
