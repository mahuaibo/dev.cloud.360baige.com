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
	timestamp := time.Now().UnixNano() / 1e6
	args2 := []schoolfee.Record{
		schoolfee.Record{
			CreateTime: timestamp,
			UpdateTime: timestamp,
			CompanyId:  replyAccessToken.CompanyId,
			ProjectId:  project_id,
			Name:       "张1三",
			ClassName:  "三1",
			IdCard:     "2004",
			Num:        "2004",
			Phone:      "18911545460",
			Status:     0,
			Price:      1.01,
			IsFee:      0,
			FeeTime:    0,
			Desc:       "",
		},
		schoolfee.Record{
			CreateTime: timestamp,
			UpdateTime: timestamp,
			CompanyId:  replyAccessToken.CompanyId,
			ProjectId:  project_id,
			Name:       "张2三",
			ClassName:  "三1",
			IdCard:     "2004",
			Num:        "2004",
			Phone:      "18911545460",
			Status:     0,
			Price:      1.02,
			IsFee:      0,
			FeeTime:    0,
			Desc:       "",

		},
		schoolfee.Record{
			CreateTime: timestamp,
			UpdateTime: timestamp,
			CompanyId:  replyAccessToken.CompanyId,
			ProjectId:  project_id,
			Name:       "张3三",
			ClassName:  "三1",
			IdCard:     "2004",
			Num:        "2004",
			Phone:      "18911545460",
			Status:     0,
			Price:      1.01,
			IsFee:      0,
			FeeTime:    0,
			Desc:       "",
		},
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
	res.Messgae = "下载缴费记录成功"
	res.Data.List = listOfRecord
	c.Data["json"] = res
	c.ServeJSON()
	return

}
