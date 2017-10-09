package schoolfee

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/schoolfee"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/schoolfee"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/utils"
	"strings"
	"os"
	"io"
	"github.com/xuri/excelize"
	"strconv"
	"github.com/tealeg/xlsx"
	"net/http"
	"sort"
	"fmt"
	"dev.cloud.360baige.com/log"
	"encoding/json"
)

// Record API
type RecordController struct {
	beego.Controller
}

// @Title 校园收费记录列表接口
// @Description Project List 校园收费记录列表接口
// @Success 200 {"code":200,"message":"获取缴费项目记录成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   projectId     query   int true       "项目ID"
// @Failure 400 {"code":400,"message":"获取缴费项目记录失败"}
// @router /list [post]
func (c *RecordController) ListOfRecord() {

	type data ListOfRecordResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	projectId := c.GetString("projectId")
	pageSize, _ := c.GetInt64("pageSize", 50)
	currentPage, _ := c.GetInt64("current", 1)
	err := utils.Unable(map[string]string{"accessToken": "string:true", "projectId": "int:true"}, c.Ctx.Input)
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

	var replyPageByCond action.PageByCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "PageByCond", action.PageByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId},
			action.CondValue{Type: "And", Key: "project_id", Val: projectId},
			action.CondValue{Type: "And", Key: "status__gt", Val: -1},
		},
		Cols:     []string{"id", "create_time", "company_id", "project_id", "name", "class_name", "id_card", "num", "phone", "price", "is_fee", "fee_time", "desc", "status" },
		OrderBy:  []string{"id"},
		PageSize: pageSize,
		Current:  currentPage,
	}, &replyPageByCond)

	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "获取缴费项目记录失败"}
		c.ServeJSON()
		return
	}

	var replyRecord []schoolfee.Record
	err = json.Unmarshal([]byte(replyPageByCond.Json), &replyRecord)
	listOfRecord := make([]Record, len(replyRecord), len(replyRecord))
	for index, rec := range replyRecord {
		resFeeTime := "无"
		if rec.IsFee == 1 {
			resFeeTime = utils.Datetime(rec.FeeTime, "2006-01-02 15:04")
		}
		listOfRecord[index] = Record{
			Id:         rec.Id,
			CreateTime: utils.Datetime(rec.CreateTime, "2006-01-02 15:04"),
			CompanyId:  rec.CompanyId,
			ProjectId:  rec.ProjectId,
			Name:       rec.Name,
			ClassName:  rec.ClassName,
			IdCard:     rec.IdCard,
			Num:        rec.Num,
			Phone:      rec.Phone,
			Price:      rec.Price,
			IsFee:      rec.IsFee,
			FeeTime:    resFeeTime,
			Desc:       rec.Desc,
			Status:     rec.Status,
		}
	}
	c.Data["json"] = data{Code: Normal, Message: "获取缴费项目记录成功", Data: ListOfRecord{
		List:     listOfRecord,
		Total:    replyPageByCond.Total,
		PageSize: pageSize,
		Current:  currentPage,
	}}
	c.ServeJSON()
	return
}

// @Title 添加收费名单接口
// @Description Record Add 添加收费名单接口
// @Success 200 {"code":200,"message":"添加收费名单成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   name     query   string true       "项目名称"
// @Param   isLimit     query   string true       "是否限制缴费"
// @Param   desc     query   string true       "描述"
// @Param   link     query   string true       "描述链接"
// @Param   status     query   string true       "状态 -1注销 0正常"
// @Failure 400 {"code":400,"message":"添加收费名单失败"}
// @router /add [post]
func (c *RecordController) AddRecord() {
	type data AddRecordResponse
	currentTimestamp := utils.CurrentTimestamp()

	accessToken := c.GetString("accessToken")
	projectId, _ := c.GetInt64("projectId", 0)
	name := c.GetString("name")
	className := c.GetString("className")
	idCard := c.GetString("idCard")
	num := c.GetString("num")
	phone := c.GetString("phone")
	price, _ := c.GetInt64("price")
	isFee, _ := c.GetInt("isFee")
	feeTime, _ := c.GetInt64("feeTime")
	desc := c.GetString("desc")
	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{CondList: []action.CondValue{action.CondValue{Type: "And", Key: "access_token", Val: accessToken }, action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTimestamp }, }, Fileds: []string{"id", "user_id", "company_id", "type"}, }, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌失效"}
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
		ClassName:  className,
		IdCard:     idCard,
		Num:        num,
		Phone:      phone,
		Desc:       desc,
		Price:      price,
		IsFee:      isFee,
		FeeTime:    feeTime,
		Status:     0,
	}, &replyRecord)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "添加收费名单失败"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = data{Code: Normal, Message: "添加收费名单成功", Data: AddRecord{
		Id: replyRecord.Id,
	}}
	c.ServeJSON()
	return
}

// @Title 查看缴费记录接口
// @Description Project Add 查看缴费记录接口
// @Success 200 {"code":200,"message":"查看缴费记录成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   recordId     query   string true       "记录ID"
// @Failure 400 {"code":400,"message":"查看缴费记录失败"}
// @router /detail [post]
func (c *RecordController) DetailRecord() {
	type data DetailRecordResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	recordId, _ := c.GetInt64("recordId", 0)
	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{CondList: []action.CondValue{action.CondValue{Type: "And", Key: "access_token", Val: accessToken }, action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTimestamp }, }, Fileds: []string{"id", "user_id", "company_id", "type"}, }, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	var replyRecord schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "FindById", &schoolfee.Record{
		Id: recordId,
	}, &replyRecord)

	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "修改缴费项目失败"}
		c.ServeJSON()
		return
	}

	resFeeTime := "无"
	if replyRecord.IsFee == 1 {
		resFeeTime = utils.Datetime(replyRecord.FeeTime, "2006-01-02 15:04")
	}

	log.Println("resFeeTime:", resFeeTime)

	c.Data["json"] = data{Code: Normal, Message: "修改缴费项目成功", Data: DetailRecord{
		Data: Record{
			Id:         replyRecord.Id,
			CreateTime: utils.Datetime(replyRecord.CreateTime, "2006-01-02 15:04"),
			UpdateTime: replyRecord.UpdateTime,
			CompanyId:  replyRecord.CompanyId,
			ProjectId:  replyRecord.ProjectId,
			Name:       replyRecord.Name,
			ClassName:  replyRecord.ClassName,
			IdCard:     replyRecord.IdCard,
			Num:        replyRecord.Num,
			Phone:      replyRecord.Phone,
			Status:     replyRecord.Status,
			Price:      replyRecord.Price,
			IsFee:      replyRecord.IsFee,
			FeeTime:    resFeeTime,
			Desc:       replyRecord.Desc,
		},
	}}
	c.ServeJSON()
	return
}

// @Title 修改缴费记录接口
// @Description Project Add 修改缴费记录接口
// @Success 200 {"code":200,"message":"修改缴费记录成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   recordId     query   string true       "记录ID"
// @Param   className     query   string true       "班级名称"
// @Failure 400 {"code":400,"message":"修改缴费记录失败"}
// @router /modify [post]
func (c *RecordController) ModifyRecord() {
	type data ModifyRecordResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	recordId, _ := c.GetInt64("recordId", 0)
	name := c.GetString("name")
	className := c.GetString("className")
	idCard := c.GetString("idCard")
	num := c.GetString("num")
	phone := c.GetString("phone")
	price, _ := c.GetFloat("price")
	isFee, _ := c.GetInt("isFee")
	feeTime, _ := c.GetInt64("feeTime")
	desc := c.GetString("desc")
	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{CondList: []action.CondValue{action.CondValue{Type: "And", Key: "access_token", Val: accessToken }, action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTimestamp }, }, Fileds: []string{"id", "user_id", "company_id", "type"}, }, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	var replyRecord schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "FindById", &schoolfee.Record{
		Id: recordId,
	}, &replyRecord)

	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "修改缴费项目失败"}
		c.ServeJSON()
		return
	}

	if replyRecord.CompanyId != replyUserPosition.CompanyId {
		c.Data["json"] = data{Code: ErrorLogic, Message: "非法操作"}
		c.ServeJSON()
		return
	}

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "UpdateById", &action.UpdateByIdCond{
		Id: []int64{replyRecord.Id},
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "update_time", Val: currentTimestamp},
			action.UpdateValue{Key: "name", Val: name},
			action.UpdateValue{Key: "class_name", Val: className},
			action.UpdateValue{Key: "id_card", Val: idCard},
			action.UpdateValue{Key: "num", Val: num},
			action.UpdateValue{Key: "phone", Val: phone},
			action.UpdateValue{Key: "price", Val: price},
			action.UpdateValue{Key: "is_fee", Val: isFee},
			action.UpdateValue{Key: "fee_time", Val: feeTime},
			action.UpdateValue{Key: "desc", Val: desc},
		},
	}, &replyNum)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "修改缴费记录失败"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = data{Code: Normal, Message: "修改缴费记录成功", Data: ModifyRecord{
		Count: replyNum.Value,
	}}
	c.ServeJSON()
	return
}

// @Title 删除校园收费记录接口
// @Description Delete Record 删除校园收费记录接口
// @Success 200 {"code":200,"message":"删除缴费项目记录成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   recordIds     query   int true       "项目记录IDs"
// @Failure 400 {"code":400,"message":"删除缴费项目记录失败"}
// @router /delete [post]
func (c *RecordController) DeleteRecord() {
	type data DeleteRecordResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	recordIds := c.GetString("recordIds")
	log.Println("recordIds:", recordIds)

	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{CondList: []action.CondValue{action.CondValue{Type: "And", Key: "access_token", Val: accessToken }, action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTimestamp }, }, Fileds: []string{"id", "user_id", "company_id", "type"}, }, &replyUserPosition)

	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyRecord action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "DeleteById", &action.DeleteByIdCond{
		Value: utils.StrArrToInt64Arr(strings.Split(recordIds, ",")),
	}, &replyRecord)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "删除缴费项目记录失败"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = data{Code: Normal, Message: "删除缴费项目记录成功", Data: DeleteRecord{
		Count: replyRecord.Value,
	}}
	c.ServeJSON()
	return
}

// @Title 上传缴费名单接口
// @Description Delete Record 上传缴费名单接口
// @Success 200 {"code":200,"message":"上传缴费名单成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   projectId     query   string true       "项目Id"
// @Failure 400 {"code":400,"message":"上传缴费名单失败"}
// @router /upload [options,post]
func (c *RecordController) UploadRecord() {
	type data UploadRecordResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	projectId, _ := c.GetInt64("projectId")
	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{CondList: []action.CondValue{action.CondValue{Type: "And", Key: "access_token", Val: accessToken }, action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTimestamp }, }, Fileds: []string{"id", "user_id", "company_id", "type"}, }, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	var replyRecord action.Num
	if c.Ctx.Request.Method == "POST" {
		formFile, header, err := c.Ctx.Request.FormFile("uploadFile")
		if err != nil {
			c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效"}
			return
		}
		defer formFile.Close()
		objectKey := utils.RandomName("file_", header.Filename)
		// 创建保存文件
		destFile, err := os.Create(objectKey)
		if err != nil {
			c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效"}
			return
		}
		defer destFile.Close()
		// 读取表单文件，写入保存文件
		_, err = io.Copy(destFile, formFile)
		if err != nil {
			c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效"}
			return
		}
		xlsx, err := excelize.OpenFile(objectKey)
		if err != nil {
			os.Exit(1)
		}
		// Get sheet index.
		index := xlsx.GetSheetIndex("Sheet1")
		rows := xlsx.GetRows("sheet" + strconv.Itoa(index))
		var argsRecordList []schoolfee.Record = make([]schoolfee.Record, len(rows) - 1)
		for key, row := range rows {
			if key > 0 {
				Price, err := strconv.ParseInt(row[5], 10, 64)
				if err != nil {
					c.Data["json"] = data{Code: ErrorLogic, Message: "上传缴费名单失败"}
					c.ServeJSON()
					return
				}
				argsRecordList[key - 1] = schoolfee.Record{
					CreateTime: currentTimestamp,
					UpdateTime: currentTimestamp,
					CompanyId:  replyUserPosition.CompanyId,
					ProjectId:  projectId,
					Name:       row[0],
					ClassName:  row[1],
					IdCard:     row[2],
					Num:        row[3],
					Phone:      row[4],
					Price:      Price,
					Desc:       row[6],
					Status:     0,
					IsFee:      0,
					FeeTime:    0,
				}
			}
		}

		err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "AddMultiple", &argsRecordList, &replyRecord)
		if err != nil {
			c.Data["json"] = data{Code: ErrorLogic, Message: "上传缴费名单失败"}
			c.ServeJSON()
			return
		}
		defer func() {
			os.Remove(objectKey)
		}()
	}

	c.Data["json"] = data{Code: Normal, Message: "上传缴费名单成功", Data: UploadRecord{
		Count: replyRecord.Value,
	}}
	c.ServeJSON()
	return
}

// @Title 下载缴费记录接口
// @Description Delete Record 下载缴费记录接口
// @Success 200 {"code":200,"message":"下载缴费记录成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"下载缴费记录失败"}
// @router /download [get,post]
func (c *RecordController) DownloadRecord() {
	type data DownloadRecordResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	projectId, _ := c.GetInt64("projectId")
	classNames := c.GetString("classNames")
	isFees := c.GetString("isFees")
	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{
		CondList:[]action.CondValue{action.CondValue{Type: "And", Key: "access_token", Val: accessToken },
			action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTimestamp },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}
	var replyRecord []schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "ListByCond", &action.ListByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId},
			action.CondValue{Type: "And", Key: "class_name__in", Val: strings.Split(classNames, ",")},
			action.CondValue{Type: "And", Key: "is_fee__in", Val: strings.Split(isFees, ",")},
			action.CondValue{Type: "And", Key: "project_id", Val: projectId},
		},
	}, &replyRecord)

	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "下载缴费记录失败"}
		c.ServeJSON()
		return
	}

	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	row := sheet.AddRow()
	row.SetHeightCM(1) //设置每行的高度
	row.AddCell().Value = "姓名"
	row.AddCell().Value = "班级"
	row.AddCell().Value = "身份证号码"
	row.AddCell().Value = "编号"
	row.AddCell().Value = "联系电话"
	row.AddCell().Value = "应缴费用"
	row.AddCell().Value = "是否缴费"
	row.AddCell().Value = "缴费时间"
	row.AddCell().Value = "备注"
	//姓名、班级、身份证号码、编号、联系电话、应缴费用、是否缴费、缴费时间、备注
	for _, rec := range replyRecord {
		log.Println("姓名")
		row := sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度
		row.AddCell().Value = rec.Name
		row.AddCell().Value = rec.ClassName
		row.AddCell().Value = rec.IdCard
		row.AddCell().Value = rec.Num
		row.AddCell().Value = rec.Phone
		row.AddCell().Value = strconv.FormatInt(rec.Price, 10)
		if rec.IsFee == 1 {
			row.AddCell().Value = "已缴费"
			row.AddCell().Value = utils.Datetime(rec.FeeTime, "2006-01-02 15:04")
		} else {
			row.AddCell().Value = "未缴费"
			row.AddCell().Value = "无"
		}
		row.AddCell().Value = rec.Desc
	}
	objectKey := utils.RandomName("file_", ".xlsx")
	err = file.Save(objectKey)
	if err != nil {
		panic(err)
	}

	c.Ctx.Output.Header("Accept-Ranges", "bytes")
	c.Ctx.Output.Header("Content-Disposition", "attachment; filename=" + fmt.Sprintf("%s", objectKey)) //文件名
	c.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Ctx.Output.Header("Pragma", "no-cache")
	c.Ctx.Output.Header("Expires", "0")
	//最主要的一句
	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, objectKey)
	c.Data["json"] = data{Code: Normal, Message: "下载缴费记录成功"}
	c.ServeJSON()
	defer func() {
		os.Remove(objectKey)
	}()
	return
}

// @Title 班级列表接口
// @Description Project List 班级列表接口
// @Success 200 {"code":200,"message":"获取班级列表成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   projectId     query   int true       "项目ID"
// @Failure 400 {"code":400,"message":"获取班级列表失败"}
// @router /classList [post]
func (c *RecordController) ClassList() {
	type data ClassListOfRecordResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	projectId := c.GetString("projectId")
	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{CondList: []action.CondValue{action.CondValue{Type: "And", Key: "access_token", Val: accessToken }, action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTimestamp }, }, Fileds: []string{"id", "user_id", "company_id", "type"}, }, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	var replyRecord []schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "ListByCond", &action.PageByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId},
			action.CondValue{Type: "And", Key: "project_id", Val: projectId},
		},
		Cols: []string{"class_name" },
	}, &replyRecord)

	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "获取班级列表失败"}
		c.ServeJSON()
		return
	}

	var listOfRecord []string
	for _, rec := range replyRecord {
		listOfRecord = append(listOfRecord, rec.ClassName)
	}

	sort.Strings(listOfRecord)
	c.Data["json"] = data{Code: Normal, Message: "获取班级列表成功", Data: RemoveDuplicatesAndEmpty(listOfRecord)}
	c.ServeJSON()
	return
}

func RemoveDuplicatesAndEmpty(a []string) (ret []map[string]string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i - 1] == a[i]) || len(a[i]) == 0 {
			continue;
		}
		data := make(map[string]string)
		data["value"] = a[i]
		ret = append(ret, data)
	}
	return
}
