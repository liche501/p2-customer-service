package extends

import (
	"best/p2-customer-service/logs"
	"encoding/json"
	"io/ioutil"
)

var (
	ErrorList = map[string]string{}
)

func InitErrorList() {
	bytes, err := ioutil.ReadFile("config/error_list.json")
	if err != nil {
		logs.Error.Println("Read error_list error: ", err)
	}

	if err := json.Unmarshal(bytes, &ErrorList); err != nil {
		logs.Error.Println("Unmarshal error_list error: ", err)
	}
}
