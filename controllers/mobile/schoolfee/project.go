package schoolfee

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/mobile/schoolfee"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/schoolfee"
	"dev.model.360baige.com/action"
)

// Project API
type ProjectController struct {
	beego.Controller
}

// @Title 获取缴费项目列表接口
// @Description No Limit Project List 获取缴费项目列表接口
// @Success 200 {"code":200,"message":"获取缴费项目列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取缴费项目列表失败"}
// @router /nolimitlist [get]
func (c *ProjectController) ListOfNoLimitProject() {
	res := ListOfNoLimitProjectResponse{}
	access_token := c.GetString("access_token")
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
	var args2 action.FindByCond
	args2.CondList = append(args2.CondList, action.CondValue{
		Type: "And",
		Key:  "company_id",
		Val:  replyAccessToken.CompanyId,
	}, action.CondValue{
		Type: "And",
		Key:  "is_limit",
		Val:  1,
	}, action.CondValue{
		Type: "And",
		Key:  "status",
		Val:  0,
	})
	var replyProject []schoolfee.Project
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "ListByCond", args2, &replyProject)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "获取非限制缴费项目列表失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var listOfProject []Project = make([]Project, len(replyProject), len(replyProject))
	for index, pro := range replyProject {
		listOfProject[index] = Project{
			Id:         pro.Id,
			CreateTime: pro.CreateTime,
			UpdateTime: pro.UpdateTime,
			CompanyId:  pro.CompanyId,
			Name:       pro.Name,
			IsLimit:    pro.IsLimit,
			Desc:       pro.Desc,
			Link:       pro.Link,
			Status:     pro.Status,
		}
	}
	res.Code = ResponseNormal
	res.Message = "获取非限制缴费项目列表成功"
	res.Data.List = listOfProject
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 查询缴费信息接口
// @Description No Limit Project List 查询缴费信息接口
// @Success 200 {"code":200,"message":"查询缴费信息成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"查询缴费信息失败"}
// @router /search [get]
func (c *ProjectController) SearchProjectInfo() {
	res := SearchProjectInfoResponse{}
	access_token := c.GetString("access_token")
	search_type := c.GetString("search_type") // 1：身份证号码 其他：编码
	search_key := c.GetString("search_key")
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

	var args2 action.FindByCond
	args2.CondList = append(args2.CondList,
		action.CondValue{Type: "And", Key: "company_id", Val: replyAccessToken.CompanyId},
		action.CondValue{Type: "And", Key: "is_fee", Val: 0},
		action.CondValue{Type: "And", Key: "status", Val: 0},
	)
	if search_type == "1" {
		args2.CondList = append(args2.CondList, action.CondValue{Type: "And", Key: "id_card", Val: search_key})
	} else {
		args2.CondList = append(args2.CondList, action.CondValue{Type: "And", Key: "num", Val: search_key})
	}

	// 2.
	var replyRecord []schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "ListByCond", args2, &replyRecord)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "获取缴费项目列表失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	var project_id_s []int64 = make([]int64, len(replyRecord), len(replyRecord))
	for index, record := range replyRecord {
		project_id_s[index] = record.ProjectId
	}

	var args3 action.FindByCond
	args3.CondList = append(args3.CondList,
		action.CondValue{Type: "And", Key: "company_id", Val: replyAccessToken.CompanyId},
		action.CondValue{Type: "And", Key: "id__in", Val: project_id_s},
		//action.CondValue{Type: "And", Key: "status", Val: 0},
	)
	var replyProject []schoolfee.Project
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "ListByCond", args3, &replyProject)

	var map_project map[int64]schoolfee.Project = make(map[int64]schoolfee.Project)
	for _, pro := range replyProject {
		map_project[pro.Id] = pro
	}

	var list []RecordProject = make([]RecordProject, len(replyRecord), len(replyRecord))
	for index, record := range replyRecord {
		list[index] = RecordProject{
			Id:         record.Id,
			CreateTime: record.CreateTime,
			UpdateTime: record.UpdateTime,
			CompanyId:  record.CompanyId,
			ProjectId:  record.ProjectId,
			Name:       record.Name,
			ClassName:  record.ClassName,
			IdCard:     record.IdCard,
			Num:        record.Num,
			Phone:      record.Phone,
			Status:     record.Status,
			Price:      record.Price,
			IsFee:      record.IsFee,
			FeeTime:    record.FeeTime,
			Desc:       record.Desc,
			Project: Project{
				Id:         map_project[record.ProjectId].Id,
				CreateTime: map_project[record.ProjectId].CreateTime,
				UpdateTime: map_project[record.ProjectId].UpdateTime,
				CompanyId:  map_project[record.ProjectId].CompanyId,
				Name:       map_project[record.ProjectId].Name,
				IsLimit:    map_project[record.ProjectId].IsLimit,
				Desc:       map_project[record.ProjectId].Desc,
				Link:       map_project[record.ProjectId].Link,
				Status:     map_project[record.ProjectId].Status,
			},
		}
	}

	res.Code = ResponseNormal
	res.Message = "获取缴费项目列表成功"
	res.Data.List = list
	c.Data["json"] = res
	c.ServeJSON()
}
