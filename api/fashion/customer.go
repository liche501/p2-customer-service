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
func APICheckMobileAvailableForRegister(c echo.Context) error {
	mobile := c.QueryParam("mobile")
	brandCode := c.QueryParam("brandCode")

	logs.Debug.Println(mobile, brandCode)
	if mobile == "" || brandCode == "" {
		return c.JSON(http.StatusBadRequest, APIResult{
			Error: APIError{
				Code:    10012,
				Message: "Mobile, WxOpenID and BrandCode are required parameter.",
			}})
	}

	_, err := model.FashionBrandCustomer{}.GetByMobile(brandCode, mobile)
	if err != nil && err != model.CustomerNotExistError {
		return echo.NewHTTPError(http.StatusInternalServerError, 10012)
	}
	if err == model.CustomerNotExistError {
		return c.JSON(http.StatusOK, APIResult{Success: true})
	}

	return c.JSON(http.StatusOK, APIResult{Success: false})
}

// curl -X POST -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJicmFuZENvZGUiOiJyYyIsIm9wZW5JZCI6Im9ZaVI2d1R6NmFucjVLcGlSSC1tUmNwdnZMUGMiLCJtb2JpbGUiOiIxMzY5MTE5NDIyMyIsImN1c3RObyI6IjAwMDE4NTIzNTkiLCJleHAiOjE0ODAwNjQ2NDIsImlzcyI6ImxpY2hlIn0.huxuvLITetwHzdpZHX-T_sfZe0rEeMM_2DOnugdUjRo" -H "Cache-Control: no-cache" -H "Postman-Token: 28bf6f4d-9809-26a2-3229-4a177c8d29cf" "http://localhost:9000/api/v1/fashion/user/register?mobile=13691194223&verCode=1234"
func APIRegister(c echo.Context) error {

	mobile := c.QueryParam("mobile")
	verCode := c.QueryParam("verCode")
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	openId := c.Get("user").(*extends.AuthClaims).OpenId

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
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	//UpdateStatus BrandCustomerInitiated
	// bc := model.BrandCustomer{}
	// bc.BrandCode = brandCode
	// bc.CustomerId = e.CustomerID
	// bc.Status = "BrandCustomerInitiated"
	// err := bc.UpdateStatus()
	// if err != nil {
	// 	logs.Error.Println(err)
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err)
	// }

	// WillDo:: sendEvent => BrandCustomerConfirmed

	return c.JSON(http.StatusOK, APIResult{Success: true})

}

func APILogin(c echo.Context) error {
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	openId := c.Get("user").(*extends.AuthClaims).OpenId
	if openId == "" || brandCode == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, 10012)
	}

	us, err := model.FashionBrandCustomerInfo{}.GetSuccessUserByWxOpenID(brandCode, openId)
	if err == model.CustomerNotExistError {
		return c.JSON(http.StatusOK, APIResult{Error: APIError{Code: 20003, Message: ""}})
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}

	jsonWebToken, err := extends.AuthHandler(brandCode, openId, us.Customer.Mobile, us.FashionBrandCustomer.CustNo)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, 20002)
	}
	logs.Warning.Println("Create JsonWebToken: ", jsonWebToken)
	rs := APIResult{Success: true, Result: map[string]string{"token": jsonWebToken, "status": us.Status()}}
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

	ui, err := model.FashionBrandCustomerInfo{}.GetByWxOpenID(brandCode, openId)
	if err != nil {
		logs.Error.Println("GetByOpenidWithoutResgitStatus error: ", err)
		return c.String(http.StatusOK, "data:nouser\n\nexist")
	}

	if ui == nil {
		return c.String(http.StatusOK, "data:nouser\n\n")
	}

	logs.Debug.Printf("RegistStatus for openId %v: %v", openId, ui.Status)

	if ui.Status() == "BrandCustomerDuplicated" {
		return c.String(http.StatusOK, "data:nouser\n\nError01")
	}
	if ui.Status() == "BrandCustomerFailed" {
		return c.String(http.StatusOK, "data:nouser\n\nOtherError")

	}
	if len(ui.CustNo) > 0 && ui.Status() == "BrandCustomerCreated" {
		return c.String(http.StatusOK, "data:nouser\n\nexist")
	}
	return c.String(http.StatusOK, "data:nouser\n\n")

	// ===============SSE END===============
}

func APIGetUserInfo(c echo.Context) error {

	openId := c.Get("user").(*extends.AuthClaims).OpenId
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode

	ui, err := model.FashionBrandCustomerInfo{}.GetSuccessUserByWxOpenID(brandCode, openId)
	if err != nil {
		logs.Error.Println("GetByOpenid error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}

	if ui == nil {
		logs.Error.Println("UserShop not exist")
		return c.JSON(http.StatusOK, APIResult{Error: APIError{Code: 20003}})
	}

	// If need perfect user info
	bc, err := model.BrandCustomer{}.Get(ui.FashionBrandCustomer.BrandCode, ui.Customer.Mobile)
	if err != nil {
		logs.Error.Println("GetUserDetail error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10003)
	}
	needPerfect := bc == nil || !bc.HasFilled

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

	bc, err := model.BrandCustomer{}.Get(brandCode, mobile)
	if err != nil {
		logs.Error.Println("GetUserDetail error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}
	// 完善顾客信息上线前注册的用户
	if bc == nil {
		bc = &model.BrandCustomer{}
		bc.HasFilled = false
	}
	bc = getUserDetailFromCSL(brandCode, openId, bc, bc.HasFilled)

	result := make(map[string]interface{})
	if bc != nil {
		result["mobile"] = mobile
		result["name"] = bc.Name
		result["gender"] = bc.Gender
		result["birthday"] = bc.Birthday
		result["address"] = bc.Address
		result["detailAddress"] = bc.DetailAddress
		result["email"] = bc.Email
		result["isMarried"] = bc.IsMarried
		result["hasFilled"] = bc.HasFilled
	}

	return c.JSON(http.StatusOK, APIResult{Success: true, Result: result})

}

func APIUpdatePerfectInfo(c echo.Context) error {
	// var brandCustomer model.BrandCustomer
	oldMobile := c.Get("user").(*extends.AuthClaims).Mobile
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	openId := c.Get("user").(*extends.AuthClaims).OpenId
	custNo := c.Get("user").(*extends.AuthClaims).CustNo

	brandCustomer, err := model.BrandCustomer{}.Get(brandCode, oldMobile)
	if err != nil {
		logs.Error.Println(err)
	}

	brandCustomer.WxOpenID = openId
	brandCustomer.CustNo = custNo
	brandCustomer.BrandCode = brandCode
	brandCustomer.Name = c.FormValue("name")
	brandCustomer.Gender = c.FormValue("gender")
	phoneNo := c.FormValue("mobile")
	//logs.Error.Println(userDetail.Name)
	//logs.Error.Println(phoneNo)
	if phoneNo != oldMobile {
		brandCustomer.Mobile = phoneNo
	} else {
		brandCustomer.Mobile = oldMobile
	}

	birthday := c.FormValue("birthday")
	if birthday != "" {
		brandCustomer.Birthday = birthday
	}

	brandCustomer.Address = c.FormValue("address")
	brandCustomer.DetailAddress = c.FormValue("detailAddress")
	brandCustomer.Email = c.FormValue("email")

	isMarried := c.FormValue("isMarried")
	// var err error
	if isMarried != "" {
		brandCustomer.IsMarried, err = strconv.ParseBool(isMarried)
		if err != nil {
			logs.Error.Println("[UpdatePerfectInfo]isMarried convert bool error", err)
			return echo.NewHTTPError(http.StatusBadRequest, 10011)
		}
	}
	logs.Warning.Println(phoneNo, " ", oldMobile)
	if phoneNo != oldMobile {
		err = brandCustomer.ChangeMobileWithOld(oldMobile, phoneNo)
		if err != nil {
			logs.Error.Println("[UpdatePerfectInfo]saveUserDetail", err)
			return echo.NewHTTPError(http.StatusInternalServerError, 10013)
		}
	} else {
		err = brandCustomer.Save()
		if err != nil {
			logs.Error.Println("[UpdatePerfectInfo]saveUserDetail", err)
			return echo.NewHTTPError(http.StatusInternalServerError, 10013)
		}
	}

	brandCustomer.HasFilled = true
	err = brandCustomer.UpdateHasFilled()

	err = saveUserDetailToCSL(brandCustomer, openId)
	if err != nil {
		logs.Error.Println("[SaveUserDetailToCSL]saveUserDetail", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}

	return c.JSON(http.StatusOK, APIResult{Success: true, Result: map[string]string{"openId": openId}})

}

// 从CSL获取用户信息
func getUserDetailFromCSL(brandCode, openId string, bc *model.BrandCustomer, splitAddress bool) *model.BrandCustomer {
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
		bc.Address = string([]byte(data.Address)[0:n])
		bc.DetailAddress = string([]byte(data.Address)[n+1:])
	} else {
		bc.DetailAddress = data.Address
	}
	bc.Name = data.CustName
	bc.Mobile = data.TelephoneNo
	bc.Birthday = data.Birthday
	bc.BrandCode = data.BrandCode
	bc.Gender = data.SexCode
	bc.Email = data.EmailAddress
	bc.IsMarried = data.WeddingChk
	bc.IsNewCust = data.IsNewCust
	return bc
}

func saveUserDetailToCSL(bc *model.BrandCustomer, openId string) error {
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

	data.Address = bc.Address + "," + bc.DetailAddress
	data.Birthday = bc.Birthday
	data.BrandCode = strings.ToUpper(bc.BrandCode)
	data.CustName = bc.Name
	data.EmailAddress = bc.Email
	data.SexCode = bc.Gender
	data.TelephoneNo = bc.Mobile
	data.OpenId = openId
	if bc.IsMarried {
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
