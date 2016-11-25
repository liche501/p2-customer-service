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
	BrandCustomer        BrandCustomer `xorm:"extends"`
}

func (fbci *FashionBrandCustomerInfo) Create() error {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	//WillDo: this logic have to change
	// return if exist
	fashionBrandCustomer, err := FashionBrandCustomer{}.GetByWxOpenID(fbci.FashionBrandCustomer.BrandCode, fbci.FashionBrandCustomer.WxOpenID)
	if err != nil && err != CustomerNotExistError {
		return err
	}
	if fashionBrandCustomer != nil {
		return BrandCustomerAlreadyExistError
	}

	// check Customer exist
	customer, err := Customer{}.GetByMobile(fbci.Customer.Mobile)
	if err != nil && err != CustomerNotExistError {
		return err
	}
	if customer == nil {
		if err := fbci.Customer.Create(); err != nil && err != CustomerNotExistError {
			return err
		}
		customer = &fbci.Customer
	}

	// create BrandCustomer
	brandCustomer, err := BrandCustomer{}.Get(fbci.FashionBrandCustomer.BrandCode, fbci.Mobile)
	if err != nil && err != CustomerNotExistError {
		return err
	}

	if brandCustomer == nil {
		brandCustomer = &BrandCustomer{
			CustomerId: customer.Id,
			Mobile:     customer.Mobile,
			WxOpenID:   fbci.FashionBrandCustomer.WxOpenID,
			BrandCode:  fbci.FashionBrandCustomer.BrandCode,
			Status:     "CustomerCreated",
		}
		if err := brandCustomer.Save(); err != nil {
			return err
		}
	}

	// create FashionBrandCustomer
	fbci.FashionBrandCustomer.CustomerId = customer.Id
	affected, err := db.InsertOne(&fbci.FashionBrandCustomer)
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

func (FashionBrandCustomerInfo) GetByWxOpenID(brandCode, wxOpenId string) (*FashionBrandCustomerInfo, error) {
	var c FashionBrandCustomerInfo
	has, err := db.Table("user").Join("INNER", "user_detail", "user_detail.user_id = user.id").
		Join("INNER", "user_shop", "user_shop.user_id = user.id").
		Where("user_shop.open_id = ?", wxOpenId).And("user_shop.brand_code = ?", brandCode).
		Get(&c)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, CustomerNotExistError
	}
	return &c, nil
}

func (FashionBrandCustomerInfo) GetSuccessUserByWxOpenID(brandCode, wxOpenId string) (*FashionBrandCustomerInfo, error) {
	const successState = "BrandCustomerCreated"
	var c FashionBrandCustomerInfo
	has, err := db.Table("user").Join("INNER", "user_detail", "user_detail.user_id = user.id").
		Join("INNER", "user_shop", "user_shop.user_id = user.id").
		Where("user_shop.open_id = ?", wxOpenId).And("user_shop.brand_code = ?", brandCode).And("user_detail.status = ?", successState).
		Get(&c)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, CustomerNotExistError
	}
	return &c, nil
}

func (c *FashionBrandCustomerInfo) Status() string {
	return c.BrandCustomer.Status
}
