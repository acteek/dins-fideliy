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

type DinsApi struct {
	Endpoint     string
	client       *http.Client
	CurrentMeals map[string]Meal
}

func NewDinsApi(apiEndpoint string) *DinsApi {
	mealStore, err := currentMeals(apiEndpoint)
	if err != nil {
		log.Fatal("Failed connect to dins ", apiEndpoint, err)
	}

	return &DinsApi{
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

func (d *DinsApi) GetMenu(u User) []Meal {
	resp, err := d.client.Get(d.Endpoint + "/cafe-new/tomorrow_get_menu_array.php?user_id=" + u.ID)
	if err != nil {
		log.Panic(err)
	}

	data := ParseResponse(resp)
	if !data.isAbleToOrder || len(data.Orders) > 0 {
		return []Meal{}
	} else {
		return data.Menu
	}

}

func (d *DinsApi) GetOrders(u User) []Order {
	resp, err := d.client.Get(d.Endpoint + "/cafe-new/tomorrow_get_menu_array.php?user_id=" + u.ID)
	if err != nil {
		log.Panic(err)
	}

	data := ParseResponse(resp)
	return data.Orders
}

//TODO refactor auth method
func (d *DinsApi) GetUser(token string) (User, error) {
	cookie := http.Cookie{Name: "mydins-auth", Value: token}
	req, err := http.NewRequest(http.MethodGet, d.Endpoint+"/?page=fidel", nil)
	if err != nil {
		return User{}, err
	}
	req.AddCookie(&cookie)
	resp, err := d.client.Do(req)

	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	expName, _ := regexp.Compile(`full_name \s*=\s*"([\S\s]+)" ;`)
	expID, _ := regexp.Compile(`user_id\s*=\s*([0-9]+);`)

	str := string(body)
	parseName := expName.FindStringSubmatch(str)
	parseId := expID.FindStringSubmatch(str)

	if len(parseId) < 2 || len(parseName) < 2 {
		err := errors.New("failed parse user_id or user_name")
		log.Println(err)
		return User{}, err
	} else {
		return User{ID: parseId[1], Name: parseName[1]}, nil
	}

}

func (d *DinsApi) SendOrder(order []Meal, user User) error {
	orderJson, _ := json.Marshal(order)
	values := url.Values{"user_id": {user.ID}, "full_name": {user.Name}, "order": {string(orderJson)}, "make_the_order": {"Заказать"}}
	req, _ := http.NewRequest(http.MethodPost, d.Endpoint+"/cafe-new/user_order.php", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err := d.client.Do(req)
	log.Println(user.Name + " make order " + string(orderJson))

	return err
}

func (d *DinsApi) CancelOrder(orderID string, user User) error {
	values := url.Values{"user_id": {user.ID}, "full_name": {user.Name}, "cancel_the_order": {"Отменить"}, "order_id": {orderID}, "order": {""}}
	req, _ := http.NewRequest(http.MethodPost, d.Endpoint+"/cafe-new/user_order.php", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err := d.client.Do(req)
	log.Println(user.Name + " cancel order " + orderID)

	return err
}
