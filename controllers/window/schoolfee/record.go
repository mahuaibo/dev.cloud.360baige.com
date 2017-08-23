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
	"time"
	"os"
	"io"
	"github.com/xuri/excelize"
	"strconv"
	"github.com/tealeg/xlsx"
	"net/http"
	"sort"
	"fmt"
)

// Record API
type RecordController struct {
	beego.Controller
}

// @Title 校园收费记录列表接口
// @Description Project List 校园收费记录列表接口
// @Success 200 {"code":200,"message":"获取缴费项目记录成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   projectId     query   int true       "项目ID"
// @Failure 400 {"code":400,"message":"获取缴费项目记录失败"}
// @router /list [get]
func (c *RecordController) ListOfRecord() {
	type data ListOfRecordResponse
	accessToken := c.GetString("accessToken")
	projectId := c.GetString("projectId")
	pageSize, _ := c.GetInt64("pageSize")
	currentPage, _ := c.GetInt64("current")
	if accessToken == "" {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "accessToken", Val: accessToken})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var args2 action.PageByCond
	args2.CondList = append(args2.CondList,
		action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId},
		action.CondValue{Type: "And", Key: "projectId", Val: projectId},
		action.CondValue{Type: "And", Key: "status__gt", Val: -1},
	)
	args2.OrderBy = []string{"id"}
	args2.Cols = []string{"id", "create_time", "company_id", "projectId", "name", "class_name", "id_card", "num",
						  "phone", "price", "is_fee", "fee_time", "desc", "status" }
	args2.PageSize = pageSize
	args2.Current = currentPage
	var replyRecord []schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "ListByCond", args2, &replyRecord)
	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "获取缴费项目记录失败"}
		c.ServeJSON()
		return
	}

	listOfRecord := make([]Record, len(replyRecord), len(replyRecord))
	for index, rec := range replyRecord {
		listOfRecord[index] = Record{
			Id:         rec.Id,
			CreateTime: time.Unix(rec.CreateTime/1000, 0).Format("2006-01-02"),
			CompanyId:  rec.CompanyId,
			ProjectId:  rec.ProjectId,
			Name:       rec.Name,
			ClassName:  rec.ClassName,
			IdCard:     rec.IdCard,
			Num:        rec.Num,
			Phone:      rec.Phone,
			Price:      rec.Price,
			IsFee:      rec.IsFee,
			FeeTime:    rec.FeeTime,
			Desc:       rec.Desc,
			Status:     rec.Status,
		}
	}

	c.Data["json"] = data{Code: ResponseNormal, Message: "获取缴费项目记录成功", Data: ListOfRecord{
		List: listOfRecord,
	}}
	c.ServeJSON()
	return
}

// @Title 添加收费名单接口
// @Description Record Add 添加收费名单接口
// @Success 200 {"code":200,"message":"添加收费名单成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   name     query   string true       "项目名称"
// @Param   isLimit     query   string true       "是否限制缴费"
// @Param   desc     query   string true       "描述"
// @Param   link     query   string true       "描述链接"
// @Param   status     query   string true       "状态 -1注销 0正常"
// @Failure 400 {"code":400,"message":"添加收费名单失败"}
// @router /add [post]
func (c *RecordController) AddRecord() {
	res := AddRecordResponse{}
	accessToken := c.GetString("accessToken")
	projectId, _ := c.GetInt64("projectId", 0)
	name := c.GetString("name")
	class_name := c.GetString("class_name")
	id_card := c.GetString("id_card")
	num := c.GetString("num")
	phone := c.GetString("phone")
	price, _ := c.GetFloat("price")
	is_fee, _ := c.GetInt8("is_fee")
	fee_time, _ := c.GetInt64("fee_time")
	desc := c.GetString("desc")
	if accessToken == "" {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type: "And",
		Key:  "accessToken",
		Val:  accessToken,
	})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyUserPosition)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌失效"
		c.Data["json"] = res
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
		ClassName:  class_name,
		IdCard:     id_card,
		Num:        num,
		Phone:      phone,
		Desc:       desc,
		Price:      price,
		IsFee:      is_fee,
		FeeTime:    fee_time,
		Status:     0,
	}
	var replyRecord schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "Add", args2, &replyRecord)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "添加收费名单失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	res.Code = ResponseNormal
	res.Message = "添加收费名单成功"
	res.Data.Id = replyRecord.Id
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// @Title 查看缴费记录接口
// @Description Project Add 查看缴费记录接口
// @Success 200 {"code":200,"message":"查看缴费记录成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   record_id     query   string true       "记录ID"
// @Failure 400 {"code":400,"message":"查看缴费记录失败"}
// @router /detail [get]
func (c *RecordController) DetailRecord() {
	res := DetailRecordResponse{}
	accessToken := c.GetString("accessToken")
	record_id, _ := c.GetInt64("record_id", 0)
	if accessToken == "" {
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
		Key:  "accessToken",
		Val:  accessToken,
	})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyUserPosition)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	// 2.
	args2 := &schoolfee.Record{
		Id: record_id,
	}
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "FindById", args2, args2)

	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "修改缴费项目失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	res.Code = ResponseNormal
	res.Message = "修改缴费项目成功"
	res.Data.Data = Record{
		Id:         args2.Id,
		CreateTime: time.Unix(args2.CreateTime/1000, 0).Format("2006-01-02"),
		UpdateTime: args2.UpdateTime,
		CompanyId:  args2.CompanyId,
		ProjectId:  args2.ProjectId,
		Name:       args2.Name,
		ClassName:  args2.ClassName,
		IdCard:     args2.IdCard,
		Num:        args2.Num,
		Phone:      args2.Phone,
		Status:     args2.Status,
		Price:      args2.Price,
		IsFee:      args2.IsFee,
		FeeTime:    args2.FeeTime,
		Desc:       args2.Desc,
	}
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// @Title 修改缴费记录接口
// @Description Project Add 修改缴费记录接口
// @Success 200 {"code":200,"message":"修改缴费记录成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   record_id     query   string true       "记录ID"
// @Failure 400 {"code":400,"message":"修改缴费记录失败"}
// @router /modify [post]
func (c *RecordController) ModifyRecord() {
	res := ModifyRecordResponse{}
	accessToken := c.GetString("accessToken")
	record_id, _ := c.GetInt64("record_id", 0)
	name := c.GetString("name")
	class_name := c.GetString("class_name")
	id_card := c.GetString("id_card")
	num := c.GetString("num")
	phone := c.GetString("phone")
	price, _ := c.GetFloat("price")
	is_fee, _ := c.GetInt8("is_fee")
	fee_time, _ := c.GetInt64("fee_time")
	desc := c.GetString("desc")
	if accessToken == "" {
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
		Key:  "accessToken",
		Val:  accessToken,
	})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyUserPosition)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	// 2.
	args2 := &schoolfee.Record{
		Id: record_id,
	}
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "FindById", args2, args2)

	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "修改缴费项目失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	if args2.CompanyId != replyUserPosition.CompanyId {
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
	operationTime := time.Now().UnixNano() / 10e6
	args3.UpdateList = append(args3.UpdateList,
		action.UpdateValue{Key: "update_time", Val: operationTime},
		action.UpdateValue{Key: "name", Val: name},
		action.UpdateValue{Key: "class_name", Val: class_name},
		action.UpdateValue{Key: "id_card", Val: id_card},
		action.UpdateValue{Key: "num", Val: num},
		action.UpdateValue{Key: "phone", Val: phone},
		action.UpdateValue{Key: "price", Val: price},
		action.UpdateValue{Key: "is_fee", Val: is_fee},
		action.UpdateValue{Key: "fee_time", Val: fee_time},
		action.UpdateValue{Key: "desc", Val: desc},
	)

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "UpdateById", args3, &replyNum)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "修改缴费记录失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	res.Code = ResponseNormal
	res.Message = "修改缴费记录成功"
	res.Data.Count = replyNum.Value
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// @Title 删除校园收费记录接口
// @Description Delete Record 删除校园收费记录接口
// @Success 200 {"code":200,"message":"删除缴费项目记录成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   recordIds     query   int true       "项目记录IDs"
// @Failure 400 {"code":400,"message":"删除缴费项目记录失败"}
// @router /delete [post]
func (c *RecordController) DeleteRecord() {
	res := DeleteRecordResponse{}
	accessToken := c.GetString("accessToken")
	recordIds := c.GetString("recordIds")
	record_id_s := utils.StrArrToInt64Arr(strings.Split(recordIds, ","))

	if accessToken == "" {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "accessToken", Val: accessToken, })
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyUserPosition)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 2.
	args2 := action.DeleteByIdCond{Value: record_id_s}
	var replyRecord action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "DeleteById", args2, &replyRecord)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "删除缴费项目记录失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	res.Code = ResponseNormal
	res.Message = "删除缴费项目记录成功"
	res.Data.Count = replyRecord.Value
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// @Title 上传缴费名单接口
// @Description Delete Record 上传缴费名单接口
// @Success 200 {"code":200,"message":"上传缴费名单成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   projectId     query   string true       "项目Id"
// @Failure 400 {"code":400,"message":"上传缴费名单失败"}
// @router /upload [options,post]
func (c *RecordController) UploadRecord() {
	type data UploadRecordResponse
	accessToken := c.GetString("accessToken")
	projectId, _ := c.GetInt64("projectId")
	if accessToken == "" {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "accessToken", Val: accessToken })
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	requestType := c.Ctx.Request.Method
	var replyRecord action.Num
	if requestType == "POST" {
		formFile, header, err := c.Ctx.Request.FormFile("uploadFile")
		if err != nil {
			return
		}
		defer formFile.Close()
		objectKey := "./" + strconv.FormatInt(time.Now().UnixNano()/1e6, 10) + header.Filename
		// 创建保存文件
		destFile, err := os.Create(objectKey)
		if err != nil {
			return
		}
		defer destFile.Close()
		// 读取表单文件，写入保存文件
		_, err = io.Copy(destFile, formFile)
		if err != nil {
			return
		}

		xlsx, err := excelize.OpenFile(objectKey)
		if err != nil {
			os.Exit(1)
		}
		// Get sheet index.
		index := xlsx.GetSheetIndex("Sheet1")
		rows := xlsx.GetRows("sheet" + strconv.Itoa(index))
		timestamp := time.Now().UnixNano() / 1e6
		var args2 []schoolfee.Record = make([]schoolfee.Record, len(rows)-1)
		for key, row := range rows {
			if key > 0 {
				Price, err := strconv.ParseFloat(row[5], 64)
				if err != nil {
					c.Data["json"] = data{Code: ResponseLogicErr, Message: "上传缴费名单失败"}
					c.ServeJSON()
					return
				}
				args2[key-1] = schoolfee.Record{
					CreateTime: timestamp,
					UpdateTime: timestamp,
					CompanyId:  replyUserPosition.CompanyId,
					ProjectId:  projectId,
					Name:       row[0],
					ClassName:  row[1],
					IdCard:     row[2],
					Num:        row[3],
					Phone:      row[4],
					Status:     0,
					Price:      Price,
					IsFee:      0,
					FeeTime:    0,
					Desc:       row[6],
				}
			}
		}

		err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "AddMultiple", args2, &replyRecord)
		if err != nil {
			c.Data["json"] = data{Code: ResponseLogicErr, Message: "上传缴费名单失败"}
			c.ServeJSON()
			return
		}

		err = os.Remove(objectKey) // 删除文件
		if err != nil {
		} else {
		}
	}

	c.Data["json"] = data{Code: ResponseNormal, Message: "上传缴费名单成功", Data: UploadRecord{
		Count: replyRecord.Value,
	}}
	c.ServeJSON()
	return
}

// @Title 下载缴费记录接口
// @Description Delete Record 下载缴费记录接口
// @Success 200 {"code":200,"message":"下载缴费记录成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"下载缴费记录失败"}
// @router /download [get,post]
func (c *RecordController) DownloadRecord() {
	type data DownloadRecordResponse
	accessToken := c.GetString("accessToken")
	classNames := c.GetString("classNames")
	isFees := c.GetString("isFees")
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
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	var replyRecord []schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "ListByCond", action.ListByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId},
			action.CondValue{Type: "And", Key: "class_name__in", Val: strings.Split(classNames, ",")},
			action.CondValue{Type: "And", Key: "is_fee__in", Val: strings.Split(isFees, ",")},
		},
	}, &replyRecord)
	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "下载缴费记录失败"}
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
		row := sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度
		row.AddCell().Value = rec.Name
		row.AddCell().Value = rec.ClassName
		row.AddCell().Value = rec.IdCard
		row.AddCell().Value = rec.Num
		row.AddCell().Value = rec.Phone
		row.AddCell().Value = strconv.FormatFloat(rec.Price, 'f', -1, 64)
		var isFee int64 = int64(rec.IsFee)
		row.AddCell().Value = strconv.FormatInt(isFee, 10)
		row.AddCell().Value = strconv.FormatInt(rec.FeeTime, 10)
		row.AddCell().Value = rec.Desc
	}
	objectKey := strconv.FormatInt(time.Now().UnixNano()/1e6, 10) + "file.xlsx"
	err = file.Save(objectKey)
	if err != nil {
		panic(err)
	}

	c.Ctx.Output.Header("Accept-Ranges", "bytes")
	c.Ctx.Output.Header("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s", objectKey)) //文件名
	c.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Ctx.Output.Header("Pragma", "no-cache")
	c.Ctx.Output.Header("Expires", "0")
	//最主要的一句
	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, objectKey)
	c.Data["json"] = data{Code: ResponseNormal, Message: "下载缴费记录成功"}
	c.ServeJSON()
	return
}

// @Title 班级列表接口
// @Description Project List 班级列表接口
// @Success 200 {"code":200,"message":"获取班级列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   projectId     query   int true       "项目ID"
// @Failure 400 {"code":400,"message":"获取班级列表失败"}
// @router /classList [get]
func (c *RecordController) ClassList() {
	type data ClassListOfRecordResponse
	accessToken := c.GetString("accessToken")
	projectId := c.GetString("projectId")
	if accessToken == "" {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "accessToken", Val: accessToken})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	var replyRecord []schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "ListByCond", action.PageByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId},
			action.CondValue{Type: "And", Key: "project_id", Val: projectId},
		},
		Cols: []string{"class_name" },
	}, &replyRecord)

	if err != nil {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "获取班级列表失败"}
		c.ServeJSON()
		return
	}

	var listOfRecord []string
	for _, rec := range replyRecord {
		listOfRecord = append(listOfRecord, rec.ClassName)
	}

	sort.Strings(listOfRecord)
	c.Data["json"] = data{Code: ResponseNormal, Message: "获取班级列表成功", Data: RemoveDuplicatesAndEmpty(listOfRecord)}
	c.ServeJSON()
	return
}

func RemoveDuplicatesAndEmpty(a []string) (ret []map[string]string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue;
		}
		data := make(map[string]string)
		data["value"] = a[i]
		ret = append(ret, data)
	}
	return
}
