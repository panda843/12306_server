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
	data, err := p.Kyfw.GetPassenger()
	if err != nil {
		p.Fail().SetMsg(err.Error()).Send()
		return
	}
	var reData map[string]interface{}
	json.Unmarshal([]byte(data), &reData)
	p.Success().SetData(reData).Send()
}
