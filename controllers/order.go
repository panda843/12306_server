package controllers

import (

	"github.com/chuanshuo843/12306_server/utils/kyfw"
)

var (
	kyfwOrder kyfw.Order
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

	secretStr := o.GetString("secret_key")
	trainNo := o.GetString("train_no")
	trainCode := o.GetString("train_code")
	trainDate := o.GetString("train_date")
	formatDate := o.GetString("format_date")
	// formatDate := "Thu+Jan+18+2018+00%3A00%3A00+GMT%2B0800+(%E4%B8%AD%E5%9B%BD%E6%A0%87%E5%87%86%E6%97%B6%E9%97%B4)"
	startStation := o.GetString("start_station")
	endStation := o.GetString("end_station")
	startCode := o.GetString("start_code")
	endCode := o.GetString("end_code")
	ticketStr := o.GetString("ticket_str") 
	passengerStr := o.GetString("passenger_str")
	err := kyfwOrder.PlaceAnOrder(secretStr, trainNo,trainCode,startStation,startCode, endStation, endCode,trainDate,formatDate,ticketStr,passengerStr)
	if err != nil {
		o.Fail().SetMsg(err.Error()).Send()
		return
	}
	o.Success().SetMsg("任务添加成功").Send() 
}
