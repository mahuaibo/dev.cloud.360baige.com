// @APIVersion 1.0.0
// @Title 360baige.com Cloud API
// @Description 360baige.com Cloud API
// @Contact mahuaibo@360baige.com
// @TermsOfServiceUrl http://www.360baige.com
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	center_window "dev.cloud.360baige.com/controllers/window/center"
	//center_mobile "dev.cloud.360baige.com/controllers/mobile/center"
	schoolfee_mobile "dev.cloud.360baige.com/controllers/mobile/schoolfee"
	schoolfee_window "dev.cloud.360baige.com/controllers/window/schoolfee"
	personnel_window "dev.cloud.360baige.com/controllers/window/personnel"
	"github.com/astaxie/beego"
)

func init() {
	centerWindowRouter() // window->admin后台管理
	centerMobileRouter()

	schoolfeeWindowRouter() // window->缴费管理
	schoolfeeMobileRouter()

	personnelWindowRouter() // window->人事管理
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

func centerWindowRouter() {
	windowApi := beego.NewNamespace("/cloud/window/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&center_window.UserController{},
			),
		),
		beego.NSNamespace("/userPosition",
			beego.NSInclude(
				&center_window.UserPositionController{},
			),
		),
		beego.NSNamespace("/company",
			beego.NSInclude(
				&center_window.CompanyController{},
			),
		),
		beego.NSNamespace("/account",
			beego.NSInclude(
				&center_window.AccountController{},
			),
		),
		beego.NSNamespace("/account_item",
			beego.NSInclude(
				&center_window.AccountItemController{},
			),
		),
		beego.NSNamespace("/order",
			beego.NSInclude(
				&center_window.OrderController{},
			),
		),
		beego.NSNamespace("/application",
			beego.NSInclude(
				&center_window.ApplicationController{},
			),
		),
		beego.NSNamespace("/application_tpl",
			beego.NSInclude(
				&center_window.ApplicationTplController{},
			),
		),
		beego.NSNamespace("/logger",
			beego.NSInclude(
				&center_window.LoggerController{},
			),
		),
	)
	beego.AddNamespace(windowApi)
}

func centerMobileRouter() {
	mobileApi := beego.NewNamespace("/cloud/mobile/v1",
		//beego.NSNamespace("/user",
		//	beego.NSInclude(
		//		&center_mobile.UserController{},
		//	),
		//),
		//beego.NSNamespace("/userposition",
		//	beego.NSInclude(
		//		&center_mobile.UserPositionController{},
		//	),
		//),
		//beego.NSNamespace("/account",
		//	beego.NSInclude(
		//		&center_mobile.AccountController{},
		//	),
		//),
		//beego.NSNamespace("/account_item",
		//	beego.NSInclude(
		//		&center_mobile.AccountItemController{},
		//	),
		//),
		//beego.NSNamespace("/order",
		//	beego.NSInclude(
		//		&center_mobile.OrderController{},
		//	),
		//),
		//beego.NSNamespace("/application",
		//	beego.NSInclude(
		//		&center_mobile.ApplicationController{},
		//	),
		//),
		//beego.NSNamespace("/application_tpl",
		//	beego.NSInclude(
		//		&center_mobile.ApplicationTplController{},
		//	),
		//),
		//beego.NSNamespace("/message_temp",
		//	beego.NSInclude(
		//		&center_mobile.MessageTempController{},
		//	),
		//),
		//beego.NSNamespace("/logger",
		//	beego.NSInclude(
		//		&center_mobile.LoggerController{},
		//	),
		//),

	)
	beego.AddNamespace(mobileApi)
}

func personnelWindowRouter() {
	mobileApi := beego.NewNamespace("/win/personnel/v1",
		beego.NSNamespace("/person",
			beego.NSInclude(
				&personnel_window.PersonController{},
			),
		),
		beego.NSNamespace("/person_relation",
			beego.NSInclude(
				&personnel_window.PersonRelationController{},
			),
		),
		beego.NSNamespace("/person_structure",
			beego.NSInclude(
				&personnel_window.PersonStructureController{},
			),
		),
		beego.NSNamespace("/structure",
			beego.NSInclude(
				&personnel_window.StructureController{},
			),
		),
	)
	beego.AddNamespace(mobileApi)
}
