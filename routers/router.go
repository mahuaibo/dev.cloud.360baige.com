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
	windowApi := beego.NewNamespace("/cloud",
		beego.NSNamespace(
			"/window",
			beego.NSNamespace("/v1",
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
			),
		))
	beego.AddNamespace(windowApi)

	mobileApi := beego.NewNamespace("/cloud",
		beego.NSNamespace(
			"/mobile",
			beego.NSNamespace("/v1",
				beego.NSNamespace("/user",
					beego.NSInclude(
						&mobile.UserController{},
					),
				),
			),
		))
	beego.AddNamespace(mobileApi)
}
