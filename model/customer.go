package model

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Customer struct {
	Id        int64
	Mobile    string    `xorm:"unique"`
	CreatedAt time.Time `xorm:"created 'in_date_time'"`
	UpdatedAt time.Time `xorm:"updated 'modi_date_time'"`
}

func (Customer) TableName() string {
	return "user"
}

func (Customer) Get(id int64) (*Customer, error) {
	var c Customer
	has, err := db.Id(id).Get(&c)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New("User not exists")
	}
	return &c, nil
}

func (Customer) GetByMobile(mobile string) (*Customer, error) {
	var c Customer
	has, err := db.Where("mobile = ?", mobile).Get(&c)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New("User not exists")
	}
	return &c, nil
}

func (u *Customer) Create() error {
	exist, err := Customer{}.GetByMobile(u.Mobile)
	if err != nil {
		return err
	}

	if affected, err := db.InsertOne(u); err != nil {
		return err
	} else if affected == 0 {
		return errors.New("Affected rows : 0")
	}

	// TODO:: This is illogic. Have to change this logic.
	if exist != nil {
		if affected, err := db.Id(exist.Id).Cols("mobile").Update(&Customer{Mobile: strconv.FormatInt(u.Id, 10)}); err != nil {
			return err
		} else if affected == 0 {
			return errors.New("Affected rows : 0")
		}
	}
	return nil
}

func (Customer) ChangeMobileWithID(id int64, mobile string) error {
	exist, _ := Customer{}.GetByMobile(mobile)
	if exist != nil && exist.Id == id {
		return nil // do nothing
	}
	if exist != nil && exist.Id != id {
		// TODO:: This is illogic. Have to change this logic.
		if _, err := db.Id(exist.Id).Cols("mobile").Update(&Customer{Mobile: strconv.FormatInt(id, 10)}); err != nil {
			return fmt.Errorf("Cannot change exist mobile: ", err)
		}
	}

	affected, err := db.Id(id).Cols("mobile").Update(&Customer{Mobile: mobile})
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("Affected rows : 0")
	}
	return nil
}

func (Customer) ChangeMobileWithOld(oldMobile, newMobile string) error {
	if oldMobile == newMobile {
		return nil
	}

	exist, err := Customer{}.GetByMobile(oldMobile)
	if err != nil {
		return err
	}
	if exist == nil {
		return errors.New("Affected rows : 0")
	}

	if err := (Customer{}.ChangeMobileWithID(exist.Id, newMobile)); err != nil {
		return err
	}

	return nil
}

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

	exist, err := CustomerInfo{}.Get(u.Mobile, u.BrandCode)
	if err != nil {
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

func (CustomerInfo) Get(mobile, brandCode string) (*CustomerInfo, error) {
	user := CustomerInfo{}
	has, err := db.Where("mobile = ?", mobile).And("brand_code = ?", brandCode).Get(&user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, nil
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
