package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

type _Header struct {
	JwtHead string `json:"type"`
	JwtAlg  string `json:"alg"`
}

type _Payload struct {
	Iss  string `json:"iss"`  //Issuer，发行者
	Sub  string `json:"sub"`  //Subject，主题
	Aud  string `json:"aud"`  //Audience，观众
	Data string `json:"data"` //请求数据
	Exp  int64  `json:"exp"`  //Expiration time，过期时间
	Nbf  int64  `json:"nbf"`  //Not before
	Iat  int64  `json:"iat"`  //Issued at，发行时间
	Jti  int64  `json:"jti"`  //JWT ID
}

type Jwt struct{}

var initPayload _Payload

var initHeader _Header

var secretKey string

func init() {
	//初始化秘钥
	secretKey = "jwt_key"
	//设置header头
	initHeader.JwtHead = "JWT"
	initHeader.JwtAlg = "HS256"
	//设置payload
	initPayload.Iss = "https://www.ganktools.com"
	initPayload.Sub = "https://www.ganktools.com"
	initPayload.Aud = "https://www.ganktools.com"
}

//编码JWT的Header头
func (_header *_Header) EncodeHeader() string {
	json_data, _ := json.Marshal(_header)
	return base64.StdEncoding.EncodeToString(json_data)
}

//解码JWT的Header头
func (_header *_Header) DecodeHeader(data string) bool {
	decode_header, _ := base64.StdEncoding.DecodeString(data)
	err_header := json.Unmarshal(decode_header, &_header)
	if err_header != nil {
		return false
	}
	return true
}

//编码payload部分
func (_payload *_Payload) EncodePayload() string {
	json_data, _ := json.Marshal(_payload)
	return base64.StdEncoding.EncodeToString(json_data)
}

//解码payload部分
func (_payload *_Payload) DecodePayload(data string) bool {
	decode_payload, _ := base64.StdEncoding.DecodeString(data)
	err_payload := json.Unmarshal(decode_payload, &_payload)
	if err_payload != nil {
		return false
	}
	return true
}

//JWT的secret部分加密
func signature(jwt, key string) string {
	secret := []byte(key)
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(jwt))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

//设置Issuer
func (jwt *Jwt) SetIssuer(iss string) *Jwt {
	initPayload.Iss = iss
	return jwt
}

//设置Subject
func (jwt *Jwt) SetSubject(sub string) *Jwt {
	initPayload.Sub = sub
	return jwt
}

//设置Audience
func (jwt *Jwt) SetAudience(aud string) *Jwt {
	initPayload.Aud = aud
	return jwt
}

//设置Key
func (jwt *Jwt) SetSecretKey(key string) {
	secretKey = key
}

//JWT加密
func (jwt *Jwt) Encode(exp int64, data string) string {
	current_time := time.Now().Unix()
	initPayload.Jti, initPayload.Iat, initPayload.Nbf = current_time, current_time, current_time
	initPayload.Exp = exp
	initPayload.Data = data
	encode_header := initHeader.EncodeHeader()
	encode_payload := initPayload.EncodePayload()
	encode_jwt := encode_header + "." + encode_payload
	secret := signature(encode_jwt, secretKey)
	return encode_jwt + "." + secret
}

//JWT检测
func (jwt *Jwt) Checkd(token string) bool {
	data := strings.Split(token, ".")
	//检测长度
	if len(data) != 3 {
		return false
	}
	//检测Hash是否一致
	secret := signature(string(data[0])+"."+string(data[1]), secretKey)
	if secret != string(data[2]) {
		return false
	}
	//解码Payload
	if !initPayload.DecodePayload(string(data[1])) {
		return false
	}
	//检测JWT是否过期
	if initPayload.Exp <= time.Now().Unix() {
		return false
	}
	//检测什么时间之后可用
	if initPayload.Nbf >= time.Now().Unix() {
		return false
	}
	return true
}

func (jwt *Jwt) GetData() string {
	return initPayload.Data
}
