package config

import (
	"fmt"

	"github.com/jinzhu/configor"
)

var Config = struct {
	DB struct {
		Conn string `json:"Conn"`
	} `json:"DB"`
	Coupon struct {
		Contents struct {
			ServiceHost     string `json:"ServiceHost"`
			Port            string `json:"Port"`
			CouponAPIPrefix string `json:"CouponApiPrefix"`
		} `json:"Contents"`
		Marketing struct {
			EventCoupon string `json:"EventCoupon"`
		} `json:"Marketing"`
	} `json:"Coupon"`
	Adapter struct {
		CSL struct {
			CslWebService        string `json:"CslWebService"`
			CustomerInterfaceAPI string `json:"CustomerInterfaceApi"`
		} `json:"CSL"`
		MH struct {
			ServiceHost     string `json:"ServiceHost"`
			Port            string `json:"Port"`
			MemberAPIPrefix string `json:"MemberApiPrefix"`
		} `json:"MH"`
	} `json:"Adapter"`
	CaptchaService struct {
		ServiceHost string `json:"ServiceHost"`
		Port        string `json:"Port"`
		APIPrefix   string `json:"ApiPrefix"`
	} `json:"CaptchaService"`
}{}

var UrlCaptcha string
var UrlMhService string

func InitConfig() {
	configor.Load(&Config, "config/config.json")

	UrlMhService = fmt.Sprintf("http://%s:%s/%s", Config.Adapter.MH.ServiceHost, Config.Adapter.MH.Port, Config.Adapter.MH.MemberAPIPrefix)
	UrlCaptcha = fmt.Sprintf("http://%s:%s/%s", Config.CaptchaService.ServiceHost, Config.CaptchaService.Port, Config.CaptchaService.APIPrefix)
}
