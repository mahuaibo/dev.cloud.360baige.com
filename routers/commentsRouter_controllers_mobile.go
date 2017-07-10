package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:UserController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:UserPositionController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:UserPositionController"],
		beego.ControllerComments{
			Method: "List",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
