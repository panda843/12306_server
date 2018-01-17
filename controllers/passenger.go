package controllers

import "encoding/json"

// "encoding/json"
// "github.com/astaxie/beego"

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
	req := p.tokenReq()
	data, err := kyfwQuery.GetPassenger(req)
	if err != nil {
		p.Fail().SetMsg(err.Error()).Send()
		return
	}
	var reData map[string]interface{}
	json.Unmarshal([]byte(data), &reData)
	p.Success().SetData(reData).Send()
	//获取乘客列表
	// request.SetHeader("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	// request.SetHeader("X-Requested-With","XMLHttpRequest")
	// isGetOk, passengerData := request.SetURL(beego.AppConfig.String("12306::URLGetPassgener")).Get()
	// beego.Info("获取乘客列表 ----->  ", isGetOk)
	// if !isGetOk {
	// 	p.Fail().SetMsg("乘客列表获取失败").Send()
	// 	return
	// }
	// var reData map[string]interface{}
	// json.Unmarshal([]byte(passengerData), &reData)
	// p.Success().SetData(reData).Send()
}
