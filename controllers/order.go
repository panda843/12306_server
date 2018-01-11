package controllers

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

// PassengerController Operations about object
type OrderController struct {
	BaseController
}

// @Title Post
// @Description 添加订单
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router / [post]
func (o *OrderController) Post() {
	//获取提交信息
	secretStr := o.GetString("secret_key")
	trainDate := o.GetString("train_date")
	startStation := o.GetString("start_station")
	endStation := o.GetString("end_station")
	//检测登录用户
	// checkUserData := &url.Values{}
	// checkUserData.Set("_json_att","")
	// request.SetHeader("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	// request.SetHeader("X-Requested-With","XMLHttpRequest")
	// isCheckOk, checkResData := request.SetURL(beego.AppConfig.String("12306::URLCheckLoginUser")).Post(checkUserData)
	// beego.Info("检测登录用户 ----->  ", isCheckOk,checkResData)
	// if !isCheckOk {
	// 	o.Fail().SetMsg("检测登录用户失败").Send()
	// 	return
	// }
	var cookies []*http.Cookie

	expiration := time.Now().Add(10 * time.Minute)
	cookie := &http.Cookie{Name: "_jc_save_formDate", Value: trainDate, Expires: expiration}
	cookies = append(cookies, cookie)
	cookie = &http.Cookie{Name: "_jc_save_fromStation", Value: url.QueryEscape(startStation), Expires: expiration}
	cookies = append(cookies, cookie)
	cookie = &http.Cookie{Name: "_jc_save_toStation", Value: url.QueryEscape(endStation), Expires: expiration}
	cookies = append(cookies, cookie)
	cookie = &http.Cookie{Name: "_jc_save_toDate", Value: trainDate, Expires: expiration}
	cookies = append(cookies, cookie)
	cookie = &http.Cookie{Name: "_jc_save_wfdc_flag", Value: "dc", Expires: expiration}
	cookies = append(cookies, cookie)
	cookie = &http.Cookie{Name: "current_captcha_type", Value: "Z", Expires: expiration}
	cookies = append(cookies, cookie)
	request.SetCookie(cookies)
	//提交订单
	str := strings.Split(startStation, ",")
	end := strings.Split(endStation, ",")
	submitOrderData := &url.Values{}
	submitOrderData.Set("secretStr", secretStr)
	submitOrderData.Set("train_date", trainDate)
	submitOrderData.Set("back_train_date", trainDate)
	submitOrderData.Set("tour_flag", "dc")
	submitOrderData.Set("purpose_codes", "ADULT")
	submitOrderData.Set("query_from_station_name", string(str[0]))
	submitOrderData.Set("query_to_station_name", string(end[0]))
	submitOrderData.Set("undefined", "")
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	request.SetHeader("X-Requested-With", "XMLHttpRequest")
	isSubmitOk, submitOrderResData := request.SetURL(beego.AppConfig.String("12306::URLSubmitOrder")).Post(submitOrderData)
	beego.Info("提交用户订单 ----->  ", isSubmitOk, submitOrderResData)
	if !isSubmitOk {
		o.Fail().SetMsg("提交用户订单失败").Send()
		return
	}

	//下单页面init
	initOrderData := &url.Values{}
	initOrderData.Set("_json_att", "")
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	request.SetHeader("X-Requested-With", "XMLHttpRequest")
	isInitOrderOk, initOrderResData := request.SetURL(beego.AppConfig.String("12306::URLSubmitOrderInit")).Post(initOrderData)
	beego.Info("下单页面init ----->  ", isInitOrderOk, initOrderResData)
	if !isInitOrderOk {
		o.Fail().SetMsg("下单init失败").Send()
		return
	}
	splData := strings.Split(initOrderResData, "\n")
	var submitToken []byte
	beego.Debug(splData[8])
	if len(splData) > 64 {
		submitToken = []byte(splData[8])[32:64]
		beego.Debug("SubmitToken:", string(submitToken))
	} else {
		beego.Debug("获取Token出错了")
	}

	//获取乘客列表
	passengerData := &url.Values{}
	passengerData.Set("_json_att", "")
	passengerData.Set("REPEAT_SUBMIT_TOKEN", string(submitToken))
	request.SetHeader("X-Requested-With", "XMLHttpRequest")
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	isPassOk, passengerResData := request.SetURL(beego.AppConfig.String("12306::URLGetPassgener")).Post(passengerData)
	beego.Info("获取乘客列表 ----->  ", isPassOk, passengerResData)
	if !isPassOk {
		o.Fail().SetMsg("乘客列表获取失败").Send()
		return
	}

	//https://kyfw.12306.cn/otn/confirmPassenger/checkOrderInfo
	//cancel_flag=2
	//&bed_level_order_num=000000000000000000000000000000
	//&passengerTicketStr=3,0,1,段恩建,1,510723199209121772,18780597049,N
	//&oldPassengerStr=段恩建,1,510723199209121772,1_
	//&tour_flag=dc
	//&randCode=
	//&whatsSelect=1
	//&_json_att=
	//&REPEAT_SUBMIT_TOKEN=8ff8976a1014762e2af9f72ab9c516ab
	//确认订单
	confirmData := &url.Values{}
	confirmData.Set("cancel_flag", "")
	confirmData.Set("bed_level_order_num", "000000000000000000000000000000")
	confirmData.Set("passengerTicketStr", "3,0,1,段恩建,1,510723199209121772,18780597049,N")
	confirmData.Set("oldPassengerStr", "段恩建,1,510723199209121772,1_")
	confirmData.Set("tour_flag", "dc")
	confirmData.Set("randCode", "")
	confirmData.Set("whatsSelect", "1")
	confirmData.Set("_json_att", "")
	confirmData.Set("REPEAT_SUBMIT_TOKEN", string(submitToken))
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	request.SetHeader("X-Requested-With", "XMLHttpRequest")
	isconfirmOk, confirmResData := request.SetURL("https://kyfw.12306.cn/otn/confirmPassenger/checkOrderInfo").Post(confirmData)
	beego.Info("确认订单 ----->  ", isconfirmOk, confirmResData)
	if !isconfirmOk {
		o.Fail().SetMsg("确认订单失败").Send()
		return
	}

	o.Fail().SetMsg("TTTTTT").Send()
}
