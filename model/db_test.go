package model

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func initDB() error {
	return InitDB("mysql", "root:aabb1122@tcp(127.0.0.1:3306)/wxshop?charset=utf8&parseTime=True&loc=Local")
}

func TestUserShopCURD(t *testing.T) {
	initDB()
	Convey("UserShopCURD", t, func() {
		Convey("Create", func() {
			i := FashionBrandCustomerInfo{}
			i.Customer.Mobile = "13161955000"
			i.FashionBrandCustomer.BrandCode = "tt"
			i.FashionBrandCustomer.CustNo = "NO99999999"
			i.FashionBrandCustomer.WxOpenID = "oYiR6wV9swnxcaXrEkXDPLzt27Wg"
			err := i.Create()
			if err != nil {
				fmt.Println("Create error: ", err)
			}
			if err != BrandCustomerAlreadyExistError {
				So(err, ShouldBeNil)
			}
			Convey("Check Customer", func() {
				c, err := Customer{}.GetByMobile("13161955000")
				if err != nil {
					fmt.Println("GetByMobile error: ", err)
				} else {
					fmt.Println("GetByMobile: ", c)
				}
				So(err, ShouldBeNil)
			})
			Convey("Check BrandCustomer", func() {
				c, err := BrandCustomer{}.Get("tt", "13161955000")
				if err != nil {
					fmt.Println("GetByMobile error: ", err)
				} else {
					fmt.Println("GetByMobile: ", c)
				}
				So(err, ShouldBeNil)
			})
			Convey("Check FashionBrandCustomer", func() {
				c, err := FashionBrandCustomer{}.GetByMobile("tt", "13161955000")
				if err != nil {
					fmt.Println("GetByMobile error: ", err)
				} else {
					fmt.Println("GetByMobile: ", c)
				}
				So(err, ShouldBeNil)
			})
			Convey("Check FashionBrandCustomerInfo", func() {
				c, err := FashionBrandCustomerInfo{}.GetByWxOpenID("tt", "oYiR6wV9swnxcaXrEkXDPLzt27Wg")
				if err != nil {
					fmt.Println("Check FashionBrandCustomerInfo Error: ", err)
				} else {
					fmt.Println("Check FashionBrandCustomerInfo: ", c)
				}
				fmt.Println("c.Customer:", c.Customer)
				fmt.Println("c.FashionBrandCustomer:", c.FashionBrandCustomer)
				fmt.Println("c.BC:", c.BC)
				So(err, ShouldBeNil)
			})
		})

		Convey("Change Mobile", func() {
			i := FashionBrandCustomerInfo{}
			i.Customer.Mobile = "13161955001"
			i.FashionBrandCustomer.BrandCode = "tt"
			i.FashionBrandCustomer.CustNo = "1"
			i.FashionBrandCustomer.WxOpenID = "oYiR6wV9swnxcaXrEkXDPLzt27Wg111"
			i.Create()

			i2 := FashionBrandCustomerInfo{}
			i2.Customer.Mobile = "13161955002"
			i2.FashionBrandCustomer.BrandCode = "tt"
			i2.FashionBrandCustomer.CustNo = "2"
			i2.FashionBrandCustomer.WxOpenID = "oYiR6wV9swnxcaXrEkXDPLzt27Wg211"
			i2.Create()

			i3 := FashionBrandCustomerInfo{}
			i3.Customer.Mobile = "13161955003"
			i3.FashionBrandCustomer.BrandCode = "tt"
			i3.FashionBrandCustomer.CustNo = "3"
			i3.FashionBrandCustomer.WxOpenID = "oYiR6wV9swnxcaXrEkXDPLzt27Wg311"
			i3.Create()

			err := i3.ChangeMobileWithOld("13161955003", "13161955002")
			So(err, ShouldBeNil)

			c, err := FashionBrandCustomer{}.GetByMobile("tt", "13161955002")
			So(err, ShouldBeNil)
			So(c.CustNo == "3", ShouldBeTrue)

			c2, err := Customer{}.Get(i2.CustomerId)
			So(err, ShouldBeNil)
			fmt.Println("c2.Mobile:", c2.Mobile)
			// t.Log(c2.Mobile)

		})

		Convey("UpdatePassword", func() {
			i := FashionBrandCustomerInfo{}
			i.Customer.Mobile = "13161955000"
			i.FashionBrandCustomer.BrandCode = "tt"
			err := i.UpdatePassword()
			if err != nil {
				fmt.Println("UpdatePassword error: ", err)
			}
			So(err, ShouldBeNil)
		})

		Convey("UpdateForGame", func() {
			i := FashionBrandCustomerInfo{}
			i.FashionBrandCustomer.WxOpenID = "oYiR6wV9swnxcaXrEkXDPLzt27Wg"
			i.FashionBrandCustomer.BrandCode = "tt"
			i.ReceiveAddress = "北京朝阳"
			i.ReceiveName = "强哥"
			i.ReceiveSize = "1000"
			i.ReceiveTelephone = "13161955008"
			err := i.UpdateForGame()
			if err != nil {
				fmt.Println("UpdateForGame error: ", err)
			}
			So(err, ShouldBeNil)
		})

		// Convey("UpdateCouponErrorTimes", func() {
		// 	err := FashionBrandCustomer{}.UpdateCouponErrorTimes("13161955000", "tt", false)
		// 	if err != nil {
		// 		fmt.Println("UpdateCouponErrorTimes error: ", err)
		// 	}
		// 	So(err, ShouldBeNil)
		// })

		Convey("GetByMobile", func() {
			userShop, err := FashionBrandCustomer{}.GetByMobile("tt", "13161955000")
			if err != nil {
				fmt.Println("GetByMobile error: ", err)
			} else {
				fmt.Println("GetByMobile: ", userShop)
			}
			So(err, ShouldBeNil)
		})

		Convey("GetByOpenid", func() {
			userShop, err := FashionBrandCustomer{}.GetByWxOpenID("tt", "oYiR6wV9swnxcaXrEkXDPLzt27Wg")
			if err != nil {
				fmt.Println("GetByOpenid error: ", err)
			} else {
				fmt.Println("GetByOpenid: ", userShop)
			}
			So(err, ShouldBeNil)
		})

		Convey("CheckMobile", func() {
			flag, err := FashionBrandCustomer{}.GetByMobile("tt", "13161955000")
			if err != nil {
				fmt.Println("CheckMobile error: ", err)
			}
			fmt.Println("CheckMobile flag:", flag)
			So(err, ShouldBeNil)
		})

		// Convey("GetCouponErrorTimes", func() {
		// 	times, err := FashionBrandCustomer{}.GetCouponErrorTimes("13161955000", "tt")
		// 	if err != nil {
		// 		fmt.Println("GetCouponErrorTimes error: ", err)
		// 	}
		// 	fmt.Println("GetCouponErrorTimes times:", times)
		// 	So(err, ShouldBeNil)
		// })

		Convey("Delete", func() {
			i, err := FashionBrandCustomer{}.GetByMobile("tt", "13161955000")
			So(err, ShouldBeNil)
			So(i, ShouldNotBeNil)
			if i != nil {
				err = FashionBrandCustomerInfo{}.Delete(i.BrandCode, i.CustomerId)
				if err != nil {
					fmt.Println("Delete error: ", err)
				}
				So(err, ShouldBeNil)
			}

		})
	})
}

// func TestErrorUserC(t *testing.T) {
// 	initDB()
// 	Convey("ErrorUserC", t, func() {
// 		Convey("SaveErrorUser", func() {
// 			i := ErrorUser{}
// 			i.BrandCode = "tt"
// 			i.ErrorMessage = "Testing error user"
// 			i.Mobile = "13161955000"
// 			i.OpenId = "oYiR6wV9swnxcaXrEkXDPLzt27Wg"
// 			err := i.SaveErrorUser()
// 			if err != nil {
// 				fmt.Println("SaveErrorUser error: ", err)
// 			}
// 			So(err, ShouldBeNil)

// 			// Delete
// 			user := ErrorUser{}
// 			_, err = db.Where("mobile = ?", "13161955000").And("brand_code = ?", "tt").Delete(&user)
// 			if err != nil {
// 				fmt.Println("DeleteErrorUser error: ", err)
// 			}
// 		})
// 	})
// }

func TestSmsCR(t *testing.T) {
	initDB()
	Convey("TestSmsCR", t, func() {
		Convey("Create", func() {
			s := Sms{}
			s.Mobile = "13161955000"
			s.Type = "Testing"
			s.VerCode = "654321"
			s.BrandCode = "tt"
			err := s.Create()
			if err != nil {
				fmt.Println("Create error: ", err)
			}
			So(err, ShouldBeNil)
		})

		Convey("CheckVerCode", func() {
			s := Sms{}
			flag, err := s.CheckVerCode("13161955000", "654321")
			if err != nil {
				fmt.Println("CheckVerCode error: ", err)
			}
			fmt.Println("CheckVerCode flag:", flag)
			So(err, ShouldBeNil)

			// Delete
			s = Sms{}
			_, err = db.Where("mobile = ?", "13161955000").Delete(&s)
			if err != nil {
				fmt.Println("DeleteSms error: ", err)
			}
		})
	})
}

func TestUserDetailCUR(t *testing.T) {
	initDB()
	Convey("TestUserDetailCUR", t, func() {
		Convey("SaveUserDetail", func() {
			u := BrandCustomer{}
			u.Address = "北京 北京 朝阳区"
			u.Birthday = "19900614"
			u.BrandCode = "tt"
			u.DetailAddress = "酒仙桥中路恒通商务园B51座"
			u.Email = "elvinchan@foxmail.com"
			u.Gender = "M"
			u.HasFilled = false
			u.IsMarried = false
			u.IsNewCust = 1
			u.Mobile = "13161955000"
			u.Name = "强哥"
			err := u.Save()
			if err != nil {
				fmt.Println("SaveUserDetail error: ", err)
			}
			So(err, ShouldBeNil)
		})
		Convey("GetUserDetail", func() {
			user, err := BrandCustomer{}.Get("tt", "13161955000")
			if err != nil {
				fmt.Println("GetUserDetail error: ", err)
			}
			fmt.Println("GetUserDetail: ", user)
			So(err, ShouldBeNil)
		})
		Convey("UpdateHasFilled", func() {
			u := BrandCustomer{}
			u.Mobile = "13161955000"
			u.BrandCode = "tt"
			u.HasFilled = true
			err := u.UpdateHasFilled()
			if err != nil {
				fmt.Println("UpdateHasFilled error: ", err)
			}
			So(err, ShouldBeNil)

			// Delete
			u = BrandCustomer{}
			_, err = db.Where("mobile = ?", "13161955000").And("brand_code = ?", "tt").Delete(&u)
			if err != nil {
				fmt.Println("DeleteUserDetail error: ", err)
			}
		})
	})
}

// func TestUserMhCURD(t *testing.T) {
// 	initDB()
// 	Convey("UserMhCURD", t, func() {
// 		Convey("Create", func() {
// 			i := FashionBrandCustomerInfoMh{}
// 			i.Customer.Mobile = "13161955000"
// 			i.UserMh.BrandCode = "tt"
// 			i.UserMh.VipCode = "juifnewiuncwnencw"
// 			err := i.Create()
// 			if err != nil {
// 				fmt.Println("Create error: ", err)
// 			}
// 			So(err, ShouldBeNil)
// 		})

// 		Convey("GetByMobile", func() {
// 			userMh, err := UserMh{}.GetByMobile("13161955000", "tt")
// 			if err != nil {
// 				fmt.Println("GetByMobile error: ", err)
// 			} else {
// 				fmt.Println("GetByMobile: ", userMh)
// 			}
// 			So(err, ShouldBeNil)
// 		})

// 		Convey("CheckMobile", func() {
// 			flag, err := UserMh{}.CheckMobile("13161955000", "tt")
// 			if err != nil {
// 				fmt.Println("CheckMobile error: ", err)
// 			}
// 			fmt.Println("CheckMobile flag:", flag)
// 			So(err, ShouldBeNil)
// 		})

// 		Convey("Update", func() {
// 			userMh, err := UserMh{}.GetByMobile("13161955000", "tt")
// 			if err != nil {
// 				fmt.Println("GetByMobile error: ", err)
// 			} else {
// 				fmt.Println("GetByMobile: ", userMh)
// 			}

// 			userMh.VipCode = "jioasnicncnasn"
// 			err = userMh.Update()
// 			if err != nil {
// 				fmt.Println("Update error: ", err)
// 			}
// 			So(err, ShouldBeNil)
// 		})

// 		Convey("Delete", func() {
// 			i := FashionBrandCustomerInfoMh{}
// 			i.Customer.Mobile = "13161955000"
// 			i.UserMh.BrandCode = "tt"
// 			err := i.Delete()
// 			if err != nil {
// 				fmt.Println("Delete error: ", err)
// 			}
// 			So(err, ShouldBeNil)
// 		})
// 	})
// }
