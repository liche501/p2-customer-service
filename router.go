package main

import (
	"best/p2-customer-service/api/common"
	"best/p2-customer-service/api/fashion"
	"best/p2-customer-service/event"
	"best/p2-customer-service/logs"

	. "best/p2-customer-service/dto"
	"best/p2-customer-service/extends"

	"net/http"

	"github.com/labstack/echo"
)

func RouterInit() {
	// Login route
	e.POST("/createtoken", CreateToken)
	// Unauthenticated route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "accessible")
	})
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	e.GET("/success", func(c echo.Context) error {
		return c.JSON(http.StatusOK, APIResult{Success: true, Result: "aaaaaaa"})
	})
	e.GET("/error", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, 10003)
	})
	e.GET("/error2", func(c echo.Context) error {
		return c.JSON(http.StatusOK, APIResult{Error: APIError{Code: 400, Message: "StatusBadRequest"}})
	})
	e.GET("/event/:eventType", func(c echo.Context) error {
		aa := new(event.EventSender)
		// url := fmt.Sprintf("/v1/streams/%v/events/%v", "marketing", "BrandCustomerInitiated")
		aa.EventBrokerUrl = "http://staging.p2shop.cn:50110"
		var payload interface{}
		c.Bind(&payload)
		err := aa.SendEvent("marketing", "BrandCustomerInitiated", payload)
		if err != nil {
			logs.Error.Println(err)
		}
		return c.JSON(http.StatusOK, APIResult{Success: true})
	})
	e.GET("/token", func(c echo.Context) error {
		token, err := extends.AuthHandler("rc", "oYiR6wTz6anr5KpiRH-mRcpvvLPc", "13691194223", "0001852359")
		if err == nil {
			return err
		}
		return c.JSON(http.StatusOK, APIResult{Success: true, Result: token})
	})
	t := e.Group("/jwt")
	t.Use(extends.JWTMiddleware)
	t.GET("", func(c echo.Context) error {

		return c.String(http.StatusOK, "Welcome")
	})

	api := e.Group("/api")
	v1 := api.Group("/v1")

	// Common
	c := v1.Group("/common")
	c.POST("/events", event.ApiHandleEvent)

	s := c.Group("/sms")
	s.GET("/code", common.ApiSendSms)
	s.GET("/success", common.ApiCheckSms)
	s.GET("/active", common.Active)

	p := c.Group("/captcha")
	p.GET("/key", common.APIGetCaptchaKey)
	p.GET("/image", demo)
	p.GET("/success", common.ApiCheckCaptcha)

	a := c.Group("/auth")
	a.GET("/set_auth", demo)

	fa := v1.Group("/fashion")
	fa.Use(extends.JWTMiddleware) //********** important check JWT prefix /fashion  *****
	//User
	user := fa.Group("/user")

	user.POST("/register", fashion.APIRegister)
	user.GET("/login", fashion.APILogin)
	user.GET("/get_customer_info", fashion.APIGetCustomerInfo)
	user.GET("/get_user_info", fashion.APIGetUserInfo)
	user.POST("/update_perfect_info", fashion.APIUpdatePerfectInfo)
	user.GET("/check_mobile", fashion.APICheckMobileAvailableForRegister)
	user.GET("/get_member_info", fashion.APIGetMemberInfo)

	//Coupon
	co := fa.Group("/coupon")
	co.GET("/get_coupon_list", fashion.APIGetCouponList)

	//Integral
	in := fa.Group("/integral")
	in.GET("/current", fashion.ApiGetCurrentIntegral)
	in.GET("/history", fashion.ApiGetIntegralHistory)
	in.GET("/grade", fashion.ApiGetVipGrade)
	in.GET("/update_integral_exchange", fashion.ApiUpdateIntegralExchange)

}
