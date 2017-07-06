package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountController"],
		beego.ControllerComments{
			Method: "Account",
			Router: `/account`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountController"],
		beego.ControllerComments{
			Method: "List",
			Router: `/list`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountController"],
		beego.ControllerComments{
			Method: "UpdateByIds",
			Router: `/updateByIds`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountController"],
		beego.ControllerComments{
			Method: "AddMultiple",
			Router: `/addMultiple`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountItemController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountItemController"],
		beego.ControllerComments{
			Method: "GetBillList",
			Router: `/getbilllist`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountItemController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountItemController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountItemController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountItemController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountItemController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AccountItemController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AppIntroductionController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AppIntroductionController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationArrangeController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationArrangeController"],
		beego.ControllerComments{
			Method: "GetAppArrange",
			Router: `/get-app-arrange`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationCommentController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationCommentController"],
		beego.ControllerComments{
			Method: "GetCommentList",
			Router: `/get-commentlist`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationCommentController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationCommentController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationController"],
		beego.ControllerComments{
			Method: "GetAppList",
			Router: `/getapplist`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationController"],
		beego.ControllerComments{
			Method: "ModifyApp",
			Router: `/modifyapp`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationController"],
		beego.ControllerComments{
			Method: "SearchAppList",
			Router: `/search-applist`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationController"],
		beego.ControllerComments{
			Method: "GetAppcorrelationtList",
			Router: `/get-appcorrelationtlist`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationReportController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationReportController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationTplController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationTplController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationTplController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationTplController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationTplController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationTplController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceGroupController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceGroupController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceGroupController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceGroupController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceGroupController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceGroupController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceRecordController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceRecordController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceRecordController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceRecordController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceRecordController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceRecordController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceSetupController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceSetupController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceSetupController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceSetupController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceSetupController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceSetupController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftItemController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftItemController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftItemController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftItemController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftItemController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftItemController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftRecordController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftRecordController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftRecordController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftRecordController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftRecordController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:AttendanceShiftRecordController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CardController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CardController"],
		beego.ControllerComments{
			Method: "BindIcCard",
			Router: `/bindiccard`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CardController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CardController"],
		beego.ControllerComments{
			Method: "UnBindIcCard",
			Router: `/unbindiccard`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CardController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CardController"],
		beego.ControllerComments{
			Method: "LossIcCard",
			Router: `/lossiccard`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CardController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CardController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CardController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CardController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CardController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CardController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CompanyController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CompanyController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CompanyController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CompanyController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CompanyController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CompanyController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:FeedbackController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:FeedbackController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:LoggerController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:LoggerController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:LoggerController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:LoggerController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:LoggerController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:LoggerController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:LogisticsController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:LogisticsController"],
		beego.ControllerComments{
			Method: "GetLogistics",
			Router: `/getlogistics`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:MachineController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:MachineController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:MachineController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:MachineController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:MachineController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:MachineController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:MessageController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:MessageController"],
		beego.ControllerComments{
			Method: "GetMessageList",
			Router: `/get-message-list`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:MessageReminderController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:MessageReminderController"],
		beego.ControllerComments{
			Method: "Set",
			Router: `/set`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:MessageReminderController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:MessageReminderController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:OrderCommentController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:OrderCommentController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:OrderController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:OrderController"],
		beego.ControllerComments{
			Method: "GetOrderList",
			Router: `/getorderlist`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:OrderController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:OrderController"],
		beego.ControllerComments{
			Method: "CancelOrder",
			Router: `/cancelorder`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:OrderController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:OrderController"],
		beego.ControllerComments{
			Method: "ConfirmationOfReceipt",
			Router: `/confirmationOfReceipt`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:OrderController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:OrderController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:OrderController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:OrderController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:OrderController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:OrderController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:OrderRemindshipmentController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:OrderRemindshipmentController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"],
		beego.ControllerComments{
			Method: "GetSchoolArea",
			Router: `/getschoolarea`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"],
		beego.ControllerComments{
			Method: "GetCardList",
			Router: `/getcardlist`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"],
		beego.ControllerComments{
			Method: "GetChildList",
			Router: `/getchildlist`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"],
		beego.ControllerComments{
			Method: "SetMyChild",
			Router: `/setmychild`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonRelationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonRelationController"],
		beego.ControllerComments{
			Method: "AddPersonRelation",
			Router: `/addchildparent`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonRelationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonRelationController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonRelationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonRelationController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonRelationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonRelationController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonStructureController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonStructureController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonStructureController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonStructureController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonStructureController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonStructureController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:QuietHoursController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:QuietHoursController"],
		beego.ControllerComments{
			Method: "Set",
			Router: `/set`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:QuietHoursController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:QuietHoursController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ResourceController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ResourceController"],
		beego.ControllerComments{
			Method: "StartPageList",
			Router: `/get-startpage-img`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ResourceController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ResourceController"],
		beego.ControllerComments{
			Method: "HomepageImgList",
			Router: `/get-homepage-img`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:StructureController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:StructureController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:StructureController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:StructureController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:StructureController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:StructureController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:TransactionController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:TransactionController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:TransactionController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:TransactionController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:TransactionController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:TransactionController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"],
		beego.ControllerComments{
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"],
		beego.ControllerComments{
			Method: "OtherLogin",
			Router: `/other-login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"],
		beego.ControllerComments{
			Method: "LogOut",
			Router: `/logout`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"],
		beego.ControllerComments{
			Method: "SendMessageCode",
			Router: `/send-message-code`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"],
		beego.ControllerComments{
			Method: "Register",
			Router: `/register`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"],
		beego.ControllerComments{
			Method: "ModifyPWD",
			Router: `/modify-pwd`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"],
		beego.ControllerComments{
			Method: "PersonalInformation",
			Router: `/personal-information`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserPositionController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserPositionController"],
		beego.ControllerComments{
			Method: "PositionList",
			Router: `/positionlist`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserPositionController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserPositionController"],
		beego.ControllerComments{
			Method: "SetPosition",
			Router: `/setposition`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserPositionController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserPositionController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserPositionController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserPositionController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserPositionController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserPositionController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
