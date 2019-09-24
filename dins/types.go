package dins

import (
	"encoding/json"
	"log"
)

type Meal struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Name    string      `json:"name"`
	Price   interface{} `json:"price"`
	Counter interface{} `json:"counter"`
	Check   string      `json:"check"`
	Qty     int         `json:"qty"` // number orders
}

type Order struct {
	ID     string `json:"order_id"`
	MealID string `json:"meal_id"`
	Qty    string `json:"qty"`
}
type MenuResponse struct {
	isAbleToOrder bool
	Meals         map[string]Meal
	Menu          []Meal
	Orders        []Order
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (u *User) GetBytes() []byte {
	bytes, ParsErr := json.Marshal(u)
	if ParsErr != nil {
		log.Fatal("Parse error: ", ParsErr)
	}
	return bytes
}
