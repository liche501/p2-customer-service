package main

import (
	. "best/p2-customer-service/dto"
	"best/p2-customer-service/event"
	"best/p2-customer-service/extends"
	"best/p2-customer-service/logs"
	"net/http"

	"github.com/labstack/echo"
)

func RouterDemoInit() {
	// Unauthenticated route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "accessible")
	})
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	skip := e.Group("/skip")
	skip.GET("/success", func(c echo.Context) error {
		return c.JSON(http.StatusOK, APIResult{Success: true, Result: "aaaaaaa"})
	})
	skip.GET("/error", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusInternalServerError, 10003)
	})
	skip.GET("/error2", func(c echo.Context) error {
		return c.JSON(http.StatusOK, APIResult{Error: APIError{Code: 400, Message: "StatusBadRequest"}})
	})
	skip.GET("/event/:eventType", func(c echo.Context) error {
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
	skip.GET("/token", func(c echo.Context) error {
		logs.Debug.Println("token start")

		token, err := extends.AuthHandler("rc", "oYiR6wTz6anr5KpiRH-mRcpvvLPc", "13691194223", "0001852359")
		if err != nil {
			return err
		}
		logs.Debug.Println(token)
		return c.JSON(http.StatusOK, APIResult{Success: true, Result: token})
	})
	t := e.Group("/jwt")
	// t.Use(extends.JWTMiddleware)
	t.GET("", func(c echo.Context) error {
		openId := c.Get("user").(*extends.AuthClaims).OpenId
		return c.String(http.StatusOK, "Welcome "+openId)
		// return c.String(http.StatusOK, "jwt func ")
	})
}
