package model

import (
	"errors"
	"sync"
	"time"
)

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

func (FashionBrandCustomer) GetByWxOpenID(brandCode, wxOpenID string) (*FashionBrandCustomer, error) {
	c := FashionBrandCustomer{BrandCode: brandCode, WxOpenID: wxOpenID}
	has, err := db.Get(&c)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, CustomerNotExistError
	}

	return &c, nil
}

func (FashionBrandCustomer) GetByMobile(brandCode, mobile string) (*FashionBrandCustomer, error) {
	c, err := Customer{}.GetByMobile(mobile)
	if err != nil {
		return nil, err
	}
	return FashionBrandCustomer{}.GetByCustomerID(brandCode, c.Id)
}

func (FashionBrandCustomer) GetByCustomerID(brandCode string, customerID int64) (*FashionBrandCustomer, error) {
	c := FashionBrandCustomer{BrandCode: brandCode, CustomerId: customerID}
	has, err := db.Get(&c)

	if err != nil {
		return nil, err
	}

	if !has {
		return nil, CustomerNotExistError
	}

	return &c, nil
}

type FashionBrandCustomerInfo struct {
	Customer             `xorm:"extends"`
	FashionBrandCustomer `xorm:"extends"`
}

func (u *FashionBrandCustomerInfo) Create() error {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	// return if exist
	exist, err := FashionBrandCustomer{}.GetByWxOpenID(u.FashionBrandCustomer.BrandCode, u.FashionBrandCustomer.WxOpenID)
	if err != nil && err != CustomerNotExistError {
		return err
	}
	if exist != nil {
		return BrandCustomerAlreadyExistError
	}

	// check Customer exist
	customer, err := Customer{}.GetByMobile(u.Customer.Mobile)
	if err != nil && err != CustomerNotExistError {
		return err
	}
	if customer == nil {
		if err := u.Customer.Create(); err != nil && err != CustomerNotExistError {
			return err
		}
		customer = &u.Customer
	}

	// create CustomerInfo
	customerInfo, err := CustomerInfo{}.Get(u.FashionBrandCustomer.BrandCode, u.Mobile)
	if err != nil && err != CustomerNotExistError {
		return err
	}

	if customerInfo == nil {
		customerInfo = &CustomerInfo{
			CustomerId: customer.Id,
			Name:       u.FashionBrandCustomer.ReceiveName,
			Mobile:     customer.Mobile,
			BrandCode:  u.FashionBrandCustomer.BrandCode,
		}
		if err := customerInfo.Save(); err != nil {
			return err
		}
	}

	// create FashionBrandCustomer
	u.FashionBrandCustomer.CustomerId = customer.Id
	affected, err := db.InsertOne(&u.FashionBrandCustomer)
	if err != nil {
		return err
	}
	if affected == 0 {
		err = errors.New("Affected rows : 0")
	}
	return nil
}

func (FashionBrandCustomerInfo) Delete(brandCode string, customerID int64) error {
	c, err := FashionBrandCustomer{}.GetByCustomerID(brandCode, customerID)
	if err != nil {
		return err
	}

	var deleted FashionBrandCustomer
	affected, err := db.Id(c.Id).Delete(&deleted)
	if err != nil {
		return err
	}
	if affected == 0 {
		err = errors.New("Affected rows : 0")
	}
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

func (FashionBrandCustomerInfo) GetByWxOpenIDAndStatus(brandCode, wxOpenId, status string) (*FashionBrandCustomerInfo, error) {

	return nil, nil
}
