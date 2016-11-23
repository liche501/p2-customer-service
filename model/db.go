package model

import (
	"errors"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	CustomerNotExistError          = errors.New("User not exists")
	BrandCustomerAlreadyExistError = errors.New("BrandCustomer already exists")
)

var db *xorm.Engine

func InitDB(dialect, conn string) error {
	var err error
	db, err = xorm.NewEngine(dialect, conn)
	if err != nil {
		return err
	}

	isDebug := os.Getenv("WXSHOPDEBUG")
	if len(isDebug) > 0 {
		db.ShowSQL(true)
	}

	return db.Sync2(new(Customer), new(FashionBrandCustomer), new(BrandCustomer), new(RetailBrandCustomer), new(Sms))

}
