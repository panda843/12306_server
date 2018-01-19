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
	err := s.Kyfw.InitQuerySchedule()
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
	startStation := s.GetString("start_station")
	endStation := s.GetString("end_station")
	startCode := s.GetString("start_code")
	endCode := s.GetString("end_code")
	date := s.GetString("date")
	if startStation == "" || endStation == "" || date == "" {
		s.Fail().SetMsg("请选择正确的站台信息").Send()
		return
	}
	data, err := s.Kyfw.GetSchedule(startStation, endStation, startCode, endCode, date)
	if err != nil {
		s.Fail().SetMsg(err.Error()).Send()
		return
	}
	var reData map[string]interface{}
	json.Unmarshal([]byte(data), &reData)
	s.Success().SetData(reData).Send()
}
