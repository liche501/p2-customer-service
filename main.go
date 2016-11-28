package main

import (
	"best/p2-customer-service/config"
	. "best/p2-customer-service/dto"
	"best/p2-customer-service/extends"
	"best/p2-customer-service/model"

	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	ERROR_MESSAGE_NEED_CONVERT = 5
	e                          = echo.New()
)

func main() {
	// Middleware
	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: "method=${method}, uri=${uri}, status=${status}, latency=${latency}\n",
	// }))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(middleware.CSRF())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding},
	}))
	e.Use(middleware.Secure())
	e.Use(extends.JWTMiddleware())
	e.Use(extends.JWTMiddlewareDataFormat)
	e.HTTPErrorHandler = JSONHTTPErrorHandler

	RouterDemoInit()
	RouterInit()
	config.InitConfig()
	model.InitDB("mysql", config.Config.DB.Conn)
	extends.InitErrorList()

	// e.Logger.Fatal(e.Start(":9000"))
	e.Start(":9000")
}

func demo(c echo.Context) error {
	return c.String(http.StatusOK, "demo func")
}

func JSONHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := "Internal Server Error"
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	}
	if !c.Response().Committed {
		if len(msg) == ERROR_MESSAGE_NEED_CONVERT {
			c.JSON(code, APIResult{Error: APIError{
				Code:    code,
				Message: extends.ErrorList[msg],
			}})
		} else {
			c.JSON(code, map[string]interface{}{
				"statusCode": code,
				"message":    msg,
			})
		}

	}
}
