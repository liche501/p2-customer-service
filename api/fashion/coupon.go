package fashion

import (
	"best/p2-customer-service/config"
	. "best/p2-customer-service/dto"
	"best/p2-customer-service/extends"
	"best/p2-customer-service/logs"

	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/smallnest/goreq"
)

type Coupon struct {
	CouponName string `json:"couponName"`
	CouponNo   string `json:"couponNo"`
	CouponType struct {
		CouponTypeCode string `json:"couponTypeCode"`
		CouponTypeName string `json:"couponTypeName"`
	} `json:"couponType"`
	DiscountDescription string    `json:"discountDescription"`
	EndDate             time.Time `json:"endDate"`
	RuleDescription     string    `json:"ruleDescription"`
	StartDate           time.Time `json:"startDate"`
	Used                bool      `json:"used"`
}
type CouponWxcrm struct {
	Success bool `json:"success"`
	Result  []struct {
		CouponNo    string    `json:"couponNo"`
		Title       string    `json:"title"`
		Type        string    `json:"type"`
		Description string    `json:"description"`
		Notice      string    `json:"notice"`
		StartDate   time.Time `json:"startDate"`
		EndDate     time.Time `json:"endDate"`
		UseChk      bool      `json:"useChk"`
	} `json:"result"`
	Error string `json:"error"`
}

func APIGetCouponList(c echo.Context) error {
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	custNo := c.Get("user").(*extends.AuthClaims).CustNo

	// TODO:: getCouponIfNeed(ui.CustNo, mobile, ui.BrandCode)

	var map_data []Coupon
	serviceURL := fmt.Sprintf("http://%s:%s/%s", config.Config.Coupon.Contents.ServiceHost, config.Config.Coupon.Contents.Port, config.Config.Coupon.Contents.CouponAPIPrefix)
	logs.Debug.Println("serviceUrl:", serviceURL)

	q := fmt.Sprintf("CustomerNo=%v&BrandCode=%v", custNo, strings.ToUpper(brandCode))
	_, _, reqErr := goreq.New().Get(serviceURL).Query(q).BindBody(&map_data).SetCurlCommand(true).End()
	if reqErr != nil {
		logs.Error.Println("reqErr:", reqErr)
		return echo.NewHTTPError(http.StatusInternalServerError, 10003)
	}

	wxcrmCoupons, err := Coupon{}.GetFromWxcrm(strings.ToUpper(brandCode), custNo)
	if err != nil {
		logs.Error.Println("reqErr:", reqErr)
		return echo.NewHTTPError(http.StatusInternalServerError, 10003)
	}
	map_data = append(map_data, wxcrmCoupons...)

	container := Coupon{}.Decompose(map_data)
	return c.JSON(http.StatusOK, APIResult{Success: true, Result: container})
}

func (Coupon) GetFromWxcrm(brandCode, custNo string) ([]Coupon, error) {
	var wx_data CouponWxcrm
	//"http://10.202.101.17:5700/api/v1/Coupon/EE/GetCouponList/0000114179"
	serviceUrl := fmt.Sprintf("%s/api/v1/Coupon/%s/GetCouponList/%s", config.Config.Coupon.Marketing.EventCoupon, brandCode, custNo)

	_, _, reqErr := goreq.New().Get(serviceUrl).
		BindBody(&wx_data).
		SetCurlCommand(true).
		End()
	if reqErr != nil {
		logs.Error.Println("reqErr:", reqErr)
		return nil, reqErr[0]
	}

	coupons := []Coupon{}
	for _, value := range wx_data.Result {
		tempData := Coupon{}
		tempData.CouponName = value.Title
		tempData.CouponNo = value.CouponNo
		tempData.CouponType.CouponTypeCode = value.Type
		CouponTypeName := ""
		switch value.Type {
		case "CostMinus":
			CouponTypeName = "现金券"
		case "Discount":
			CouponTypeName = "折扣券"
		case "Exchange":
			CouponTypeName = "兑换券"
		}
		tempData.CouponType.CouponTypeName = CouponTypeName
		tempData.DiscountDescription = value.Description
		tempData.RuleDescription = value.Notice
		tempData.StartDate = value.StartDate
		tempData.EndDate = value.EndDate
		tempData.Used = value.UseChk
		coupons = append(coupons, tempData)
	}
	return coupons, nil
}

func (Coupon) Decompose(map_data []Coupon) map[string][]Coupon {
	container := make(map[string][]Coupon)
	container["dueList"] = []Coupon{}
	container["noUseList"] = []Coupon{}
	container["useList"] = []Coupon{}
	var noUseList []Coupon
	var dueList []Coupon
	var useList []Coupon

	logs.Debug.Println("map_data_len:", len(map_data))

	for _, value := range map_data {
		// start_data := value.StartDate
		end_data := value.EndDate

		if time.Now().After(end_data) {
			dueList = append(dueList, value)
			container["dueList"] = dueList
		}

		if value.Used == false && time.Now().Before(end_data) {
			noUseList = append(noUseList, value)
			container["noUseList"] = noUseList
		}
		if value.Used == true {
			useList = append(useList, value)
			container["useList"] = useList
		}
	}
	logs.Debug.Println("dueList_cout:", len(container["dueList"]))
	logs.Debug.Println("noUseList_cout:", len(container["noUseList"]))
	logs.Debug.Println("useList_cout:", len(container["useList"]))
	return container
}

func SendCoupon(brandCode, custNo string) error {

	url := config.Config.Adapter.CSL.CustomerInterfaceAPI + "/" + strings.ToUpper(brandCode) + custNo + "/Coupons"
	logs.Debug.Println("[SendCoupon]url", url)
	_, textData, reqErr := goreq.New().Post(url).SetCurlCommand(true).End()
	if reqErr != nil || textData == "" {
		logs.Error.Println("reqErr", len(reqErr), " textData:", textData)
		return reqErr[0]
	}
	logs.Succ.Println("[SendCoupon] success", textData)

	return nil
}
