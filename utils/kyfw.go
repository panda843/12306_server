package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

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
	//提交订单
	OrderSubmitOrderURL = "https://kyfw.12306.cn/otn/leftTicket/submitOrderRequest"
	//初始化订单页面
	OrderInitOrderURL = "https://kyfw.12306.cn/otn/confirmPassenger/initDc"
	//检测订单
	OrderCheckedURL = "https://kyfw.12306.cn/otn/confirmPassenger/checkOrderInfo"
	//获取车票数和排队人数
	OrderGetCountURL = "https://kyfw.12306.cn/otn/confirmPassenger/getQueueCount"
	//加入购买队列
	OrderJoinBuyQueue = "https://kyfw.12306.cn/otn/confirmPassenger/confirmSingleForQueue"
	//获取取票码
	OrderGetTicketCode = "https://kyfw.12306.cn/otn/confirmPassenger/queryOrderWaitTime?random=1516252521719&tourFlag=dc&_json_att=&REPEAT_SUBMIT_TOKEN=%s"
	//获取站台信息
	QueryGetStation = "https://kyfw.12306.cn/otn/resources/js/framework/station_name.js?station_version=1.9042"
	//查询车次初始化
	QueryScheduleInit = "https://kyfw.12306.cn/otn/leftTicket/init"
	//车次查询日志
	QueryScheduleLog = "https://kyfw.12306.cn/otn/leftTicket/log?leftTicketDTO.train_date=%s&leftTicketDTO.from_station=%s&leftTicketDTO.to_station=%s&purpose_codes=ADULT"
	//查询车次信息
	QuerySchedule = "https://kyfw.12306.cn/otn/%s?leftTicketDTO.train_date=%s&leftTicketDTO.from_station=%s&leftTicketDTO.to_station=%s&purpose_codes=ADULT"
	//查询乘客信息
	QueryPassenger = "https://kyfw.12306.cn/otn/confirmPassenger/getPassengerDTOs"
)

type Kyfw struct {
	Request               *Request
	IsLogin               bool
	LoginName             string
	LoginToken            string
	QueryScheduleURL      string
	OrderSubmitToken      string
	OrderKeyCheckIsChange string
	OrderLeftTicketStr    string
	OrderTrainLocation    string
	OrderTicketCode       string
	LoginNumber           int
	CheckLoginNumber      int
	QueryNumber           int
	OrderSubmitNumber     int
	OrderCheckdNumber     int
	OrderGetCountNumber   int
	OrderJoinBuyNumber    int
	OrderQueryWaitNumber  int
	MaxLoopNumber         int
}

type KyfwList struct {
	UserKyfw *sync.Map
}

func InitKyfwList() *KyfwList {
	return &KyfwList{
		UserKyfw: &sync.Map{},
	}
}

func (kyfws *KyfwList) Create(key string) *Kyfw {
	kyf := &Kyfw{
		Request:       InitRequest(),
		MaxLoopNumber: 10,
		IsLogin:       false,
	}
	kyfws.Set(key, kyf)
	return kyf
}

func (kyfws *KyfwList) Set(key string, ky *Kyfw) {
	kyfws.UserKyfw.Store(kyfws.GetKeyHash(key), ky)
}

func (kyfws *KyfwList) Get(key string) *Kyfw {
	v, ok := kyfws.UserKyfw.Load(kyfws.GetKeyHash(key))
	if !ok {
		return nil
	}
	return v.(*Kyfw)
}

func (kyfws *KyfwList) GetKeyHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (kyfws *KyfwList) Move(forKey, toKey string) error {
	kyfws.Set(toKey, kyfws.Get(forKey))
	kyfws.Delete(forKey)
	return nil
}

func (kyfws *KyfwList) Foreach() {
	kyfws.UserKyfw.Range(func(k, v interface{}) bool {
		fmt.Println(k.(string), v)
		return true
	})
}

func (kyfws *KyfwList) Delete(key string) {
	kyfws.UserKyfw.Delete(kyfws.GetKeyHash(key))
}

//登录页面初始化
func (kyfw *Kyfw) InitLogin() ([]byte, error) {
	err := kyfw.Request.CreateHttpRequest(UserLoginInit, "GET", nil)
	if err != nil {
		return nil, err
	}
	kyfw.Request.SetHeader("Referer", "https://kyfw.12306.cn/otn/index/init")
	return kyfw.Request.Send()
}

//获取验证码
func (kyfw *Kyfw) GetVerifyImages() ([]byte, error) {
	err := kyfw.Request.CreateHttpRequest(fmt.Sprintf(UserGetVerifyImg, rand.Float64()), "GET", nil)
	if err != nil {
		return nil, err
	}
	kyfw.Request.SetHeader("Referer", "https://kyfw.12306.cn/otn/login/init")
	return kyfw.Request.Send()
}

//获取12306登录token
func (kyfw *Kyfw) Get12306Token(appToken string) error {
	err := kyfw.Request.CreateHttpRequest(UserGetToken, "POST", &url.Values{"tk": {appToken}})
	if err != nil {
		return err
	}
	kyfw.Request.SetHeader("Referer", "https://kyfw.12306.cn/otn/passport?redirect=/otn/login/userLogin")
	data, errSend := kyfw.Request.Send()
	beego.Debug("Get12306Token:", string(data))
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
		kyfw.IsLogin = false
		return errors.New(tokenRes["result_message"].(string))
	}
	kyfw.LoginName = tokenRes["username"].(string)
	kyfw.LoginToken = tokenRes["apptk"].(string)
	kyfw.IsLogin = true
	return nil
}

//检测用户是否登录
func (kyfw *Kyfw) CheckIsLogin() (string, error) {
	err := kyfw.Request.CreateHttpRequest(UserAuthUAMTK, "POST", &url.Values{"appid": {"otn"}})
	if err != nil {
		return "", err
	}
	kyfw.Request.SetHeader("Referer", "https://kyfw.12306.cn/otn/login/init")
	data, sendErr := kyfw.Request.Send()
	beego.Debug(string(data))
	if sendErr != nil {
		return "", sendErr
	}
	//没有取到数据时递归调用直到取到数据为止
	if len(data) == 0 {
		if kyfw.CheckLoginNumber <= kyfw.MaxLoopNumber {
			return kyfw.CheckIsLogin()
		}
		return "", errors.New("获取检测登录数据失败")
	}
	kyfw.CheckLoginNumber = 0
	//解析返回数据
	var checkRes map[string]interface{}
	errJosn := json.Unmarshal(data, &checkRes)
	if errJosn != nil {
		return "", errJosn
	}
	// {"result_message":"验证通过","result_code":0,"apptk":null,"newapptk":"P5e8H_FPPq-Q6kfa9uUsKC0PUdOyqGtE6OSTPKvol9Qhuc1c0"}
	if checkRes["result_code"].(float64) != 0 {
		kyfw.IsLogin = false
		return "", errors.New(checkRes["result_message"].(string))
	}
	kyfw.IsLogin = true
	return checkRes["newapptk"].(string), nil
}

//用户登录
func (kyfw *Kyfw) Login(username, password, verify string) error {
	//检测验证码
	errVer := kyfw.CheckVerifyCode(verify)
	if errVer != nil {
		return errVer
	}
	//登录12306
	errLogin := kyfw.Login12306(username, password)
	if errLogin != nil {
		return errLogin
	}
	//检测用户是否登录
	appTk, errCheck := kyfw.CheckIsLogin()
	if errCheck != nil {
		return errCheck
	}
	//获取用户Token
	errTk := kyfw.Get12306Token(appTk)
	if errTk != nil {
		return errTk
	}
	return nil
}

//登录12306
func (kyfw *Kyfw) Login12306(username, password string) error {
	kyfw.LoginNumber++
	err := kyfw.Request.CreateHttpRequest(UserLogin12306, "POST", &url.Values{"username": {username}, "password": {password}, "appid": {"otn"}})
	if err != nil {
		return err
	}
	kyfw.Request.SetHeader("Referer", "https://kyfw.12306.cn/otn/login/init")
	kyfw.Request.SetHeader("X-Requested-With", "XMLHttpRequest")
	data, errSend := kyfw.Request.Send()
	if errSend != nil {
		return errSend
	}
	//没有取到数据时递归调用直到取到数据为止
	if len(data) == 0 {
		if kyfw.LoginNumber <= kyfw.MaxLoopNumber {
			return kyfw.Login12306(username, password)
		}
		return errors.New("获取登录数据失败")
	}
	kyfw.LoginNumber = 0
	var loginRes map[string]interface{}
	errJson := json.Unmarshal(data, &loginRes)
	if errJson != nil {
		return errJson
	}
	//{"result_message":"登录成功","result_code":0,"uamtk":"tnRPMlCjrDGm3k5IbzlRKQrbmnKToZC_8WN4ePn32Mkhuc1c0"}
	if loginRes["result_code"].(float64) != 0 {
		return errors.New(loginRes["result_message"].(string))
	}
	return nil
}

//检测验证码
func (kyfw *Kyfw) CheckVerifyCode(verifyCode string) error {
	err := kyfw.Request.CreateHttpRequest(UserCheckVerify, "POST", &url.Values{"answer": {verifyCode}, "login_site": {"E"}, "rand": {"sjrand"}})
	if err != nil {
		return err
	}
	kyfw.Request.SetHeader("Referer", "https://kyfw.12306.cn/otn/login/init")
	kyfw.Request.SetHeader("X-Requested-With", "XMLHttpRequest")
	data, errSend := kyfw.Request.Send()
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

//获取站台信息 -------------------------------------------------------------------------------------Query
func (kyfw *Kyfw) GetStations() (string, error) {
	err := kyfw.Request.CreateHttpRequest(QueryGetStation, "GET", nil)
	if err != nil {
		return "", err
	}
	data, errSend := kyfw.Request.Send()
	if errSend != nil {
		return "", errSend
	}
	stationList := strings.Split(string(data), "'")
	if len(stationList) != 3 {
		return "", errors.New("获取站台信息失败")
	}
	return string(stationList[1]), nil
}

//查询车次信息Init
func (kyfw *Kyfw) InitQuerySchedule() error {
	//获取查询地址
	err := kyfw.Request.CreateHttpRequest(QueryScheduleInit, "GET", nil)
	if err != nil {
		return err
	}
	initData, _ := kyfw.Request.Send()
	splData := strings.Split(string(initData), "\n")
	kyfw.QueryScheduleURL = string([]byte(splData[13])[23:40])
	return nil
}

//添加车次查询日志
func (kyfw *Kyfw) AddQueryScheduleLog(startStation, endStation, startCode, endCode, date string) error {
	//设置Cookie
	cookies := []*http.Cookie{
		&http.Cookie{Name: "_jc_save_fromDate", Value: date},
		&http.Cookie{Name: "_jc_save_fromStation", Value: url.QueryEscape(startStation + "," + startCode)},
		&http.Cookie{Name: "_jc_save_toDate", Value: date},
		&http.Cookie{Name: "_jc_save_toStation", Value: url.QueryEscape(endStation + "," + endCode)},
		&http.Cookie{Name: "_jc_save_wfdc_flag", Value: "dc"},
		&http.Cookie{Name: "current_captcha_type", Value: "Z"},
		&http.Cookie{Name: "RAIL_DEVICEID", Value: "P-3rV0RG9eaSEhKKfLTH7F5AAyUaGB84osHAAUosnqYAJA-izUyJkRbfB0Cw-UwpQjI_pmxe3quoWQz3tWxrwVOUPA0RBhiYKCqt4hc028fOVEOlM9NUScapq2xqoMQgLFrUmVJ5-Z1f_GVWrSS9MUdfrkvnvqgR"},
		&http.Cookie{Name: "RAIL_EXPIRATION", Value: "1516096637926"},
	}
	kyfw.Request.CreateHttpRequest(fmt.Sprintf(QueryScheduleLog, date, startCode, endCode), "GET", nil)
	kyfw.Request.SetCookie(cookies)
	kyfw.Request.SetHeader("X-Requested-With", "XMLHttpRequest")
	_, err := kyfw.Request.Send()
	if err != nil {
		return errors.New("车次查询日志添加失败")
	}
	return nil
}

func (kyfw *Kyfw) QueryScheduleInfo(date, startCode, endCode string) (string, error) {
	kyfw.QueryNumber++
	//查询车次信息
	errQ := kyfw.Request.CreateHttpRequest(fmt.Sprintf(QuerySchedule, kyfw.QueryScheduleURL, date, startCode, endCode), "GET", nil)
	if errQ != nil {
		return "", errQ
	}
	kyfw.Request.SetHeader("X-Requested-With", "XMLHttpRequest")
	data, errQuery := kyfw.Request.Send()
	if errQuery != nil {
		return "", errQuery
	}
	//没有取到数据时递归调用直到取到数据为止
	if len(data) == 0 {
		if kyfw.QueryNumber <= kyfw.MaxLoopNumber {
			return kyfw.QueryScheduleInfo(date, startCode, endCode)
		}
		return "", errors.New("获取车次信息失败")
	}
	beego.Info("QuerySchedule:", string(data))
	kyfw.QueryNumber = 0
	return string(data), nil
}

//查询车次信息
func (kyfw *Kyfw) GetSchedule(startStation, endStation, startCode, endCode, date string) (string, error) {
	//检测查询URL
	if kyfw.QueryScheduleURL == "" {
		errQueryURL := kyfw.InitQuerySchedule()
		if errQueryURL != nil {
			return "", errQueryURL
		}
	}
	//添加查询日志
	errLog := kyfw.AddQueryScheduleLog(startStation, endStation, startCode, endCode, date)
	if errLog != nil {
		return "", errLog
	}
	//查询车次信息
	return kyfw.QueryScheduleInfo(date, startCode, endCode)
}

//查询乘客信息
func (kyfw *Kyfw) GetPassenger() ([]byte, error) {
	kyfw.Request.CreateHttpRequest(QueryPassenger, "GET", nil)
	kyfw.Request.SetHeader("X-Requested-With", "XMLHttpRequest")
	data, err := kyfw.Request.Send()
	if err != nil {
		return nil, err
	}
	beego.Info("GetPassenger:", string(data))
	//解析返回信息
	var passRes map[string]interface{}
	errJson := json.Unmarshal(data, &passRes)
	if errJson != nil {
		return nil, errJson
	}
	if passRes["status"].(bool) != true {
		return nil, errors.New(passRes["message"].(string))
	}
	passData := passRes["data"].(map[string]interface{})
	if passData["isExist"].(bool) != true {
		kyfw.IsLogin = false
		return nil, errors.New(passData["exMsg"].(string))
	}
	return data, nil
}

func (kyfw *Kyfw) GetSeatPrice(trainNo, startNo, endNo, seatType, date string) {
	//https://kyfw.12306.cn/otn/leftTicket/queryTicketPrice?train_no=5l0000G13061&from_station_no=01&to_station_no=12&seat_types=O9M&train_date=2017-02-05`]
}

//下单 ---------------------------------------------------------------------Order
func (kyfw *Kyfw) PlaceAnOrder(secret, trainNo, trainCode, start, startCode, end, endCode, date, formatDate, ticketStr, passengerStr string) error {
	//提交订单
	errSub := kyfw.SubmitOrder(secret, start, end, date)
	if errSub != nil {
		return errSub
	}
	//初始化订单确认页面
	errInit := kyfw.InitConfirmOrder()
	if errInit != nil {
		return errInit
	}
	//检测订单
	_, errCheckd := kyfw.CheckConfirmOrder(ticketStr, passengerStr)
	if errCheckd != nil {
		return errCheckd
	}
	ticketArr := strings.Split(ticketStr, ",")
	//获取排队信息
	errQueue := kyfw.GetOrderTicketQueueInfo(formatDate, trainNo, trainCode, ticketArr[0], startCode, endCode)
	if errQueue != nil {
		return errQueue
	}
	//加入购买队列
	errJoin := kyfw.JoinBuyTicketQueue(ticketStr, passengerStr)
	if errJoin != nil {
		return errJoin
	}
	//获取取票码
	errTicket := kyfw.GetTicketCode()
	if errTicket != nil {
		return errTicket
	}
	return nil
}

//提交订单
func (kyfw *Kyfw) SubmitOrder(secret, start, end, date string) error {
	kyfw.OrderSubmitNumber++
	params := fmt.Sprintf("secretStr=%s&train_date=%s&back_train_date=%s&tour_flag=dc&purpose_codes=ADULT&"+
		"query_from_station_name=%s&query_to_station_name=%s&undefined=", secret, date,
		date, start, end)
	err := kyfw.Request.CreateHttpRequest(OrderSubmitOrderURL, "POST", params)
	kyfw.Request.SetHeader("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	kyfw.Request.SetHeader("X-Requested-With", "XMLHttpRequest")
	if err != nil {
		return err
	}
	data, errSend := kyfw.Request.Send()
	if errSend != nil {
		return err
	}
	//没有取到数据时递归调用直到取到数据为止
	if len(data) == 0 {
		if kyfw.OrderSubmitNumber <= kyfw.MaxLoopNumber {
			return kyfw.SubmitOrder(secret, start, end, date)
		}
		return errors.New("提交订单失败")
	}
	beego.Info("SubmitOrder:", string(data))
	kyfw.OrderSubmitNumber = 0
	var subRes map[string]interface{}
	errJson := json.Unmarshal(data, &subRes)
	if errJson != nil {
		return errJson
	}
	if subRes["status"].(bool) != true {
		msg := subRes["messages"].([]interface{})
		return errors.New(string(msg[0].(string)))
	}
	return nil
}

//初始化确认订单页面
func (kyfw *Kyfw) InitConfirmOrder() error {
	err := kyfw.Request.CreateHttpRequest(OrderInitOrderURL, "POST", &url.Values{"_json_att": {""}})
	if err != nil {
		return err
	}
	kyfw.Request.SetHeader("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	data, errSend := kyfw.Request.Send()
	if errSend != nil {
		return errSend
	}
	//获取KeyCheck
	var keyCheckRegexp = regexp.MustCompile(`'key_check_isChange':'[\S]+','left`)
	keyCacheSource := keyCheckRegexp.FindAllStringSubmatch(string(data), -1)
	if len(keyCacheSource) == 0 {
		return errors.New("获取key_check_isChange失败")
	}
	if len(keyCacheSource[0]) == 0 {
		return errors.New("获取key_check_isChange失败")
	}
	keyCacheStr := strings.Split(keyCacheSource[0][0], `'`)[3]
	kyfw.OrderKeyCheckIsChange = keyCacheStr
	//获取leftTicketStr
	var ticketRegexp = regexp.MustCompile(`'leftTicketStr':'[\S]+','limit`)
	ticketSource := ticketRegexp.FindAllStringSubmatch(string(data), -1)
	if len(ticketSource) == 0 {
		return errors.New("获取leftTicketStr失败")
	}
	if len(ticketSource[0]) == 0 {
		return errors.New("获取leftTicketStr失败")
	}
	ticketStr := strings.Split(ticketSource[0][0], `'`)[3]
	kyfw.OrderLeftTicketStr = ticketStr
	//获取trian_location
	var locationRegexp = regexp.MustCompile(`'train_location':'[\S]+'};`)
	locationSource := locationRegexp.FindAllStringSubmatch(string(data), -1)
	if len(locationSource) == 0 {
		return errors.New("获取trian_location失败")
	}
	if len(locationSource[0]) == 0 {
		return errors.New("获取trian_location失败")
	}
	locationStr := strings.Split(locationSource[0][0], `'`)[3]
	kyfw.OrderTrainLocation = locationStr
	//获取submitToken
	splData := strings.Split(string(data), "\n")
	if len(splData) > 64 {
		kyfw.OrderSubmitToken = string([]byte(splData[11])[32:64])
		return nil
	} else {
		return errors.New("获取submitToken失败")
	}
	return nil
}

//检测订单
func (kyfw *Kyfw) CheckConfirmOrder(ticketStr, passengerStr string) ([]byte, error) {
	kyfw.OrderCheckdNumber++
	params := fmt.Sprintf("cancel_flag=2&bed_level_order_num=000000000000000000000000000000&passengerTicketStr=%s&oldPassengerStr=%s"+
		"&tour_flag=dc&randCode=&_json_att=&REPEAT_SUBMIT_TOKEN=%s", ticketStr, passengerStr, kyfw.OrderSubmitToken)
	err := kyfw.Request.CreateHttpRequest(OrderCheckedURL, "POST", params)
	if err != nil {
		return nil, err
	}
	kyfw.Request.SetHeader("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	kyfw.Request.SetHeader("X-Requested-With", "XMLHttpRequest")
	data, errSend := kyfw.Request.Send()
	if errSend != nil {
		return nil, errSend
	}
	//没有取到数据时递归调用直到取到数据为止
	if len(data) == 0 {
		if kyfw.OrderCheckdNumber <= kyfw.MaxLoopNumber {
			return kyfw.CheckConfirmOrder(ticketStr, passengerStr)
		}
		return nil, errors.New("获取检测订单数据失败")
	}
	beego.Info("CheckOrder:", string(data))
	kyfw.OrderCheckdNumber = 0
	var checkdRes map[string]interface{}
	errJson := json.Unmarshal(data, &checkdRes)
	if errJson != nil {
		return nil, errJson
	}
	if checkdRes["status"].(bool) != true {
		msg := checkdRes["messages"].([]interface{})
		return nil, errors.New(msg[0].(string))
	}
	checkdSubmit := checkdRes["data"].(map[string]interface{})
	if checkdSubmit["submitStatus"].(bool) != true {
		return nil, errors.New(checkdSubmit["errMsg"].(string))
	}
	return data, nil
}

//获取排队人数和Ticket信息
func (kyfw *Kyfw) GetOrderTicketQueueInfo(trainDate, trainNo, trainCode, seatType, startCode, endCode string) error {
	kyfw.OrderGetCountNumber++
	params := fmt.Sprintf("train_date=%s&train_no=%s&stationTrainCode=%s&seatType=%s&fromStationTelecode=%s"+
		"&toStationTelecode=%s&leftTicket=%s&purpose_codes=00&train_location=%s&_json_att=&REPEAT_SUBMIT_TOKEN=%s",
		url.QueryEscape(trainDate), trainNo, trainCode, seatType, startCode, endCode, kyfw.OrderLeftTicketStr, kyfw.OrderTrainLocation, kyfw.OrderSubmitToken)
	//train_date=Thu+Jan+18+2018+00%3A00%3A00+GMT%2B0800+(%E4%B8%AD%E5%9B%BD%E6%A0%87%E5%87%86%E6%97%B6%E9%97%B4)&train_no=240000K4110T&stationTrainCode=K411&seatType=3&fromStationTelecode=BJP&toStationTelecode=TXP&leftTicket=%252BG1447KdMTkeq4e7Llen0%252BmM502nk7jzf4Re3fL57omLNhBAoYLby8tmSP8%253D&purpose_codes=00&train_location=PA&_json_att=&REPEAT_SUBMIT_TOKEN=f1bea73b1e94d07fab28b4fb63f12616
	err := kyfw.Request.CreateHttpRequest(OrderGetCountURL, "POST", params)
	kyfw.Request.SetHeader("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	kyfw.Request.SetHeader("X-Requested-With", "XMLHttpRequest")
	if err != nil {
		return err
	}
	data, errSend := kyfw.Request.Send()
	if errSend != nil {
		return errSend
	}
	//没有取到数据时递归调用直到取到数据为止
	if len(data) == 0 {
		if kyfw.OrderGetCountNumber <= kyfw.MaxLoopNumber {
			time.Sleep(1 * time.Second)
			return kyfw.GetOrderTicketQueueInfo(trainDate, trainNo, trainCode, seatType, startCode, endCode)
		}
		return errors.New("获取排队人数信息失败")
	}
	beego.Info("QueueCount:", string(data))
	kyfw.OrderGetCountNumber = 0
	var queRes map[string]interface{}
	errJson := json.Unmarshal(data, &queRes)
	if errJson != nil {
		return errJson
	}
	if queRes["status"].(bool) != true {
		msg := queRes["messages"].([]interface{})
		return errors.New(msg[0].(string))
	}
	return nil
}

//加入购买队列
func (kyfw *Kyfw) JoinBuyTicketQueue(ticketStr, passengerStr string) error {
	kyfw.OrderJoinBuyNumber++
	params := fmt.Sprintf("passengerTicketStr=%s&oldPassengerStr=%s&randCode=&purpose_codes=00&key_check_isChange=%s&leftTicketStr=%s&train_location=%s&choose_seats=&seatDetailType=000&whatsSelect=1&roomType=00&dwAll=N&_json_att=&REPEAT_SUBMIT_TOKEN=%s", ticketStr, passengerStr, kyfw.OrderKeyCheckIsChange, kyfw.OrderLeftTicketStr, kyfw.OrderTrainLocation, kyfw.OrderSubmitToken)
	//passengerTicketStr=%s&oldPassengerStr=%s&randCode=&purpose_codes=00&key_check_isChange=%s&leftTicketStr=%s&train_location=%s&choose_seats=&seatDetailType=000&whatsSelect=1&roomType=00&dwAll=N&_json_att=&REPEAT_SUBMIT_TOKEN=%s
	err := kyfw.Request.CreateHttpRequest(OrderJoinBuyQueue, "POST", params)
	if err != nil {
		return err
	}
	kyfw.Request.SetHeader("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	kyfw.Request.SetHeader("X-Requested-With", "XMLHttpRequest")
	data, errSend := kyfw.Request.Send()
	if errSend != nil {
		return errSend
	}
	//没有取到数据时递归调用直到取到数据为止
	if len(data) == 0 {
		if kyfw.OrderJoinBuyNumber <= kyfw.MaxLoopNumber {
			return kyfw.JoinBuyTicketQueue(ticketStr, passengerStr)
		}
		return errors.New("加入购买队列失败")
	}
	beego.Info("BuyQueue:", string(data))
	kyfw.OrderJoinBuyNumber = 0
	var joinRes map[string]interface{}
	errJosn := json.Unmarshal(data, &joinRes)
	if errJosn != nil {
		return errJosn
	}
	if joinRes["status"].(bool) != true {
		return errors.New(joinRes["data"].(string))
	}
	joinData := joinRes["data"].(map[string]interface{})
	if joinData["submitStatus"].(bool) != true {
		return errors.New("submit false")
		//return nil,errors.New(joinData["errMsg"].(string))
	}
	return nil
}

//获取取票码
func (kyfw *Kyfw) GetTicketCode() error {
	kyfw.OrderQueryWaitNumber++
	err := kyfw.Request.CreateHttpRequest(fmt.Sprintf(OrderGetTicketCode, kyfw.OrderSubmitToken), "GET", nil)
	if err != nil {
		return err
	}
	kyfw.Request.SetHeader("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	kyfw.Request.SetHeader("X-Requested-With", "XMLHttpRequest")
	data, errSend := kyfw.Request.Send()
	if errSend != nil {
		return errSend
	}
	//没有取到数据时递归调用直到取到数据为止
	if len(data) == 0 {
		if kyfw.OrderQueryWaitNumber <= kyfw.MaxLoopNumber {
			return kyfw.GetTicketCode()
		}
		return errors.New("下单成功,获取取票码数据失败")
	}
	beego.Info("Wiat:", string(data))
	kyfw.OrderQueryWaitNumber = 0
	var tickRes map[string]interface{}
	errJosn := json.Unmarshal(data, &tickRes)
	if errJosn != nil {
		return errJosn
	}
	if tickRes["status"].(bool) != true {
		msg := tickRes["messages"].([]interface{})
		return errors.New(msg[0].(string))
	}
	tickData := tickRes["data"].(map[string]interface{})
	//检测waitTime是否小于0不小于0时递归调用
	if tickData["waitTime"].(float64) < 0 {
		switch tickData["orderId"].(type) {
		case nil:
			return errors.New(tickData["msg"].(string))
		case string:
			kyfw.OrderTicketCode = tickData["orderId"].(string)
			return nil
		}
	} else {
		return kyfw.GetTicketCode()
	}
	return nil
}
