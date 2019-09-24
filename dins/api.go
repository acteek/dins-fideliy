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
	apiEndpoint  string
	client       *http.Client
	currentMeals map[string]Meal
}

func NewDinsApi(apiEndpoint string) *DinsApi {
	return &DinsApi{
		apiEndpoint:  apiEndpoint,
		client:       &http.Client{},
		currentMeals: currentMeals(apiEndpoint),
	}
}

func currentMeals(apiEndpoint string) map[string]Meal {
	resp, err := http.Get(apiEndpoint + "/cafe-new/tomorrow_get_menu_array.php")
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	data := ParseResponse(body)

	return data.Meals
}

func (d *DinsApi) CurrentMeals() map[string]Meal {
	return d.currentMeals
}

func (d *DinsApi) GetMenu(u User) []Meal {
	resp, err := d.client.Get(d.apiEndpoint + "/cafe-new/tomorrow_get_menu_array.php?user_id=" + u.ID)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	data := ParseResponse(body)

	if !data.isAbleToOrder || len(data.Orders) > 0 {
		return []Meal{}
	} else {
		return data.Menu
	}

}

func (d *DinsApi) GetOrders(u User) []Orders {
	resp, err := d.client.Get(d.apiEndpoint + "/cafe-new/tomorrow_get_menu_array.php?user_id=" + u.ID)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	data := ParseResponse(body)

	return data.Orders
}

func (d *DinsApi) GetUser(token string) (User, error) {
	cookie := http.Cookie{Name: "mydins-auth", Value: token}
	req, err := http.NewRequest(http.MethodGet, d.apiEndpoint+"/?page=fidel", nil)
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

func (d *DinsApi) SendOrder(basket []string, user User) error {
	var orders []Order
	mealStore := d.CurrentMeals()
	for _, id := range basket {
		meal := mealStore[id]
		orders = append(orders, Order{
			ID:      meal.ID,
			Qty:     1,
			Name:    meal.Name,
			Price:   meal.Price,
			Type:    meal.Type,
			Counter: meal.Counter,
		})
	}

	orderJson, ParsErr := json.Marshal(orders)
	if ParsErr != nil {
		log.Fatal("Parse error: ", ParsErr)
	}

	values := url.Values{"user_id": {user.ID}, "full_name": {user.Name}, "order": {string(orderJson)}, "make_the_order": {"Заказать"}}
	req, _ := http.NewRequest(http.MethodPost, d.apiEndpoint+"/cafe-new/user_order.php", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err := d.client.Do(req)
	log.Println(user.Name + " make order " + string(orderJson))

	return err

}

func (d *DinsApi) CancelOrder(orderID string, user User) error {
	values := url.Values{"user_id": {user.ID}, "full_name": {user.Name}, "cancel_the_order": {"Отменить"}, "order_id": {orderID}}
	req, _ := http.NewRequest(http.MethodPost, d.apiEndpoint+"/cafe-new/user_order.php", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err := d.client.Do(req)
	log.Println(user.Name + " cancel order " + orderID)

	return err
}
