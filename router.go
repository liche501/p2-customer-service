package main

import (
	"best/p2-customer-service/api/common"
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

	e.GET("/token", func(c echo.Context) error {
		token, _ := extends.AuthHandler("rc", "openid_111111111", "13691194223", "cust_1000")
		return c.JSON(http.StatusOK, APIResult{Success: true, Result: token})
	})
	t := e.Group("/jwt")
	// t.Use(extends.JWTMiddleware)
	t.GET("", extends.JWTMiddleware(func(c echo.Context) error {
		return c.JSON(http.StatusOK, c.Get("openIdWithToken"))
	}))

	api := e.Group("/api")
	v1 := api.Group("/v1")
	shop := v1.Group("/shop")

	//User
	user := shop.Group("/user")
	user.POST("/register", demo)
	user.GET("/login", demo)
	user.GET("/get_customer_info", demo)
	user.GET("/get_user_info", demo)
	user.POST("/update_perfect_info", demo)
	user.GET("/check_mobile", demo)
	user.GET("/get_member_info", demo)

	//Coupon
	co := v1.Group("/coupon")
	co.GET("/get_coupon_list", demo)

	//Integral
	in := v1.Group("/integral")
	in.GET("/get_current_integral", demo)
	in.GET("/get_integral_history", demo)
	in.GET("/get_vip_grade", demo)
	in.GET("/update_integral_exchange", demo)

	// Common
	c := v1.Group("/common")
	s := c.Group("/sms")
	s.GET("/get_code", demo)
	s.GET("/check_sms", demo)
	s.GET("/ative", demo)

	p := c.Group("/captcha")
	p.GET("/key", common.APIGetCaptchaKey)
	p.GET("/image", demo)
	p.GET("/success", common.ApiCheckCaptcha)

	a := c.Group("/auth")
	a.GET("/set_auth", demo)
}
