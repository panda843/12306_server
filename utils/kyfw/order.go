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
)

type Order struct {
	Base
	SubmitToken      string
	KeyCheckIsChange string
	LeftTicketStr    string
}

//下单
func (order *Order) PlaceAnOrder(secret,trainNo,trainCode,seatType,start, end, date,formatDate, ticketStr, passengerStr string) error {
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
	//获取排队信息
	errQueue := order.GetOrderTicketQueueInfo(formatDate,trainNo,trainCode,seatType,start,end)
	if errQueue != nil {
		return errQueue
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
	beego.Debug(string(data))
	if errSend != nil {
		return err
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
	if len(ticketSource) == 0 {
		return errors.New("获取leftTicketStr失败")
	}
	if len(ticketSource[0]) == 0 {
		return errors.New("获取leftTicketStr失败")
	}
	ticketStr := strings.Split(ticketSource[0][0],`'`)[3]
	order.LeftTicketStr = ticketStr
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
	beego.Debug(string(data))
	if errSend != nil {
		return nil, errSend
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
	return data, nil
}

func (order *Order) GetOrderTicketQueueInfo(trainDate,trainNo,trainCode,seatType,startCode,endCode string) error{
	params := fmt.Sprintf("train_date=%s&train_no=%s&stationTrainCode=%s&seatType=%s&fromStationTelecode=%s"+
	"&toStationTelecode=%s&leftTicket=%s&purpose_codes=00&train_location=PA&_json_att=&REPEAT_SUBMIT_TOKEN=%s",
	trainDate,trainNo,trainCode,seatType,startCode,endCode,order.LeftTicketStr,order.SubmitToken)
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
	var queRes map[string]interface{}
	errJson := json.Unmarshal(data,&queRes)
	if errJson != nil {
		return errJson
	}
	if queRes["status"].(bool) != true {
		return errors.New("getCount Fail")
	}
	return nil
}
