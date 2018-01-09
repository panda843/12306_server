package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/chuanshuo843/12306_server/utils"
)

// ScheduleController Operations about object
type ScheduleController struct {
	BaseController
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router / [get]
func (s *ScheduleController) Get() {
	request := &utils.Request{}
	//"2018-01-08", "BJP", "CDW"
	strartStation := s.GetString("out")
	endStation := s.GetString("in")
	date := s.GetString("date")
	if strartStation == "" || endStation == "" || date == ""{
		s.Fail().SetMsg("请选择正确的站台信息").Send()
		return
	}
	request.IsDisableHeader(true)
	//获取cookie
	isCookieOk, cookie := request.SetURL(beego.AppConfig.String("12306::URLGetCookie")).Get()
	beego.Info("获取Cookies -----> %t", isCookieOk)
	if !isCookieOk {
		s.Fail().SetMsg(cookie).Send()
		return
	}
	//查询init
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	isOk, iniData := request.SetURL(fmt.Sprintf(beego.AppConfig.String("12306::URLTrafficInquiryInit"), date, strartStation, endStation)).Get()
	beego.Info("车次查询init -----> %t", isOk)
	if !isOk {
		s.Fail().SetMsg(iniData).Send()
		return
	}
	//查询信息
	query_str := string([]byte(cookie)[753:770])
	queryIsOk, queryData := request.SetURL(fmt.Sprintf(beego.AppConfig.String("12306::URLTrafficInquiry"),query_str, date, strartStation, endStation)).Get()
	beego.Info("查询车次信息 -----> %t", queryIsOk)
	if !queryIsOk {
		s.Fail().SetMsg(queryData).Send()
	}
	var reData map[string]interface{}
    json.Unmarshal([]byte(queryData), &reData)
	s.Success().SetData(reData).Send()
}
