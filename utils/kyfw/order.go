package kyfw

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/astaxie/beego"
)

var (
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
)

type Order struct {
	Base
	SubmitToken      string
	KeyCheckIsChange string
	LeftTicketStr    string
	TrainLocation string
}

//下单
func (order *Order) PlaceAnOrder(secret,trainNo,trainCode,start,startCode, end, endCode,date,formatDate, ticketStr, passengerStr string) error {
	//提交订单
	errSub := order.SubmitOrder(secret, start, end, date)
	if errSub != nil {
		return errSub
	}
	//初始化订单确认页面
	errInit := order.InitConfirmOrder()
	if errInit != nil {
		return errInit
	}
	//检测订单
	_, errCheckd := order.CheckConfirmOrder(ticketStr, passengerStr)
	if errCheckd != nil {
		return errCheckd
	}
	ticketArr := strings.Split(ticketStr,",")
	//获取排队信息
	errQueue := order.GetOrderTicketQueueInfo(formatDate,trainNo,trainCode,ticketArr[0],startCode,endCode)
	if errQueue != nil {
		return errQueue
	}
	//加入购买队列
	errJoin := order.JoinBuyTicketQueue(ticketStr,passengerStr)
	if errJoin != nil {
		return errJoin
	}
	//获取取票码
	errTicket := order.GetTicketCode()
	if errTicket != nil {
		return errTicket
	}
	return nil
}

//提交订单
func (order *Order) SubmitOrder(secret, start, end, date string) error {
	params := fmt.Sprintf("secretStr=%s&train_date=%s&back_train_date=%s&tour_flag=dc&purpose_codes=ADULT&"+
		"query_from_station_name=%s&query_to_station_name=%s&undefined=", secret, date,
		date, start, end)
	err := request.CreateHttpRequest(OrderSubmitOrderURL, "POST", params)
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	request.SetHeader("X-Requested-With", "XMLHttpRequest")
	if err != nil {
		return err
	}
	data, errSend := request.Send()
	if errSend != nil {
		return err
	}
	//没有取到数据时递归调用直到取到数据为止
	if len(data) == 0 {
		return order.SubmitOrder(secret, start, end, date)
	}
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
func (order *Order) InitConfirmOrder() error {
	err := request.CreateHttpRequest(OrderInitOrderURL, "POST", &url.Values{"_json_att": {""}})
	if err != nil {
		return err
	}
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	data, errSend := request.Send()
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
	keyCacheStr := strings.Split(keyCacheSource[0][0],`'`)[3]
	order.KeyCheckIsChange = keyCacheStr
	//获取leftTicketStr
	var ticketRegexp = regexp.MustCompile(`'leftTicketStr':'[\S]+','limit`)
	ticketSource := ticketRegexp.FindAllStringSubmatch(string(data), -1)
	beego.Debug(ticketSource)
	if len(ticketSource) == 0 {
		return errors.New("获取leftTicketStr失败")
	}
	if len(ticketSource[0]) == 0 {
		return errors.New("获取leftTicketStr失败")
	}
	ticketStr := strings.Split(ticketSource[0][0],`'`)[3]
	order.LeftTicketStr = ticketStr
	//获取trian_location
	var locationRegexp = regexp.MustCompile(`'train_location':'[\S]+'};`)
	locationSource := locationRegexp.FindAllStringSubmatch(string(data), -1)
	if len(locationSource) == 0 {
		return errors.New("获取trian_location失败")
	}
	if len(locationSource[0]) == 0 {
		return errors.New("获取trian_location失败")
	}
	locationStr := strings.Split(locationSource[0][0],`'`)[3]
	order.TrainLocation = locationStr
	//获取submitToken
	splData := strings.Split(string(data), "\n")
	if len(splData) > 64 {
		order.SubmitToken = string([]byte(splData[11])[32:64])
		return nil
	} else {
		return errors.New("获取submitToken失败")
	}
	return nil
}

//检测订单
func (order *Order) CheckConfirmOrder(ticketStr, passengerStr string) ([]byte, error) {
	params := fmt.Sprintf("cancel_flag=2&bed_level_order_num=000000000000000000000000000000&passengerTicketStr=%s&oldPassengerStr=%s"+
		"&tour_flag=dc&randCode=&_json_att=&REPEAT_SUBMIT_TOKEN=%s", ticketStr, passengerStr, order.SubmitToken)
	err := request.CreateHttpRequest(OrderCheckedURL, "POST", params)
	if err != nil {
		return nil, err
	}
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	request.SetHeader("X-Requested-With", "XMLHttpRequest")
	data, errSend := request.Send()
	if errSend != nil {
		return nil, errSend
	}
	//没有取到数据时递归调用直到取到数据为止
	if len(data) == 0 {
		return order.CheckConfirmOrder(ticketStr, passengerStr)
	}
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
		return nil,errors.New(checkdSubmit["errMsg"].(string))
	}
	return data, nil
}
//获取排队人数和Ticket信息
func (order *Order) GetOrderTicketQueueInfo(trainDate,trainNo,trainCode,seatType,startCode,endCode string) error{
	params := fmt.Sprintf("train_date=%s&train_no=%s&stationTrainCode=%s&seatType=%s&fromStationTelecode=%s"+
	"&toStationTelecode=%s&leftTicket=%s&purpose_codes=00&train_location=%s&_json_att=&REPEAT_SUBMIT_TOKEN=%s",
	url.QueryEscape(trainDate),trainNo,trainCode,seatType,startCode,endCode,order.LeftTicketStr,order.TrainLocation,order.SubmitToken)
	//train_date=Thu+Jan+18+2018+00%3A00%3A00+GMT%2B0800+(%E4%B8%AD%E5%9B%BD%E6%A0%87%E5%87%86%E6%97%B6%E9%97%B4)&train_no=240000K4110T&stationTrainCode=K411&seatType=3&fromStationTelecode=BJP&toStationTelecode=TXP&leftTicket=%252BG1447KdMTkeq4e7Llen0%252BmM502nk7jzf4Re3fL57omLNhBAoYLby8tmSP8%253D&purpose_codes=00&train_location=PA&_json_att=&REPEAT_SUBMIT_TOKEN=f1bea73b1e94d07fab28b4fb63f12616
	err := request.CreateHttpRequest(OrderGetCountURL,"POST",params)
	request.SetHeader("Referer","https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	request.SetHeader("X-Requested-With","XMLHttpRequest")
	if err != nil {
		return err 
	}
	data,errSend := request.Send()
	if errSend != nil {
		return errSend
	}
	//没有取到数据时递归调用直到取到数据为止
	if len(data) == 0 {
		return order.GetOrderTicketQueueInfo(trainDate,trainNo,trainCode,seatType,startCode,endCode)
	}
	var queRes map[string]interface{}
	errJson := json.Unmarshal(data,&queRes)
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
func (order *Order) JoinBuyTicketQueue(ticketStr, passengerStr string) error {
	params := fmt.Sprintf("passengerTicketStr=%s&oldPassengerStr=%s&randCode=&purpose_codes=00&key_check_isChange=%s&leftTicketStr=%s&train_location=%s&choose_seats=&seatDetailType=000&whatsSelect=1&roomType=00&dwAll=N&_json_att=&REPEAT_SUBMIT_TOKEN=%s",ticketStr,passengerStr,order.KeyCheckIsChange,order.LeftTicketStr,order.TrainLocation,order.SubmitToken)
	//passengerTicketStr=%s&oldPassengerStr=%s&randCode=&purpose_codes=00&key_check_isChange=%s&leftTicketStr=%s&train_location=%s&choose_seats=&seatDetailType=000&whatsSelect=1&roomType=00&dwAll=N&_json_att=&REPEAT_SUBMIT_TOKEN=%s
	err := request.CreateHttpRequest(OrderJoinBuyQueue,"POST",params)
	if err != nil {
		return err
	}
	request.SetHeader("Referer","https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	request.SetHeader("X-Requested-With","XMLHttpRequest")
	data,errSend := request.Send()
	if errSend != nil {
		return errSend
	}
	//没有取到数据时递归调用直到取到数据为止
	if len(data) == 0 {
		return order.JoinBuyTicketQueue(ticketStr, passengerStr)
	}
	var joinRes map[string]interface{}
	errJosn := json.Unmarshal(data,&joinRes)
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
func (order *Order) GetTicketCode() error {
	err := request.CreateHttpRequest(fmt.Sprintf(OrderGetTicketCode,order.SubmitToken),"GET",nil)
	if err != nil {
		return err
	}
	request.SetHeader("Referer","https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	request.SetHeader("X-Requested-With","XMLHttpRequest")
	data,errSend := request.Send()
	if errSend != nil {
		return errSend
	}
	//没有取到数据时递归调用直到取到数据为止
	if len(data) == 0 {
		return order.GetTicketCode()
	}
	var tickRes map[string]interface{}
	errJosn := json.Unmarshal(data,&tickRes)
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
			return nil	 
		}
	}else{
		return order.GetTicketCode()
	}
	return nil
}

