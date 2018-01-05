package controllers

import (
	tools "ganktools.com/12306_server/utils/12306/query"
	"github.com/astaxie/beego"
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
	schedule := &tools.Schedule{}
	is_ok, msg := schedule.Query("2018-01-08", "BJP", "CDW", "ADULT")
	if is_ok {

	}
	beego.Debug(msg)
	s.Success("查询成功").Send()
}
