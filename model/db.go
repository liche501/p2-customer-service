package model

import (
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
)

var (
	CustomerNotExistError          = echo.NewHTTPError(http.StatusInternalServerError, 20003)
	BrandCustomerAlreadyExistError = echo.NewHTTPError(http.StatusInternalServerError, 20014)
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
		// db.ShowExecTime(true)
		// db.Logger().SetLevel(core.LOG_WARNING)
		db.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	}

	return db.Sync2(new(Customer), new(FashionBrandCustomer), new(BrandCustomer), new(Sms))

}
