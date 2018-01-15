package kyfw

import (
	"errors"
	"strings"
)

const (
	QueryGetStation = "https://kyfw.12306.cn/otn/resources/js/framework/station_name.js?station_version=1.9042"
)

type Query struct {
	Base
}

//获取站台信息
func (query *Query) GetStations() (string,error) {
	err := request.CreateHttpRequest(QueryGetStation,"GET",nil)
	if err != nil {
		return "",err
	}
	data,errSend := request.Send()
	if errSend != nil {
		return "",errSend
	}
	stationList := strings.Split(string(data), "'")
	if len(stationList) != 3 {
		return "",errors.New("获取站台信息失败")
	}
	return string(stationList[1]),nil
}