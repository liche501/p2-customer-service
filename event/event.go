package event

import (
	"reflect"
	"time"
)

type Event struct {
	StreamName string      `json:"streamName"`
	EventType  string      `json:"eventType"`
	Payload    interface{} `json:"payload"`
	Timestamp  int64       `json:"timestamp,omitempty"`
	Partition  int32       `json:"partition,omitempty"`
	Offset     int64       `json:"offset,omitempty"`
}

var eventTypeMap = map[string]reflect.Type{
	"BrandCustomerConfirmed":  reflect.TypeOf(BrandCustomerConfirmed{}),
	"BrandCustomerCreated":    reflect.TypeOf(BrandCustomerCreated{}),
	"BrandCustomerFailed":     reflect.TypeOf(BrandCustomerFailed{}),
	"BrandCustomerDuplicated": reflect.TypeOf(BrandCustomerDuplicated{}),
	"SendCouponInitiated":     reflect.TypeOf(SendCouponInitiated{}),
}

type CustomerCreated struct {
	CustomerID int64     `json:"customerId"`
	Mobile     string    `json:"mobile"`
	CreatedAt  time.Time `json:"createdAt"`
}

type BrandCustomerInitiated struct {
	CustomerID int64     `json:"customerId"`
	Telephone  string    `json:"telephone"`
	Password   string    `json:"password"`
	BrandCode  string    `json:"brandCode"`
	WxOpenID   string    `json:"openId"`
	CreatedAt  time.Time `json:"createdAt"`
}

type BrandCustomerConfirmed struct {
	CustomerID int64     `json:"customerId"`
	CustNo     string    `json:"custNo"`
	BrandCode  string    `json:"brandCode"`
	IsNewCust  string    `json:"isNewCust"`
	CreatedAt  time.Time `json:"createdAt"`
}

type BrandCustomerCreated struct {
	CustomerID int64     `json:"customerId"`
	BrandCode  string    `json:"brandCode"`
	CustNo     string    `json:"custNo"`
	CreatedAt  time.Time `json:"createdAt"`
}

type BrandCustomerFailed struct {
	CustomerID int64     `json:"customerId"`
	BrandCode  string    `json:"brandCode"`
	CreatedAt  time.Time `json:"createdAt"`
}

type BrandCustomerDuplicated struct {
	CustomerID int64     `json:"customerId"`
	BrandCode  string    `json:"brandCode"`
	CreatedAt  time.Time `json:"createdAt"`
}

type SendCouponInitiated struct {
	BrandCode string    `json:"brandCode"`
	CustNo    string    `json:"custNo"`
	CreatedAt time.Time `json:"createdAt"`
}
