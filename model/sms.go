package model

import (
	"errors"
	"time"
)

// `xorm:"index(IDX_SMS)"`
type Sms struct {
	BrandCode string
	Mobile    string `xorm:"index"`
	Type      string
	VerCode   string

	InDateTime   time.Time `xorm:"created"`
	ModiDateTime time.Time `xorm:"updated"`
}

func (s *Sms) Create() error {
	affected, err := db.InsertOne(s)
	if err == nil && affected == 0 {
		err = errors.New("Affected rows : 0")
	}
	return err
}

func (Sms) CheckVerCode(mobile, verCode string) (bool, error) {
	s := Sms{}
	has, err := db.Where("mobile = ?", mobile).And("ver_code = ?", verCode).Desc("in_date_time").Get(&s)
	if err != nil {
		return false, err
	}

	if !has || s.ModiDateTime.Before(time.Now().Add(-time.Minute*30)) {
		return false, nil
	}

	return true, nil
}

func (Sms) CheckRepeatVerCode(verCode string) (bool, error) {
	s := Sms{}
	startTimestamp := time.Now().Unix()
	startTm := time.Unix(startTimestamp, 0)
	startTimeStr := startTm.Format("2006-01-02")

	has, err := db.Where("in_date_time > ?", startTimeStr).And("ver_code = ?", verCode).Get(&s)
	// logs.Debug.Println(has)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil

	}

	return true, nil
}
