package controllers

import (
	"encoding/json"
	"fmt"
	"strings"
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
	if strartStation == "" || endStation == "" || date == "" {
		s.Fail().SetMsg("请选择正确的站台信息").Send()
		return
	}
	//查询车次Init
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	isInitOk, initData := request.SetURL(beego.AppConfig.String("12306::URLLTrafficInquiryInit")).Get()
	beego.Info("获取Cookies -----> %t", isInitOk)
	if !isInitOk {
		s.Fail().SetMsg(initData).Send()
		return
	}
	//查询车次日志
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	request.SetHeader("X-Requested-With", "XMLHttpRequest")

	isLogOk, queryLog := request.SetURL(fmt.Sprintf(beego.AppConfig.String("12306::URLTrafficInquiryLog"), date, strartStation, endStation)).Get()
	beego.Info("车次查询init -----> %t", isLogOk)
	if !isLogOk {
		s.Fail().SetMsg(queryLog).Send()
		return
	}

	//查询信息
	splData := strings.Split(initData,"\n")
	queryStr := []byte(splData[13])[23:40]
	queryIsOk, queryData := request.SetURL(fmt.Sprintf(beego.AppConfig.String("12306::URLTrafficInquiry"), queryStr, date, strartStation, endStation)).Get()
	beego.Info("查询车次信息 -----> %t", queryIsOk)
	if !queryIsOk {
		s.Fail().SetMsg("查询失败了").Send()
	}
	var reData map[string]interface{}
	json.Unmarshal([]byte(queryData), &reData)
	s.Success().SetData(reData).Send()
}
