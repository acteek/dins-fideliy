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

func NewDinsApi(apiEndpoint string) *DinsApi {
	return &DinsApi{
		apiEndpoint: apiEndpoint,
		client:      &http.Client{},
	}
}

//TODO use user_id (my id ?user_id=1092)
func (d *DinsApi) GetMenu() interface{} {
	resp, err := d.client.Get(d.apiEndpoint + "cafe-new/tomorrow_get_menu_array.php")
	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()
	var dat map[string]interface{}

	body, _ := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &dat); err != nil {
		panic(err)
	}
	return dat["menu_array"]

}
