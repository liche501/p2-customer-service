package event

import (
	"best/p2-customer-service/logs"
	"best/p2-customer-service/model"

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
	logs.Warning.Println("BrandCustomerConfirmed ative")
	logs.Warning.Println(e)
	bc := model.BrandCustomer{}
	bc.BrandCode = e.BrandCode
	bc.CustomerId = e.CustomerID
	bc.Status = "BrandCustomerConfirmed"
	bc.CustNo = e.CustNo
	err := bc.UpdateStatusAndCustNo()
	if err != nil {
		logs.Error.Println(err)
		return err
	}
	//WillDo:: SendCoupon
	// err = fashion.SendCoupon(e.BrandCode, e.CustNo)
	// if err != nil {
	// 	logs.Error.Println(err)
	// }
	return nil
}

func (e *BrandCustomerCreated) Handle() error {
	logs.Warning.Println("BrandCustomerCreated ative")
	logs.Warning.Println(e)

	return nil
}

func (e *BrandCustomerFailed) Handle() error {
	logs.Warning.Println("BrandCustomerFailed ative")
	logs.Warning.Println(e)

	bc := model.BrandCustomer{}
	bc.BrandCode = e.BrandCode
	bc.CustomerId = e.CustomerID
	bc.Status = "BrandCustomerFailed"
	err := bc.UpdateStatus()
	if err != nil {
		logs.Error.Println(err)
		return err
	}
	return nil
}

func (e *BrandCustomerDuplicated) Handle() error {
	logs.Warning.Println("BrandCustomerDuplicated ative")
	logs.Warning.Println(e)
	bc := model.BrandCustomer{}
	bc.BrandCode = e.BrandCode
	bc.CustomerId = e.CustomerID
	bc.Status = "BrandCustomerDuplicated"
	err := bc.UpdateStatus()
	if err != nil {
		logs.Error.Println(err)
		return err
	}
	return nil
}
