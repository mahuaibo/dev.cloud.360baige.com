package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/models/app"
	. "dev.model.360baige.com/models/response"
)

type AppIntroductionController struct {
	beego.Controller
}
// @Title 信息
// @Description 信息
// @Success 200 {"code":200,"messgae":"信息查询成功", "data":{ ... ... }}
// @Param   id     query   string true       "ID"
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Param personId query   string true       "身份id"
// @Failure 400 {"code":400,"message":"..."}
// @router /detail [post]
func (c *AppIntroductionController) Detail() {
	id, _ := c.GetInt64("id")
	res := Response{}
	var reply AppIntroduction
	args := &AppIntroduction{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "AppIntroduction", "FindById", args, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "信息查询失败"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "信息查询成功"
		res.Data = reply
		c.Data["json"] = res
		c.ServeJSON()
	}
}
