package model

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

type CustomerInfo struct {
	Id            int64
	CustomerId    int64 `xorm:"'user_id'"`
	Name          string
	Mobile        string `xorm:"index"`
	BrandCode     string
	Gender        string
	Birthday      string
	Address       string
	DetailAddress string
	Email         string
	IsMarried     bool
	IsNewCust     int64
	HasFilled     bool

	CreatedAt time.Time `xorm:"created 'in_date_time'"`
	UpdatedAt time.Time `xorm:"updated 'modi_date_time'"`
}

func (CustomerInfo) TableName() string {
	return "user_detail"
}

func (u *CustomerInfo) Save() error {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	exist, err := CustomerInfo{}.Get(u.BrandCode, u.Mobile)
	if err != nil && err != CustomerNotExistError {
		return err
	}
	if exist != nil {
		affected, err := db.Id(exist.Id).AllCols().Update(u)
		if err != nil {
			return err
		}
		if affected == 0 {
			return errors.New("Affected rows : 0")
		}
		return nil
	}

	affected, err := db.InsertOne(u)
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("Affected rows : 0")
	}
	return nil
}

func (CustomerInfo) Get(brandCode, mobile string) (*CustomerInfo, error) {
	user := CustomerInfo{}
	has, err := db.Where("mobile = ?", mobile).And("brand_code = ?", brandCode).Get(&user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, CustomerNotExistError
	}
	return &user, nil
}

func (CustomerInfo) FindByMobile(mobile string) (customers []CustomerInfo, err error) {
	err = db.Where("mobile = ?", mobile).Find(&customers)
	return
}

func (u *CustomerInfo) UpdateHasFilled() error {
	affected, err := db.Where("mobile = ?", u.Mobile).And("brand_code = ?", u.BrandCode).Cols("has_filled").Update(u)
	if err == nil && affected == 0 {
		return errors.New("Affected rows : 0")
	}
	return err
}

func (CustomerInfo) ChangeMobileWithOld(oldMobile, newMobile string) error {
	if oldMobile == newMobile {
		return nil
	}

	if err := (Customer{}.ChangeMobileWithOld(oldMobile, newMobile)); err != nil {
		return nil
	}

	// TODO:: This is illogic. Have to remove this logic.
	c, err := Customer{}.GetByMobile(newMobile)
	newMobileCustomers, err := CustomerInfo{}.FindByMobile(newMobile)
	if err != nil {
		return err
	}
	if len(newMobileCustomers) > 0 {
		for _, exist := range newMobileCustomers {
			db.Id(exist.Id).Update(&CustomerInfo{Mobile: strconv.FormatInt(c.Id, 10)})
		}
	}
	oldMobileCustomers, err := CustomerInfo{}.FindByMobile(oldMobile)
	if err != nil {
		return err
	}
	if len(oldMobileCustomers) > 0 {
		for _, exist := range oldMobileCustomers {
			db.Id(exist.Id).Update(&CustomerInfo{CustomerId: c.Id, Mobile: newMobile})
		}
	}

	return nil
}
