package model

import "time"

type Sms struct {
	BrandCode string
	Mobile    string `xorm:"index"`
	Type      string
	VerCode   string

	InDateTime   time.Time `xorm:"created"`
	ModiDateTime time.Time `xorm:"updated"`
}
