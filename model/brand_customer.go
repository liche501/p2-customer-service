package model

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

// TODO:: WxOpenID, CustNo 로직 검증(jang.jaehue)
type BrandCustomer struct {
	Id            int64
	CustomerId    int64 `xorm:"'user_id'"`
	Name          string
	Mobile        string `xorm:"index"`
	WxOpenID      string `xorm:"'wx_open_id'"`
	CustNo        string
	BrandCode     string
	Gender        string
	Birthday      string
	Address       string
	DetailAddress string
	Email         string
	IsMarried     bool
	IsNewCust     int64
	HasFilled     bool

	Status    string
	CreatedAt time.Time `xorm:"created 'in_date_time'"`
	UpdatedAt time.Time `xorm:"updated 'modi_date_time'"`
}

func (BrandCustomer) TableName() string {
	return "user_detail"
}

func (bc *BrandCustomer) Create() error {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	brandCustomer, err := BrandCustomer{}.Get(bc.BrandCode, bc.Mobile)
	if err != nil {
		return err
	}

	if brandCustomer == nil {
		affected, err := db.InsertOne(bc)
		if err != nil {
			return err
		}
		if affected == 0 {
			return errors.New("Affected rows : 0")
		}
	}

	return nil
}

func (bc *BrandCustomer) Update() error {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	brandCustomer, err := BrandCustomer{}.Get(bc.BrandCode, bc.Mobile)
	if err != nil {
		return err
	}
	if brandCustomer == nil {
		return errors.New("BrandCustomer not exist")
	}

	affected, err := db.Id(brandCustomer.Id).Cols("name", "mobile", "gender", "birthday", "address", "detail_address", "email", "is_married").Update(bc)
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("Affected rows : 0")
	}
	return nil

}
func (BrandCustomer) Get(brandCode, mobile string) (*BrandCustomer, error) {
	brandCustomer := BrandCustomer{}
	has, err := db.Where("mobile = ?", mobile).And("brand_code = ?", brandCode).Get(&brandCustomer)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}
	return &brandCustomer, nil
}

func (BrandCustomer) GetByWxOpenID(brandCode, wxOpenID string) (*BrandCustomer, error) {
	brandCustomer := BrandCustomer{}
	has, err := db.Where("wx_open_id = ?", wxOpenID).And("brand_code = ?", brandCode).Get(&brandCustomer)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}
	return &brandCustomer, nil
}

func (BrandCustomer) FindByMobile(mobile string) (customers []BrandCustomer, err error) {
	err = db.Where("mobile = ?", mobile).Find(&customers)
	return customers, err
}

func (bc *BrandCustomer) UpdateHasFilled() error {
	affected, err := db.Where("mobile = ?", bc.Mobile).And("brand_code = ?", bc.BrandCode).Cols("has_filled").Update(bc)
	if err == nil && affected == 0 {
		return errors.New("Affected rows : 0")
	}
	return err
}

func (BrandCustomer) ChangeMobileWithOld(oldMobile, newMobile string) error {
	if oldMobile == newMobile {
		return nil
	}

	if err := (Customer{}.ChangeMobileWithOld(oldMobile, newMobile)); err != nil {
		return nil
	}

	// TODO:: This is illogic. Have to remove this logic.
	c, err := Customer{}.GetByMobile(newMobile)
	newMobileCustomers, err := BrandCustomer{}.FindByMobile(newMobile)
	if err != nil {
		return err
	}
	if len(newMobileCustomers) > 0 {
		for _, exist := range newMobileCustomers {
			db.Id(exist.Id).Update(&BrandCustomer{Mobile: strconv.FormatInt(c.Id, 10)})
		}
	}
	oldMobileCustomers, err := BrandCustomer{}.FindByMobile(oldMobile)
	if err != nil {
		return err
	}
	if len(oldMobileCustomers) > 0 {
		for _, exist := range oldMobileCustomers {
			db.Id(exist.Id).Update(&BrandCustomer{CustomerId: c.Id, Mobile: newMobile})
		}
	}

	return nil
}

func (bc *BrandCustomer) UpdateStatusAndCustNo() error {
	affected, err := db.Where("user_id = ?", bc.CustomerId).And("brand_code = ?", bc.BrandCode).Cols("status", "cust_no").Update(bc)
	if err == nil && affected == 0 {
		return errors.New("Affected rows : 0")
	}
	return err
}
func (bc *BrandCustomer) UpdateStatus() error {
	affected, err := db.Where("user_id = ?", bc.CustomerId).And("brand_code = ?", bc.BrandCode).Cols("status").Update(bc)
	if err == nil && affected == 0 {
		return errors.New("Affected rows : 0")
	}
	return err
}
