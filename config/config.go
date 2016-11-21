package config

import (
	"fmt"

	"github.com/jinzhu/configor"
)

var Config = struct {
	DB struct {
		Conn string
	}

	Mongo struct {
		Conn string
	}

	Contents struct {
		ServiceHost     string
		Port            string
		CouponApiPrefix string
	}

	MhService struct {
		ServiceHost     string
		Port            string
		MemberApiPrefix string
	}

	Adapter struct {
		CslWebService        string
		CustomerInterfaceApi string
	}

	UserInfoQueue struct {
		ServiceHost string
		Port        string
		ApiPrefix   string
	}
	TaskQueue struct {
		Broker string
	}

	CaptchaService struct {
		ServiceHost string
		Port        string
		ApiPrefix   string
	}
	FrontDomain struct {
		Fashion string
		RC      string
	}

	Marketing struct {
		EventCoupon string
	}
}{}

var UrlCaptcha string
var UrlMhService string

func InitConfig() {
	configor.Load(&Config, "config/config.json")

	UrlMhService = fmt.Sprintf("http://%s:%s/%s", Config.MhService.ServiceHost, Config.MhService.Port, Config.MhService.MemberApiPrefix)
	UrlCaptcha = fmt.Sprintf("http://%s:%s/%s", Config.CaptchaService.ServiceHost, Config.CaptchaService.Port, Config.CaptchaService.ApiPrefix)
}
