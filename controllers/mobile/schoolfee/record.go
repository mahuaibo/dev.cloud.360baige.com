package schoolfee

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/mobile/schoolfee"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/schoolfee"
	"dev.model.360baige.com/action"
	"encoding/json"
	"strings"
	"dev.cloud.360baige.com/utils"
	"regexp"
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
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	searchType := c.GetString("searchType", "num")
	searchKey := c.GetString("searchKey")
	pageSize, _ := c.GetInt64("pageSize", 50)
	currentPage, _ := c.GetInt64("current", 1)
	err := utils.Unable(map[string]string{"accessToken": "string:true", "searchKey": "string:true"}, c.Ctx.Input)
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

	r, _ := regexp.Compile("^[1-9]\\d{7}((0\\d)|(1[0-2]))(([0|1|2]\\d)|3[0-1])\\d{3}$|^[1-9]\\d{5}[1-9]\\d{3}((0\\d)|(1[0-2]))(([0|1|2]\\d)|3[0-1])\\d{3}([0-9]|X)$")
	if r.MatchString(searchKey) {
		searchType = "id_card"
	}
	var replyPageByCond action.PageByCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "PageByCond", action.PageByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId},
			action.CondValue{Type: "And", Key: "is_fee", Val: 1},
			action.CondValue{Type: "And", Key: "status", Val: 0},
			action.CondValue{Type: "And", Key: searchType, Val: searchKey},
		},
		OrderBy:  []string{"id"},
		PageSize: pageSize,
		Current:  currentPage,
	}, &replyPageByCond)
	jsonReplyRecord := []schoolfee.Record{}
	err = json.Unmarshal([]byte(replyPageByCond.Json), &jsonReplyRecord)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: Message(50002, "Record")}
		c.ServeJSON()
		return
	}

	var project_ids []int64
	for _, record := range jsonReplyRecord {
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

	var recordProjectList []RecordProject = make([]RecordProject, len(jsonReplyRecord), len(jsonReplyRecord))
	for index, record := range jsonReplyRecord {
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

	c.Data["json"] = data{Code: Normal, Message: Message(20003), Data: ListOfRecordProject{
		List:        recordProjectList,
		Total:       replyPageByCond.Total,
		Current:     currentPage,
		CurrentSize: replyPageByCond.CurrentSize,
		OrderBy:     replyPageByCond.OrderBy,
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
	currentTimestamp := utils.CurrentTimestamp()
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
	var isFee int = 1
	feeTime, _ := c.GetInt64("feeTime")
	desc := c.GetString("desc")

	err := utils.Unable(map[string]string{"accessToken": "string:true", "searchKey": "string:true"}, c.Ctx.Input)
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

	var replyRecord schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "Add", &schoolfee.Record{
		CreateTime: currentTimestamp,
		UpdateTime: currentTimestamp,
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
	}, &replyRecord)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: Message(50001)}
		c.ServeJSON()
		return
	}
	c.Data["json"] = data{Code: Normal, Message: Message(20004), Data: AddRecord{
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
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	recordIds := c.GetString("recordIds")
	feeTime, _ := c.GetInt64("feeTime")

	err := utils.Unable(map[string]string{"accessToken": "string:true", "recordIds": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: Message(40000, err.Error())}
		c.ServeJSON()
		return
	}

	recordIdsArr := strings.Split(recordIds, ",")
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

	var replyRecords []schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "ListByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId},
			action.CondValue{Type: "And", Key: "id__in", Val: recordIdsArr},
			action.CondValue{Type: "And", Key: "is_fee", Val: 0},
		},
	}, &replyRecords)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: Message(50001)}
		c.ServeJSON()
		return
	}
	var listOfRecord []schoolfee.Record
	for _, record := range replyRecords {
		listOfRecord = append(listOfRecord, schoolfee.Record{
			CreateTime: currentTimestamp,
			UpdateTime: currentTimestamp,
			CompanyId:  record.CompanyId,
			ProjectId:  record.ProjectId,
			Name:       record.Name,
			ClassName:  record.ClassName,
			IdCard:     record.IdCard,
			Num:        record.Num,
			Phone:      record.Phone,
			Price:      record.Price,
			FeeTime:    feeTime,
			IsFee:      1,
			Desc:       record.Desc,
		})
	}

	var replyRecord action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "AddMultiple", &listOfRecord, &replyRecord)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: Message(50001)}
		c.ServeJSON()
		return
	}
	c.Data["json"] = data{Code: Normal, Message: Message(20005), Data: AddMultipleRecord{
		Num: replyRecord.Value,
	}}
	c.ServeJSON()
	return
}
