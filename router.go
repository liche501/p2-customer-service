package main

import (
	"best/p2-customer-service/api/common"
	"best/p2-customer-service/api/fashion"
	"best/p2-customer-service/event"
)

func RouterInit() {

	api := e.Group("/api")
	v1 := api.Group("/v1")

	// Common
	c := v1.Group("/common")
	c.POST("/events", event.ApiHandleEvent)

	s := c.Group("/sms")
	s.GET("/code", common.ApiSendSms)
	s.GET("/success", common.ApiCheckSms)
	s.GET("/active", common.Active)

	p := c.Group("/captcha")
	p.GET("/key", common.APIGetCaptchaKey)
	p.GET("/image", demo)
	p.GET("/success", common.ApiCheckCaptcha)

	a := c.Group("/auth")
	a.GET("/set_auth", demo)

	fa := v1.Group("/fashion")
	//User
	user := fa.Group("/user")

	user.POST("/register", fashion.APIRegister)
	user.GET("/login", fashion.APILogin)
	user.GET("/brand_customer_status", fashion.APIBrandCustomerStatus)
	user.GET("/get_user_info", fashion.APIGetUserInfo)
	user.POST("/update_perfect_info", fashion.APIUpdatePerfectInfo)
	user.GET("/check_mobile", fashion.APICheckMobileAvailableForRegister)
	user.GET("/get_member_info", fashion.APIGetMemberInfo)

	//Coupon
	co := fa.Group("/coupon")
	co.GET("/get_coupon_list", fashion.APIGetCouponList)

	//Integral
	in := fa.Group("/integral")
	in.GET("/current", fashion.ApiGetCurrentIntegral)
	in.GET("/history", fashion.ApiGetIntegralHistory)
	in.GET("/grade", fashion.ApiGetVipGrade)
	in.GET("/update_integral_exchange", fashion.ApiUpdateIntegralExchange)

}
