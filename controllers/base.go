package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/chuanshuo843/12306_server/utils"
)

type _ResData struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// BaseController Operations about Users
type BaseController struct {
	beego.Controller
	req *utils.Request
}

var _res _ResData

func init() {
	_res.Status = true
	_res.Message = "success"
	_res.Data = ""
}

// Prepare .
// func (b *BaseController) Prepare() {
// 	authString := b.Ctx.Input.Header("Authorization")
// 	token := utils.Token(authString)
// 	if token != "" {
// 		req := kyfw.Load(token)
// 		if req != nil {
// 			b.req = req
// 		}
// 	}
// }

// Success .
func (b *BaseController) Success() *BaseController {
	_res.Status = true
	return b
}

// SetMsg .
func (b *BaseController) SetMsg(message string) *BaseController {
	_res.Message = message
	return b
}

// Fail .
func (b *BaseController) Fail() *BaseController {
	_res.Status = false
	return b
}

// SetData .
func (b *BaseController) SetData(data interface{}) *BaseController {
	_res.Data = data
	return b
}

// Send .
func (b *BaseController) Send() {
	json_data, _ := json.Marshal(_res)
	b.Data["json"] = string(json_data)
	//初始化数据
	_res.Status = true
	_res.Message = "success"
	_res.Data = ""
	b.ServeJSON()
}
