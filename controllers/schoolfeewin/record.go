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
	"os"
	"io"
	"github.com/xuri/excelize"
	"strconv"
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
	pageSize, _ := c.GetInt64("page_size")
	currentPage, _ := c.GetInt64("current")
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
	// 2.
	var args2 action.PageByCond
	args2.CondList = append(args2.CondList,
		action.CondValue{Type: "And", Key: "company_id", Val: replyAccessToken.CompanyId},
		action.CondValue{Type: "And", Key: "project_id", Val: project_id},
		action.CondValue{Type: "And", Key: "status__gt", Val: -1},
	)
	args2.OrderBy = []string{"id"}
	args2.Cols = []string{"id", "create_time", "company_id", "project_id", "name", "class_name", "id_card", "num",
		"phone", "price", "is_fee", "fee_time", "desc", "status" }
	args2.PageSize = pageSize
	args2.Current = currentPage
	var replyRecord []schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "ListByCond", args2, &replyRecord)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "获取缴费项目记录失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	// 3.
	listOfRecord := make([]Record, len(replyRecord), len(replyRecord))
	fmt.Println(replyRecord)
	for index, rec := range replyRecord {
		listOfRecord[index] = Record{
			Id:         rec.Id,
			CreateTime: time.Unix(rec.CreateTime / 1000, 0).Format("2006-01-02"),
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
// @router /upload [options,post]
func (c *RecordController) UploadRecord() {
	res := UploadRecordResponse{}
	access_token := c.GetString("a")
	projectId, _ := c.GetInt64("i")
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

	requestType := c.Ctx.Request.Method
	var replyRecord action.Num
	if requestType == "POST" {
		formFile, header, err := c.Ctx.Request.FormFile("uploadFile")
		if err != nil {
			fmt.Println("Get form file failed: %s\n", err)
			return
		}
		defer formFile.Close()
		objectKey := "./" + strconv.FormatInt(time.Now().UnixNano() / 1e6, 10) + header.Filename
		// 创建保存文件
		destFile, err := os.Create(objectKey)
		if err != nil {
			fmt.Println("Create failed: %s\n", err)
			return
		}
		defer destFile.Close()
		// 读取表单文件，写入保存文件
		_, err = io.Copy(destFile, formFile)
		if err != nil {
			fmt.Println("Write file failed: %s\n", err)
			return
		}

		xlsx, err := excelize.OpenFile(objectKey)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Get sheet index.
		index := xlsx.GetSheetIndex("Sheet1")
		rows := xlsx.GetRows("sheet" + strconv.Itoa(index))
		fmt.Println("rows:", rows)
		timestamp := time.Now().UnixNano() / 1e6
		var args2 []schoolfee.Record = make([]schoolfee.Record, len(rows) - 1)
		for key, row := range rows {
			if key > 0 {
				Price, err := strconv.ParseFloat(row[5], 64)
				fmt.Println("Price:", Price)
				fmt.Println("err:", err)
				if err != nil {
					res.Code = ResponseLogicErr
					res.Messgae = "上传缴费名单失败"
					c.Data["json"] = res
					c.ServeJSON()
					return
				}
				args2[key - 1] = schoolfee.Record{
					CreateTime: timestamp,
					UpdateTime: timestamp,
					CompanyId:  replyAccessToken.CompanyId,
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
		fmt.Println("args2", args2)
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "AddMultiple", args2, &replyRecord)
		fmt.Println("replyRecord", replyRecord)
		if err != nil {
			res.Code = ResponseLogicErr
			res.Messgae = "上传缴费名单失败"
			c.Data["json"] = res
			c.ServeJSON()
			return
		}

		err = os.Remove(objectKey) // 删除文件
		if err != nil {
			fmt.Println("file remove Error!", err)
		} else {
			fmt.Print("file remove OK!")
		}
	}

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
// @router /download [post]
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
	listOfRecord := make([]Record, len(replyRecord))
	for index, rec := range replyRecord {
		listOfRecord[index] = Record{
			Id:         rec.Id,
			CreateTime: time.Unix(rec.CreateTime / 1000, 0).Format("2006-01-02"),
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
	res.Messgae = "下载缴费记录成功"
	res.Data.List = listOfRecord
	c.Data["json"] = res
	c.ServeJSON()
	return
}
