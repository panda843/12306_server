package kyfw

import (
	"sync"

	"github.com/chuanshuo843/12306_server/utils"
)

var (
	request *utils.Request
)

// SetRequest .
func SetRequest(req *utils.Request) {
	if !req.InitBool {
		req.InitRequest()
	}

	request = req
}

// UserRequest .
var UserRequest *sync.Map

func init() {
	UserRequest = &sync.Map{}
}

// Base .
type Base struct {
}

// func init() {
// 	request.InitRequest()
// }

// Store .
func Store(key string, req *utils.Request) {
	if !req.InitBool {
		req.InitRequest()
	}
	UserRequest.Store(key, req)
}

// Load .
func Load(key string) *utils.Request {
	v, ok := UserRequest.Load(key)
	if !ok {
		return nil
	}
	req, ok := v.(*utils.Request)
	if !ok {
		return nil
	}
	return req
}
