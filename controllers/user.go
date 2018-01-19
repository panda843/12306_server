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

func (u *UserController) InitLogin() {
	sid := u.Ctx.Input.CruSession.SessionID()
	u.Kyfw = kyfws.Create(sid)
	_, err := u.Kyfw.InitLogin()
	if err != nil {
		u.Fail().SetMsg(err.Error()).Send()
		return
	}
	u.Success().SetData(`{"app_id":"`+sid+`"}`).Send()
}

//登录12306
func (u *UserController) Login() {
	verify := u.GetString("verify")
	username := u.GetString("username")
	password := u.GetString("password")
	// key := u.GetString("key")
	errLogin := u.Kyfw.Login(username, password, verify)
	if errLogin != nil {
		u.Fail().SetMsg(errLogin.Error()).Send()
	}
	//生成JWT
	jwt := &utils.Jwt{}
	jwt.InitJwt()
	jwt.Payload.Jti = time.Now().Unix()
	jwt.Payload.Iat = time.Now().Unix()
	jwt.Payload.Nbf = time.Now().Unix()
	jwt.Payload.Exp = time.Now().Unix() + 70000
	jwt.Payload.Data = `{"username":"` + u.Kyfw.LoginName + `"}`
	token := jwt.Encode()
	kyfws.Move(u.AppID,token)
	reJson := map[string]string{"access_token": token}
	u.Success().SetMsg("登录成功").SetData(reJson).Send()
}

//获取12306登录验证码
func (u *UserController) VerifyCode() {
	//获取验证码
	data, errVer := u.Kyfw.GetVerifyImages()
	if errVer != nil {
		http.Error(u.Ctx.ResponseWriter, "Not Found", 404)
		return 
	}
	u.Ctx.Output.ContentType("png")
	u.Ctx.Output.Body(data)
}
