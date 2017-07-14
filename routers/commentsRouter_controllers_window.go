package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:AccountController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:AccountController"],
		beego.ControllerComments{
			Method: "AccountStatistics",
			Router: `/accountstatistics`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:AccountItemController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:AccountItemController"],
		beego.ControllerComments{
			Method: "List",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:AccountItemController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:AccountItemController"],
		beego.ControllerComments{
			Method: "TradingList",
			Router: `/tradinglist`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:AccountItemController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:AccountItemController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

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
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:LoggerController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:LoggerController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:LoggerController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:LoggerController"],
		beego.ControllerComments{
			Method: "GetList",
			Router: `/getList`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:OrderController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:OrderController"],
		beego.ControllerComments{
			Method: "List",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:OrderController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:OrderController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:OrderController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:OrderController"],
		beego.ControllerComments{
			Method: "DetailByCode",
			Router: `/detailbycode`,
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
			Router: `/logout`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserController"],
		beego.ControllerComments{
			Method: "ModifyPassword",
			Router: `/modifypassword`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserPositionController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserPositionController"],
		beego.ControllerComments{
			Method: "PositionList",
			Router: `/positionlist`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserPositionController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/window:UserPositionController"],
		beego.ControllerComments{
			Method: "PositionToken",
			Router: `/positiontoken`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
