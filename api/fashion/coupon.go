package fashion

import (
	"best/p2-customer-service/config"
	. "best/p2-customer-service/dto"
	"best/p2-customer-service/extends"
	"best/p2-customer-service/logs"
	"best/p2-customer-service/model"

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
	mobile := c.Get("user").(*extends.AuthClaims).Mobile
	// openId := c.Get("user").(*extends.AuthClaims).OpenId
	brandCode := c.Get("user").(*extends.AuthClaims).BrandCode
	custNo := c.Get("user").(*extends.AuthClaims).CustNo

	ui, err := model.FashionBrandCustomer{}.GetByMobile(brandCode, mobile)
	if err != nil {
		logs.Error.Println("call GetCurrentUser err:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, 10013)
	}
	getCouponIfNeed(ui.CustNo, mobile, ui.BrandCode)

	var map_data []Coupon
	serviceURL := fmt.Sprintf("http://%s:%s/%s", config.Config.Coupon.Contents.ServiceHost, config.Config.Coupon.Contents.Port, config.Config.Coupon.Contents.CouponAPIPrefix)
	logs.Debug.Println("serviceUrl:", serviceURL)

	q := fmt.Sprintf("CustomerNo=%v&BrandCode=%v", ui.CustNo, strings.ToUpper(brandCode))
	_, _, reqErr := goreq.New().Get(serviceURL).Query(q).BindBody(&map_data).SetCurlCommand(true).End()
	if reqErr != nil {
		logs.Error.Println("reqErr:", reqErr)
		// extends.ReturnJsonFailure(w, http.StatusInternalServerError, 10003)
		// return
		return echo.NewHTTPError(http.StatusInternalServerError, 10003)
	}
	// logs.Debug.Println("custNo==", custNo)

	map_data, errWxcrm := Coupon{}.AddCouponFromWxcrm(c, map_data, strings.ToUpper(brandCode), custNo)
	if errWxcrm != nil {
		logs.Error.Println("errWxcrm:", errWxcrm)
		// extends.ReturnJsonFailure(w, http.StatusInternalServerError, 10003)
		// return
		return echo.NewHTTPError(http.StatusInternalServerError, 10003)
	}
	container := Coupon{}.DecomposeCoupon(map_data)
	//extends.ReturnJsonSuccess(w, http.StatusOK, container)
	return c.JSON(http.StatusOK, APIResult{Success: true, Result: container})
}
func (Coupon) AddCouponFromWxcrm(c echo.Context, map_data []Coupon, brandCode, custNo string) ([]Coupon, error) {
	var wx_data CouponWxcrm
	//"http://10.202.101.17:5700/api/v1/Coupon/EE/GetCouponList/0000114179"
	serviceUrl := fmt.Sprintf("%s/api/v1/Coupon/%s/GetCouponList/%s", config.Config.Coupon.Marketing.EventCoupon, brandCode, custNo)

	_, _, reqErr := goreq.New().Get(serviceUrl).
		BindBody(&wx_data).
		SetCurlCommand(true).
		End()
	if reqErr != nil {
		logs.Error.Println("reqErr:", reqErr)
		// extends.ReturnJsonFailure(w, http.StatusInternalServerError, 10003)
		// return nil, reqErr[0]
		//   return c.JSON(http.StatusOK, APIResult{Success: true, Result: container})
		return nil, reqErr[0]
	}
	// logs.Warning.Println("body==", body)
	// logs.Succ.Println("wx_data==", wx_data)

	for _, value := range wx_data.Result {
		temp_data := Coupon{}
		temp_data.CouponName = value.Title
		temp_data.CouponNo = value.CouponNo
		temp_data.CouponType.CouponTypeCode = value.Type
		CouponTypeName := ""
		switch value.Type {
		case "CostMinus":
			CouponTypeName = "现金券"
		case "Discount":
			CouponTypeName = "折扣券"
		case "Exchange":
			CouponTypeName = "兑换券"
		}
		temp_data.CouponType.CouponTypeName = CouponTypeName
		temp_data.DiscountDescription = value.Description
		temp_data.RuleDescription = value.Notice
		temp_data.StartDate = value.StartDate
		temp_data.EndDate = value.EndDate
		temp_data.Used = value.UseChk
		map_data = append(map_data, temp_data)
	}
	// logs.Succ.Println("map_data==", map_data)
	return map_data, nil
}

func (Coupon) DecomposeCoupon(map_data []Coupon) map[string][]Coupon {
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

func getCouponIfNeed(custNo, mobile, brandCode string) {
	// 判断是否需要发券
	// us := model.FashionBrandCustomer{}
	// times, err := us.GetCouponErrorTimes(mobile, brandCode)
	// if err != nil {
	// 	logs.Error.Println("Call GetCouponErrorTimes err:", err)
	// 	return
	// }
	// if times > 0 && times < 4 {
	// 	//task.SendCoupon(mobile, custNo, brandCode, config.Config.Adapter.CustomerInterfaceApi)
	// }
}
