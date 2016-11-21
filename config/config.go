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

	EE struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	TT struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	FO struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	PC struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	RC struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	EK struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	PO struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	SF struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	SM struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	PR struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	CZ struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	MH struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	EH struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	EC struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	BC struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	WA struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	NC struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	CO struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
	MD struct {
		AppId       string
		AppSecret   string
		RedirectUrl string
	}
}{}

var UrlCaptcha string
var UrlMhService string

func InitConfig() {
	configor.Load(&Config, "config/config.json")

	UrlMhService = fmt.Sprintf("http://%s:%s/%s", Config.MhService.ServiceHost, Config.MhService.Port, Config.MhService.MemberApiPrefix)
	UrlCaptcha = fmt.Sprintf("http://%s:%s/%s", Config.CaptchaService.ServiceHost, Config.CaptchaService.Port, Config.CaptchaService.ApiPrefix)
}
