package controllers

import (
	"net/url"
	"strings"
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

	//提交订单
	submitOrderData := &url.Values{}
	submitOrderData.Set("secretStr",secretStr)
	submitOrderData.Set("train_date",trainDate)
	submitOrderData.Set("back_train_date",trainDate)
	submitOrderData.Set("tour_flag","dc")
	submitOrderData.Set("purpose_codes","ADULT")
	submitOrderData.Set("query_from_station_name",startStation)
	submitOrderData.Set("query_to_station_name",endStation)
	submitOrderData.Set("undefined","")
	request.SetHeader("X-Requested-With","XMLHttpRequest")
	isSubmitOk, submitOrderResData := request.SetURL(beego.AppConfig.String("12306::URLSubmitOrder")).Post(submitOrderData)
	beego.Info("提交用户订单 ----->  ", isSubmitOk,submitOrderResData)
	if !isSubmitOk {
		o.Fail().SetMsg("提交用户订单失败").Send()
		return
	}

	
	//下单页面init
	initOrderData := &url.Values{}
	initOrderData.Set("_json_att","")
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	request.SetHeader("X-Requested-With","XMLHttpRequest")
	isInitOrderOk, initOrderResData := request.SetURL(beego.AppConfig.String("12306::URLSubmitOrderInit")).Post(initOrderData)
	beego.Info("下单页面init ----->  ", isInitOrderOk,initOrderResData)
	if !isInitOrderOk {
		o.Fail().SetMsg("下单init失败").Send()
		return
	}
	splData := strings.Split(initOrderResData,"\n")
	var submitToken []byte
	beego.Debug(splData[8])
	if len(splData) > 32 {
		submitToken = []byte(splData[8])[32:64]
		beego.Debug("SubmitToken:",string(submitToken))
	}else{
		beego.Debug("获取Token出错了")
	}
	//检测登录用户
	checkUserData := &url.Values{}
	checkUserData.Set("_json_att","")
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	request.SetHeader("X-Requested-With","XMLHttpRequest")
	isCheckOk, checkResData := request.SetURL(beego.AppConfig.String("12306::URLCheckLoginUser")).Post(checkUserData)
	beego.Info("检测登录用户 ----->  ", isCheckOk,checkResData)
	if !isCheckOk {
		o.Fail().SetMsg("检测登录用户失败").Send()
		return
	}


	//获取乘客列表
	passengerData := &url.Values{}
	passengerData.Set("_json_att","")
	passengerData.Set("REPEAT_SUBMIT_TOKEN",string(submitToken))
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	request.SetHeader("X-Requested-With","XMLHttpRequest")
	isPassOk, passengerResData := request.SetURL(beego.AppConfig.String("12306::URLGetPassgener")).Post(passengerData)
	beego.Info("获取乘客列表 ----->  ", isPassOk,passengerResData)
	if !isPassOk {
		o.Fail().SetMsg("乘客列表获取失败").Send()
		return
	}


	

	o.Fail().SetMsg("TTTTTT").Send()
}