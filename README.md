# 12306API接口

```golang
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
```

