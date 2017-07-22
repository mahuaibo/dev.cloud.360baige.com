package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:AccountController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:AccountController"],
		beego.ControllerComments{
			Method: "AccountStatistics",
			Router: `/accountstatistics`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:AccountItemController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:AccountItemController"],
		beego.ControllerComments{
			Method: "List",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:AccountItemController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:AccountItemController"],
		beego.ControllerComments{
			Method: "TradingList",
			Router: `/tradinglist`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:AccountItemController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:AccountItemController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:ApplicationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:ApplicationController"],
		beego.ControllerComments{
			Method: "List",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:ApplicationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:ApplicationController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:ApplicationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:ApplicationController"],
		beego.ControllerComments{
			Method: "ModifyStatus",
			Router: `/modifystatus`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:ApplicationTplController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:ApplicationTplController"],
		beego.ControllerComments{
			Method: "List",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:ApplicationTplController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:ApplicationTplController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:ApplicationTplController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:ApplicationTplController"],
		beego.ControllerComments{
			Method: "Subscription",
			Router: `/subscription`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:ApplicationTplController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:ApplicationTplController"],
		beego.ControllerComments{
			Method: "ModifyStatus",
			Router: `/modifystatus`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:LoggerController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:LoggerController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:MessageTempController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:MessageTempController"],
		beego.ControllerComments{
			Method: "List",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:OrderController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:OrderController"],
		beego.ControllerComments{
			Method: "List",
			Router: `/list`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:OrderController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:OrderController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:OrderController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:OrderController"],
		beego.ControllerComments{
			Method: "DetailByCode",
			Router: `/detailbycode`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:UserController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:UserController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/logout`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:UserController"],
		beego.ControllerComments{
			Method: "ModifyPassword",
			Router: `/modifypassword`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:UserPositionController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:UserPositionController"],
		beego.ControllerComments{
			Method: "PositionList",
			Router: `/positionlist`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:UserPositionController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers/mobile:UserPositionController"],
		beego.ControllerComments{
			Method: "PositionToken",
			Router: `/positiontoken`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
