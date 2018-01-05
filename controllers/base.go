package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
)

type _ResData struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Operations about Users
type BaseController struct {
	beego.Controller
}

var _res _ResData

func init() {
	_res.Status = true
	_res.Message = "success"
	_res.Data = ""
}

func (b *BaseController) Success(message string) *BaseController {
	_res.Status = true
	if message != "" {
		_res.Message = message
	}
	return b
}

func (b *BaseController) Fail(message string) *BaseController {
	_res.Status = false
	if message != "" {
		_res.Message = message
	}
	return b
}

func (b *BaseController) SetData(data interface{}) *BaseController {
	_res.Data = data
	return b
}

func (b *BaseController) Send() {
	b.AllowCross()
	json_data, _ := json.Marshal(_res)
	b.Data["json"] = string(json_data)
	b.ServeJSON()
}

func (b *BaseController) Options() {
    b.AllowCross() //允许跨域
    b.Data["json"] = ""
    b.ServeJSON()
}

//跨域
func (b *BaseController) AllowCross() {
    b.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")       //允许访问源
    b.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")    //允许post访问
    b.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization") //header的类型
    b.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
    b.Ctx.ResponseWriter.Header().Set("content-type", "application/json") //返回数据格式是json
	b.Ctx.Output.Header("Cache-Control", "no-store")
}
