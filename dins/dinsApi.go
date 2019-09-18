package dins

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type DinsApi struct {
	apiEndpoint string
	client      *http.Client
}

//TODO all strict
type Meal struct {
	ID            string     `json:"id"`
	Type          string      `json:"type"`
	Name          string     `json:"name"`
	Price         interface{} `json:"price"`
	Counter       interface{} `json:"counter"`
	Check         string      `json:"check"`
	CheckTomorrow string      `json:"check_tomorrow"`
}

type MenuResponse struct {
	CheckOrders   string          `json:"check_orders"`
	MealArray     map[string]Meal `json:"meal_array"`
	MenuArray     map[string]Meal `json:"menu_array"`
	Orders        []string        `json:"orders"`
	OrdersContent []string        `json:"orders_content"`
}

func (r *MenuResponse) GetMenuMeals() []Meal {
	var meals []Meal
	for _, m := range r.MenuArray {
		meals = append(meals, m)
	}

	return meals
}

func NewDinsApi(apiEndpoint string) *DinsApi {
	return &DinsApi{
		apiEndpoint: apiEndpoint,
		client:      &http.Client{},
	}
}

//TODO use user_id (my id ?user_id=1092)
func (d *DinsApi) GetMenu() []Meal {
	resp, err := d.client.Get(d.apiEndpoint + "/cafe-new/tomorrow_get_menu_array.php")
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	data := MenuResponse{}

	body, _ := ioutil.ReadAll(resp.Body)

	if parseErr := json.Unmarshal(body, &data); parseErr != nil {
		log.Fatal("Parse error: ", parseErr)
	}

	return data.GetMenuMeals()

}
