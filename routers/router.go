// @APIVersion 1.0.0
// @Title 360baige.com Cloud API
// @Description 360baige.com Cloud API
// @Contact mahuaibo@360baige.com
// @TermsOfServiceUrl http://www.360baige.com
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	centerWindow "dev.cloud.360baige.com/controllers/window/center"
	centerMobile "dev.cloud.360baige.com/controllers/mobile/center"
	schoolFeeMobile "dev.cloud.360baige.com/controllers/mobile/schoolfee"
	schoolFeeWindow "dev.cloud.360baige.com/controllers/window/schoolfee"
	personnelWindow "dev.cloud.360baige.com/controllers/window/personnel"
	safeWindow "dev.cloud.360baige.com/controllers/window/safe"
)

func init() {
	centerWindowRouter() // window->后台管理
	centerMobileRouter() // mobile->后台管理

	schoolFeeWindowRouter() // window->缴费管理
	schoolFeeMobileRouter() // mobile->缴费管理

	personnelWindowRouter() // window->人事管理
	safeWindowRouter() //windiw->安全中心
}

func schoolFeeWindowRouter() {
	schoolFeeWindowRouter := beego.NewNamespace("/cloud/window/schoolFee/v1",
		beego.NSNamespace("/project",
			beego.NSInclude(
				&schoolFeeWindow.ProjectController{},
			),
		),
		beego.NSNamespace("/record",
			beego.NSInclude(
				&schoolFeeWindow.RecordController{},
			),
		),
	)
	beego.AddNamespace(schoolFeeWindowRouter)
}

func schoolFeeMobileRouter() {
	schoolFeeMobileRouter := beego.NewNamespace("/cloud/mobile/schoolFee/v1",
		beego.NSNamespace("/project",
			beego.NSInclude(
				&schoolFeeMobile.ProjectController{},
			),
		),
		beego.NSNamespace("/record",
			beego.NSInclude(
				&schoolFeeMobile.RecordController{},
			),
		),
	)
	beego.AddNamespace(schoolFeeMobileRouter)
}

func centerWindowRouter() {
	centerWindowRouter := beego.NewNamespace("/cloud/window/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&centerWindow.UserController{},
			),
		),
		beego.NSNamespace("/userPosition",
			beego.NSInclude(
				&centerWindow.UserPositionController{},
			),
		),
		beego.NSNamespace("/company",
			beego.NSInclude(
				&centerWindow.CompanyController{},
			),
		),
		beego.NSNamespace("/account",
			beego.NSInclude(
				&centerWindow.AccountController{},
			),
		),
		beego.NSNamespace("/accountItem",
			beego.NSInclude(
				&centerWindow.AccountItemController{},
			),
		),
		beego.NSNamespace("/order",
			beego.NSInclude(
				&centerWindow.OrderController{},
			),
		),
		beego.NSNamespace("/application",
			beego.NSInclude(
				&centerWindow.ApplicationController{},
			),
		),
		beego.NSNamespace("/applicationTpl",
			beego.NSInclude(
				&centerWindow.ApplicationTplController{},
			),
		),
		beego.NSNamespace("/logger",
			beego.NSInclude(
				&centerWindow.LoggerController{},
			),
		),
	)
	beego.AddNamespace(centerWindowRouter)
}

func centerMobileRouter() {
	centerMobileRouter := beego.NewNamespace("/cloud/mobile/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&centerMobile.UserController{},
			),
		),
		beego.NSNamespace("/userPosition",
			beego.NSInclude(
				&centerMobile.UserPositionController{},
			),
		),
		beego.NSNamespace("/company",
			beego.NSInclude(
				&centerMobile.CompanyController{},
			),
		),
		beego.NSNamespace("/account",
			beego.NSInclude(
				&centerMobile.AccountController{},
			),
		),
		beego.NSNamespace("/accountItem",
			beego.NSInclude(
				&centerMobile.AccountItemController{},
			),
		),
		beego.NSNamespace("/order",
			beego.NSInclude(
				&centerMobile.OrderController{},
			),
		),
		beego.NSNamespace("/application",
			beego.NSInclude(
				&centerMobile.ApplicationController{},
			),
		),
		beego.NSNamespace("/applicationTpl",
			beego.NSInclude(
				&centerMobile.ApplicationTplController{},
			),
		),
		beego.NSNamespace("/logger",
			beego.NSInclude(
				&centerMobile.LoggerController{},
			),
		),
	)
	beego.AddNamespace(centerMobileRouter)
}

func personnelWindowRouter() {
	personnelWindowRouter := beego.NewNamespace("/cloud/window/personnel/v1",
		beego.NSNamespace("/person",
			beego.NSInclude(
				&personnelWindow.PersonController{},
			),
		),
		beego.NSNamespace("/personRelation",
			beego.NSInclude(
				&personnelWindow.PersonRelationController{},
			),
		),
		beego.NSNamespace("/personStructure",
			beego.NSInclude(
				&personnelWindow.PersonStructureController{},
			),
		),
		beego.NSNamespace("/structure",
			beego.NSInclude(
				&personnelWindow.StructureController{},
			),
		),
	)
	beego.AddNamespace(personnelWindowRouter)
}

func safeWindowRouter() {
	safeWindowRouter := beego.NewNamespace("/cloud/window/safe/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&safeWindow.UserController{},
			),
		),
		beego.NSNamespace("/userPosition",
			beego.NSInclude(
				&safeWindow.UserPositionController{},
			),
		),
	)
	beego.AddNamespace(safeWindowRouter)
}
