package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:ScheduleController"] = append(beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:ScheduleController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:UserController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:UserController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:UserController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:UserController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"put"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:UserController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:uid`,
			AllowHTTPMethods: []string{"delete"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:UserController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/logout`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
