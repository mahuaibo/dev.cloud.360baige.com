package schoolfee

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.cloud.360baige.com/utils"
	. "dev.model.360baige.com/http/window/schoolfee"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/schoolfee"
	"dev.model.360baige.com/action"
	"encoding/json"
)

// Project API
type ProjectController struct {
	beego.Controller
}

// @Title 校园收费列表
// @Description Project List 校园收费列表
// @Success 200 {"code":200,"message":"获取校园收费列表成功"}
// @Failure 400 {"code":400,"message":"获取校园收费列表失败"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   pageSize     query   string true       "页码数量 默认50"
// @Param   current     query   string true       "页码 默认1"
// @Param   accessToken     query   string true       "访问令牌"
// @router /list [post]
func (c *ProjectController) ListOfProject() {
	type data ListOfProjectResponse
	accessToken := c.GetString("accessToken")
	currentTimestamp := utils.CurrentTimestamp()
	pageSize, _ := c.GetInt64("pageSize", 50)
	currentPage, _ := c.GetInt64("current", 1)

	err := utils.Unable(map[string]string{"accessToken": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: Message(40000, err.Error())}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{CondList: []action.CondValue{action.CondValue{Type: "And", Key: "access_token", Val: accessToken }, action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTimestamp }, }, Fileds: []string{"id", "user_id", "company_id", "type"}, }, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: Message(50000)}
		c.ServeJSON()
		return
	}
	if replyUserPosition.UserId == 0 {
		c.Data["json"] = data{Code: ErrorPower, Message: Message(30000)}
		c.ServeJSON()
		return
	}

	var replyPageByCond action.PageByCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "PageByCond", &action.PageByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId },
			action.CondValue{Type: "And", Key: "status__gt", Val: -1},
		},
		Cols:     []string{"id", "create_time", "company_id", "name", "is_limit", "desc", "link", "status" },
		OrderBy:  []string{"id"},
		PageSize: pageSize,
		Current:  currentPage,
	}, &replyPageByCond)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: Message(50001, "Project")}
		c.ServeJSON()
		return
	}

	var jsonProjectList []schoolfee.Project
	err = json.Unmarshal([]byte(replyPageByCond.Json), &jsonProjectList)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: Message(50002, "Project")}
		c.ServeJSON()
		return
	}

	var projectList []Project = make([]Project, len(jsonProjectList), len(jsonProjectList))
	for index, pro := range jsonProjectList {
		projectList[index] = Project{
			Id:         pro.Id,
			CreateTime: utils.Datetime(pro.CreateTime, "2006-01-02 03:04"),
			CompanyId:  pro.CompanyId,
			Name:       pro.Name,
			IsLimit:    pro.IsLimit,
			Desc:       pro.Desc,
			Link:       pro.Link,
			Status:     pro.Status,
		}
	}

	c.Data["json"] = data{Code: Normal, Message: Message(20000), Data: ListOfProject{
		List:     projectList,
		Total:    replyPageByCond.Total,
		PageSize: pageSize,
		Current:  currentPage,
	}}
	c.ServeJSON()
	return
}

// @Title 添加校园收费项目接口
// @Description Project Add 添加校园收费项目接口
// @Success 200 {"code":200,"message":"添加缴费项目成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   name     query   string true       "项目名称"
// @Param   isLimit     query   string true       "是否限制缴费"
// @Param   desc     query   string true       "描述"
// @Param   link     query   string true       "描述链接"
// @Param   status     query   string true       "状态 -1注销 0正常"
// @Failure 400 {"code":400,"message":"添加缴费项目失败"}
// @router /add [post]
func (c *ProjectController) AddProject() {
	type data AddProjectResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	name := c.GetString("name")
	isLimit, _ := c.GetInt8("isLimit", 0)
	desc := c.GetString("desc")
	link := c.GetString("link")
	status, _ := c.GetInt8("status", 0)

	err := utils.Unable(map[string]string{
		"accessToken": "string:true",
		"name":        "string:true",
	}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: Message(40000, err.Error())}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{CondList: []action.CondValue{action.CondValue{Type: "And", Key: "access_token", Val: accessToken }, action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTimestamp }, }, Fileds: []string{"id", "user_id", "company_id", "type"}, }, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyProject schoolfee.Project
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "Add", &schoolfee.Project{
		CreateTime: currentTimestamp,
		UpdateTime: currentTimestamp,
		CompanyId:  replyUserPosition.CompanyId,
		Name:       name,
		IsLimit:    isLimit,
		Desc:       desc,
		Link:       link,
		Status:     status,
	}, &replyProject)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "添加缴费项目失败"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = data{Code: Normal, Message: "添加缴费项目成功", Data: AddProject{
		Id: replyProject.Id,
	}}
	c.ServeJSON()
	return
}

// @Title 修改校园收费项目接口
// @Description Project Add 修改校园收费项目接口
// @Success 200 {"code":200,"message":"修改缴费项目成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   id     query   string true       "项目id"
// @Param   name     query   string true       "项目名称"
// @Param   isLimit     query   string true       "是否限制缴费"
// @Param   desc     query   string true       "描述"
// @Param   link     query   string true       "描述链接"
// @Param   status     query   string true       "状态 -1注销 0正常"
// @Failure 400 {"code":400,"message":"修改缴费项目失败"}
// @router /modify [post]
func (c *ProjectController) ModifyProject() {
	type data ModifyProjectResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	id, _ := c.GetInt64("id", 0)
	name := c.GetString("name")
	isLimit, _ := c.GetInt8("isLimit", 0)
	desc := c.GetString("desc")
	link := c.GetString("list")
	status, _ := c.GetInt8("status", 0)

	err := utils.Unable(map[string]string{"accessToken": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: Message(40000, err.Error())}
		c.ServeJSON()
		return
	}
	var replyUserPosition user.UserPosition
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{CondList: []action.CondValue{action.CondValue{Type: "And", Key: "access_token", Val: accessToken }, action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTimestamp }, }, Fileds: []string{"id", "user_id", "company_id", "type"}, }, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效2"}
		c.ServeJSON()
		return
	}

	var replyProject schoolfee.Project
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "FindById", &schoolfee.Project{
		Id:        id,
		CompanyId: replyUserPosition.CompanyId,
	}, &replyProject)

	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效3"}
		c.ServeJSON()
		return
	}
	if replyProject.Id == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效4"}
		c.ServeJSON()
		return
	}

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "UpdateById", &action.UpdateByIdCond{
		Id: []int64{replyProject.Id},
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "UpdateTime", Val: currentTimestamp},
			action.UpdateValue{Key: "Name", Val: name},
			action.UpdateValue{Key: "IsLimit", Val: isLimit},
			action.UpdateValue{Key: "Desc", Val: desc},
			action.UpdateValue{Key: "Link", Val: link},
			action.UpdateValue{Key: "Status", Val: status},
		},
	}, &replyNum)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "修改缴费项目失败"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = data{Code: Normal, Message: "修改缴费项目成功", Data: ModifyProject{
		Count: replyNum.Value,
	}}
	c.ServeJSON()
	return
}
