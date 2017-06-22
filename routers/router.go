// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"dev.cloud.360baige.com/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/company",
			beego.NSInclude(
				&controllers.CompanyController{},
			),
		),
		beego.NSNamespace("/person",
			beego.NSInclude(
				&controllers.PersonController{},
			),
		),
		beego.NSNamespace("/personRelation",
			beego.NSInclude(
				&controllers.PersonRelationController{},
			),
		),
		beego.NSNamespace("/personStructure",
			beego.NSInclude(
				&controllers.PersonStructureController{},
			),
		),
		beego.NSNamespace("/structure",
			beego.NSInclude(
				&controllers.StructureController{},
			),
		),
		beego.NSNamespace("/account",
			beego.NSInclude(
				&controllers.AccountController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
