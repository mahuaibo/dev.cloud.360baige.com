package schoolfeewin

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/schoolfeewin"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/schoolfee"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/utils"
	"fmt"
	"strings"
	"time"
	"github.com/tealeg/xlsx"
	"strconv"
	"net/http"
)

// Record API
type RecordController struct {
	beego.Controller
}

// @Title 校园收费记录列表接口
// @Description Project List 校园收费记录列表接口
// @Success 200 {"code":200,"messgae":"获取缴费项目记录成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   project_id     query   int true       "项目ID"
// @Failure 400 {"code":400,"message":"获取缴费项目记录失败"}
// @router /list [get]
func (c *RecordController) ListOfRecord() {
	res := ListOfRecordResponse{}
	access_token := c.GetString("access_token")
	project_id := c.GetString("project_id")
	if access_token == "" {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "access_token", Val: access_token})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	fmt.Println("1:", replyAccessToken)
	// 2.
	var args2 action.FindByCond
	args2.CondList = append(args2.CondList,
		action.CondValue{Type: "And", Key: "company_id", Val: replyAccessToken.CompanyId},
		action.CondValue{Type: "And", Key: "project_id", Val: project_id},
	)

	var replyRecord []schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "ListByCond", args2, &replyRecord)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "获取缴费项目记录失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	fmt.Println("2:", replyRecord)

	// 3.
	listOfRecord := make([]Record, len(replyRecord), len(replyRecord))
	for index, rec := range replyRecord {
		listOfRecord[index] = Record{
			Id:         rec.Id,
			CreateTime: rec.CreateTime,
			UpdateTime: rec.UpdateTime,
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

	res.Code = ResponseNormal
	res.Messgae = "获取缴费项目记录成功"
	res.Data.List = listOfRecord
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// @Title 删除校园收费记录接口
// @Description Delete Record 删除校园收费记录接口
// @Success 200 {"code":200,"messgae":"删除缴费项目记录成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   record_ids     query   int true       "项目记录IDs"
// @Failure 400 {"code":400,"message":"删除缴费项目记录失败"}
// @router /delete [post]
func (c *RecordController) DeleteRecord() {
	res := DeleteRecordResponse{}
	access_token := c.GetString("access_token")
	record_ids := c.GetString("record_ids")
	record_id_s := utils.StrArrToInt64Arr(strings.Split(record_ids, ","))

	if access_token == "" {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "access_token", Val: access_token, })
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	fmt.Println("1:", replyAccessToken)
	// 2.
	args2 := action.DeleteByIdCond{Value: record_id_s}
	var replyRecord action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "DeleteById", args2, &replyRecord)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "删除缴费项目记录失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	fmt.Println("2:", replyRecord)
	res.Code = ResponseNormal
	res.Messgae = "删除缴费项目记录成功"
	res.Data.Count = replyRecord.Value
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// @Title 上传缴费名单接口
// @Description Delete Record 上传缴费名单接口
// @Success 200 {"code":200,"messgae":"上传缴费名单成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"上传缴费名单失败"}
// @router /upload [get]
func (c *RecordController) UploadRecord() {
	excelFileName := "file.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		return
	}
	res := UploadRecordResponse{}
	access_token := c.GetString("access_token")
	project_id, _ := c.GetInt64("project_id", 1)

	if access_token == "" {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "access_token", Val: access_token })
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken user.UserPosition
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	fmt.Println("1:", replyAccessToken)

	// 2.
	timestamp := time.Now().UnixNano() / 1e6
	var args2 []schoolfee.Record
	for _, sheet := range xlFile.Sheets {
		for index, row := range sheet.Rows {
			if index != 0 {
				feeTime, _ := strconv.ParseInt(row.Cells[7].Value, 10, 64)
				price, _ := strconv.ParseFloat(row.Cells[5].Value, 64) //row.Cells[5].Value.(float64)
				isFee,_ :=  strconv.ParseInt(row.Cells[6].Value, 10, 64)
				args2 = append(args2, schoolfee.Record{
					CreateTime: timestamp,
					UpdateTime: timestamp,
					CompanyId:  replyAccessToken.CompanyId,
					ProjectId:  project_id,
					Name:       row.Cells[0].Value,
					ClassName:  row.Cells[1].Value,
					IdCard:     row.Cells[2].Value,
					Num:        row.Cells[3].Value,
					Phone:      row.Cells[4].Value,
					Price:      price,
					IsFee:      int8(isFee),
					FeeTime:    feeTime,
					Desc:       row.Cells[8].Value,
					Status:     0,
				})
			}
		}
	}

	var replyRecord action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "AddMultiple", args2, &replyRecord)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "上传缴费名单失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	fmt.Println("2:", replyRecord)
	// 3.
	res.Code = ResponseNormal
	res.Messgae = "上传缴费名单成功"
	res.Data.Count = replyRecord.Value
	c.Data["json"] = res
	c.ServeJSON()
	return

}

// @Title 下载缴费记录接口
// @Description Delete Record 下载缴费记录接口
// @Success 200 {"code":200,"messgae":"下载缴费记录成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   record_ids     query   int true       "项目记录IDs"
// @Failure 400 {"code":400,"message":"下载缴费记录失败"}
// @router /download [get,post]
func (c *RecordController) DownloadRecord() {
	res := DownloadRecordResponse{}
	access_token := c.GetString("access_token")
	class_name_s := c.GetString("class_names")
	is_fee_s := c.GetString("is_fees")

	if access_token == "" {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "access_token", Val: access_token })
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	fmt.Println("1:", replyAccessToken)
	// 2.
	args2 := action.ListByCond{CondList: []action.CondValue{
		action.CondValue{Type: "And", Key: "company_id", Val: replyAccessToken.CompanyId},
		action.CondValue{Type: "And", Key: "class_name__in", Val: strings.Split(class_name_s, ",")},
		action.CondValue{Type: "And", Key: "is_fee__in", Val: strings.Split(is_fee_s, ",")},
	}}
	var replyRecord []schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "ListByCond", args2, &replyRecord)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "下载缴费记录失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	fmt.Println("2:", replyRecord)
	// 3.
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	row := sheet.AddRow()
	row.SetHeightCM(1) //设置每行的高度
	cell := row.AddCell()
	cell.Value = "姓名"
	cell = row.AddCell()
	cell.Value = "班级"
	cell = row.AddCell()
	cell.Value = "身份证号码"
	cell = row.AddCell()
	cell.Value = "编号"
	cell = row.AddCell()
	cell.Value = "联系电话"
	cell = row.AddCell()
	cell.Value = "应缴费用"
	cell = row.AddCell()
	cell.Value = "是否缴费"
	cell = row.AddCell()
	cell.Value = "缴费时间"
	cell = row.AddCell()
	cell.Value = "备注"
	//姓名、班级、身份证号码、编号、联系电话、应缴费用、是否缴费、缴费时间、备注
	for _, rec := range replyRecord {
		row := sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度
		cell = row.AddCell()
		cell.Value = rec.Name
		cell = row.AddCell()
		cell.Value = rec.ClassName
		cell = row.AddCell()
		cell.Value = rec.IdCard
		cell = row.AddCell()
		cell.Value = rec.Num
		cell = row.AddCell()
		cell.Value = rec.Phone
		cell = row.AddCell()
		cell.Value = strconv.FormatFloat(rec.Price, 'f', -1, 64)
		cell = row.AddCell()
		var isfee int64 = int64(rec.IsFee)
		cell.Value = strconv.FormatInt(isfee, 10)
		cell = row.AddCell()
		cell.Value = strconv.FormatInt(rec.FeeTime, 10)
		cell = row.AddCell()
		cell.Value = rec.Desc
	}
	err = file.Save("file.xlsx")
	if err != nil {
		panic(err)
	}

	c.Ctx.Output.Header("Accept-Ranges", "bytes")
	c.Ctx.Output.Header("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s", "file.xls")) //文件名
	c.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Ctx.Output.Header("Pragma", "no-cache")
	c.Ctx.Output.Header("Expires", "0")
	//最主要的一句
	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, "file.xlsx")
	res.Code = ResponseNormal
	res.Messgae = "下载缴费记录成功"
	c.Data["json"] = res
	c.ServeJSON()
	return

}
