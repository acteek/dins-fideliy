package dins

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type DinsApi struct {
	apiEndpoint string
	client      *http.Client
}

func NewDinsApi(apiEndpoint string) *DinsApi {
	return &DinsApi{
		apiEndpoint: apiEndpoint,
		client:      &http.Client{},
	}
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
