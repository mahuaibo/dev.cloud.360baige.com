package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/models/company"
	. "dev.model.360baige.com/models/response"
	"time"
)

// COMPANY API
type CompanyController struct {
	beego.Controller
}

// @Title 企业信息接口
// @Description 企业信息接口
// @Success 200 {"code":200,"messgae":"获取企业信息成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   company_id     query   string true       "企业ID"
// @Failure 400 {"code":400,"message":"获取企业信息失败"}
// @router /detail [get]
func (c *CompanyController) Detail() {
	res := Response{}
	company_id, _ := c.GetInt64("company_id", 0)
	if company_id == 0 {
		res.Code = ResponseSystemErr
		res.Messgae = "获取企业信息失败"
		c.Data["json"] = res
		c.ServeJSON()
	}
	var reply Company
	var err error

	err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "FindById", &Company{
		Id: company_id,
	}, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取企业信息失败"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "获取企业信息成功"
		replyData := Company{
			Id:         reply.Id,
			CreateTime: reply.CreateTime,
			UpdateTime: reply.UpdateTime,
			Type:       reply.Type,
			Level:      reply.Level,
			Logo:       reply.Logo,
			Name:       reply.Name,
			ShortName:  reply.ShortName,
			SubDomain:  reply.SubDomain,
			CityId:     reply.CityId,
			Address:    reply.Address,
			PositionX:  reply.PositionX,
			PositionY:  reply.PositionY,
			Remark:     reply.Remark,
			Brief:      reply.Brief,
			Status:     reply.Status,
		}
		res.Data = replyData
		c.Data["json"] = res
		c.ServeJSON()
	}
}

// @Title 企业信息修改接口
// @Description 企业信息修改接口
// @Success 200 {"code":200,"messgae":"获取企业信息修改成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   company_id     query   string true       "企业ID"
// @Failure 400 {"code":400,"message":"获取企业信息修改失败"}
// @router /modify [get]
func (c *CompanyController) Modify() {
	res := Response{}
	company_id, _ := c.GetInt64("company_id", 0)
	name := c.GetString("name")
	short_name := c.GetString("short_name")
	logo := c.GetString("logo")
	if company_id == 0 {
		res.Code = ResponseSystemErr
		res.Messgae = "获取企业信息失败"
		c.Data["json"] = res
		c.ServeJSON()
	}
	var reply Company
	var err error
	timestamp := time.Now().UnixNano()/10e6

	err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "UpdateById", &Company{
		Id: company_id,
		UpdateTime:timestamp,
		Name:name,
		ShortName:short_name,
		Logo:logo,
	}, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取企业信息失败"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "获取企业信息成功"
		replyData := Company{
			Id:         reply.Id,
			CreateTime: reply.CreateTime,
			UpdateTime: reply.UpdateTime,
			Type:       reply.Type,
			Level:      reply.Level,
			Logo:       reply.Logo,
			Name:       reply.Name,
			ShortName:  reply.ShortName,
			SubDomain:  reply.SubDomain,
			CityId:     reply.CityId,
			Address:    reply.Address,
			PositionX:  reply.PositionX,
			PositionY:  reply.PositionY,
			Remark:     reply.Remark,
			Brief:      reply.Brief,
			Status:     reply.Status,
		}
		res.Data = replyData
		c.Data["json"] = res
		c.ServeJSON()
	}
}