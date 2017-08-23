package schoolfee

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/schoolfee"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/schoolfee"
	"dev.model.360baige.com/action"
	"time"
)

// Project API
type ProjectController struct {
	beego.Controller
}

// @Title 校园收费列表接口
// @Description Project List 校园收费列表接口
// @Success 200 {"code":200,"message":"获取缴费项目成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取缴费项目失败"}
// @router /list [post]
func (c *ProjectController) ListOfProject() {
	type data ListOfProjectResponse
	accessToken := c.GetString("accessToken")
	pageSize, _ := c.GetInt64("pageSize")
	currentPage, _ := c.GetInt64("current")
	if accessToken == "" {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "accessToken", Val: accessToken },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	var replyProject []schoolfee.Project
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "ListByCond", &action.PageByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId },
		},
		OrderBy:  []string{"id"},
		Cols:     []string{"id", "create_time", "company_id", "name", "isLimit", "desc", "link", "status" },
		PageSize: pageSize,
		Current:  currentPage,
	}, &replyProject)
	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "获取缴费项目失败"}
		c.ServeJSON()
		return
	}

	var projectList []Project = make([]Project, len(replyProject), len(replyProject))
	for index, pro := range replyProject {
		projectList[index] = Project{
			Id:         pro.Id,
			CreateTime: time.Unix(pro.CreateTime/1000, 0).Format("2006-01-02"),
			CompanyId:  pro.CompanyId,
			Name:       pro.Name,
			IsLimit:    pro.IsLimit,
			Desc:       pro.Desc,
			Link:       pro.Link,
			Status:     pro.Status,
		}
	}
	c.Data["json"] = data{Code: ResponseNormal, Message: "获取缴费项目成功", Data: ListOfProject{
		List: projectList,
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
	accessToken := c.GetString("accessToken")
	name := c.GetString("name")
	isLimit, _ := c.GetInt8("isLimit", 0)
	desc := c.GetString("desc")
	link := c.GetString("link")
	status, _ := c.GetInt8("status", 0)
	if accessToken == "" {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList:
		[]action.CondValue{
			action.CondValue{Type: "And", Key: "accessToken", Val: accessToken },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	operationTime := time.Now().UnixNano() / 1e6
	var replyProject schoolfee.Project
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "Add", schoolfee.Project{
		CreateTime: operationTime,
		UpdateTime: operationTime,
		CompanyId:  replyUserPosition.CompanyId,
		Name:       name,
		IsLimit:    isLimit,
		Desc:       desc,
		Link:       link,
		Status:     status,
	}, &replyProject)
	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "添加缴费项目失败"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = data{Code: ResponseNormal, Message: "添加缴费项目成功", Data: AddProject{
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
	accessToken := c.GetString("accessToken")
	id, _ := c.GetInt64("id", 0)
	name := c.GetString("name")
	isLimit, _ := c.GetInt8("isLimit", 0)
	desc := c.GetString("desc")
	link := c.GetString("list")
	status, _ := c.GetInt8("status", 0)

	if accessToken == "" {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌无效1"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "accessToken", Val: accessToken },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌无效2"}
		c.ServeJSON()
		return
	}

	var replyProject schoolfee.Project
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "FindById", schoolfee.Project{
		Id:        id,
		CompanyId: replyUserPosition.CompanyId,
	}, &replyProject)

	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌无效3"}
		c.ServeJSON()
		return
	}
	if replyProject.Id == 0 {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌无效4"}
		c.ServeJSON()
		return
	}

	operationTime := time.Now().UnixNano() / 1e6
	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "UpdateById", action.UpdateByIdCond{
		Id: []int64{replyProject.Id},
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "UpdateTime", Val: operationTime},
			action.UpdateValue{Key: "Name", Val: name},
			action.UpdateValue{Key: "IsLimit", Val: isLimit},
			action.UpdateValue{Key: "Desc", Val: desc},
			action.UpdateValue{Key: "Link", Val: link},
			action.UpdateValue{Key: "Status", Val: status},
		},
	}, &replyNum)
	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "修改缴费项目失败"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = data{Code: ResponseNormal, Message: "修改缴费项目成功", Data: ModifyProject{
		Count: replyNum.Value,
	}}
	c.ServeJSON()
	return
}
