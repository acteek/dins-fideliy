package main

import (
	"encoding/json"
	"fideliy/dins"
	"fmt"
	"log"
)

func main() {

	v := []byte(`{
		"10": {
			"id": "10",
				"type": "Гарнир",
				"name": "Картофельное пюре",
				"price": null,
				"counter": 0,
				"check": "0",
				"check_tomorrow": "1"
		},
		"71": {
			"id": "71",
				"type": "Гарнир",
				"name": "Рис с кукурузой",
				"price": null,
				"counter": null,
				"check": "0",
				"check_tomorrow": "1"
		}
     }`)
	var data map[string]dins.Meal

	if parseErr := json.Unmarshal(v, &data); parseErr != nil {
		log.Fatal("Parse error")
	}

	var dishes []dins.Meal

	for _, dish := range data {
		dishes =  append(dishes, dish)
	}

	fmt.Println(dishes)

}
