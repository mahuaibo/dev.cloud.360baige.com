package schoolfeewin

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/schoolfeewin"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/schoolfee"
	"dev.model.360baige.com/action"
	"fmt"
)

// Record API
type RecordController struct {
	beego.Controller
}

// @Title 校园收费列表接口
// @Description Project List 校园收费列表接口
// @Success 200 {"code":200,"messgae":"获取缴费项目成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   project_id     query   int true       "项目ID"
// @Failure 400 {"code":400,"message":"获取缴费项目失败"}
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
	var args2 action.FindByCond
	args2.CondList = append(args2.CondList,
		action.CondValue{Type: "And", Key: "company_id", Val: replyAccessToken.CompanyId},
		action.CondValue{Type: "And", Key: "project_id", Val: project_id},
	)

	var replyRecord []schoolfee.Record
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "ListByCond", args2, &replyRecord)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "获取缴费项目失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	fmt.Println("2:", replyRecord)

	// 3.
	var listOfRecord []Record = make([]Record, len(replyRecord), len(replyRecord))
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
	res.Messgae = "获取缴费项目成功"
	res.Data.List = listOfRecord
	c.Data["json"] = res
	c.ServeJSON()
	return
}
