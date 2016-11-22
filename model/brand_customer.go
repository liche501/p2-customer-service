package model

import "time"

type FashionBrandCustomer struct {
	Id        int64
	UserId    int64 `xorm:"index"`
	BrandCode string
	CustNo    string
	WxOpenId  string

	ReceiveAddress   string
	ReceiveTelephone string
	ReceiveName      string
	ReceiveSize      string

	Status    string    `xorm:"varchar(40)"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated 'mode_date_time'"`
}

type RetailBrandCustomer struct {
	Id        int64
	UserId    int64 `xorm:"index"`
	BrandCode string
	VipCode   string
	WxOpenId  string

	Status    string    `xorm:"varchar(40)"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated 'mode_date_time'"`
}
