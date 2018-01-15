package controllers

import (
	// "encoding/json"
	// "fmt"
	// "math/rand"
	// "net/url"
	"time"

	"net/http"
	// "github.com/astaxie/beego"
	"github.com/chuanshuo843/12306_server/utils"
	"github.com/chuanshuo843/12306_server/utils/kyfw"
)

var (
	kyfwUser kyfw.User
)

// var request utils.Request

// //{"result_message":"验证码校验成功","result_code":"4"}
// type _VerifyRes struct {
// 	ResultMessage string `json:"result_message"`
// 	ResultCode    string `json:"result_code"`
// }

// //{"result_message":"登录成功","result_code":0,"uamtk":"tnRPMlCjrDGm3k5IbzlRKQrbmnKToZC_8WN4ePn32Mkhuc1c0"}
// type _LoginRes struct {
// 	ResultMessage string `json:"result_message"`
// 	ResultCode    int    `json:"result_code"`
// 	Uamtk         string `json:"uamtk"`
// }

// // {"result_message":"验证通过","result_code":0,"apptk":null,"newapptk":"P5e8H_FPPq-Q6kfa9uUsKC0PUdOyqGtE6OSTPKvol9Qhuc1c0"}
// type _UaMtkRes struct {
// 	ResultCode    int    `json:"result_code"`
// 	ResultMessage string `json:"result_message"`
// 	AppTk         string `json:"apptk"`
// 	NewAppTK      string `json:"newapptk"`
// }

// //{"apptk":"6fgxwb7avXwqubqIZr5kHbmHZY2wxV2RqUjDkX0xs8Etyc2c0","result_code":0,"result_message":"验证通过","username":"YouName"}
// type _AuthOk struct {
// 	ResultCode    int    `json:"result_code"`
// 	ResultMessage string `json:"result_message"`
// 	AppTk         string `json:"apptk"`
// 	UserName      string `json:"username"`
// }

// Operations about Users
type UserController struct {
	BaseController
}

func (u *UserController) InitLogin() {
	_, err := kyfwUser.InitLogin()
	if err != nil {
		u.Fail().SetMsg(err.Error()).Send()
		return
	}
	u.Success().Send()
}

//登录12306
func (u *UserController) Login() {
	verify := u.GetString("verify")
	username := u.GetString("username")
	password := u.GetString("password")
	// key := u.GetString("key")
	errLogin := kyfwUser.Login(username, password, verify)
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
	jwt.Payload.Data = `{"username":"` + kyfwUser.UserName + `"}`
	token := jwt.Encode()
	reJson := map[string]string{"access_token": token}
	u.Success().SetMsg("登录成功").SetData(reJson).Send()

	// //获取基本信息
	// verify := u.GetString("verify")
	// username := u.GetString("username")
	// password := u.GetString("password")
	// // key := u.GetString("key")
	// //检测验证码
	// data := &url.Values{}
	// data.Set("answer", verify)
	// data.Set("login_site", `E`)
	// data.Set("rand", `sjrand`)
	// request.SetHeader("Referer", "https://kyfw.12306.cn/otn/login/init")
	// request.SetHeader("X-Requested-With","XMLHttpRequest")
	// isInit, checkd := request.SetURL(beego.AppConfig.String("12306::URLCheckVerifyCode")).Post(data)
	// beego.Info("验证码检测 ----->  ", isInit, checkd)
	// if !isInit {
	// 	u.Fail().SetMsg("验证码检测调用失败").Send()
	// 	return
	// }
	// verifyRes := &_VerifyRes{}
	// errJsonVerify := json.Unmarshal([]byte(checkd), verifyRes)
	// if errJsonVerify != nil {
	// 	u.Fail().SetMsg("验证码检测解析失败").Send()
	// 	return
	// }

	// if verifyRes.ResultCode != "4" {
	// 	u.Fail().SetMsg(verifyRes.ResultMessage).Send()
	// 	return
	// }

	// //用户登录
	// loginData := &url.Values{}
	// loginData.Set("username", username)
	// loginData.Set("password", password)
	// loginData.Set("appid", `otn`)
	// request.SetHeader("Referer", "https://kyfw.12306.cn/otn/login/init")
	// request.SetHeader("X-Requested-With","XMLHttpRequest")
	// isLogin, login := request.SetURL(beego.AppConfig.String("12306::URLGetUserLogin")).Post(loginData)
	// beego.Info("用户登录调用 ----->  ", isLogin,login)
	// if !isLogin {
	// 	u.Fail().SetMsg("用户登录调用失败").Send()
	// 	return
	// }
	// loginRes := &_LoginRes{}
	// errJsonLogin := json.Unmarshal([]byte(login), loginRes)
	// if errJsonLogin != nil {
	// 	u.Fail().SetMsg("用户登录解析失败").Send()
	// 	return
	// }
	// if loginRes.ResultCode != 0 {
	// 	u.Fail().SetMsg(loginRes.ResultMessage).Send()
	// 	return
	// }

	// //检测用户是否登录
	// uaMtkData := &url.Values{}
	// uaMtkData.Set("appid", "otn")
	// request.SetHeader("Referer", "https://kyfw.12306.cn/otn/passport?redirect=/otn/login/userLogin")
	// isUamtk, uamtk := request.SetURL(beego.AppConfig.String("12306::URLGetUAMTK")).Post(uaMtkData)
	// beego.Info("用户登录检测调用 ----->  ", isUamtk)
	// if !isUamtk {
	// 	u.Fail().SetMsg("用户登录检测调用失败").Send()
	// 	return
	// }
	// uaMtkRes := &_UaMtkRes{}
	// errJsonMtk := json.Unmarshal([]byte(uamtk), uaMtkRes)
	// if errJsonMtk != nil {
	// 	u.Fail().SetMsg("用户登录检测解析失败").Send()
	// }
	// //获取登录用户信息
	// uaAuthData := &url.Values{}
	// uaAuthData.Set("tk", uaMtkRes.NewAppTK)
	// request.SetHeader("Referer", "https://kyfw.12306.cn/otn/passport?redirect=/otn/login/userLogin")
	// isUaAuth, uaAuthPage := request.SetURL("https://kyfw.12306.cn/otn/uamauthclient").Post(uaAuthData)
	// beego.Info("获取登录信息调用 ----->  ", isUaAuth)
	// if !isUaAuth {
	// 	u.Fail().SetMsg("获取登录信息调用失败").Send()
	// 	return
	// }
	// authInfo := &_AuthOk{}
	// errJsonAuth := json.Unmarshal([]byte(uaAuthPage), authInfo)
	// if errJsonAuth != nil {
	// 	u.Fail().SetMsg("解析登录信息失败").Send()
	// 	return
	// }
	// //生成JWT
	// jwt := &utils.Jwt{}
	// jwt.Secret.Key = beego.AppConfig.String("JwtKey")
	// jwt.Payload.Jti = time.Now().Unix()
	// jwt.Payload.Iat = time.Now().Unix()
	// jwt.Payload.Nbf = time.Now().Unix()
	// jwt.Payload.Exp = time.Now().Unix()+70000
	// jwt.Payload.Data = `{"username":"`+authInfo.UserName+`"}`
	// token := jwt.Encode()
	// reJson := map[string]string{"access_token": token}
	// u.Success().SetMsg("登录成功").SetData(reJson).Send()
}

//获取12306登录验证码
func (u *UserController) VerifyCode() {
	//初始化登录页面
	// _, errInit := kyfwUser.InitLogin()
	// if errInit != nil {
	// 	http.Error(u.Ctx.ResponseWriter, "Not Found", 404)
	// 	return
	// }
	//获取验证码
	data, errVer := kyfwUser.GetVerifyImages()
	if errVer != nil {
		http.Error(u.Ctx.ResponseWriter, "Not Found", 404)
		return
	}
	u.Ctx.Output.ContentType("png")
	u.Ctx.Output.Body(data)

	//登录页面初始化
	// request.SetHeader("Referer", "https://kyfw.12306.cn/otn/index/init")
	// isInit, _ := request.SetURL(beego.AppConfig.String("12306::URLGetLoginInit")).Get()
	// beego.Info("登录页面init ----->  ", isInit)
	// if !isInit {
	// 	http.Error(u.Ctx.ResponseWriter, "Not Found", 404)
	// 	return
	// }
	// //检测用户是否登录
	// data := &url.Values{}
	// data.Set("appid", "otn")
	// request.SetHeader("Referer", "https://kyfw.12306.cn/otn/login/init")
	// isUamtk, _ := request.SetURL(beego.AppConfig.String("12306::URLGetUAMTK")).Post(data)
	// beego.Info("用户登录检测调用 ----->  ", isUamtk)
	// if !isInit {
	// 	http.Error(u.Ctx.ResponseWriter, "Not Found", 404)
	// 	return
	// }
	// //获取验证码
	// request.SetHeader("Referer", "https://kyfw.12306.cn/otn/login/init")
	// isGet, dataImg := request.SetURL(fmt.Sprintf(beego.AppConfig.String("12306::URLGetLoginCodeImg"), rand.Float64())).Download()
	// beego.Info("验证码下载 ----->  ", isGet)
	// if !isGet {
	// 	http.Error(u.Ctx.ResponseWriter, "Not Found", 404)
	// 	return
	// }
	// u.Ctx.Output.ContentType("png")
	// u.Ctx.Output.Body(dataImg)
}
