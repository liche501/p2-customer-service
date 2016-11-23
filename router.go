package main

import (
	"best/p2-customer-service/api/common"
	"best/p2-customer-service/api/fashion"
	"best/p2-customer-service/event"

	. "best/p2-customer-service/dto"
	"best/p2-customer-service/extends"

	"net/http"

	"github.com/labstack/echo"
)

func RouterInit() {
	// Login route
	e.POST("/login", login)
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
		return echo.NewHTTPError(http.StatusInternalServerError, "10003", " error01")
	})
	e.GET("/error2", func(c echo.Context) error {
		return c.JSON(http.StatusOK, APIResult{Error: APIError{Code: 10001, Message: "StatusBadRequest"}})
	})
	e.GET("/event",func(c echo.Context)error{
		 aa := new(event.EventSender)
		 aa.EventBrokerUrl = "http://localhost:9000"
		 aa.SendEvent()
		return c.JSON(http.StatusOK, APIResult{Success:true}})
	})
	e.GET("/token", func(c echo.Context) error {
		token, _ := extends.AuthHandler("rc", "oYiR6wTz6anr5KpiRH-mRcpvvLPc", "13691194223", "0001852359")
		return c.JSON(http.StatusOK, APIResult{Success: true, Result: token})
	})
	t := e.Group("/jwt")
	// t.Use(extends.JWTMiddleware)
	t.GET("", extends.JWTMiddleware(func(c echo.Context) error {
		return c.JSON(http.StatusOK, c.Get("openIdWithToken"))
	}))

	api := e.Group("/api")
	v1 := api.Group("/v1")
	fa := v1.Group("/fashion")

	//User
	user := fa.Group("/user")
	user.POST("/register", extends.JWTMiddleware(fashion.APIRegister))
	user.GET("/login", extends.JWTMiddleware(fashion.APILogin))
	user.GET("/get_customer_info", extends.JWTMiddleware(fashion.APIGetCustomerInfo))
	user.GET("/get_user_info", extends.JWTMiddleware(fashion.APIGetUserInfo))
	user.POST("/update_perfect_info", extends.JWTMiddleware(fashion.APIUpdatePerfectInfo))
	user.GET("/check_mobile", extends.JWTMiddleware(fashion.APICheckMobile))
	user.GET("/get_member_info", extends.JWTMiddleware(fashion.APIGetMemberInfo))

	//Coupon
	co := fa.Group("/coupon")
	co.GET("/get_coupon_list", extends.JWTMiddleware(fashion.APIGetCouponList))

	//Integral
	in := fa.Group("/integral")
	in.GET("/get_current_integral", demo)
	in.GET("/get_integral_history", demo)
	in.GET("/get_vip_grade", demo)
	in.GET("/update_integral_exchange", demo)

	// Common
	c := v1.Group("/common")
	s := c.Group("/sms")
	s.GET("/code", common.ApiSendSms)
	s.GET("/success", common.ApiCheckSms)
	s.GET("/ative", common.Active)

	p := c.Group("/captcha")
	p.GET("/key", common.APIGetCaptchaKey)
	p.GET("/image", demo)
	p.GET("/success", common.ApiCheckCaptcha)

	a := c.Group("/auth")
	a.GET("/set_auth", demo)

	v1.POST("/events", event.ApiHandleEvent)
}
