package fashion

import (
	"best/p2-customer-service/config"
	. "best/p2-customer-service/dto"
	"best/p2-customer-service/extends"
	"best/p2-customer-service/logs"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"net/http"

	"github.com/labstack/echo"
	"github.com/smallnest/goreq"
)

// ApiGetCurrentIntegral 获取总积分
func ApiGetCurrentIntegral(c echo.Context) error {

	mobile := c.Get("user").(*extends.AuthClaims).Mobile
	openId := c.Get("user").(*extends.AuthClaims).OpenId
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	logs.Debug.Println(mobile, openId, brandCode)
	if mobile == "" || openId == "" || brandCode == "" {
		return echo.NewHTTPError(http.StatusForbidden, 20001)
	}

	type IntegrationDto struct {
		Xmlns         string `xml:"xmlns,attr"`
		SellingAmt    string `xml:"SellingAmt"`
		UseMileage    string `xml:"UseMileage"`
		ObtainMileage string `xml:"ObtainMileage"`
		RetainMileage string `xml:"RetainMileage"`
		Xsi           string `xml:"xsi,attr"`
		Xsd           string `xml:"xsd,attr"`
	}

	url := config.Config.Adapter.CSL.CslWebService + "/GetCurrentIntegral"
	q := fmt.Sprintf("strCallUserCode=%v&strCallPassword=%v&brandCode=%v&wechatID=%v", "Eland", "Eland1234", strings.ToUpper(brandCode), openId)
	_, body, reqErr := goreq.New().Get(url).Query(q).ContentType("xml").SetCurlCommand(true).EndBytes()
	if reqErr != nil {
		logs.Error.Println("GetCurrentIntegral error: ", reqErr)
		// extends.ReturnJsonFailure(w, http.StatusInternalServerError, 10003, reqErr[0].Error())
		// return
		return echo.NewHTTPError(http.StatusInternalServerError, 10003)
	}

	var structData IntegrationDto
	err := xml.Unmarshal(body, &structData)
	if err != nil {
		logs.Error.Println("GetCurrentIntegral unmarshal error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10003)
	}

	result := make(map[string]string)
	result["RetainMileage"] = structData.RetainMileage
	//extends.ReturnJsonSuccess(w, http.StatusOK, result)
	return c.JSON(http.StatusOK, APIResult{Success: true, Result: result})
}

// ApiGetIntegralHistory 获取积分记录
func ApiGetIntegralHistory(c echo.Context) error {

	mobile := c.Get("user").(*extends.AuthClaims).Mobile
	openId := c.Get("user").(*extends.AuthClaims).OpenId
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	logs.Debug.Println(mobile, openId, brandCode)
	if mobile == "" || openId == "" || brandCode == "" {
		return echo.NewHTTPError(http.StatusForbidden, 20001)
	}

	// openId = "oS1CewapWOtyanVNPJsd4DesSwys"
	// brandCode = "BC"

	// A:获取积分记录 B:使用积分记录
	integralType := c.FormValue("integralType")
	if integralType == "" {

		return echo.NewHTTPError(http.StatusBadRequest, 10012)
	}

	type ArrayOfIntegrationDto struct {
		Xsi           string   `xml:"xsi,attr"`
		UseMileage    []string `xml:"IntegrationDto>UseMileage"`
		Date          []string `xml:"IntegrationDto>Date"`
		SellingAmt    []string `xml:"IntegrationDto>SellingAmt"`
		ObtainMileage []string `xml:"IntegrationDto>ObtainMileage"`
		RetainMileage []string `xml:"IntegrationDto>RetainMileage"`
		Xsd           string   `xml:"xsd,attr"`
		Xmlns         string   `xml:"xmlns,attr"`
		SaleNo        []string `xml:"IntegrationDto>SaleNo"`
	}

	url := config.Config.Adapter.CSL.CslWebService + "/GetIntegralHistoryForALL"
	q := fmt.Sprintf("strCallUserCode=%v&strCallPassword=%v&brandCode=%v&wechatID=%v&type=%v", "Eland", "Eland1234", strings.ToUpper(brandCode), openId, integralType)
	_, body, reqErr := goreq.New().Get(url).Query(q).ContentType("xml").SetCurlCommand(true).EndBytes()
	if reqErr != nil {
		logs.Error.Println("GetIntegralHistoryForALL error: ", reqErr)

		return echo.NewHTTPError(http.StatusInternalServerError, 10003)
	}

	structData := ArrayOfIntegrationDto{}
	err := xml.Unmarshal(body, &structData)
	if err != nil {
		logs.Error.Println("GetIntegralHistoryForALL unmarshal error: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10003)
	}

	var listData []map[string]string
	for i := 0; i < len(structData.SaleNo); i++ {
		m := make(map[string]string)

		m["SaleNo"] = structData.SaleNo[i]
		m["SellingAmt"] = structData.SellingAmt[i]
		m["UseMileage"] = structData.UseMileage[i]
		m["ObtainMileage"] = structData.ObtainMileage[i]
		m["RetainMileage"] = structData.RetainMileage[i]
		m["Date"] = structData.Date[i]
		listData = append(listData, m)
	}
	//extends.ReturnJsonSuccess(w, http.StatusOK, listData)

	return c.JSON(http.StatusOK, APIResult{Success: true, Result: listData})
}

type DataFromXML struct {
	Text string `xml:",chardata"`
}
type VipGradeDto struct {
	Message          string
	Reuslt           int64
	Status           int64
	TotalPurchaseAmt float64
}

func ApiGetVipGrade(c echo.Context) error {

	custNo := c.Get("user").(*extends.AuthClaims).CustNo
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	logs.Debug.Println(custNo, brandCode)
	if custNo == "" || brandCode == "" {
		return echo.NewHTTPError(http.StatusForbidden, 20001)
	}

	// url := "http://10.202.101.10:8022/CustomerServiceWS.asmx/GetTotalPurchaseAmt?strCallUserCode=Eland&strCallPassword=Eland1234&brandCode=rc&custNo=0001852352"
	url := config.Config.Adapter.CSL.CslWebService + "/GetTotalPurchaseAmt"

	q := fmt.Sprintf("strCallUserCode=%v&strCallPassword=%v&brandCode=%v&custNo=%v", "Eland", "Eland1234", brandCode, custNo)

	_, body, reqErr := goreq.New().Get(url).Query(q).ContentType("xml").SetCurlCommand(true).End()
	if reqErr != nil {
		logs.Error.Println("GetMemberInfo error: ", reqErr)
	}
	DataFromXML := new(DataFromXML)
	err := xml.Unmarshal([]byte(body), DataFromXML)
	if err != nil {
		logs.Error.Println("Unmarshal xml error: ", err)
	}
	// logs.Warning.Println(DataFromXML.Text)

	vipGradeDto := new(VipGradeDto)
	err = json.Unmarshal([]byte(DataFromXML.Text), vipGradeDto)
	if err != nil {
		logs.Error.Println("Unmarshal xml error: ", err)
	}
	// logs.Succ.Println(vipGradeDto)
	//extends.ReturnJsonSuccess(w, http.StatusOK, vipGradeDto)
	return c.JSON(http.StatusOK, APIResult{Success: true, Result: vipGradeDto})
}

func ApiUpdateIntegralExchange(c echo.Context) error {

	custNo := c.Get("user").(*extends.AuthClaims).CustNo
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	logs.Debug.Println(custNo, brandCode)
	if custNo == "" || brandCode == "" {
		return echo.NewHTTPError(http.StatusForbidden, 20001)
	}

	integral := c.FormValue("integral")

	// url := "http://10.202.101.10:8022/CustomerServiceWS.asmx/IntegralExchange"
	url := config.Config.Adapter.CSL.CslWebService + "/IntegralExchange"

	q := fmt.Sprintf("strCallUserCode=%v&strCallPassword=%v&brandCode=%v&custNo=%v&integral=%v&etcMileageTypeCode=%v&etcMileageTypeClassCode=%v", "Eland", "Eland1234", brandCode, custNo, integral, "w22", "w22001")
	_, body, reqErr := goreq.New().Get(url).Query(q).ContentType("xml").SetCurlCommand(true).End()
	if reqErr != nil {
		logs.Error.Println("GetMemberInfo error: ", reqErr)
	}
	DataFromXML := new(DataFromXML)
	err := xml.Unmarshal([]byte(body), DataFromXML)
	if err != nil {
		logs.Error.Println("Unmarshal xml error: ", err)
	}
	logs.Warning.Println(DataFromXML.Text)
	vipGradeDto := new(VipGradeDto)
	err = json.Unmarshal([]byte(DataFromXML.Text), vipGradeDto)
	if err != nil {
		logs.Error.Println("Unmarshal xml error: ", err)
		// return nil, nil
	}
	// logs.Succ.Println(vipGradeDto)
	//extends.ReturnJsonSuccess(w, http.StatusOK, vipGradeDto)

	return c.JSON(http.StatusOK, APIResult{Success: true, Result: vipGradeDto})

}
