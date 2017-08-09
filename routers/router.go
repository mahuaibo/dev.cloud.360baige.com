// @APIVersion 1.0.0
// @Title 360baige.com Cloud API
// @Description 360baige.com Cloud API
// @Contact mahuaibo@360baige.com
// @TermsOfServiceUrl http://www.360baige.com
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	admin_window "dev.cloud.360baige.com/controllers/window/admin"
	center_mobile "dev.cloud.360baige.com/controllers/mobile/center"
	schoolfee_mobile "dev.cloud.360baige.com/controllers/mobile/schoolfee"
	schoolfee_window "dev.cloud.360baige.com/controllers/window/schoolfee"
	"github.com/astaxie/beego"
)

func init() {
	adminWindowRouter()
	//centerMobileRouter()
	schoolfeeWindowRouter()
	//schoolfeeMobileRouter()
}

func schoolfeeWindowRouter() {
	schoolfeewinApi := beego.NewNamespace("/win/schoolfee/v1",
		beego.NSNamespace("/project",
			beego.NSInclude(
				&schoolfee_window.ProjectController{},
			),
		),
		beego.NSNamespace("/record",
			beego.NSInclude(
				&schoolfee_window.RecordController{},
			),
		),
	)
	beego.AddNamespace(schoolfeewinApi)
}

func schoolfeeMobileRouter() {
	schoolfeeappApi := beego.NewNamespace("/app/schoolfee/v1",
		beego.NSNamespace("/project",
			beego.NSInclude(
				&schoolfee_mobile.ProjectController{},
			),
		),
		beego.NSNamespace("/record",
			beego.NSInclude(
				&schoolfee_mobile.RecordController{},
			),
		),
	)
	beego.AddNamespace(schoolfeeappApi)
}

func adminWindowRouter() {
	windowApi := beego.NewNamespace("/cloud/window/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&admin_window.UserController{},
			),
		),
		beego.NSNamespace("/userposition",
			beego.NSInclude(
				&admin_window.UserPositionController{},
			),
		),
		beego.NSNamespace("/company",
			beego.NSInclude(
				&admin_window.CompanyController{},
			),
		),
		beego.NSNamespace("/account",
			beego.NSInclude(
				&admin_window.AccountController{},
			),
		),
		beego.NSNamespace("/account_item",
			beego.NSInclude(
				&admin_window.AccountItemController{},
			),
		),
		beego.NSNamespace("/order",
			beego.NSInclude(
				&admin_window.OrderController{},
			),
		),
		beego.NSNamespace("/application",
			beego.NSInclude(
				&admin_window.ApplicationController{},
			),
		),
		beego.NSNamespace("/application_tpl",
			beego.NSInclude(
				&admin_window.ApplicationTplController{},
			),
		),
		beego.NSNamespace("/logger",
			beego.NSInclude(
				&admin_window.LoggerController{},
			),
		),
	)
	beego.AddNamespace(windowApi)
}

func centerMobileRouter() {
	mobileApi := beego.NewNamespace("/cloud/mobile/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&center_mobile.UserController{},
			),
		),
		beego.NSNamespace("/userposition",
			beego.NSInclude(
				&center_mobile.UserPositionController{},
			),
		),
		beego.NSNamespace("/account",
			beego.NSInclude(
				&center_mobile.AccountController{},
			),
		),
		beego.NSNamespace("/account_item",
			beego.NSInclude(
				&center_mobile.AccountItemController{},
			),
		),
		beego.NSNamespace("/order",
			beego.NSInclude(
				&center_mobile.OrderController{},
			),
		),
		beego.NSNamespace("/application",
			beego.NSInclude(
				&center_mobile.ApplicationController{},
			),
		),
		beego.NSNamespace("/application_tpl",
			beego.NSInclude(
				&center_mobile.ApplicationTplController{},
			),
		),
		beego.NSNamespace("/message_temp",
			beego.NSInclude(
				&center_mobile.MessageTempController{},
			),
		),
		beego.NSNamespace("/logger",
			beego.NSInclude(
				&center_mobile.LoggerController{},
			),
		),

	)
	beego.AddNamespace(mobileApi)
}
