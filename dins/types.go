package dins

import (
	"encoding/json"
	"log"
	t "time"
)

//Meal describes a meal from dinsAPI
type Meal struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Name    string      `json:"name"`
	Price   interface{} `json:"price"`
	Counter interface{} `json:"counter"`
	Check   string      `json:"check"`
	Qty     int         `json:"qty"` // number orders
}

//Order describes a order from dinsAPI
type Order struct {
	ID     string `json:"order_id"`
	MealID string `json:"meal_id"`
	Qty    string `json:"qty"`
}

//MenuResponse common struct for parsse response
type MenuResponse struct {
	isAbleToOrder bool
	Meals         map[string]Meal
	Menu          []Meal
	Orders        []Order
}

/*
User has user data field.
Save to store by telegram chatId

ID   - Dins Domain ID
Name - Dins Domain Name
Subs - Active subscriptions for user
*/
type User struct {
	ID   string            `json:"id"`
	Name string            `json:"name"`
	Subs map[string]t.Time `json:"subs"`
}

//GetBytes wraper for Serialization User struct
func (u *User) GetBytes() []byte {
	bytes, ParsErr := json.Marshal(u)
	if ParsErr != nil {
		log.Fatal("Parse error: ", ParsErr)
	}
	return bytes
}
