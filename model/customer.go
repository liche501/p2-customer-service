package model

import (
	"errors"
	"fmt"
	"strconv"
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
		return nil, nil
	}
	return &c, nil
}

func (Customer) GetByMobile(mobile string) (*Customer, error) {
	var c Customer
	has, err := db.Where("mobile = ?", mobile).Get(&c)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}
	return &c, nil
}

func (c *Customer) Create() error {
	customer, err := Customer{}.GetByMobile(c.Mobile)
	if err != nil {
		return err
	}

	// WillDO:: This is illogic. Have to change this logic.
	if customer != nil {
		if affected, err := db.Id(customer.Id).Cols("mobile").Update(&Customer{Mobile: strconv.FormatInt(c.Id, 10)}); err != nil {
			return err
		} else if affected == 0 {
			return errors.New("Affected rows : 0")
		}
	} else {
		if affected, err := db.InsertOne(c); err != nil {
			return err
		} else if affected == 0 {
			return errors.New("Affected rows : 0")
		}
	}
	return nil
}

func (Customer) ChangeMobileWithID(id int64, mobile string) error {
	customer, _ := Customer{}.GetByMobile(mobile)
	if customer != nil && customer.Id == id {
		return nil // do nothing
	}
	if customer != nil && customer.Id != id {
		// TODO:: This is illogic. Have to change this logic.
		if _, err := db.Id(customer.Id).Cols("mobile").Update(&Customer{Mobile: strconv.FormatInt(id, 10)}); err != nil {
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

	customer, err := Customer{}.GetByMobile(oldMobile)
	if err != nil {
		return err
	}
	if customer == nil {
		return errors.New("Affected rows : 0")
	}

	if err := (Customer{}.ChangeMobileWithID(customer.Id, newMobile)); err != nil {
		return err
	}

	return nil
}
