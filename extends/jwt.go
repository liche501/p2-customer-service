package extends

import (
	"best/p2-customer-service/logs"
	"strings"

	"io/ioutil"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var privKey []byte

func init() {
	key, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		log.Fatal(err)
	}
	privKey = key
}

type AuthClaims struct {
	BrandCode string `json:"brandCode"`
	OpenId    string `json:"openId"`
	Mobile    string `json:"mobile"`
	CustNo    string `json:"custNo"`

	jwt.StandardClaims
}

// AuthHandler create the Claims
func AuthHandler(brandCode, openId, mobile, custNo string) (string, error) {
	claims := AuthClaims{
		brandCode,
		openId,
		mobile,
		custNo,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1172).Unix(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jsonWebToken, err := jwtToken.SignedString([]byte("secret"))
	if err != nil {
		logs.Error.Println("create jwtToken: ", err)
		return "", err
	}
	return jsonWebToken, nil
}

func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("secret"),
		Claims:     &AuthClaims{},
		Skipper: func(c echo.Context) bool {
			switch {
			case strings.HasPrefix(c.Path(), "/ping"):
				return true
			case strings.HasPrefix(c.Path(), "/skip"):
				return true
			case strings.HasPrefix(c.Path(), "/api/va/fashion/user/check_mobile"):
				return true
			case strings.HasPrefix(c.Path(), "/api/v1/common"):
				return true
			}
			return false
		},
	})
}

func JWTMiddlewareDataFormat(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// logs.Warning.Println(c.Get("user"))
		if c.Get("user") != nil {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*AuthClaims)
			c.Set("user", claims)
		}
		return next(c)
	}
}
