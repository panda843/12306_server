package kyfw

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/astaxie/beego"
	"github.com/chuanshuo843/12306_server/utils"
)

var (
	//提交订单
	OrderSubmitOrderURL = "https://kyfw.12306.cn/otn/leftTicket/submitOrderRequest"
	//初始化订单页面
	OrderInitOrderURL = "https://kyfw.12306.cn/otn/confirmPassenger/initDc"
	//检测订单
	OrderCheckedURL = "https://kyfw.12306.cn/otn/confirmPassenger/checkOrderInfo"
)

type Order struct {
	Base
	Token string
}

//下单
func (order *Order) PlaceAnOrder(request *utils.Request, secret, start, end, date string) error {
	//提交订单
	errSub := order.SubmitOrder(request, secret, start, end, date)
	if errSub != nil {
		return errSub
	}
	//初始化订单确认页面
	errInit := order.InitConfirmOrder(request)
	if errInit != nil {
		return errInit
	}
	_, errCheckd := order.CheckConfirmOrder(request, "", "")
	if errCheckd != nil {
		return errCheckd
	}
	return nil
}

//提交订单
func (order *Order) SubmitOrder(request *utils.Request, secret, start, end, date string) error {
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
func (order *Order) InitConfirmOrder(request *utils.Request) error {
	err := request.CreateHttpRequest(OrderInitOrderURL, "POST", &url.Values{"_json_att": {""}})
	if err != nil {
		return err
	}
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/leftTicket/init")
	data, errSend := request.Send()
	if errSend != nil {
		return errSend
	}
	splData := strings.Split(string(data), "\n")
	beego.Debug(string(splData[11]))
	if len(splData) > 64 {
		order.Token = string([]byte(splData[11])[32:64])
		return nil
	}

	return errors.New("获取订单token出错")

}

//检测订单
func (order *Order) CheckConfirmOrder(request *utils.Request, ticketStr, passengerStr string) ([]byte, error) {
	paramsMap := &url.Values{
		"cancel_flag":         {"2"},
		"bed_level_order_num": {"000000000000000000000000000000"},
		"passengerTicketStr":  {ticketStr},
		"oldPassengerStr":     {passengerStr},
		"tour_flag":           {"dc"},
		"randCode":            {""},
		"whatsSelect":         {"1"},
		"_json_att":           {""},
		"REPEAT_SUBMIT_TOKEN": {order.Token},
	}
	err := request.CreateHttpRequest(OrderCheckedURL, "POST", paramsMap)
	if err != nil {
		return nil, err
	}
	request.SetHeader("Referer", "https://kyfw.12306.cn/otn/confirmPassenger/initDc")
	request.SetHeader("X-Requested-With", "XMLHttpRequest")
	data, errSend := request.Send()
	if errSend != nil {
		return nil, errSend
	}
	var checkdRes map[string]interface{}
	errJson := json.Unmarshal(data, &checkdRes)
	if errJson != nil {
		return nil, errJson
	}
	if checkdRes["status"].(bool) != true {
		return nil, errors.New(string(data))
	}
	return data, nil
}
