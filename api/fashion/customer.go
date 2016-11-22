package fashion

import (
	"best/p2-customer-service/config"
	. "best/p2-customer-service/dto"
	"best/p2-customer-service/extends"
	"best/p2-customer-service/logs"
	"best/p2-customer-service/model"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"net/http"

	"github.com/labstack/echo"
	"github.com/smallnest/goreq"
)

// APICheckMobile: 注册/修改手机号时检查注册状态
func APICheckMobile(c echo.Context) error {
	mobile := c.Get("user").(*extends.AuthClaims).Mobile
	openId := c.Get("user").(*extends.AuthClaims).OpenId
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode

	logs.Debug.Println(mobile, openId, brandCode)
	if mobile == "" || openId == "" || brandCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, 10012)
	}
	// status, err := models.UserShop{}.CheckMobileForRegister(mobile, openId, brandCode)
	// if err != nil {
	// 	logs.Error.Println("CheckMobile error: ", err)

	// 	return echo.NewHTTPError(http.StatusInternalServerError, 10013)

	// }
	status := "heihei"
	return c.JSON(http.StatusOK, APIResult{Success: true, Result: map[string]string{"status": status}})

}

// curl -X POST -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJicmFuZENvZGUiOiJyYyIsIm9wZW5JZCI6Im9ZaVI2d1R6NmFucjVLcGlSSC1tUmNwdnZMUGMiLCJtb2JpbGUiOiIxMzY5MTE5NDIyMyIsImN1c3RObyI6IjAwMDE4NTIzNTkiLCJleHAiOjE0ODAwNjQ2NDIsImlzcyI6ImxpY2hlIn0.huxuvLITetwHzdpZHX-T_sfZe0rEeMM_2DOnugdUjRo" -H "Cache-Control: no-cache" -H "Postman-Token: 28bf6f4d-9809-26a2-3229-4a177c8d29cf" "http://localhost:9000/api/v1/fashion/user/register?mobile=13691194223&verCode=1234"
func APIRegister(c echo.Context) error {

	mobile := c.FormValue("mobile")
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	openId := c.Get("user").(*extends.AuthClaims).OpenId
	verCode := c.FormValue("verCode")

	if verCode == "" || mobile == "" || openId == "" || brandCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, 10012)
	}

	//CheckSms
	// flag, err := model.Sms{}.CheckVerCode(mobile, verCode)
	// if err != nil {
	// 	logs.Error.Println("Check sms error: ", err)
	// 	return echo.NewHTTPError(http.StatusInternalServerError, 10013)

	// } else if !flag {
	// 	logs.Error.Println("Sms code invalid")
	// 	return c.JSON(http.StatusOK, APIResult{Error: APIError{Code: 20006, Message: "20006"}})

	// }
	logs.Debug.Println(brandCode, openId)

	fashionBrandCustomerInfo := new(model.FashionBrandCustomerInfo)
	fashionBrandCustomerInfo.Customer.Mobile = mobile
	fashionBrandCustomerInfo.FashionBrandCustomer.BrandCode = brandCode
	fashionBrandCustomerInfo.FashionBrandCustomer.WxOpenID = openId
	logs.Debug.Println(fashionBrandCustomerInfo)

	//customer regist
	if err := fashionBrandCustomerInfo.Create(); err != nil {
		logs.Error.Println(err)
	}

	return c.JSON(http.StatusOK, APIResult{Success: true})

}

func APILogin(c echo.Context) error {
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	openId := c.Get("user").(*extends.AuthClaims).OpenId
	if openId == "" || brandCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, 10012)
	}
	us, err := model.FashionBrandCustomerInfo{}.GetByWxOpenIDAndStatus(brandCode, openId, "Success")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}
	if us == nil {

		return c.JSON(http.StatusOK, APIResult{Error: APIError{Code: 20003, Message: ""}})
	}

	jsonWebToken, err := extends.AuthHandler(brandCode, openId, us.Customer.Mobile, us.FashionBrandCustomer.CustNo)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, 20002)
	}
	logs.Warning.Println("Create JsonWebToken: ", jsonWebToken)
	rs := APIResult{Success: true, Result: map[string]string{"token": jsonWebToken, "status": us.FashionBrandCustomer.Status}}
	return c.JSON(http.StatusOK, rs)

}

// APIGetCustomerInfo SSE 获取顾客信息
func APIGetCustomerInfo(c echo.Context) error {
	openId := c.Get("user").(*extends.AuthClaims).OpenId
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	if openId == "" || brandCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, 10012)
	}

	// ===============SSE START===============

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")

	ui, err := model.FashionBrandCustomer{}.GetByWxOpenID(brandCode, openId)
	if err != nil {
		logs.Error.Println("GetByOpenidWithoutResgitStatus error: ", err)
		return c.String(http.StatusOK, "data:nouser\n\nexist")
	}

	if ui == nil {
		return c.String(http.StatusOK, "data:nouser\n\n")
	}

	logs.Debug.Printf("RegistStatus for openId %v: %v", openId, ui.Status)

	if ui.Status == "Error01" {
		return c.String(http.StatusOK, "data:nouser\n\nError01")

	} else if ui.Status == "OtherError" {
		return c.String(http.StatusOK, "data:nouser\n\nOtherError")

	} else if len(ui.CustNo) > 0 && ui.Status == "Success" {
		return c.String(http.StatusOK, "data:nouser\n\nexist")
	}
	return c.String(http.StatusOK, "data:nouser\n\n")

	// ===============SSE END===============
}

func APIGetUserInfo(c echo.Context) error {

	openId := c.Get("user").(*extends.AuthClaims).OpenId
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode

	ui, err := model.FashionBrandCustomerInfo{}.GetByWxOpenIDAndStatus(brandCode, openId, "Success")
	if err != nil {
		logs.Error.Println("GetByOpenid error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}

	if ui == nil {
		logs.Error.Println("UserShop not exist")
		return c.JSON(http.StatusOK, APIResult{Error: APIError{Code: 20003}})
	}

	// If need perfect user info
	ud, err := model.CustomerInfo{}.Get(ui.FashionBrandCustomer.BrandCode, ui.Customer.Mobile)
	if err != nil {
		logs.Error.Println("GetUserDetail error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10003)
	}
	needPerfect := ud == nil || !ud.HasFilled

	custNo := ui.FashionBrandCustomer.CustNo
	mobile := ui.Customer.Mobile

	totalCode := ""
	if len(custNo) > 0 {
		totalCode = strings.ToUpper(brandCode) + custNo
	}

	displayIntegral := setMyScoreDisplay(brandCode)

	rs := map[string]interface{}{"custNo": totalCode, "originCustNo": custNo, "mobile": mobile, "displayIntegral": displayIntegral, "needPerfect": needPerfect}
	return c.JSON(http.StatusOK, APIResult{Success: true, Result: rs})
}

func APIGetMemberInfo(c echo.Context) error {
	mobile := c.Get("user").(*extends.AuthClaims).Mobile
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	openId := c.Get("user").(*extends.AuthClaims).OpenId

	ud, err := model.CustomerInfo{}.Get(brandCode, mobile)
	if err != nil {
		logs.Error.Println("GetUserDetail error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}
	// 完善顾客信息上线前注册的用户
	if ud == nil {
		ud = &model.CustomerInfo{}
		ud.HasFilled = false
	}
	ud = getUserDetailFromCSL(brandCode, openId, ud, ud.HasFilled)

	result := make(map[string]interface{})
	if ud != nil {
		result["mobile"] = mobile
		result["name"] = ud.Name
		result["gender"] = ud.Gender
		result["birthday"] = ud.Birthday
		result["address"] = ud.Address
		result["detailAddress"] = ud.DetailAddress
		result["email"] = ud.Email
		result["isMarried"] = ud.IsMarried
		result["hasFilled"] = ud.HasFilled
	}

	return c.JSON(http.StatusOK, APIResult{Success: true, Result: result})

}

func APIUpdatePerfectInfo(c echo.Context) error {
	var userDetail model.CustomerInfo
	var oldmobile string
	mobile := c.Get("user").(*extends.AuthClaims).Mobile
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	openId := c.Get("user").(*extends.AuthClaims).OpenId

	userDetail.BrandCode = brandCode
	userDetail.Name = c.FormValue("name")
	userDetail.Gender = c.FormValue("gender")
	phoneNo := c.FormValue("mobile")
	//logs.Error.Println(userDetail.Name)
	//logs.Error.Println(phoneNo)
	if phoneNo != mobile {
		userDetail.Mobile = phoneNo
		mobile = phoneNo
		oldmobile = mobile
	} else {
		userDetail.Mobile = mobile
	}

	birthday := c.FormValue("birthday")
	if birthday != "" {
		userDetail.Birthday = birthday
	}

	userDetail.Address = c.FormValue("address")
	userDetail.DetailAddress = c.FormValue("detailAddress")
	userDetail.Email = c.FormValue("email")

	isMarried := c.FormValue("isMarried")
	var err error
	if isMarried != "" {
		userDetail.IsMarried, err = strconv.ParseBool(isMarried)
		if err != nil {
			logs.Error.Println("[UpdatePerfectInfo]isMarried convert bool error", err)
			return echo.NewHTTPError(http.StatusBadRequest, 10011)
		}
	}
	if phoneNo != mobile {
		err = userDetail.ChangeMobileWithOld(oldmobile, mobile)
		if err != nil {
			logs.Error.Println("[UpdatePerfectInfo]saveUserDetail", err)
			return echo.NewHTTPError(http.StatusInternalServerError, 10013)
		}
	} else {
		err = userDetail.Save()
		if err != nil {
			logs.Error.Println("[UpdatePerfectInfo]saveUserDetail", err)
			return echo.NewHTTPError(http.StatusInternalServerError, 10013)
		}
	}

	userDetail.HasFilled = true
	err = userDetail.UpdateHasFilled()

	// err = saveUserDetailToCSL(&userDetail, openId)
	// if err != nil {
	// 	logs.Error.Println("[SaveUserDetailToCSL]saveUserDetail", err)
	// 	return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	// }

	return c.JSON(http.StatusOK, APIResult{Success: true, Result: map[string]string{"openId": openId}})

}

// 从CSL获取用户信息
func getUserDetailFromCSL(brandCode, openId string, ud *model.CustomerInfo, splitAddress bool) *model.CustomerInfo {
	url := fmt.Sprintf("%v?BrandCode=%v&OpenId=%v", config.Config.Adapter.CSL.CustomerInterfaceAPI, brandCode, openId)
	logs.Debug.Println("Start GetUserInfoDetail:", url)

	type UserDetailCSL struct {
		CustNo        string `json:"custNo"`
		TelephoneNo   string `json:"telephoneNo"`
		CustName      string `json:"custName"`
		SexCode       string `json:"sexCode"`
		Birthday      string `json:"birthday"`
		Address       string `json:"address"`
		CustGradeName string `json:"custGradeName"`
		WeddingChk    bool   `json:"weddingChk"`
		BrandCode     string `json:"brandCode"`
		EmailAddress  string `json:"emailAddress"`
		LeaveChk      bool   `json:"leaveChk"`
		IsNewCust     int64  `json:"isNewCust"`
	}
	var data UserDetailCSL
	resp, _, reqErr := goreq.New().Get(url).ContentType("json").BindBody(&data).SetCurlCommand(true).End()
	if reqErr != nil {
		logs.Error.Println("ReqErr:", reqErr, " StatusCode:", resp.StatusCode)
		return nil
	}
	logs.Debug.Println("Success GetUserInfoDetail: ", data)

	n := strings.Index(data.Address, ",")
	if splitAddress && n >= 0 {
		ud.Address = string([]byte(data.Address)[0:n])
		ud.DetailAddress = string([]byte(data.Address)[n+1:])
	} else {
		ud.DetailAddress = data.Address
	}
	ud.Name = data.CustName
	ud.Mobile = data.TelephoneNo
	ud.Birthday = data.Birthday
	ud.BrandCode = data.BrandCode
	ud.Gender = data.SexCode
	ud.Email = data.EmailAddress
	ud.IsMarried = data.WeddingChk
	ud.IsNewCust = data.IsNewCust
	return ud
}

func saveUserDetailToCSL(ud *model.CustomerInfo, openId string) error {
	type UserDetailCSL struct {
		Address      string `json:"Address"`
		Birthday     string `json:"Birthday"`
		BrandCode    string `json:"BrandCode"`
		CustName     string `json:"CustName"`
		EmailAddress string `json:"EmailAddress"`
		SexCode      string `json:"SexCode"`
		TelephoneNo  string `json:"TelephoneNo"`
		Marry        string `json:"Marry"`
		OpenId       string `json:"OpenId"`
		Intergal     int    `json:"Intergal"`
	}
	var data UserDetailCSL

	url := config.Config.Adapter.CSL.CustomerInterfaceAPI

	data.Address = ud.Address + "," + ud.DetailAddress
	data.Birthday = ud.Birthday
	data.BrandCode = strings.ToUpper(ud.BrandCode)
	data.CustName = ud.Name
	data.EmailAddress = ud.Email
	data.SexCode = ud.Gender
	data.TelephoneNo = ud.Mobile
	data.OpenId = openId
	if ud.IsMarried {
		data.Marry = "1"
	} else {
		data.Marry = "0"
	}
	data.Intergal = 0

	resp, _, reqErr := goreq.New().Put(url).SendStruct(data).ContentType("json").SetCurlCommand(true).End()
	if reqErr != nil || resp.StatusCode != 204 {
		logs.Error.Println("ReqErr:", reqErr, " StatusCode:", resp.StatusCode)
		return errors.New("Save userinfo failed!")
	}

	return nil
}
func setMyScoreDisplay(brandCode string) bool {
	if brandCode == "sf" || brandCode == "sm" || brandCode == "cz" {
		return false
	}
	return true
}
