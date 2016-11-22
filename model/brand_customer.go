package model

import "time"

type FashionBrandCustomer struct {
	Id         int64
	CustomerId int64 `xorm:"index"`
	BrandCode  string
	CustNo     string
	WxOpenId   string

	ReceiveAddress   string
	ReceiveTelephone string
	ReceiveName      string
	ReceiveSize      string

	Status    string    `xorm:"varchar(40)"`
	CreatedAt time.Time `xorm:"created 'in_date_time'"`
	UpdatedAt time.Time `xorm:"updated 'modi_date_time'"`
}

type RetailBrandCustomer struct {
	Id         int64
	CustomerId int64 `xorm:"index"`
	BrandCode  string
	VipCode    string
	WxOpenId   string

	Status    string    `xorm:"varchar(40)"`
	CreatedAt time.Time `xorm:"created 'in_date_time'"`
	UpdatedAt time.Time `xorm:"updated 'modi_date_time'"`
}

type FashionBrandCustomerInfo struct {
	Customer             `xorm:"extends"`
	FashionBrandCustomer `xorm:"extends"`
}

type RetailBrandCustomerInfo struct {
	Customer            `xorm:"extends"`
	RetailBrandCustomer `xorm:"extends"`
}
