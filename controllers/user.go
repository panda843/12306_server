package controllers

import (
	"time"

	"net/http"

	"github.com/chuanshuo843/12306_server/utils"
)

// Operations about Users
type UserController struct {
	BaseController
}

func (u *UserController) Prepare() {
	//初始化返回数据
	u.res.Status = true
	u.res.Message = "success"
	u.res.Data = ""
	//获取用户对应的信息
	u.GetUserKyfw()
}

func (u *UserController) InitLogin() {
	sid := u.Ctx.Input.CruSession.SessionID()
	u.Kyfw = kyfws.Create(sid)
	_, err := u.Kyfw.InitLogin()
	if err != nil {
		u.Fail().SetMsg(err.Error()).Send()
		return
	}
	u.Success().SetData(`{"app_id":"` + sid + `"}`).Send()
}

//登录12306
func (u *UserController) Login() {
	if u.Kyfw == nil {
		u.Fail().SetMsg("获取用户数据失败").Send()
		return
	}
	verify := u.GetString("verify")
	username := u.GetString("username")
	password := u.GetString("password")
	// key := u.GetString("key")
	errLogin := u.Kyfw.Login(username, password, verify)
	if errLogin != nil {
		u.Fail().SetMsg(errLogin.Error()).Send()
	}
	//生成JWT
	jwt := utils.InitJwt()
	jwt.Payload.Jti = time.Now().Unix()
	jwt.Payload.Iat = time.Now().Unix() - 30    //减30秒以防请求过快
	jwt.Payload.Nbf = time.Now().Unix() - 30    //减30秒以防请求过快
	jwt.Payload.Exp = time.Now().Unix() + 43200 //有效期十二个小时
	jwt.Payload.Data = `{"username":"` + u.Kyfw.LoginName + `"}`
	token := jwt.Encode()
	kyfws.Move(u.AppID, token)
	u.res.IsLogin = u.Kyfw.IsLogin
	reJson := map[string]string{"access_token": token}
	u.Success().SetMsg("登录成功").SetData(reJson).Send()
}

//获取12306登录验证码
func (u *UserController) VerifyCode() {
	if u.Kyfw == nil {
		u.Fail().SetMsg("获取用户数据失败").Send()
		return
	}
	//获取验证码
	data, errVer := u.Kyfw.GetVerifyImages()
	if errVer != nil {
		http.Error(u.Ctx.ResponseWriter, "Not Found", 404)
		return
	}
	u.Ctx.Output.ContentType("png")
	u.Ctx.Output.Body(data)
}
