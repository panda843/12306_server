package utils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

//Request
type Request struct {
	HttpClient           *http.Client
	HttpRequest          *http.Request
	HttpResponse         *http.Response
	DisableDefaultHeader bool
}

//初始化Request
func InitRequest() *Request {
	req := &Request{
		HttpClient:   &http.Client{},
		HttpRequest:  &http.Request{},
		HttpResponse: &http.Response{},
	}
	//初始化CookieJar
	req.HttpClient.Jar, _ = cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	//设置CheckRedirect
	req.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	//设置Transport
	req.HttpClient.Transport = &http.Transport{
		// 12306 https CA认证
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		// TLSClientConfig: &tls.Config{RootCAs: pool},
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 60 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 60 * time.Second,
	}
	return req
}

//创建请求
func (request *Request) CreateHttpRequest(requestUrl, method string, val interface{}) error {
	var nReq *http.Request
	var nErr error
	switch data := val.(type) {
	case string:
		nReq, nErr = http.NewRequest(method, requestUrl, io.Reader(bytes.NewBuffer([]byte(data))))
	case []byte:
		nReq, nErr = http.NewRequest(method, requestUrl, io.Reader(bytes.NewBuffer(data)))
	case url.Values:
		nReq, nErr = http.NewRequest(method, requestUrl, strings.NewReader(data.Encode()))
	case *url.Values:
		nReq, nErr = http.NewRequest(method, requestUrl, strings.NewReader(data.Encode()))
	default:
		nReq, nErr = http.NewRequest(method, requestUrl, nil)
	}
	if nErr != nil {
		return errors.New("Request请求创建失败")
	}
	request.HttpRequest = nReq
	return nil
}

//设置Header头
func (request *Request) SetHeader(key, val string) error {
	if request.HttpRequest != nil {
		request.HttpRequest.Header.Set(key, val)
		return nil
	}
	return errors.New("HttpRequest为空")
}

//设置Cookie
func (request *Request) SetCookie(cookie []*http.Cookie) error {
	if request.HttpClient.Jar != nil && request.HttpRequest != nil {
		request.HttpClient.Jar.SetCookies(request.HttpRequest.URL, cookie)
		return nil
	}
	return errors.New("HttpClient为空")
}

//设置默认请求头
func SetRequestDefaultHeader(request *Request) {
	//设置默认Header
	if !request.DisableDefaultHeader {
		request.SetHeader("Accept", "*/*")
		request.SetHeader("Accept-Encoding", "deflate, br")
		request.SetHeader("Accept-Language", "zh-CN,zh;q=0.9")
		request.SetHeader("Cache-Control", "no-cache")
		request.SetHeader("Connection", "keep-alive")
		request.SetHeader("Host", "kyfw.12306.cn")
		request.SetHeader("Origin", "https://kyfw.12306.cn")
		request.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
	}
}

//读取返回数据
func ReadHttpResponseData(request *Request) ([]byte, error) {
	//关闭Response body
	defer request.HttpResponse.Body.Close()
	//返回读取数据
	return ioutil.ReadAll(request.HttpResponse.Body)
}

//发送Request请求
func SendHttpRequest(request *Request) error {
	//判断Post请求,设置Header头Content-Type
	if strings.ToUpper(request.HttpRequest.Method) == "POST" {
		request.SetHeader("Content-Type", "application/x-www-form-urlencoded; charset=utf-")
	}
	nRsp, errSend := request.HttpClient.Do(request.HttpRequest)
	if errSend != nil {
		return errors.New("发送Request请求失败")
	}
	request.HttpResponse = nRsp
	return nil
}

//发送请求
func (request *Request) Send() ([]byte, error) {
	if request.HttpRequest != nil {
		//设置header
		SetRequestDefaultHeader(request)
		//发送请求
		err := SendHttpRequest(request)
		if err != nil {
			return nil, err
		}
		//读取数据
		return ReadHttpResponseData(request)
	}
	return nil, errors.New("HttpRequest为空")
}
