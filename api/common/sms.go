package common

import (
	"best/p2-customer-service/config"
	. "best/p2-customer-service/dto"
	"best/p2-customer-service/logs"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"wxshop/extends"
	"wxshop/models"

	"github.com/labstack/echo"
	"github.com/smallnest/goreq"
)

type SmsResult struct {
	Result string `json:"result"`
	Msg    string `json:"msg"`
	Data   Data   `json:"data"`
}

type Data struct {
	TelNo   string `json:"telNo"`
	SmsCode string `json:"smsCode"`
	Flag    int    `json:"flag"`
}

func ApiSendSms(c echo.Context) error {
	brandCode := ""
	if c.FormValue("brandCode") == "mh" {
		brandCode = "mh"
	}
	//  else {
	// 	brandCode = strings.ToLower(mux.Vars(c)["tenantCode"])
	// }

	mobile := c.FormValue("mobile")
	if mobile == "" {
		// extends.ReturnJsonFailure(w, http.StatusBadRequest, 10012)
		// return
		return echo.NewHTTPError(http.StatusBadRequest, 10012)
	}
	verCode := createVerCode()
	smsStr := &SmsResult{}
	url := config.Config.Adapter.CustomerInterfaceApi + "?telephoneNumber=" + mobile + "&verCode=" + verCode
	_, _, err := goreq.New().Get(url).BindBody(smsStr).SetCurlCommand(true).End()
	if err != nil {
		logs.Error.Println("Call sms api error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10003)

		// extends.ReturnJsonFailure(w, http.StatusInternalServerError, 10003, err[0].Error())
		// return
	} else if smsStr.Data.Flag != 0 {
		logs.Error.Println("Call sms api error: ", smsStr.Data.Flag)
		// extends.ReturnJsonFailure(w, http.StatusInternalServerError, 10003, fmt.Sprintf("Flag == %v", smsStr.Data.Flag))
		// return
		return echo.NewHTTPError(http.StatusInternalServerError, 10003)
	}
	msg := fmt.Sprintf("SMS %v %v %v ", brandCode, mobile, smsStr.Data.SmsCode)
	logs.Succ.Println(msg)
	var sms models.Sms
	sms.BrandCode = brandCode
	sms.Mobile = mobile
	sms.VerCode = smsStr.Data.SmsCode
	sms.Type = "register"
	sms.Create()
	return c.JSON(http.StatusOK, APIResult{Success: true, Result: map[string]string{"status": "OK"}})
	//extends.ReturnJsonSuccess(w, http.StatusOK, map[string]string{"status": "OK"})
}

func ApiCheckSms(c echo.Context) error {
	// Verify SmsCode
	var sms models.Sms
	mobile := r.FormValue("mobile")
	verCode := r.FormValue("verCode")
	if mobile == "" || verCode == "" {
		extends.ReturnJsonFailure(w, http.StatusBadRequest, 10012)
		return
	}
	result, err := sms.CheckVerCode(mobile, verCode)
	if err != nil {
		logs.Error.Println("Check sms error: ", err)
		extends.ReturnJsonFailure(w, http.StatusInternalServerError, 10013)
		return
	}
	extends.ReturnJsonSuccess(w, http.StatusOK, map[string]bool{"flag": result})
}

var smsChannel chan int64 = make(chan int64, 32)

func init() {
	go func() {
		var old int64
		for {
			o := rand.New(rand.NewSource(time.Now().UnixNano())).Int63()
			if old != o {
				old = o
				select {
				case smsChannel <- o:
				}
			}
		}
	}()
}
func RandInt64() (r int64) {
	select {
	case rand := <-smsChannel:
		r = rand
	}
	return
}

func Active(c echo.Context) error {

	aa := createVerCode()
	logs.Succ.Println(aa)
}
func createVerCode() string {
	verCode := ""
	for i := 0; i < 5; i++ {
		verCode = cutRandomCode()
		result := checkRepeatSms(verCode)
		if !result {
			// logs.Debug.Println("ok ", verCode)
			break
		}
		logs.Debug.Println(verCode, " is repeated")
	}

	return verCode
}
func cutRandomCode() string {
	rd := RandInt64()
	str := strconv.FormatInt(rd, 10)
	strCut := str[0:6]
	return strCut
}

func checkRepeatSms(verCode string) bool {
	var sms models.Sms
	result, err := sms.CheckRepeatVerCode(verCode)
	if err != nil {
		logs.Error.Println("CheckRepeatVerCode sms error: ", err)
		return false
	}
	return result
}
