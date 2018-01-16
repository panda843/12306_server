package kyfw

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
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

type Query struct {
	Base
	ScheduleQueryURL string
}

//获取站台信息
func (query *Query) GetStations() (string, error) {
	err := request.CreateHttpRequest(QueryGetStation, "GET", nil)
	if err != nil {
		return "", err
	}
	data, errSend := request.Send()
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
func (query *Query) InitQuerySchedule() error {
	//获取查询地址
	err := request.CreateHttpRequest(QueryScheduleInit, "GET", nil)
	if err != nil {
		return err
	}
	initData, _ := request.Send()
	splData := strings.Split(string(initData), "\n")
	query.ScheduleQueryURL = string([]byte(splData[13])[23:40])
	return nil
}

//添加车次查询日志
func (query *Query) AddQueryScheduleLog(startStation, endStation, startCode, endCode, date string) error {
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
	request.CreateHttpRequest(fmt.Sprintf(QueryScheduleLog, date, startCode, endCode), "GET", nil)
	request.SetCookie(cookies)
	request.SetHeader("X-Requested-With", "XMLHttpRequest")
	_, err := request.Send()
	if err != nil {
		return errors.New("车次查询日志添加失败")
	}
	return nil
}

//查询车次信息
func (query *Query) GetSchedule(startStation, endStation, startCode, endCode, date string) (string, error) {
	//检测查询URL
	if query.ScheduleQueryURL == "" {
		errQueryURL := query.InitQuerySchedule()
		if errQueryURL != nil {
			return "", errQueryURL
		}
	}
	//添加查询日志
	errLog := query.AddQueryScheduleLog(startStation, endStation, startCode, endCode, date)
	if errLog != nil {
		return "", errLog
	}
	//查询车次信息
	errQ := request.CreateHttpRequest(fmt.Sprintf(QuerySchedule, query.ScheduleQueryURL, date, startCode, endCode), "GET", nil)
	if errQ != nil {
		return "", errQ
	}
	request.SetHeader("X-Requested-With", "XMLHttpRequest")
	data, errQuery := request.Send()
	if len(data) == 0 {
		return "", errors.New("车次查询失败,返回数据为空")
	}
	if errQuery != nil {
		return "", errQuery
	}
	return string(data), nil
}

//查询乘客信息
func (query *Query) GetPassenger() ([]byte, error) {
	request.CreateHttpRequest(QueryPassenger, "GET", nil)
	request.SetHeader("X-Requested-With", "XMLHttpRequest")
	data, err := request.Send()
	if err != nil {
		return nil, err
	}
	//解析返回信息
	var passRes map[string]interface{}
	errJson := json.Unmarshal(data, &passRes)
	if errJson != nil {
		return nil, errJson
	}
	if passRes["status"].(bool) != true {
		return nil, errors.New(passRes["message"].(string))
	}
	return data, nil
}

func (query *Query) GetSeatPrice(trainNo,startNo,endNo,seatType,date string){
	//https://kyfw.12306.cn/otn/leftTicket/queryTicketPrice?train_no=5l0000G13061&from_station_no=01&to_station_no=12&seat_types=O9M&train_date=2017-02-05`]
}
