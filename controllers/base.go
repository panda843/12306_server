package controllers

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
	"github.com/chuanshuo843/12306_server/utils"
	"github.com/chuanshuo843/12306_server/utils/kyfw"
)

type _ResData struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// BaseController Operations about Users
type BaseController struct {
	beego.Controller

	UserReq *utils.Request
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

// req .
func (b *BaseController) req() *utils.Request {
	cookie := b.Ctx.Input.Cookie(beego.BConfig.WebConfig.Session.SessionName)
	if cookie == "" {
		cookie = b.Ctx.Input.CruSession.SessionID()
	}
	req := kyfw.Load(cookie)

	return req
}

// tokenReq .
func (b *BaseController) tokenReq() *utils.Request {
	//只检测OPTIONS以外的请求
	if !b.Ctx.Input.Is("OPTIONS") {
		authString := b.Ctx.Input.Header("Authorization")
		if authString == "" {
			return nil
		}
		kv := strings.Split(authString, " ")
		if len(kv) != 2 || kv[0] != "Bearer" {
			return nil
		}
		token := kv[1]

		return kyfw.Load(token)
	}

	return nil
}

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
