// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"net/http"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/chuanshuo843/12306_server/controllers"
	"github.com/chuanshuo843/12306_server/utils"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSRouter("/auth/login", &controllers.UserController{}, "Post:Login"),
		beego.NSNamespace("/user",
			beego.NSBefore(Auth),
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/schedule",
			beego.NSBefore(Auth),
			beego.NSInclude(
				&controllers.ScheduleController{},
			),
		),
	)
	beego.AddNamespace(ns)
}

func Auth(ctx *context.Context) {
	authString := ctx.Input.Header("Authorization")
	beego.Debug("AuthString:", authString)
	if authString == "" {
		ctx.Output.Header("Cache-Control", "no-store")
		ctx.Output.Header("WWW-Authenticate", "Bearer realm=\""+beego.AppConfig.String("HostName")+"\" error=\"Authorization\" error_description=\"invalid Authorization\"")
		http.Error(ctx.ResponseWriter, "Unauthorized", 401)
		return
	}
	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		ctx.Output.Header("Cache-Control", "no-store")
		ctx.Output.Header("WWW-Authenticate", "Bearer realm=\""+beego.AppConfig.String("HostName")+"\" error=\"Authorization\" error_description=\"invalid Authorization\"")
		http.Error(ctx.ResponseWriter, "Unauthorized", 401)
		return
	}
	token := kv[1]
	beego.Debug(token)
	jwt := &utils.Jwt{}
	jwt.SetSecretKey(beego.AppConfig.String("JwtKey"))
	if !jwt.Checkd(token) {
		ctx.Output.Header("Cache-Control", "no-store")
		ctx.Output.Header("WWW-Authenticate", "Bearer realm=\""+beego.AppConfig.String("HostName")+"\" error=\"Authorization\" error_description=\"invalid Authorization\"")
		http.Error(ctx.ResponseWriter, "Unauthorized", 401)
		return
	}
}
