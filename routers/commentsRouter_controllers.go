package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CompanyController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:CompanyController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/detail`,
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

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonController"],
		beego.ControllerComments{
			Method: "PersonList",
			Router: `/personList`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonRelationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonRelationController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonRelationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonRelationController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonRelationController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonRelationController"],
		beego.ControllerComments{
			Method: "AssociatedList",
			Router: `/associatedList`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonStructureController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:PersonStructureController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:StructureController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:StructureController"],
		beego.ControllerComments{
			Method: "Add",
			Router: `/add`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:StructureController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:StructureController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:StructureController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:StructureController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:StructureController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:StructureController"],
		beego.ControllerComments{
			Method: "StructureList",
			Router: `/structureList`,
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
			Method: "Login",
			Router: `/login`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"],
		beego.ControllerComments{
			Method: "Logout",
			Router: `/logout`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"],
		beego.ControllerComments{
			Method: "Cancel",
			Router: `/cancel`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"],
		beego.ControllerComments{
			Method: "Activation",
			Router: `/activation/:username`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"],
		beego.ControllerComments{
			Method: "Detail",
			Router: `/detail`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"],
		beego.ControllerComments{
			Method: "Modify",
			Router: `/modify`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"] = append(beego.GlobalControllerRouter["dev.cloud.360baige.com/controllers:UserController"],
		beego.ControllerComments{
			Method: "ModifyPassword",
			Router: `/modifyPassword`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
