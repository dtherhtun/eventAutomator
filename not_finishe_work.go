// https://go.dev/play/p/JOISWO9gIki 

// You can edit this code!
// Click here and start typing.
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

func main() {
	var datas []data
	test := [][]string{{"asc-app-c1", "500m", "2Gi"}, {"asc-api-c1", "1", "2Gi"}}
	name := []string{"app", "cpu", "mem"}
	m := make(map[string]string)
	var datamap []map[string]string
	for i := 0; i < len(test); i++ {
		for j := 0; j < len(test[i]); j++ {
			m[name[j]] = test[i][j]
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
	fmt.Println(aeiou)
	u, _ := json.Marshal(aeiou)
	fmt.Println(string(u))
}

