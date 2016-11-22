package model

import "time"

type FashionBrandCustomer struct {
	Id         int64
	CustomerId int64 `xorm:"index 'user_id'"`
	BrandCode  string
	CustNo     string
	WxOpenID   string `xorm:"'open_id'"`

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

type FashionBrandCustomerInfo struct {
	Customer             `xorm:"extends"`
	FashionBrandCustomer `xorm:"extends"`
}

type RetailBrandCustomerInfo struct {
	Customer            `xorm:"extends"`
	RetailBrandCustomer `xorm:"extends"`
}

func (u *FashionBrandCustomerInfo) Create() error {
	return nil
}

func (FashionBrandCustomerInfo) Delete(mobile, brandCode string) error {
	return nil
}
func (u *FashionBrandCustomerInfo) UpdateCustNo() error {
	return nil
}

func (u *FashionBrandCustomerInfo) UpdatePassword() error {
	return nil
}

func (u *FashionBrandCustomerInfo) UpdateForGame() error {
	return nil
}

func (FashionBrandCustomerInfo) GetByMobile(mobile, brandCode string) (*FashionBrandCustomer, error) {
	return nil, nil
}

func (FashionBrandCustomerInfo) GetByWxOpenIDAndStatus(brandCode, openId, status string) (*FashionBrandCustomerInfo, error) {

	return nil, nil
}

func (FashionBrandCustomer) CheckWxOpenID(brandCode, openId string) (bool, error) {
	return true, nil
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
