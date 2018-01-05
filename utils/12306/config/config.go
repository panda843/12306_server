package config

// 12306 URLs
const (
	//根Url
	BaseUrl = "https://kyfw.12306.cn/otn/"
	//查询车次准备
	UrlTrafficInquiryInit = BaseUrl + "/leftTicket/log?leftTicketDTO.train_date=%s&leftTicketDTO.from_station=%s&leftTicketDTO.to_station=%s&purpose_codes=%s"
	//查询车次
	UrlTrafficInquiry = BaseUrl + "leftTicket/queryA"
)