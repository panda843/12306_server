package controllers

import (
	// "strings"
	// "github.com/astaxie/beego"
	"github.com/chuanshuo843/12306_server/utils/kyfw"
)

var (
	kyfwQuery kyfw.Query
)

// StationController Operations about object
type StationController struct {
	BaseController
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router / [get]
func (s *StationController) Get() {
	data,err := kyfwQuery.GetStations()
	if err != nil {
		s.Fail().SetMsg(err.Error()).Send()
		return
	}
	s.Ctx.Output.Header("Cache-Control:", "public,max-age=43200")
	s.Success().SetData(data).Send()

	// request := &utils.Request{}
	// isOk, data := request.SetURL(beego.AppConfig.String("12306::URLGetStationList")).IsDisableHeader(false).Get()
	// beego.Info("获取站台信息 ----->  ", isOk)
	// if !isOk {
	// 	s.Fail().SetMsg("站台信息获取失败").Send()
	// 	return
	// }
	// formatData := strings.Split(data, "'")
	// if len(formatData) != 3 {
	// 	s.Fail().SetMsg("获取车站信息失败").Send()
	// 	return
	// }
	// //缓存12个小时
	// s.Ctx.Output.Header("Cache-Control:", "public,max-age=43200")
	// s.Success().SetData(string(formatData[1])).Send()
}
