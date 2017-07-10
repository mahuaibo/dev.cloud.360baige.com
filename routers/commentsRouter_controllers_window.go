package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:CompanyController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:CompanyController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:CompanyController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:CompanyController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserController"],
		beego.ControllerComments{
			Method: "ModifyPassword",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserPositionController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserPositionController"],
		beego.ControllerComments{
			Method: "List",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserPositionController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserPositionController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
