package controllers

import (
	"github.com/astaxie/beego"
	"encoding/json"
)

type _ResData struct{
	Status bool `json:"status"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

// Operations about Users
type BaseController struct {
	beego.Controller
}

var _res _ResData

func init(){
	_res.Status = true
	_res.Message = "success"
	_res.Data = ""
}

func (b *BaseController) Success(message string) *BaseController{
	_res.Status = true
	if(message != ""){
		_res.Message = message
	}
	return b
}

func (b *BaseController) Fail(message string) *BaseController{
	_res.Status = false
	if(message != ""){
		_res.Message = message
	}
	return b
}

func (b *BaseController) SetData(data interface{}) *BaseController{
	_res.Data = data
	return b
}

func (b *BaseController) Send(){
	json_data,_ := json.Marshal(_res)
	b.Data["json"] = string(json_data)
	b.ServeJSON()
}
