package fashion

import (
	"best/p2-customer-service/config"
	. "best/p2-customer-service/dto"
	"best/p2-customer-service/event"
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
		return echo.NewHTTPError(http.StatusBadRequest, 10012)
	}

	fbc, err := model.FashionBrandCustomer{}.GetByMobile(brandCode, mobile)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}
	if fbc == nil {
		return c.JSON(http.StatusOK, APIResult{Success: true})
	}
	return c.JSON(http.StatusOK, APIResult{Error: APIError{Code: 20007, Message: "Mobile already registered"}})
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

	fbci := new(model.FashionBrandCustomerInfo)
	fbci.Customer.Mobile = mobile
	fbci.FashionBrandCustomer.BrandCode = brandCode
	fbci.FashionBrandCustomer.WxOpenID = openId
	logs.Debug.Println(fbci)

	//customer regist
	if err := fbci.Create(); err != nil {
		logs.Error.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	//check customer regist success
	brandCustomer, err := model.BrandCustomer{}.Get(brandCode, mobile)
	if err != nil {
		logs.Error.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}
	if brandCustomer == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, 20015)
	}
	//sendEvent BrandCustomerInitiated
	et := new(event.EventSender)
	// url := fmt.Sprintf("/v1/streams/%v/events/%v", "marketing", "BrandCustomerInitiated")
	et.EventBrokerUrl = "http://staging.p2shop.cn:50110"
	obj := event.BrandCustomerInitiated{}
	obj.CustomerID = brandCustomer.CustomerId
	obj.Telephone = brandCustomer.Mobile
	obj.Password = "123456"
	obj.BrandCode = brandCustomer.BrandCode
	obj.WxOpenID = brandCustomer.WxOpenID

	err = et.SendEvent("marketing", "BrandCustomerInitiated", obj)
	if err != nil {
		logs.Error.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	//UpdateStatus BrandCustomerInitiated
	brandCustomer.Status = "BrandCustomerInitiated"

	if err := brandCustomer.UpdateStatus(); err != nil {
		logs.Error.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}
	logs.Succ.Println("BrandCustomerInitiated is send")

	return c.JSON(http.StatusOK, APIResult{Success: true})

}

func APILogin(c echo.Context) error {
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	openId := c.Get("user").(*extends.AuthClaims).OpenId
	if openId == "" || brandCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, 10012)
	}
	logs.Debug.Println(brandCode, openId)

	fbci, err := model.FashionBrandCustomerInfo{}.GetSuccessUserByWxOpenID(brandCode, openId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}
	if fbci == nil {
		return c.JSON(http.StatusOK, APIResult{Success: false, Error: APIError{Code: 200, Message: "customer not exist"}})
	}

	jsonWebToken, err := extends.AuthHandler(brandCode, openId, fbci.Customer.Mobile, fbci.FashionBrandCustomer.CustNo)
	if err != nil {
		return err
	}
	logs.Succ.Println("Create JsonWebToken: ", jsonWebToken)
	rs := APIResult{Success: true, Result: map[string]string{"token": jsonWebToken, "status": fbci.Status()}}
	return c.JSON(http.StatusOK, rs)
}

// APIBrandCustomerStatus SSE
func APIBrandCustomerStatus(c echo.Context) error {
	openId := c.Get("user").(*extends.AuthClaims).OpenId
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	if openId == "" || brandCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, 10012)
	}

	// ===============SSE START===============

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")

	fbci, err := model.FashionBrandCustomerInfo{}.GetByWxOpenID(brandCode, openId)
	if err != nil {
		logs.Error.Println("GetByOpenidWithoutResgitStatus error: ", err)
		return c.String(http.StatusOK, "data:nouser\n\nexist")
	}
	if fbci == nil {
		return c.String(http.StatusOK, "data:nouser\n\nexist")
	}

	logs.Debug.Printf("RegistStatus for openId %v: %v", openId, fbci.Status)

	if fbci.Status() == "BrandCustomerDuplicated" {
		return c.String(http.StatusOK, "data:nouser\n\nError01")
	}
	if fbci.Status() == "BrandCustomerFailed" {
		return c.String(http.StatusOK, "data:nouser\n\nOtherError")

	}
	logs.Warning.Println(fbci.BrandCustomer.CustNo, " ", fbci.Status())
	if len(fbci.BrandCustomer.CustNo) > 0 && fbci.Status() == "BrandCustomerCreated" {
		return c.String(http.StatusOK, "data:nouser\n\nexist")
	}
	return c.String(http.StatusOK, "data:nouser\n\n")

	// ===============SSE END===============
}

func APIGetUserInfo(c echo.Context) error {

	openId := c.Get("user").(*extends.AuthClaims).OpenId
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode

	fbci, err := model.FashionBrandCustomerInfo{}.GetSuccessUserByWxOpenID(brandCode, openId)
	if err != nil {
		logs.Error.Println("GetByOpenid error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}
	if fbci == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, 20003)
	}

	// If need perfect user info
	bc, err := model.BrandCustomer{}.Get(fbci.FashionBrandCustomer.BrandCode, fbci.Customer.Mobile)
	if err != nil {
		logs.Error.Println("GetUserDetail error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}
	if bc == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, 20015)
	}
	needPerfect := bc == nil || !bc.HasFilled

	custNo := fbci.FashionBrandCustomer.CustNo
	mobile := fbci.Customer.Mobile

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

	bcFromDB, err := model.BrandCustomer{}.Get(brandCode, mobile)
	if err != nil {
		logs.Error.Println("GetUserDetail error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}
	if bcFromDB == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, 20015)
	}

	bc, err := getUserDetailFromCSL(bcFromDB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, APIResult{Success: true, Result: bc})

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
		echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}
	if brandCustomer == nil {
		echo.NewHTTPError(http.StatusInternalServerError, 20015)

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
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	} else {
		err = brandCustomer.Update()
		if err != nil {
			logs.Error.Println("[UpdatePerfectInfo]saveUserDetail", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	brandCustomer.HasFilled = true
	err = brandCustomer.UpdateHasFilled()

	err = saveUserDetailToCSL(brandCustomer, openId)
	if err != nil {
		logs.Error.Println("[SaveUserDetailToCSL]saveUserDetail", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, APIResult{Success: true})

}

// 从CSL获取用户信息
func getUserDetailFromCSL(bc *model.BrandCustomer) (*model.BrandCustomer, error) {
	url := fmt.Sprintf("%v?BrandCode=%v&OpenId=%v", config.Config.Adapter.CSL.CustomerInterfaceAPI, bc.BrandCode, bc.WxOpenID)
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
		return nil, reqErr[0]
	}
	logs.Debug.Println("Success GetUserInfoDetail: ", data)

	n := strings.Index(data.Address, ",")
	if bc.HasFilled && n >= 0 {
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
	return bc, nil
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
		return errors.New("Save userinfo to CSL failed!")
	}

	return nil
}
func setMyScoreDisplay(brandCode string) bool {
	if brandCode == "sf" || brandCode == "sm" || brandCode == "cz" {
		return false
	}
	return true
}
