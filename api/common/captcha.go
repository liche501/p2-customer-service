package common

import (
	"best/p2-customer-service/config"
	. "best/p2-customer-service/dto"

	"best/p2-customer-service/logs"

	"net/http"

	"github.com/labstack/echo"
	"github.com/smallnest/goreq"
)

type CaptchaResult struct {
	ErrorNo  int64  `json:"error_no"`
	ErrorMsg string `json:"error_msg"`
	Key      string `json:"key"`
}

func APIGetCaptchaKey(c echo.Context) error {

	serviceURL := config.UrlCaptcha + "/get_key"
	// serviceURL := "http://139.196.228.246:9094/captcha/get_key"
	logs.Debug.Println(serviceURL)
	result := &CaptchaResult{}
	_, _, err := goreq.New().Get(serviceURL).BindBody(result).SetCurlCommand(true).End()
	if err != nil {
		logs.Error.Println("Get captcha key error: ", err)
		// extends.ReturnJsonFailure(w, http.StatusInternalServerError, 10003, err[0].Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err[0].Error())
	}
	return c.JSON(http.StatusOK, APIResult{Success: true, Result: result})
	// extends.ReturnJsonSuccess(w, http.StatusOK, map[string]string{"key": result.Key})
}
