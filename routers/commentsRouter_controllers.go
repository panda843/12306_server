package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:PassengerController"] = append(beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:PassengerController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:ScheduleController"] = append(beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:ScheduleController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:StationController"] = append(beego.GlobalControllerRouter["github.com/chuanshuo843/12306_server/controllers:StationController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
