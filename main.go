package main

import (
	"best/p2-customer-service/config"
	"best/p2-customer-service/model"

	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	e = echo.New()
)

func main() {
	// Middleware
	// e.Use(middleware.Logger())
	// e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
	// 	SigningKey: privKey,
	// 	Claims:     extends.AuthClaims{},
	// 	Skipper: func(c echo.Context) bool {
	// 		switch {
	// 		case strings.HasPrefix(c.Path(), "/createtoken"):
	// 			return true
	// 		case strings.HasPrefix(c.Path(), "/api/v1/common"):
	// 			return true
	// 		}
	// 		return false
	// 	},
	// }))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, latency=${latency}\n",
	}))
	e.Use(middleware.Recover())
	// e.Use(middleware.CSRF())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding},
	}))
	e.Use(middleware.Secure())
	e.Use(MyMwServerHeader)
	e.HTTPErrorHandler = JSONHTTPErrorHandler

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract the credentials from HTTP request header and perform a security
			// check

			// For invalid credentials
			// return echo.NewHTTPError(http.StatusUnauthorized)

			// For valid credentials call next
			return next(c)
		}
	})
	RouterInit()
	config.InitConfig()
	model.InitDB("mysql", config.Config.DB.Conn)

	// e.Logger.Fatal(e.Start(":9000"))
	e.Start(":9000")
}

func demo(c echo.Context) error {

	fmt.Println("deme 2222")
	// time.Sleep(time.Second * 1)
	fmt.Println(c.Request().Host)
	return c.String(http.StatusOK, "test")

}

// ServerHeader middleware adds a `Server` header to the response.
func MyMwServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// fmt.Println("22222")
		c.Response().Header().Set(echo.HeaderServer, "Echo/3.0")
		return next(c)
	}
}

func JSONHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := "Internal Server Error"
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	}
	if !c.Response().Committed {
		c.JSON(code, map[string]interface{}{
			"statusCode": code,
			"message":    msg,
		})
	}
}
