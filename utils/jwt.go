package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

type JwtHeader struct {
	JwtHead string `json:"type"`
	JwtAlg  string `json:"alg"`
}

type JwtPayload struct {
	Iss  string `json:"iss"`  //Issuer，发行者
	Sub  string `json:"sub"`  //Subject，主题
	Aud  string `json:"aud"`  //Audience，观众
	Data string `json:"data"` //请求数据
	Exp  int64  `json:"exp"`  //Expiration time，过期时间
	Nbf  int64  `json:"nbf"`  //Not before
	Iat  int64  `json:"iat"`  //Issued at，发行时间
	Jti  int64  `json:"jti"`  //JWT ID
}

type JwtSecret struct {
	Key string `json:"key"` //秘钥
}

type Jwt struct {
	Header  JwtHeader
	Payload JwtPayload
	Secret  JwtSecret
}

//JWT 初始化
func (jwt *Jwt) InitJwt() {
	//设置Header
	jwt.Header.JwtAlg = "HS256"
	jwt.Header.JwtHead = "JWT"
	//设置Payload
	jwt.Payload.Iss = "https://www.ganktools.com"
	jwt.Payload.Sub = "https://www.ganktools.com"
	jwt.Payload.Aud = "https://www.ganktools.com"
	//设置加密秘钥
	jwt.Secret.Key = "jwt_key"
}

//编码JWT的Header头
func (header *JwtHeader) Encode() string {
	json_data, _ := json.Marshal(header)
	return base64.StdEncoding.EncodeToString(json_data)
}

//解码JWT的Header头
func (header *JwtHeader) Decode(data string) error {
	headerStr, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	errJson := json.Unmarshal(headerStr, header)
	if errJson != nil {
		return errJson
	}
	return nil
}

//解码payload部分
func (payload *JwtPayload) Decode(data string) error {
	payloadStr, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	errJson := json.Unmarshal(payloadStr, payload)
	if errJson != nil {
		return errJson
	}
	return nil
}

//编码JWT的payload部分
func (payload *JwtPayload) Encode() string {
	json_data, _ := json.Marshal(payload)
	return base64.StdEncoding.EncodeToString(json_data)
}

//JWT的secret部分加密
func (secret *JwtSecret) Signature(header, payload string) string {
	encode_jwt := header + "." + payload
	h := hmac.New(sha256.New, []byte(secret.Key))
	h.Write([]byte(encode_jwt))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

//JWT加密
func (jwt *Jwt) Encode() string {
	headerStr := jwt.Header.Encode()
	payloadStr := jwt.Payload.Encode()
	return headerStr + "." + payloadStr + "." + jwt.Secret.Signature(headerStr, payloadStr)
}

//JWT解码
func (jwt *Jwt) Decode(token string) error {
	data := strings.Split(token, ".")
	errHeader := jwt.Header.Decode(string(data[0]))
	if errHeader != nil {
		return errHeader
	}
	errPayload := jwt.Payload.Decode(string(data[1]))
	if errPayload != nil {
		return errPayload
	}
	return nil
}

//JWT检测
func (jwt *Jwt) Checkd(token string) bool {
	data := strings.Split(token, ".")
	//检测长度
	if len(data) != 3 {
		return false
	}
	//解码Token
	errDeCode := jwt.Decode(token)
	if errDeCode != nil {
		return false
	}
	//检测Hash是否一致
	secret := jwt.Secret.Signature(string(data[0]), string(data[1]))
	if secret != string(data[2]) {
		return false
	}
	//检测JWT是否过期
	if jwt.Payload.Exp <= time.Now().Unix() {
		return false
	}
	//检测什么时间之后可用
	if jwt.Payload.Nbf >= time.Now().Unix() {
		return false
	}
	return true
}

// Token .
func Token(authString string) string {
	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		return ""
	}
	token := kv[1]
	return token
}
