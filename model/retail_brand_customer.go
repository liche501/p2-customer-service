package model

import "time"


type RetailBrandCustomer struct {
	Id         int64
	CustomerId int64  `xorm:"index 'user_id'"`
	WxOpenID   string `xorm:"'open_id'"`
	BrandCode  string
	VipCode    string

	Status    string    `xorm:"varchar(40)"`
	CreatedAt time.Time `xorm:"created 'in_date_time'"`
	UpdatedAt time.Time `xorm:"updated 'modi_date_time'"`
}

func (RetailBrandCustomer) TableName() string {
	return "user_mh"
}


type RetailBrandCustomerInfo struct {
	Customer            `xorm:"extends"`
	RetailBrandCustomer `xorm:"extends"`
}


func (u *RetailBrandCustomer) Create() error {
	return nil
}

func (RetailBrandCustomer) GetByMobile(mobile, brandCode string) (*RetailBrandCustomer, error) {
	return nil, nil
}

func (u RetailBrandCustomer) Delete() error {
	return nil
}

func (u *RetailBrandCustomer) UpdateIsOldCust(isOldCust bool) error {
	return nil
}
