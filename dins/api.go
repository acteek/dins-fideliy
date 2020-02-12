package dins

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

/*
API struct for communicate with Dins

Endpoint     - Dins server endpoint
client       - Common client for each request
CurrentMeals - Cached list of meals
*/
type API struct {
	Endpoint     string
	client       *http.Client
	CurrentMeals map[string]Meal
}

// NewDinsAPI constructor for create new API instance
func NewDinsAPI(apiEndpoint string) *API {
	mealStore, err := currentMeals(apiEndpoint)
	if err != nil {
		log.Fatal("Failed connect to dins ", apiEndpoint, err)
	}

	return &API{
		Endpoint:     apiEndpoint,
		client:       &http.Client{},
		CurrentMeals: mealStore,
	}
}

func currentMeals(apiEndpoint string) (map[string]Meal, error) {
	resp, err := http.Get(apiEndpoint + "/cafe-new/tomorrow_get_menu_array.php")
	if err != nil {
		return map[string]Meal{}, err
	}
	data := ParseResponse(resp)
	return data.Meals, nil
}

//GetMenu returns list of meals for user and flag does user has orders or doesn't
func (d *API) GetMenu(u User) ([]Meal, bool) {
	resp, err := d.client.Get(d.Endpoint + "/cafe-new/tomorrow_get_menu_array.php?user_id=" + u.ID)
	if err != nil {
		log.Println(err)
	}

	data := ParseResponse(resp)
	if !data.isAbleToOrder {
		return []Meal{}, false
	}

	return data.Menu, len(data.Orders) > 0

}

//GetSubList returns list of meals avalible for subscription 
func (d *API) GetSubList() []Meal {
	resp, err := d.client.Get(d.Endpoint + "/cafe-new/tomorrow_get_menu_array.php")
	if err != nil {
		log.Println(err)
	}

	data := ParseResponse(resp)

	return data.Menu

}

//GetOrders returns list of orders for user
func (d *API) GetOrders(u User) []Order {
	resp, err := d.client.Get(d.Endpoint + "/cafe-new/tomorrow_get_menu_array.php?user_id=" + u.ID)
	if err != nil {
		log.Println(err)
	}

	data := ParseResponse(resp)
	return data.Orders
}

// GetUser it's auth method, return user by token -TODO refactor
func (d *API) GetUser(token string) (User, error) {
	cookie := http.Cookie{Name: "mydins-auth", Value: token}
	req, err := http.NewRequest(http.MethodGet, d.Endpoint+"/?page=fidel", nil)
	if err != nil {
		return User{}, err
	}
	req.AddCookie(&cookie)
	resp, err := d.client.Do(req)

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	expName, _ := regexp.Compile(`full_name \s*=\s*"([\S\s]+)" ;`)
	expID, _ := regexp.Compile(`user_id\s*=\s*([0-9]+);`)

	str := string(body)
	parseName := expName.FindStringSubmatch(str)
	parseID := expID.FindStringSubmatch(str)

	if len(parseID) < 2 || len(parseName) < 2 {
		err := errors.New("failed parse user_id or user_name")
		log.Println(err)
		return User{}, err
	}

	return User{ID: parseID[1], Name: parseName[1]}, nil

}

//SendOrder make order to Dins
func (d *API) SendOrder(order []Meal, user User) error {
	orderJSON, _ := json.Marshal(order)
	values := url.Values{"user_id": {user.ID}, "full_name": {user.Name}, "order": {string(orderJSON)}, "make_the_order": {"Заказать"}}
	req, _ := http.NewRequest(http.MethodPost, d.Endpoint+"/cafe-new/user_order.php", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err := d.client.Do(req)
	log.Println(user.Name + " make order " + string(orderJSON))

	return err
}

//CancelOrder cancel Order to Dins
func (d *API) CancelOrder(orderID string, user User) error {
	values := url.Values{"user_id": {user.ID}, "full_name": {user.Name}, "cancel_the_order": {"Отменить"}, "order_id": {orderID}, "order": {""}}
	req, _ := http.NewRequest(http.MethodPost, d.Endpoint+"/cafe-new/user_order.php", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err := d.client.Do(req)
	log.Println(user.Name + " cancel order " + orderID)

	return err
}
