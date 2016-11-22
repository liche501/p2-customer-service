package model

import "time"

type FashionBrandCustomer struct {
	Id         int64
	CustomerId int64 `xorm:"index 'user_id'"`
	BrandCode  string
	CustNo     string
	WxOpenId   string `xorm:"'open_id'"`

	ReceiveAddress   string `xorm:"'receiv_address'"`
	ReceiveTelephone string `xorm:"'receiv_telephone'"`
	ReceiveName      string `xorm:"'receiv_name'"`
	ReceiveSize      string `xorm:"'receiv_size'"`

	Status    string    `xorm:"varchar(40)"`
	CreatedAt time.Time `xorm:"created 'in_date_time'"`
	UpdatedAt time.Time `xorm:"updated 'modi_date_time'"`
}

func (FashionBrandCustomer) TableName() string {
	return "user_shop"
}

type RetailBrandCustomer struct {
	Id         int64
	CustomerId int64  `xorm:"index 'user_id'"`
	WxOpenId   string `xorm:"'open_id'"`
	BrandCode  string
	VipCode    string

	Status    string    `xorm:"varchar(40)"`
	CreatedAt time.Time `xorm:"created 'in_date_time'"`
	UpdatedAt time.Time `xorm:"updated 'modi_date_time'"`
}

func (RetailBrandCustomer) TableName() string {
	return "user_mh"
}

type FashionBrandCustomerInfo struct {
	Customer             `xorm:"extends"`
	FashionBrandCustomer `xorm:"extends"`
}

type RetailBrandCustomerInfo struct {
	Customer            `xorm:"extends"`
	RetailBrandCustomer `xorm:"extends"`
}
