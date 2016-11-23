package common

import (
	"best/p2-customer-service/config"
	. "best/p2-customer-service/dto"
	"best/p2-customer-service/logs"
	"best/p2-customer-service/model"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/smallnest/goreq"
)

const (
	SMS_VERICATION_CODE_LENGTH               = 6
	SMS_VERICATION_CODE_GENERATE_RETRY_COUNT = 5
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
	brandCode := c.QueryParam("brandCode")
	mobile := c.QueryParam("mobile")
	if mobile == "" {
		return echo.NewHTTPError(http.StatusBadRequest, 10012)
	}
	verCode := createVerCode()
	smsStr := &SmsResult{}
	url := config.Config.Adapter.CSL.CustomerInterfaceAPI + "?telephoneNumber=" + mobile + "&verCode=" + verCode
	_, _, err := goreq.New().Get(url).BindBody(smsStr).SetCurlCommand(true).End()
	if err != nil {
		logs.Error.Println("Call sms api error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10003)
	}

	if smsStr.Data.Flag != 0 {
		logs.Error.Println("Call sms api error: ", smsStr.Data.Flag)
		return echo.NewHTTPError(http.StatusInternalServerError, 10003)
	}

	msg := fmt.Sprintf("SMS %v %v %v ", brandCode, mobile, smsStr.Data.SmsCode)
	logs.Succ.Println(msg)

	var sms model.Sms
	sms.BrandCode = brandCode
	sms.Mobile = mobile
	sms.VerCode = smsStr.Data.SmsCode
	sms.Type = "register"
	sms.Create()
	return c.JSON(http.StatusOK, APIResult{Success: true, Result: map[string]string{"status": "OK"}})
}

func ApiCheckSms(c echo.Context) error {
	// Verify SmsCode
	var sms model.Sms
	mobile := c.QueryParam("mobile")
	verCode := c.QueryParam("verCode")
	if mobile == "" || verCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, 10012)
	}
	result, err := sms.CheckVerCode(mobile, verCode)
	if err != nil {
		logs.Error.Println("Check sms error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}
	return c.JSON(http.StatusOK, APIResult{Success: true, Result: map[string]bool{"flag": result}})
}

type SmsVerificationCode string

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
	return nil
}
func createVerCode() string {
	verCode := ""
	for i := 0; i < SMS_VERICATION_CODE_GENERATE_RETRY_COUNT; i++ {
		verCode = cutRandomCode()
		result := checkRepeatSms(verCode)
		if !result {
			break
		}
		logs.Debug.Println(verCode, " is repeated")
	}

	return verCode
}
func cutRandomCode() string {
	rd := RandInt64()
	str := strconv.FormatInt(rd, 10)
	strCut := str[0:SMS_VERICATION_CODE_LENGTH]
	return strCut
}

func checkRepeatSms(verCode string) bool {
	var sms model.Sms
	result, err := sms.CheckRepeatVerCode(verCode)
	if err != nil {
		logs.Error.Println("CheckRepeatVerCode sms error: ", err)
		return false
	}
	return result
}
