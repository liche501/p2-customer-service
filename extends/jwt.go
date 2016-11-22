package extends

import (
	"best/p2-customer-service/logs"

	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
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
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Issuer:    "liche",
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jsonWebToken, err := jwtToken.SignedString(privKey)
	if err != nil {
		return "", err
	}
	return jsonWebToken, err
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		jwtValue := c.Request().Header.Get("Authorization")

		if len(jwtValue) == 0 || jwtValue == "null" {
			logs.Warning.Println("jwt_token is empty")
			return echo.NewHTTPError(http.StatusForbidden, "20001")
		}
		var defaultKeyFunc jwt.Keyfunc = func(*jwt.Token) (interface{}, error) {
			return privKey, nil
		}
		jsonWebTokenParsed, err := jwt.ParseWithClaims(jwtValue, &AuthClaims{}, defaultKeyFunc)
		logs.Debug.Println("jsonWebTokenParsed.Valid==", jsonWebTokenParsed.Valid)
		if err != nil || !jsonWebTokenParsed.Valid {
			return echo.NewHTTPError(http.StatusForbidden, "20001")

		}
		ac := jsonWebTokenParsed.Claims.(*AuthClaims)
		brandCode := ac.BrandCode
		openId := ac.OpenId
		mobile := ac.Mobile
		custNo := ac.CustNo
		logs.Debug.Printf("brandCode:%v openId:%v mobile:%v", brandCode, openId, mobile)
		c.Set("brandCodeWithToken", brandCode)
		c.Set("openIdWithToken", openId)
		c.Set("mobileWithToken", mobile)
		c.Set("custNoWithToken", custNo)

		c.Set("user", ac)
		// c.Get("user").(*AuthClaims).OpenId

		return next(c)
	}
}
