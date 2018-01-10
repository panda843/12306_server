package controllers

import (
	"encoding/json"
	// "fmt"
	// "strings"
	"github.com/astaxie/beego"
	// "github.com/chuanshuo843/12306_server/utils"
)

// PassengerController Operations about object
type PassengerController struct {
	BaseController
}

// @Title Get
// @Description 获取乘客列表
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router / [get]
func (p *PassengerController) Get() {
	//获取乘客列表
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	request.SetHeader("X-Requested-With","XMLHttpRequest")
	isGetOk, passengerData := request.SetURL(beego.AppConfig.String("12306::URLGetPassgener")).Get()
	beego.Info("检测登录用户 -----> %t", isGetOk)
	if !isGetOk {
		p.Fail().SetMsg(passengerData).Send()
		return
	}
	var reData map[string]interface{}
	json.Unmarshal([]byte(passengerData), &reData)
	p.Success().SetData(reData).Send()
}