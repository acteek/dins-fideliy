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
type Order struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type MenuResponse struct {
	CheckOrders   string          `json:"check_orders"`
	MealArray     map[string]Meal `json:"meal_array"`
	MenuArray     map[string]Meal `json:"menu_array"`
	Orders        []string        `json:"orders"`
	OrdersContent []string        `json:"orders_content"`
}

type User struct {
	ID    string
	Name  string
	Token string
}

func (r *MenuResponse) GetCurrentMeals() []Meal {
	var meals []Meal
	for _, m := range r.MenuArray {
		meals = append(meals, m)
	}

	if r.CheckOrders == "false" {
		return meals
	} else {
		return []Meal{}
	}

}
