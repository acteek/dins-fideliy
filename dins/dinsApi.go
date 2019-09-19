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

	data := MenuResponse{}

	body, _ := ioutil.ReadAll(resp.Body)

	if parseErr := json.Unmarshal(body, &data); parseErr != nil {
		log.Fatal("Parse error: ", parseErr)
	}

	return data.MealArray
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

	data := MenuResponse{}

	body, _ := ioutil.ReadAll(resp.Body)

	if parseErr := json.Unmarshal(body, &data); parseErr != nil {
		log.Fatal("Parse error: ", parseErr)
	}

	return data.GetCurrentMeals()
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
		return User{ID: parseId[1], Name: parseName[1], Token: token}, nil
	}

}

//TODO implement send request
func (d *DinsApi) SendOrder(order Order, user User) error {
	cookie := http.Cookie{Name: "mydins-auth", Value: user.Token}
	values := url.Values{"user_id": {user.ID}, "full_name": {user.Name}, "orders": {"TODO"}}

	req, _ := http.NewRequest(http.MethodGet, d.apiEndpoint+"/TODO/", strings.NewReader(values.Encode()))
	req.AddCookie(&cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err := d.client.Do(req)

	return err
}
