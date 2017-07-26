// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"dev.cloud.360baige.com/controllers/window"
	"dev.cloud.360baige.com/controllers/mobile"
	"github.com/astaxie/beego"
)

func init() {
	windowRouter()
	//mobileRouter()
}

func windowRouter() {
	windowApi := beego.NewNamespace("/cloud/window/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&window.UserController{},
			),
		),
		beego.NSNamespace("/userposition",
			beego.NSInclude(
				&window.UserPositionController{},
			),
		),
		beego.NSNamespace("/company",
			beego.NSInclude(
				&window.CompanyController{},
			),
		),
		beego.NSNamespace("/account",
			beego.NSInclude(
				&window.AccountController{},
			),
		),
		beego.NSNamespace("/account_item",
			beego.NSInclude(
				&window.AccountItemController{},
			),
		),
		beego.NSNamespace("/order",
			beego.NSInclude(
				&window.OrderController{},
			),
		),
		beego.NSNamespace("/application",
			beego.NSInclude(
				&window.ApplicationController{},
			),
		),
		beego.NSNamespace("/application_tpl",
			beego.NSInclude(
				&window.ApplicationTplController{},
			),
		),
		beego.NSNamespace("/logger",
			beego.NSInclude(
				&window.LoggerController{},
			),
		),
	)
	beego.AddNamespace(windowApi)
}

func mobileRouter() {
	mobileApi := beego.NewNamespace("/cloud/mobile/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&mobile.UserController{},
			),
		),
		beego.NSNamespace("/userposition",
			beego.NSInclude(
				&mobile.UserPositionController{},
			),
		),
		beego.NSNamespace("/account",
			beego.NSInclude(
				&mobile.AccountController{},
			),
		),
		beego.NSNamespace("/account_item",
			beego.NSInclude(
				&mobile.AccountItemController{},
			),
		),
		beego.NSNamespace("/order",
			beego.NSInclude(
				&mobile.OrderController{},
			),
		),
		beego.NSNamespace("/application",
			beego.NSInclude(
				&mobile.ApplicationController{},
			),
		),
		beego.NSNamespace("/application_tpl",
			beego.NSInclude(
				&mobile.ApplicationTplController{},
			),
		),
		beego.NSNamespace("/message_temp",
			beego.NSInclude(
				&mobile.MessageTempController{},
			),
		),
		beego.NSNamespace("/logger",
			beego.NSInclude(
				&mobile.LoggerController{},
			),
		),

	)
	beego.AddNamespace(mobileApi)
}
