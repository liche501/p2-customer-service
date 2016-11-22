package model

import (
	"errors"
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
	user := CustomerInfo{}

	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	has, err := db.Where("mobile = ?", u.Mobile).And("brand_code = ?", u.BrandCode).Get(&user)
	if err != nil {
		return err
	} else if !has {
		// Insert
		affected, err := db.InsertOne(u)
		if err == nil && affected == 0 {
			return errors.New("Affected rows : 0")
		}
		return err
	}

	affected, err := db.Where("mobile = ?", u.Mobile).And("brand_code = ?", u.BrandCode).AllCols().Update(u)
	if err == nil && affected == 0 {
		return errors.New("Affected rows : 0")
	}
	return err
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

func (u *CustomerInfo) UpdateHasFilled() error {
	affected, err := db.Where("mobile = ?", u.Mobile).And("brand_code = ?", u.BrandCode).Cols("has_filled").Update(u)
	if err == nil && affected == 0 {
		return errors.New("Affected rows : 0")
	}
	return err
}

func (ud *CustomerInfo) SaveCustomerInfoWithMobile(mobile, oldMobile, openId string) error {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	user := Customer{}
	has, err := db.Where("mobile = ?", oldMobile).Get(&user)
	if err != nil {
		return err
	} else if !has {
		return nil
	}

	// 判断mobile是否已存在
	u := Customer{}
	has, err = db.Where("mobile = ?", mobile).Get(&u)
	if err != nil {
		return err
	}
	if has {
		// 1修改主表
		u.Mobile = strconv.FormatInt(u.Id, 10)
		affected, err := db.Where("mobile = ?", mobile).Cols("mobile").Update(&u)
		if err != nil {
			return err
		} else if affected == 0 {
			return errors.New("Affected rows : 0")
		}

		// 2修改Detail表
		udn := CustomerInfo{}
		udn.Mobile = strconv.FormatInt(u.Id, 10)
		affected, err = db.Where("mobile = ?", mobile).Cols("mobile").Update(&udn)
		if err != nil {
			return err
		}

		// 3修改ModernHouse
		mh := RetailBrandCustomer{}
		_, err = db.Where("user_id = ?", user.Id).Delete(&mh)
		if err != nil {
			return err
		}

		mhn := RetailBrandCustomer{}
		mhn.CustomerId = user.Id
		_, err = db.Where("user_id = ?", u.Id).Cols("user_id").Update(&mhn)
		if err != nil {
			return err
		}
	}

	// 改User表
	u.Mobile = mobile
	_, err = db.Where("mobile = ?", oldMobile).Cols("mobile").Update(&u)
	if err != nil {
		return err
	}

	// 改当前品牌
	udd := CustomerInfo{}
	has, err = db.Where("mobile = ?", oldMobile).And("brand_code = ?", ud.BrandCode).Get(&udd)
	if err != nil {
		return err
	} else if !has {
		// Insert
		affected, err := db.InsertOne(ud)
		if err == nil && affected == 0 {
			return errors.New("Affected rows : 0")
		}
		return err
	}

	_, err = db.Where("mobile = ?", oldMobile).And("brand_code = ?", ud.BrandCode).AllCols().Update(ud)
	if err != nil {
		return err
	}

	// 改所有品牌
	_, err = db.Where("mobile = ?", oldMobile).Cols("mobile").Update(ud)
	if err != nil {
		return err
	}
	return err
}
