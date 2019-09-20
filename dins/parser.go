package dins

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"log"
	"strconv"
)

func ParseMenu(menu gjson.Result) []Meal {
	var parsed = make(map[string]Meal)
	var data []Meal
	if menu.IsObject() {
		if parseErr := json.Unmarshal([]byte(menu.Raw), &parsed); parseErr != nil {
			log.Panic("Parse Menu error: ", parseErr)
		}
		for _, m := range parsed {
			data = append(data, m)
		}
	}

	return data
}

func ParseMeals(meals gjson.Result) map[string]Meal {
	var parsed = make(map[string]Meal)
	if parseErr := json.Unmarshal([]byte(meals.Raw), &parsed); parseErr != nil {
		log.Panic("Parse Meals error: ", parseErr)
	}
	return parsed
}

func ParseOrders(orders gjson.Result) []OrderContent {
	var data []OrderContent
	if parseErr := json.Unmarshal([]byte(orders.Raw), &data); parseErr != nil {
		log.Fatal("Parse Orders error: ", parseErr)
	}
	return data
}

func ParseAbleToOrder(check gjson.Result) bool {
	value, parseErr := strconv.ParseBool(check.Str)
	if parseErr != nil {
		log.Fatal("Parse AbleToOrder error: ", parseErr)
	}
	return value
}

func ParseResponse(bytes []byte) MenuResponse {
	p := gjson.ParseBytes(bytes)
	return MenuResponse{
		isAbleToOrder: ParseAbleToOrder(p.Get("check_orders")),
		Meals:         ParseMeals(p.Get("meal_array")),
		Menu:          ParseMenu(p.Get("menu_array")),
		Orders:        ParseOrders(p.Get("orders_content")),
	}
}
