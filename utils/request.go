package utils

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"

	"github.com/astaxie/beego"
)

var requestCookie []*http.Cookie
var requestCookieJar *cookiejar.Jar

//url
var requestURL string

var isDisableHeader bool

//request header
var requestHeader map[string]string

//response header
var responseHeader map[string]string

//Request
type Request struct {
}

func init() {
	isDisableHeader = false

	requestCookie = nil

	requestCookieJar, _ = cookiejar.New(nil)

	requestHeader = make(map[string]string)

	responseHeader = make(map[string]string)
}

//设置Url
func (request *Request) SetURL(url string) *Request {
	requestURL = url
	return request
}

//设置启用停用header
func (request *Request) IsDisableHeader(enable bool) *Request {
	isDisableHeader = enable
	return request
}

//设置默认header
func (request *Request) _SetDefaultHeader(clientRequest *http.Request) *Request {
	//设置header头
	clientRequest.Header.Set("Accept", "*/*")
	clientRequest.Header.Set("Accept-Encoding", "deflate, br")
	clientRequest.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	clientRequest.Header.Set("Cache-Control", "no-cache")
	clientRequest.Header.Set("Connection", "keep-alive")
	clientRequest.Header.Set("Host", "kyfw.12306.cn")
	clientRequest.Header.Set("If-Modified-Since", "0")
	clientRequest.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
	return request
}

//设置用户自定义header
func (request *Request) _SetDefineHeader(clientRequest *http.Request) *Request {
	//设置header头
	for k, v := range requestHeader {
		clientRequest.Header.Set(k, v)
	}
	return request
}

func (request *Request) _ResetRequestDefaultData() {
	// 恢复默认使用header
	isDisableHeader = true
	// 	清空cookies
	requestCookie = nil
	// 清空 requestHeader map
	for k, _ := range requestHeader {
		delete(requestHeader, k)
	}
	// 清空 responseHeader map
	for k, _ := range responseHeader {
		delete(responseHeader, k)
	}
}

//设置header
func (request *Request) SetHeader(key, val string) *Request {
	requestHeader[key] = val
	return request
}

//获取header
func (request *Request) GetHeader(key string) string {
	return requestHeader[key]
}

//发送Get请求
func (request *Request) Get() (bool, string) {
	beego.Debug(requestURL)
	client := &http.Client{CheckRedirect: nil, Jar: requestCookieJar}
	//新建请求
	clientRequest, errNew := http.NewRequest("GET", requestURL, nil)
	//设置header
	if isDisableHeader {
		//设置默认Header
		request._SetDefaultHeader(clientRequest)
		//设置用户自定义Header
		request._SetDefineHeader(clientRequest)
	}
	if errNew != nil {
		return false, "create get request fail"
	}
	//发送请求
	clientResponse, errSend := client.Do(clientRequest)
	beego.Debug(clientResponse)
	if errSend != nil {
		return false, "send get request fail"
	}
	//重置数据
	request._ResetRequestDefaultData()
	//关闭Response body
	defer clientResponse.Body.Close()
	//获取cookies
	requestCookie = requestCookieJar.Cookies(clientRequest.URL)
	//读取数据
	body, errRead := ioutil.ReadAll(clientResponse.Body)
	if clientResponse.StatusCode != 200 {
		return false, string(body)
	}
	if errRead != nil {
		return false, "read get request body fail"
	}
	return true, string(body)
}
