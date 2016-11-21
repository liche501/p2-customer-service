// package main

// import (
// 	"net/http"
// 	"os"

// 	"github.com/Sirupsen/logrus"
// 	"github.com/labstack/echo"
// )

// var (
// 	log = logrus.New()
// )

// func init() {
// 	log.Out = os.Stderr
// 	log.Level = logrus.DebugLevel

// }

// type User struct {
// 	Name string `json:"name" xml:"name" form:"name"`
// 	Age  string `json:"age" xml:"age" form:"age"`
// }

// func main() {
// 	e := echo.New()
// 	e.POST("/user/:id", func(c echo.Context) error {
// 		// name := c.QueryParam("name")
// 		// no := c.Param("id")
// 		// age := c.FormValue("age")
// 		// log.Debugln(age)

// 		u := new(User)
// 		if err := c.Bind(u); err != nil {
// 			return err
// 		}
// 		return c.JSON(http.StatusCreated, u)
// 	})
// 	e.GET("/ping", func(c echo.Context) error {
// 		return c.JSON(http.StatusOK, "pong")

// 	})
// 	e.Logger.Fatal(e.Start(":9000"))
// }

package main

import (
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
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, latency=${latency}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CSRF())
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
