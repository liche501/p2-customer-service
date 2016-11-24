package event

import (
	"best/p2-customer-service/logs"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpError struct {
	Endpoint   string
	StatusCode int
	Body       string
}

func (e HttpError) Error() string {
	return fmt.Sprintf("Endpoint:%s, StatusCode:%d, Body:%s", e.Endpoint, e.StatusCode, e.Body)
}

type EventSender struct {
	EventBrokerUrl string
}

func (s *EventSender) SendEvent(streamName, eventType string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{}
	url := fmt.Sprintf("%s/v1/streams/%s/events/%s", s.EventBrokerUrl, streamName, eventType)
	logs.Debug.Println(url)
	logs.Debug.Println("GetUserDetail error: ", err)
	request, err := http.NewRequest("POST", url, bytes.NewReader(data))
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		logs.Error.Println("GetUserDetail error: ", err)
		return err
	}

	defer response.Body.Close()

	b, _ := ioutil.ReadAll(response.Body)
	logs.Debug.Println(time.Now(), string(b))

	if response.StatusCode >= 500 {
		return HttpError{
			Endpoint:   s.EventBrokerUrl,
			StatusCode: response.StatusCode,
			Body:       string(b),
		}
	}

	return nil
}
