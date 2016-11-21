package common

import (
	"best/wxshop/logs"
	"net/http"

	"github.com/labstack/echo"
	"github.com/smallnest/goreq"
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

type CaptchaResult struct {
	ErrorNo  int64  `json:"error_no"`
	ErrorMsg string `json:"error_msg"`
	Key      string `json:"key"`
}

func ApiGetCaptchaKey(c echo.Context) error {
	// serviceUrl := config.UrlCaptcha + "/get_key"
	serviceUrl := "http://139.196.228.246:9094/captcha/get_key"
	logs.Debug.Println(serviceUrl)
	result := &CaptchaResult{}
	_, body, err := goreq.New().Get(serviceUrl).BindBody(result).SetCurlCommand(true).End()
	if err != nil {
		logs.Error.Println("Get captcha key error: ", err)
		// extends.ReturnJsonFailure(w, http.StatusInternalServerError, 10003, err[0].Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err[0].Error())

	}
	logs.Debug.Println(body)
	return c.JSON(http.StatusOK, apiResult{Success: true, Result: result})
	// extends.ReturnJsonSuccess(w, http.StatusOK, map[string]string{"key": result.Key})
}
