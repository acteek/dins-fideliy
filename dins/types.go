package dins

import (
	"encoding/json"
	"log"
)

type Meal struct {
	ID            string      `json:"id"`
	Type          string      `json:"type"`
	Name          string      `json:"name"`
	Price         interface{} `json:"price"`
	Counter       interface{} `json:"counter"`
	Check         string      `json:"check"`
	CheckTomorrow string      `json:"check_tomorrow"`
}

type OrderContent struct {
	ID      string `json:"id"`
	OrderID string `json:"order_id"`
	MealID  string `json:"meal_id"`
	Qty     string `json:"qty"`
}
type MenuResponse struct {
	isAbleToOrder bool
	Meals         map[string]Meal
	Menu          []Meal
	Orders        []OrderContent
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Order struct {
	ID      string      `json:"id"`
	Qty     int         `json:"qty"` // value 1 order, value 0 not order
	Name    string      `json:"name"`
	Price   interface{} `json:"price"`
	Type    string      `json:"type"`
	Counter interface{} `json:"counter"`
}

func (u *User) GetBytes() []byte {
	bytes, ParsErr := json.Marshal(u)
	if ParsErr != nil {
		log.Fatal("Parse error: ", ParsErr)
	}
	return bytes
}
