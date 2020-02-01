package dins

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func parseMenu(menu gjson.Result) []Meal {
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

func parseMeals(meals gjson.Result) map[string]Meal {
	var parsed = make(map[string]Meal)
	if parseErr := json.Unmarshal([]byte(meals.Raw), &parsed); parseErr != nil {
		log.Panic("Parse Meals error: ", parseErr)
	}
	return parsed
}

func parseOrders(orders gjson.Result) []Order {
	var data []Order
	if parseErr := json.Unmarshal([]byte(orders.Raw), &data); parseErr != nil {
		log.Fatal("Parse Order error: ", parseErr)
	}
	return data
}

func parseAbleToOrder(check gjson.Result) bool {
	value, parseErr := strconv.ParseBool(check.Str)
	if parseErr != nil {
		log.Fatal("Parse AbleToOrder error: ", parseErr)
	}
	return value
}

// ParseResponse parse response fro Dins Api
func ParseResponse(resp *http.Response) MenuResponse {
	bytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	p := gjson.ParseBytes(bytes)
	return MenuResponse{
		isAbleToOrder: parseAbleToOrder(p.Get("check_orders")),
		Meals:         parseMeals(p.Get("meal_array")),
		Menu:          parseMenu(p.Get("menu_array")),
		Orders:        parseOrders(p.Get("orders_content")),
	}

}
