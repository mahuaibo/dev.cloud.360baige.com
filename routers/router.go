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
	ns := beego.NewNamespace("/cloud",
		beego.NSNamespace(
			"/window",
			beego.NSNamespace("/v1",
				beego.NSNamespace("/user",
					beego.NSInclude(
						&window.UserController{},
					),
				),
			),
		),beego.NSNamespace(
			"/mobile",
			beego.NSNamespace("/v1",
				beego.NSNamespace("/user",
					beego.NSInclude(
						&mobile.UserController{},
					),
				),
			),
		))
	beego.AddNamespace(ns)
	//ns := beego.NewNamespace("/v1",
	//	beego.NSNamespace("/user",
	//		beego.NSInclude(
	//			&controllers.UserController{},
	//		),
	//	),
	//	beego.NSNamespace("/company",
	//		beego.NSInclude(
	//			&controllers.CompanyController{},
	//		),
	//	),
	//	beego.NSNamespace("/person",
	//		beego.NSInclude(
	//			&controllers.PersonController{},
	//		),
	//	),
	//	beego.NSNamespace("/personRelation",
	//		beego.NSInclude(
	//			&controllers.PersonRelationController{},
	//		),
	//	),
	//	beego.NSNamespace("/personStructure",
	//		beego.NSInclude(
	//			&controllers.PersonStructureController{},
	//		),
	//	),
	//	beego.NSNamespace("/structure",
	//		beego.NSInclude(
	//			&controllers.StructureController{},
	//		),
	//	),
	//	beego.NSNamespace("/account",
	//		beego.NSInclude(
	//			&controllers.AccountController{},
	//		),
	//	),
	//	beego.NSNamespace("/accountItem",
	//		beego.NSInclude(
	//			&controllers.AccountItemController{},
	//		),
	//	),
	//	beego.NSNamespace("/appIntroduction",
	//		beego.NSInclude(
	//			&controllers.AppIntroductionController{},
	//		),
	//	),
	//	beego.NSNamespace("/application",
	//		beego.NSInclude(
	//			&controllers.ApplicationController{},
	//		),
	//	),
	//	beego.NSNamespace("/applicationArrange",
	//		beego.NSInclude(
	//			&controllers.ApplicationArrangeController{},
	//		),
	//	),
	//	beego.NSNamespace("/applicationComment",
	//		beego.NSInclude(
	//			&controllers.ApplicationCommentController{},
	//		),
	//	),
	//	beego.NSNamespace("/applicationReport",
	//		beego.NSInclude(
	//			&controllers.ApplicationReportController{},
	//		),
	//	),
	//	beego.NSNamespace("/applicationTpl",
	//		beego.NSInclude(
	//			&controllers.ApplicationTplController{},
	//		),
	//	),
	//	beego.NSNamespace("/attendanceGroup",
	//		beego.NSInclude(
	//			&controllers.AttendanceGroupController{},
	//		),
	//	),
	//	beego.NSNamespace("/attendanceRecord",
	//		beego.NSInclude(
	//			&controllers.AttendanceRecordController{},
	//		),
	//	),
	//	beego.NSNamespace("/attendanceSetup",
	//		beego.NSInclude(
	//			&controllers.AttendanceSetupController{},
	//		),
	//	),
	//	beego.NSNamespace("/attendanceShift",
	//		beego.NSInclude(
	//			&controllers.AttendanceShiftController{},
	//		),
	//	),
	//	beego.NSNamespace("/attendanceShiftItem",
	//		beego.NSInclude(
	//			&controllers.AttendanceShiftItemController{},
	//		),
	//	),
	//	beego.NSNamespace("/attendanceShiftRecord",
	//		beego.NSInclude(
	//			&controllers.AttendanceShiftRecordController{},
	//		),
	//	),
	//	beego.NSNamespace("/card",
	//		beego.NSInclude(
	//			&controllers.CardController{},
	//		),
	//	),
	//	beego.NSNamespace("/feedback",
	//		beego.NSInclude(
	//			&controllers.FeedbackController{},
	//		),
	//	),
	//	beego.NSNamespace("/logger",
	//		beego.NSInclude(
	//			&controllers.LoggerController{},
	//		),
	//	),
	//	beego.NSNamespace("/logistics",
	//		beego.NSInclude(
	//			&controllers.LogisticsController{},
	//		),
	//	),
	//	beego.NSNamespace("/machine",
	//		beego.NSInclude(
	//			&controllers.MachineController{},
	//		),
	//	),
	//	beego.NSNamespace("/message",
	//		beego.NSInclude(
	//			&controllers.MessageController{},
	//		),
	//	),
	//	beego.NSNamespace("/messageReminder",
	//		beego.NSInclude(
	//			&controllers.MessageReminderController{},
	//		),
	//	),
	//	beego.NSNamespace("/order",
	//		beego.NSInclude(
	//			&controllers.OrderController{},
	//		),
	//	),
	//	beego.NSNamespace("/orderComment",
	//		beego.NSInclude(
	//			&controllers.OrderCommentController{},
	//		),
	//	),
	//	beego.NSNamespace("/orderRemindshipment ",
	//		beego.NSInclude(
	//			&controllers.OrderRemindshipmentController{},
	//		),
	//	),
	//	beego.NSNamespace("/quietHours",
	//		beego.NSInclude(
	//			&controllers.QuietHoursController{},
	//		),
	//	),
	//	beego.NSNamespace("/resource",
	//		beego.NSInclude(
	//			&controllers.ResourceController{},
	//		),
	//	),
	//	beego.NSNamespace("/transaction",
	//		beego.NSInclude(
	//			&controllers.TransactionController{},
	//		),
	//	),
	//	beego.NSNamespace("/userPosition",
	//		beego.NSInclude(
	//			&controllers.UserPositionController{},
	//		),
	//	),
	//)
	//beego.AddNamespace(ns)
}
