package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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
func CreateToken(c echo.Context) error {

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "liche"
	claims["admin"] = true
	claims["cust_no"] = "1001"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString(privKey)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})

	return echo.ErrUnauthorized
}
