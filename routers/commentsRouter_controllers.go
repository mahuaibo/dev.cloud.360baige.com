package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

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
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:ApplicationController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
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
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"get"},
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
