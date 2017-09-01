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
// @Success 200 {"code":200,"message":"获取缴费项目列表成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取缴费项目列表失败"}
// @router /noLimitList [*]
func (c *ProjectController) ListOfNoLimitProject() {
	type data ListOfNoLimitProjectResponse
	accessToken := c.GetString("accessToken")
	if accessToken == "" {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌无效"}
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
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	var replyProject []schoolfee.Project
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "ListByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId },
			action.CondValue{Type: "And", Key: "is_limit", Val: 1 },
			action.CondValue{Type: "And", Key: "status", Val: 0 },
		},
	}, &replyProject)
	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "获取非限制缴费项目列表失败"}
		c.ServeJSON()
		return
	}

	var projectList []Project = make([]Project, len(replyProject), len(replyProject))
	for index, pro := range replyProject {
		projectList[index] = Project{
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
	c.Data["json"] = data{Code: ResponseNormal, Message: "获取非限制缴费项目列表成功", Data: ListOfNoLimitProject{
		List: projectList,
	}}
	c.ServeJSON()
}

// @Title 查询缴费信息接口
// @Description No Limit Project List 查询缴费信息接口
// @Success 200 {"code":200,"message":"查询缴费信息成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   searchKey     query   string true       "查询值"
// @Failure 400 {"code":400,"message":"查询缴费信息失败"}
// @router /search [*]
func (c *ProjectController) SearchProjectInfo() {
	type data SearchProjectInfoResponse
	accessToken := c.GetString("accessToken")
	//searchType := c.GetString("searchType", "0") // 1：身份证号码 其他：编码
	searchKey := c.GetString("searchKey")
	if accessToken == "" {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌无效"}
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
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}
	searchTypeKey := "num"
	if len(searchKey) >= 15 {
		searchTypeKey = "id_card"
	}
	var replyRecord []schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "ListByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId},
			action.CondValue{Type: "And", Key: "is_fee", Val: 0},
			action.CondValue{Type: "And", Key: "status", Val: 0},
			action.CondValue{Type: "And", Key: searchTypeKey, Val: searchKey},
		},
	}, &replyRecord)
	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "获取缴费项目列表失败"}
		c.ServeJSON()
		return
	}
	var project_ids []int64
	for _, record := range replyRecord {
		project_ids = append(project_ids, record.ProjectId)
	}
	var replyFeeRecord []schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "ListByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId},
			action.CondValue{Type: "And", Key: "project_id__in", Val: project_ids},
			action.CondValue{Type: "And", Key: "is_fee", Val: 1},
			action.CondValue{Type: "And", Key: "status", Val: 0},
			action.CondValue{Type: "And", Key: searchTypeKey, Val: searchKey},
		},
	}, &replyFeeRecord)
	var feeRecordIdsList map[int64]int64 = make(map[int64]int64)
	for _, feeRe := range replyFeeRecord {
		feeRecordIdsList[feeRe.ProjectId] = feeRe.Id
	}
	var replyProject []schoolfee.Project
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "ListByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId},
			action.CondValue{Type: "And", Key: "id__in", Val: project_ids},
			//action.CondValue{Type: "And", Key: "status", Val: 0},
		},
	}, &replyProject)

	var projectList map[int64]schoolfee.Project = make(map[int64]schoolfee.Project)
	for _, pro := range replyProject {
		projectList[pro.Id] = pro
	}

	var recordProjectList []RecordProject
	for _, record := range replyRecord {
		if feeRecordIdsList[record.ProjectId] > 0 {
		} else {
			recordProjectList = append(recordProjectList, RecordProject{
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
					Id:         projectList[record.ProjectId].Id,
					CreateTime: projectList[record.ProjectId].CreateTime,
					UpdateTime: projectList[record.ProjectId].UpdateTime,
					CompanyId:  projectList[record.ProjectId].CompanyId,
					Name:       projectList[record.ProjectId].Name,
					IsLimit:    projectList[record.ProjectId].IsLimit,
					Desc:       projectList[record.ProjectId].Desc,
					Link:       projectList[record.ProjectId].Link,
					Status:     projectList[record.ProjectId].Status,
				},
			})
		}

	}

	c.Data["json"] = data{Code: ResponseNormal, Message: "", Data: ListOfRecordProject{
		List: recordProjectList,
	}}
	c.ServeJSON()
}

// @Title 收费详情接口
// @Description 收费详情接口
// @Success 200 {"code":200,"message":"获取信息详情成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"获取信息详情失败"}
// @router /detail [*]
func (c *ProjectController) Detail() {
	type data ProjectDetailResponse
	accessToken := c.GetString("accessToken")
	if accessToken == "" {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌无效"}
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
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}
	projectId, _ := c.GetInt64("id")
	if projectId == 0 {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "获取信息失败"}
		c.ServeJSON()
		return
	}
	var reply schoolfee.Project
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "FindById", &schoolfee.Project{
		Id: projectId,
	}, &reply)

	if err != nil {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取信息失败"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: ResponseNormal, Message: "获取信息成功", Data: Project{
		Id:         reply.Id,
		CreateTime: reply.CreateTime,
		UpdateTime: reply.UpdateTime,
		CompanyId:  reply.CompanyId,
		Name:       reply.Name,
		IsLimit:    reply.IsLimit,
		Desc:       reply.Desc,
		Link:       reply.Link,
		Status:     reply.Status,
	}}
	c.ServeJSON()
	return
}
