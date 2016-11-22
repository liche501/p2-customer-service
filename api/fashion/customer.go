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

// curl -X POST -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJicmFuZENvZGUiOiJyYyIsIm9wZW5JZCI6Im9ZaVI2d1R6NmFucjVLcGlSSC1tUmNwdnZMUGMiLCJtb2JpbGUiOiIxMzY5MTE5NDIyMyIsImN1c3RObyI6IjAwMDE4NTIzNTkiLCJleHAiOjE0ODAwNjQ2NDIsImlzcyI6ImxpY2hlIn0.huxuvLITetwHzdpZHX-T_sfZe0rEeMM_2DOnugdUjRo" -H "Cache-Control: no-cache" -H "Postman-Token: 28bf6f4d-9809-26a2-3229-4a177c8d29cf" "http://localhost:9000/api/v1/fashion/user/register?mobile=13691194223&verCode=1234"
func APIRegister(c echo.Context) error {

	mobile := c.FormValue("mobile")
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	openId := c.Get("user").(*extends.AuthClaims).OpenId
	verCode := c.FormValue("verCode")

	if verCode == "" || mobile == "" || openId == "" || brandCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, 10012)
	}

	//CheckSms
	// flag, err := model.Sms{}.CheckVerCode(mobile, verCode)
	// if err != nil {
	// 	logs.Error.Println("Check sms error: ", err)
	// 	return echo.NewHTTPError(http.StatusInternalServerError, 10013)

	// } else if !flag {
	// 	logs.Error.Println("Sms code invalid")
	// 	return c.JSON(http.StatusOK, APIResult{Error: APIError{Code: 20006, Message: "20006"}})

	// }
	logs.Debug.Println(brandCode, openId)

	fashionBrandCustomerInfo := new(model.FashionBrandCustomerInfo)
	fashionBrandCustomerInfo.Customer.Mobile = mobile
	fashionBrandCustomerInfo.FashionBrandCustomer.BrandCode = brandCode
	fashionBrandCustomerInfo.FashionBrandCustomer.WxOpenID = openId
	logs.Debug.Println(fashionBrandCustomerInfo)

    //customer regist 
	if err := fashionBrandCustomerInfo.Create(); err != nil {
		logs.Error.Println(err)

	}
    

	return echo.NewHTTPError(http.StatusBadRequest, 10012)

}
