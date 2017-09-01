package schoolfee

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/mobile/schoolfee"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/schoolfee"
	"dev.model.360baige.com/action"
	"encoding/json"
	"time"
	"strings"
	//"fmt"
)

// Record API
type RecordController struct {
	beego.Controller
}

// @Title 查询缴费历史接口
// @Description Limit Project List 查询缴费历史接口
// @Success 200 {"code":200,"message":"查询缴费历史成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   searchKey     query   string true       "查询值"
// @Failure 400 {"code":400,"message":"查询缴费历史失败"}
// @router /history [*]
func (c *RecordController) ListOfRecord() {
	type data RecordHistoryResponse
	accessToken := c.GetString("accessToken")
	searchKey := c.GetString("searchKey")
	pageSize, _ := c.GetInt64("pageSize")
	currentPage, _ := c.GetInt64("current")
	if accessToken == "" {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}
	//fmt.Println("--------------------")
	//fmt.Println(currentPage)
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
	var reply action.PageByCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "PageByCond", action.PageByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId},
			action.CondValue{Type: "And", Key: "is_fee", Val: 1},
			action.CondValue{Type: "And", Key: "status", Val: 0},
			action.CondValue{Type: "And", Key: searchTypeKey, Val: searchKey},
		},
		OrderBy:  []string{"id"},
		PageSize: pageSize,
		Current:  currentPage,
	}, &reply)
	replyRecord := []schoolfee.Record{}
	err = json.Unmarshal([]byte(reply.Json), &replyRecord)
	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "获取缴费项目列表失败"}
		c.ServeJSON()
		return
	}
	var project_ids []int64
	for _, record := range replyRecord {
		project_ids = append(project_ids, record.ProjectId)
	}

	var replyProject []schoolfee.Project
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Project", "ListByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId},
			action.CondValue{Type: "And", Key: "id__in", Val: project_ids},
		},
	}, &replyProject)

	var projectList map[int64]schoolfee.Project = make(map[int64]schoolfee.Project)
	for _, pro := range replyProject {
		projectList[pro.Id] = pro
	}

	var recordProjectList []RecordProject = make([]RecordProject, len(replyRecord), len(replyRecord))
	for index, record := range replyRecord {
		recordProjectList[index] = RecordProject{
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
		}
	}

	c.Data["json"] = data{Code: ResponseNormal, Message: "", Data: ListOfRecordProject{
		List:        recordProjectList,
		Total:       reply.Total,
		Current:     currentPage,
		CurrentSize: reply.CurrentSize,
		OrderBy:     reply.OrderBy,
		PageSize:    pageSize,
		SearchKey:   searchKey,
	}}
	c.ServeJSON()
}

// @Title 添加缴费接口
// @Description No limit Record Add 添加缴费接口
// @Success 200 {"code":200,"message":"添加缴费成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   name     query   string true       "姓名"
// @Param   projectId     query   string true       "项目id"
// @Param   phone     query   string true       "联系方式"
// @Param   price     query   string true       ""费用
// @Param   feeTime     query   string true       "缴费时间"
// @Param   desc     query   string true       "备注"
// @Param   key     query   string true       "编号、身份证号"
// @Failure 400 {"code":400,"message":"添加缴费失败"}
// @router /add [post]
func (c *RecordController) AddRecord() {
	type data AddRecordResponse
	accessToken := c.GetString("accessToken")
	projectId, _ := c.GetInt64("projectId", 0)
	name := c.GetString("name")
	no := c.GetString("key")
	num := ""
	id_card := ""
	if len(no) >= 15 {
		id_card = no
	} else {
		num = no
	}
	phone := c.GetString("phone")
	price, _ := c.GetInt64("price")
	var isFee int8 = 1
	feeTime, _ := c.GetInt64("feeTime")
	desc := c.GetString("desc")
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

	operationTime := time.Now().UnixNano() / 1e6
	args2 := &schoolfee.Record{
		CreateTime: operationTime,
		UpdateTime: operationTime,
		CompanyId:  replyUserPosition.CompanyId,
		ProjectId:  projectId,
		Name:       name,
		IdCard:     id_card,
		Num:        num,
		Phone:      phone,
		Desc:       desc,
		Price:      price,
		IsFee:      isFee,
		FeeTime:    feeTime,
		Status:     0,
	}
	var replyRecord schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "Add", args2, &replyRecord)
	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "添加失败"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = data{Code: ResponseNormal, Message: "添加成功", Data: AddRecord{
		Id: replyRecord.Id,
	}}
	c.ServeJSON()
	return
}

// @Title 添加缴费(多个)接口
// @Description  limit Record Add 添加缴费接口
// @Success 200 {"code":200,"message":"添加缴费成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   recordIds     query   string true       "缴费ids"
// @Param   feeTime     query   string true       "缴费时间"
// @Failure 400 {"code":400,"message":"添加缴费失败"}
// @router /addMulti [post]
func (c *RecordController) AddMultiRecord() {
	type data AddMultipleRecordResponse
	accessToken := c.GetString("accessToken")
	recordIds := c.GetString("recordIds")
	var isFee int8 = 1
	feeTime, _ := c.GetInt64("feeTime")
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
	if recordIds == "" {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "添加缴费失败"}
		c.ServeJSON()
		return
	}
	recordIdsArr := strings.Split(recordIds, ",")
	if len(recordIdsArr) == 0 {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "添加缴费失败"}
		c.ServeJSON()
		return
	}
	var replyRecords []schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "ListByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId},
			action.CondValue{Type: "And", Key: "id__in", Val: recordIdsArr},
			action.CondValue{Type: "And", Key: "is_fee", Val: 0},
		},
	}, &replyRecords)
	var args2 []schoolfee.Record
	operationTime := time.Now().UnixNano() / 1e6
	if len(replyRecords) == 0 {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "添加缴费失败"}
		c.ServeJSON()
		return
	}
	var temp schoolfee.Record
	for _, record := range replyRecords {
		temp.CreateTime = operationTime
		temp.UpdateTime = operationTime
		temp.CompanyId = record.CompanyId
		temp.ProjectId = record.ProjectId
		temp.Name = record.Name
		temp.ClassName = record.ClassName
		temp.IdCard = record.IdCard
		temp.Num = record.Num
		temp.Phone = record.Phone
		temp.Price = 0
		temp.Price = record.Price
		temp.FeeTime = feeTime
		temp.IsFee = isFee
		temp.FeeTime = feeTime
		temp.Desc = record.Desc
		args2 = append(args2, temp)
	}

	var replyRecord action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "AddMultiple", args2, &replyRecord)
	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "添加失败"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = data{Code: ResponseNormal, Message: "添加成功", Data: AddMultipleRecord{
		Num: replyRecord.Value,
	}}
	c.ServeJSON()
	return
}
