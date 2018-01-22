package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
	"github.com/chuanshuo843/12306_server/utils"
)

type _ResData struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	IsLogin bool        `json:"login"`
}

// Operations about Users
type BaseController struct {
	beego.Controller
	res    _ResData
	Kyfw   *utils.Kyfw
	AppID  string
	UserID string
}

var (
	kyfws *utils.KyfwList
	tasks *utils.TaskList
)

func init() {
	kyfws = utils.InitKyfwList()
	tasks = utils.InitTaskList()
}

func (b *BaseController) Prepare() {
	//初始化返回数据
	b.res.Status = true
	b.res.Message = "success"
	b.res.Data = ""
	//获取用户对应的信息
	b.GetUserKyfw()
	if b.Kyfw == nil {
		b.Fail().SetMsg("登录信息失效,请重新登录").Send()
		return
	}
	hasher := md5.New()
	hasher.Write([]byte(b.UserID))
	b.UserID = hex.EncodeToString(hasher.Sum(nil))
}

// 获取用户数据 .
func (b *BaseController) GetUserKyfw() {
	//Options的不获取
	if b.Ctx.Input.Is("OPTIONS") {
		return
	}
	//获取用户数据,没登录使用APPID,登录使用Token,同时存在则以Token为准
	b.AppID = b.GetString("app_id")
	b.UserID = b.AppID
	if b.AppID != "" {
		b.Kyfw = kyfws.Get(b.AppID)
	}
	authString := b.Ctx.Input.Header("Authorization")
	if authString != "" {
		kv := strings.Split(authString, " ")
		if len(kv) == 2 || kv[0] == "Bearer" {
			b.UserID = kv[1]
			b.Kyfw = kyfws.Get(kv[1])
		}
	}
}

func (b *BaseController) Success() *BaseController {
	b.res.Status = true
	return b
}

func (b *BaseController) SetMsg(message string) *BaseController {
	b.res.Message = message
	return b
}

func (b *BaseController) Fail() *BaseController {
	b.res.Status = false
	return b
}

func (b *BaseController) SetData(data interface{}) *BaseController {
	b.res.Data = data
	return b
}

func (b *BaseController) Send() {
	if b.Kyfw != nil {
		b.res.IsLogin = b.Kyfw.IsLogin
	}
	json_data, _ := json.Marshal(b.res)
	b.Data["json"] = string(json_data)
	//初始化数据
	b.res.Status = true
	b.res.Message = "success"
	b.res.Data = ""
	b.ServeJSON()
}
