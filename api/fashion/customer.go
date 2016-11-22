package fashion

import (
	. "best/p2-customer-service/dto"
	"best/p2-customer-service/extends"
	"best/p2-customer-service/logs"
	"best/p2-customer-service/model"

	"net/http"

	"github.com/labstack/echo"
)

// APICheckMobile: 注册/修改手机号时检查注册状态
func APICheckMobile(c echo.Context) error {
	mobile := c.Get("user").(*extends.AuthClaims).Mobile
	openId := c.Get("user").(*extends.AuthClaims).OpenId
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode

	logs.Debug.Println(mobile, openId, brandCode)
	if mobile == "" || openId == "" || brandCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, 10012)
	}
	// status, err := models.UserShop{}.CheckMobileForRegister(mobile, openId, brandCode)
	// if err != nil {
	// 	logs.Error.Println("CheckMobile error: ", err)

	// 	return echo.NewHTTPError(http.StatusInternalServerError, 10013)

	// }
	status := "heihei"
	return c.JSON(http.StatusOK, APIResult{Success: true, Result: map[string]string{"status": status}})

}
func APIRegister(c echo.Context) error {

	mobile := c.FormValue("mobile")
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	openId := c.Get("user").(*extends.AuthClaims).OpenId
	verCode := c.FormValue("verCode")

	if verCode == "" || mobile == "" || openId == "" || brandCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, 10012)
	}

	//CheckSms
	flag, err := model.Sms{}.CheckVerCode(mobile, verCode)
	if err != nil {
		logs.Error.Println("Check sms error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)

	} else if !flag {
		logs.Error.Println("Sms code invalid")
		return c.JSON(http.StatusOK, APIResult{Error: APIError{Code: 20006, Message: "20006"}})

	}
	return echo.NewHTTPError(http.StatusBadRequest, 10012)

}
