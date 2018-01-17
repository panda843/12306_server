package controllers

import (
	"net/url"

	"github.com/chuanshuo843/12306_server/utils/kyfw"
	// "net/http"
	// "net/url"
	// "strings"
	// "fmt"
	// "github.com/astaxie/beego"
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
	seatType := o.GetString("seat_type")
	trainDate := o.GetString("train_date")
	formatDate := o.GetString("format_date")
	startStation := o.GetString("start_station")
	endStation := o.GetString("end_station")
	ticketStr := o.GetString("ticket_str")
	passengerStr := o.GetString("passenger_str")
	pasSec, _ := url.Parse(secretStr)
	err := kyfwOrder.PlaceAnOrder(pasSec.String(), trainNo,trainCode,seatType,startStation, endStation, trainDate,formatDate,ticketStr,passengerStr)
	if err != nil {
		o.Fail().SetMsg(err.Error()).Send()
		return
	}
	o.Fail().SetMsg("DDDTTTTSSSSS").Send()
}
