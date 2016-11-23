package event

import (
	"best/p2-customer-service/logs"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/labstack/echo"
)

type EventHandler interface {
	Handle() error
}

func ApiHandleEvent(c echo.Context) error {
	ignore := []string{"CustomerCreated", "BrandCustomerInitiated"}
	var e Event
	c.Bind(&e)
	logs.Debug.Println(e)
	for _, s := range ignore {
		if s == e.EventType {
			return nil
		}
	}

	eventType, exist := eventTypeMap[e.EventType]
	if !exist {
		return fmt.Errorf("%s is not supported event\n", e.EventType)
	}

	b, err := json.Marshal(e.Payload)
	if err != nil {
		return nil
	}

	payload := reflect.New(eventType).Interface()
	if err := json.Unmarshal(b, payload); err != nil {
		return nil
	}
	logs.Debug.Println(payload)
	handler, ok := payload.(EventHandler)
	if !ok {
		return fmt.Errorf("%s is not implemented\n", e.EventType)
	}

	if err := handler.Handle(); err != nil {
		return err
	}

	return nil
}

func (e *BrandCustomerConfirmed) Handle() error {
	return nil
}

func (e *BrandCustomerCreated) Handle() error {
	return nil
}

func (e *BrandCustomerFailed) Handle() error {
	return nil
}

func (e *BrandCustomerDuplicated) Handle() error {
	return nil
}
