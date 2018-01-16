package kyfw

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"

	"github.com/astaxie/beego"
)

const (
	//用户登录Init
	UserLoginInit = "https://kyfw.12306.cn/otn/login/init"
	//检测用户是否登录
	UserAuthUAMTK = "https://kyfw.12306.cn/passport/web/auth/uamtk"
	//获取登录验证码
	UserGetVerifyImg = "https://kyfw.12306.cn/passport/captcha/captcha-image?login_site=E&module=login&rand=sjrand&%g"
	//检测登录验证码
	UserCheckVerify = "https://kyfw.12306.cn/passport/captcha/captcha-check"
	//登录12306
	UserLogin12306 = "https://kyfw.12306.cn/passport/web/login"
	//获取登录信息
	UserGetToken = "https://kyfw.12306.cn/otn/uamauthclient"
)

type User struct {
	Base
	IsLogin  bool
	UserName string
	Token    string
}

//登录页面初始化
func (user *User) InitLogin() ([]byte, error) {
	err := request.CreateHttpRequest(UserLoginInit, "GET", nil)
	if err != nil {
		return nil, err
	}
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/index/init")
	log.Printf("InitLogin = %p \n", request)
	return request.Send()
}

//获取验证码
func (user *User) GetVerifyImages() ([]byte, error) {
	err := request.CreateHttpRequest(fmt.Sprintf(UserGetVerifyImg, rand.Float64()), "GET", nil)
	if err != nil {
		return nil, err
	}
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/login/init")
	log.Printf("GetVerifyImages = %p \n", request)
	return request.Send()
}

//获取12306登录token
func (user *User) Get12306Token(appToken string) error {
	err := request.CreateHttpRequest(UserGetToken, "POST", &url.Values{"tk": {appToken}})
	if err != nil {
		return err
	}
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/passport?redirect=/otn/login/userLogin")
	data, errSend := request.Send()
	if errSend != nil {
		return errSend
	}
	//解析返回数据
	var tokenRes map[string]interface{}
	errJson := json.Unmarshal(data, &tokenRes)
	if errJson != nil {
		return errJson
	}
	//检测操作是否成功
	if tokenRes["result_code"].(float64) != 0 {
		user.IsLogin = false
		return errors.New(tokenRes["result_message"].(string))
	}
	user.UserName = tokenRes["username"].(string)
	user.Token = tokenRes["apptk"].(string)
	user.IsLogin = true
	return nil
}

//检测用户是否登录
func (user *User) CheckIsLogin() (string, error) {
	err := request.CreateHttpRequest(UserAuthUAMTK, "POST", &url.Values{"appid": {"otn"}})
	if err != nil {
		return "", err
	}
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/login/init")
	data, sendErr := request.Send()
	if sendErr != nil {
		return "", sendErr
	}
	if len(data) == 0 {
		return "", errors.New("检测登录返回数据为空")
	}
	//解析返回数据
	var checkRes map[string]interface{}
	errJosn := json.Unmarshal(data, &checkRes)
	if errJosn != nil {
		return "", errJosn
	}
	// {"result_message":"验证通过","result_code":0,"apptk":null,"newapptk":"P5e8H_FPPq-Q6kfa9uUsKC0PUdOyqGtE6OSTPKvol9Qhuc1c0"}
	if checkRes["result_code"].(float64) != 0 {
		user.IsLogin = false
		return "", errors.New(checkRes["result_message"].(string))
	}
	user.IsLogin = true
	return checkRes["newapptk"].(string), nil
}

//用户登录
func (user *User) Login(username, password, verify string) error {
	log.Printf("Login = %p \n", request)
	//检测验证码
	errVer := user.CheckVerifyCode(verify)
	if errVer != nil {
		return errVer
	}
	//登录12306
	errLogin := user.Login12306(username, password)
	if errLogin != nil {
		return errLogin
	}
	//检测用户是否登录
	appTk, errCheck := user.CheckIsLogin()
	if errCheck != nil {
		return errCheck
	}
	//获取用户Token
	errTk := user.Get12306Token(appTk)
	if errTk != nil {
		return errTk
	}

	return nil
}

//登录12306
func (user *User) Login12306(username, password string) error {
	err := request.CreateHttpRequest(UserLogin12306, "POST", &url.Values{"username": {username}, "password": {password}, "appid": {"otn"}})
	if err != nil {
		return err
	}
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/login/init")
	request.SetHeader("X-Requested-With", "XMLHttpRequest")
	data, errSend := request.Send()
	if errSend != nil {
		return errSend
	}
	if len(data) == 0 {
		return errors.New("登录返回数据为空")
	}
	var loginRes map[string]interface{}
	errJson := json.Unmarshal(data, &loginRes)
	if errJson != nil {
		return errJson
	}
	beego.Debug(string(data))
	//{"result_message":"登录成功","result_code":0,"uamtk":"tnRPMlCjrDGm3k5IbzlRKQrbmnKToZC_8WN4ePn32Mkhuc1c0"}
	if loginRes["result_code"].(float64) != 0 {
		return errors.New(loginRes["result_message"].(string))
	}
	return nil
}

//检测验证码
func (user *User) CheckVerifyCode(verifyCode string) error {
	log.Printf("CheckVerifyCode = %p \n", request)
	err := request.CreateHttpRequest(UserCheckVerify, "POST", &url.Values{"answer": {verifyCode}, "login_site": {"E"}, "rand": {"sjrand"}})
	if err != nil {
		return err
	}
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/login/init")
	request.SetHeader("X-Requested-With", "XMLHttpRequest")
	data, errSend := request.Send()
	if errSend != nil {
		return errSend
	}
	var verRes map[string]interface{}
	errJson := json.Unmarshal(data, &verRes)
	if errJson != nil {
		return errJson
	}
	//{"result_message":"验证码校验成功","result_code":"4"}
	if verRes["result_code"].(string) != "4" {
		return errors.New(verRes["result_message"].(string))
	}
	return nil
}
