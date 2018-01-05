package query

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/astaxie/beego"
	config "12306/utils/12306/config"
)

type Schedule struct{

}
func init(){

}

func(schedule *Schedule)Query(date,strart_station,end_station,code string)  (bool,string) {
	is_ok,msg := schedule.initQuery(date,strart_station,end_station,code)
	return is_ok,msg
}

func (schedule *Schedule)initQuery(date,strart_station,end_station,code string) (bool,string) {
	client := &http.Client{}
    reqest, _ := http.NewRequest("GET",fmt.Sprintf(config.UrlTrafficInquiryInit,date, strart_station,end_station,code), nil)
    //设置header头           
    reqest.Header.Set("Accept","*/*")
    reqest.Header.Set("Accept-Encoding","gzip, deflate, br")
    reqest.Header.Set("Accept-Language","zh-CN,zh;q=0.9")
    reqest.Header.Set("Cache-Control","no-cache")
    reqest.Header.Set("Connection","keep-alive")
    reqest.Header.Set("Host","kyfw.12306.cn")
    reqest.Header.Set("If-Modified-Since","0")
    reqest.Header.Set("Referer","https://kyfw.12306.cn/otn/leftTicket/init")
    reqest.Header.Set("User-Agent","Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
    reqest.Header.Set("X-Requested-With","XMLHttpRequest")
	//发送请求
	response,err := client.Do(reqest)
	defer response.Body.Close()
	if err != nil{
		return false,"Schedule Query init 请求失败"
	}
    if response.StatusCode != 200 {
		return false,"Schedule Query init 请求失败"
	}
	body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        // handle error
    }
	beego.Debug(body)
	return true,"Schedule Query init 请求成功"
}