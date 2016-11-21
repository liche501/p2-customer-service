package main

import (
	"best/p2-customer-service/api/common"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	apiResult struct {
		Result  interface{} `json:"result"`
		Success bool        `json:"success"`
		Error   apiError    `json:"error"`
	}

	apiError struct {
		Code    int         `json:"code"`
		Details interface{} `json:"details"`
		Message string      `json:"message"`
	}
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
		return c.JSON(http.StatusOK, apiResult{Success: true, Result: "aaaaaaa"})
	})
	e.GET("/error", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, "10003", " error01")
	})
	e.GET("/error2", func(c echo.Context) error {
		return c.JSON(http.StatusOK, apiResult{Error: apiError{Code: http.StatusBadRequest, Message: "StatusBadRequest"}})
	})

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", restricted)

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
	p.GET("/get_key", common.ApiGetCaptchaKey)
	p.GET("/get_image", demo)
	p.GET("/check", demo)

	a := c.Group("/auth")
	a.GET("/set_auth", demo)
}
