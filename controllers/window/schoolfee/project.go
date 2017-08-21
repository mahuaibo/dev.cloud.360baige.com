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
// @Success 200 {"code":200,"message":"获取缴费项目成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取缴费项目失败"}
// @router /list [get]
func (c *ProjectController) ListOfProject() {
	res := ListOfProjectResponse{}
	access_token := c.GetString("access_token")
	pageSize, _ := c.GetInt64("page_size")
	currentPage, _ := c.GetInt64("current")
	if access_token == "" {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type: "And",
		Key:  "access_token",
		Val:  access_token,
	})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	// 2.
	var args2 action.PageByCond
	args2.CondList = append(args2.CondList, action.CondValue{
		Type: "And",
		Key:  "company_id",
		Val:  replyAccessToken.CompanyId,
	})
	args2.OrderBy = []string{"id"}
	args2.Cols = []string{"id", "create_time", "company_id", "name", "is_limit", "desc", "link", "status" }
	args2.PageSize = pageSize
	args2.Current = currentPage
	var replyProject []schoolfee.Project
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "ListByCond", args2, &replyProject)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "获取缴费项目失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	// 3.
	var listOfProject []Project = make([]Project, len(replyProject), len(replyProject))
	for index, pro := range replyProject {
		listOfProject[index] = Project{
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
	res.Code = ResponseNormal
	res.Message = "获取缴费项目成功"
	res.Data.List = listOfProject
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// @Title 添加校园收费项目接口
// @Description Project Add 添加校园收费项目接口
// @Success 200 {"code":200,"message":"添加缴费项目成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   name     query   string true       "项目名称"
// @Param   is_limit     query   string true       "是否限制缴费"
// @Param   desc     query   string true       "描述"
// @Param   link     query   string true       "描述链接"
// @Param   status     query   string true       "状态 -1注销 0正常"
// @Failure 400 {"code":400,"message":"添加缴费项目失败"}
// @router /add [post]
func (c *ProjectController) AddProject() {
	res := AddProjectResponse{}
	access_token := c.GetString("access_token")
	name := c.GetString("name")
	is_limit, _ := c.GetInt8("is_limit", 0)
	desc := c.GetString("desc")
	link := c.GetString("link")
	status, _ := c.GetInt8("status", 0)

	if access_token == "" {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type: "And",
		Key:  "access_token",
		Val:  access_token,
	})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	// 2.
	operationTime := time.Now().UnixNano() / 1e6
	args2 := &schoolfee.Project{
		CreateTime: operationTime,
		UpdateTime: operationTime,
		CompanyId:  replyAccessToken.CompanyId,
		Name:       name,
		IsLimit:    is_limit,
		Desc:       desc,
		Link:       link,
		Status:     status,
	}
	var replyProject schoolfee.Project
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "Add", args2, &replyProject)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "添加缴费项目失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	res.Code = ResponseNormal
	res.Message = "添加缴费项目成功"
	res.Data.Id = replyProject.Id
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// @Title 修改校园收费项目接口
// @Description Project Add 修改校园收费项目接口
// @Success 200 {"code":200,"message":"修改缴费项目成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   name     query   string true       "项目名称"
// @Param   is_limit     query   string true       "是否限制缴费"
// @Param   desc     query   string true       "描述"
// @Param   link     query   string true       "描述链接"
// @Param   status     query   string true       "状态 -1注销 0正常"
// @Failure 400 {"code":400,"message":"修改缴费项目失败"}
// @router /modify [post]
func (c *ProjectController) ModifyProject() {
	res := ModifyProjectResponse{}
	access_token := c.GetString("access_token")
	id, _ := c.GetInt64("id", 0)
	name := c.GetString("name")
	is_limit, _ := c.GetInt8("is_limit", 0)
	desc := c.GetString("desc")
	link := c.GetString("list")
	status, _ := c.GetInt8("status", 0)

	if access_token == "" {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type: "And",
		Key:  "access_token",
		Val:  access_token,
	})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	// 2.
	args2 := &schoolfee.Project{
		Id: id,
	}
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "FindById", args2, args2)

	if args2.CompanyId != replyAccessToken.CompanyId {
		res.Code = ResponseLogicErr
		res.Message = "非法操作"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	// 3.
	args3 := action.UpdateByIdCond{
		Id: []int64{args2.Id},
	}
	operationTime := time.Now().UnixNano() / 1e6
	args3.UpdateList = append(args3.UpdateList,
		action.UpdateValue{Key: "UpdateTime", Val: operationTime},
		action.UpdateValue{Key: "Name", Val: name},
		action.UpdateValue{Key: "IsLimit", Val: is_limit},
		action.UpdateValue{Key: "Desc", Val: desc},
		action.UpdateValue{Key: "Link", Val: link},
		action.UpdateValue{Key: "Status", Val: status},
	)

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "UpdateById", args3, &replyNum)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "修改缴费项目失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	res.Code = ResponseNormal
	res.Message = "修改缴费项目成功"
	res.Data.Count = replyNum.Value
	c.Data["json"] = res
	c.ServeJSON()
	return
}
