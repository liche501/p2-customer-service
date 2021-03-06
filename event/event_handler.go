package event

import (
	"best/p2-customer-service/config"
	"best/p2-customer-service/logs"
	"best/p2-customer-service/model"
	"strconv"
	"strings"

	"encoding/json"
	"fmt"
	"reflect"

	"github.com/labstack/echo"
	"github.com/smallnest/goreq"
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
	numIsNewCust, _ := strconv.ParseInt(e.IsNewCust, 10, 64)
	bc.IsNewCust = numIsNewCust
	if err := bc.UpdateStatusAndCustNo(); err != nil {
		logs.Error.Println(err)
		return err
	}

	//SendCoupon
	// if err := SendCoupon(e.BrandCode, e.CustNo); err != nil {
	// 	logs.Error.Println(err)
	// }

	//sendEvent SendCouponInitiated (crm and wcs)
	sendCouponInitiated := SendCouponInitiated{}
	sendCouponInitiated.BrandCode = e.BrandCode
	sendCouponInitiated.CustNo = e.BrandCode
	es := new(EventSender)
	es.EventBrokerUrl = config.Config.Adapter.EventBrokerURL
	err := es.SendEvent("marketing", "SendCouponInitiated", sendCouponInitiated)
	if err != nil {
		logs.Error.Println(err)
	}

	bc2 := model.BrandCustomer{}
	bc2.BrandCode = e.BrandCode
	bc2.CustomerId = e.CustomerID
	bc2.Status = "BrandCustomerCreated"
	if err := bc2.UpdateStatus(); err != nil {
		logs.Error.Println(err)
		return err
	}

	return nil
}

func (e *BrandCustomerCreated) Handle() error {
	// logs.Warning.Println("BrandCustomerCreated ative")
	// logs.Warning.Println(e)

	// bc := model.BrandCustomer{}
	// bc.BrandCode = e.BrandCode
	// bc.CustomerId = e.CustomerID
	// bc.Status = "BrandCustomerCreated"
	// err := bc.UpdateStatus()
	// if err != nil {
	// 	logs.Error.Println(err)
	// 	return err
	// }
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

func SendCoupon(brandCode, custNo string) error {
	url := config.Config.Adapter.CSL.CustomerInterfaceAPI + "/" + strings.ToUpper(brandCode) + custNo + "/Coupons"
	logs.Debug.Println("[SendCoupon]url", url)
	_, textData, reqErr := goreq.New().Post(url).SetCurlCommand(true).End()
	if reqErr != nil || textData == "" {
		logs.Error.Println("reqErr", len(reqErr), " textData:", textData)
		return reqErr[0]
	}
	logs.Succ.Println("[SendCoupon] success", textData)

	return nil
}
