package model

import "time"

type Customer struct {
	Id        int64
	Mobile    string    `xorm:"unique"`
	CreatedAt time.Time `xorm:"created 'in_date_time'"`
	UpdatedAt time.Time `xorm:"updated 'modi_date_time'"`
}

func (Customer) TableName() string {
	return "user"
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
