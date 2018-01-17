package controllers

import "encoding/json"

// "encoding/json"
// "fmt"
// "strings"
// "github.com/astaxie/beego"

// ScheduleController Operations about object
type ScheduleController struct {
	BaseController
}

// @router /init [get]
func (s *ScheduleController) InitQuery() {
	req := s.req()
	err := kyfwQuery.InitQuerySchedule(req)
	if err != nil {
		s.Fail().SetMsg(err.Error()).Send()
		return
	}
	s.Success().Send()
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router / [get]
func (s *ScheduleController) Get() {
	req := s.tokenReq()
	startStation := s.GetString("start_station")
	endStation := s.GetString("end_station")
	startCode := s.GetString("start_code")
	endCode := s.GetString("end_code")
	date := s.GetString("date")
	if startStation == "" || endStation == "" || date == "" {
		s.Fail().SetMsg("请选择正确的站台信息").Send()
		return
	}
	data, err := kyfwQuery.GetSchedule(req, startStation, endStation, startCode, endCode, date)
	if err != nil {
		s.Fail().SetMsg(err.Error()).Send()
		return
	}
	var reData map[string]interface{}
	json.Unmarshal([]byte(data), &reData)
	s.Success().SetData(reData).Send()

	//"2018-01-08", "BJP", "CDW"
	// strartStation := s.GetString("out")
	// endStation := s.GetString("in")
	// date := s.GetString("date")
	// if strartStation == "" || endStation == "" || date == "" {
	// 	s.Fail().SetMsg("请选择正确的站台信息").Send()
	// 	return
	// }
	// //查询车次Init
	// request.SetHeader("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	// isInitOk, initData := request.SetURL(beego.AppConfig.String("12306::URLLTrafficInquiryInit")).Get()
	// beego.Info("获取Cookies ----->  ", isInitOk)
	// if !isInitOk {
	// 	s.Fail().SetMsg("车次查询初始化失败").Send()
	// 	return
	// }
	// //查询车次日志
	// request.SetHeader("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	// request.SetHeader("X-Requested-With", "XMLHttpRequest")

	// isLogOk, _ := request.SetURL(fmt.Sprintf(beego.AppConfig.String("12306::URLTrafficInquiryLog"), date, strartStation, endStation)).Get()
	// beego.Info("车次查询init ----->  ", isLogOk)
	// if !isLogOk {
	// 	s.Fail().SetMsg("车次查询日志调用失败").Send()
	// 	return
	// }

	// //查询信息
	// splData := strings.Split(initData,"\n")
	// queryStr := []byte(splData[13])[23:40]
	// queryIsOk, queryData := request.SetURL(fmt.Sprintf(beego.AppConfig.String("12306::URLTrafficInquiry"), queryStr, date, strartStation, endStation)).Get()
	// beego.Info("查询车次信息 ----->  ", queryIsOk)
	// if !queryIsOk {
	// 	s.Fail().SetMsg("车次查询失败了").Send()
	// }
	// var reData map[string]interface{}
	// json.Unmarshal([]byte(queryData), &reData)
	// s.Success().SetData(reData).Send()
}
