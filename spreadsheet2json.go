package main

import "encoding/json"

type data struct {
	App interface{} `json:"app"`
	Cpu interface{} `json:"cpu"`
	Mem interface{} `json:"mem"`
}

type scale struct {
	Data   []data `json:"data"`
	Commit string `json:"commit"`
}

func spreadsheet2json(sheetdata [][]interface{}) []byte {
	title := []string{"app", "cpu", "mem"}
	var allData []data
	m := make(map[string]interface{})
	var datamap []map[string]interface{}

	for i := 0; i < len(sheetdata); i++ {
		for j := 0; j < len(sheetdata[i]); j++ {
			m[title[j]] = sheetdata[i][j]
		}
		datamap = append(datamap, m)
	}

	for _, v := range datamap {
		temp := data{
			App: v["app"],
			Cpu: v["cpu"],
			Mem: v["mem"],
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
