package dins

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
	CheckOrders   string          `json:"check_orders"`
	MealArray     map[string]Meal `json:"meal_array"`
	MenuArray     map[string]Meal `json:"menu_array"`
	OrdersContent []OrderContent  `json:"orders_content"`
	//Orders        []string        `json:"orders"` //TODO
}

type User struct {
	ID    string
	Name  string
	Token string
}

type Order struct {
	ID      string      `json:"id"`
	Qty     int         `json:"qty"` // value 1 order, value 0 not order
	Name    string      `json:"name"`
	Price   interface{} `json:"price"`
	Type    string      `json:"type"`
	Counter interface{} `json:"counter"`
}

func (r *MenuResponse) GetCurrentMeals() []Meal {
	var meals []Meal
	for _, m := range r.MenuArray {
		meals = append(meals, m)
	}

	if r.CheckOrders == "true" {
		return meals
	} else {
		return []Meal{}
	}

}
