package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/models/app"
	_"fmt"
	"time"
)

type QuietHoursController struct {
	beego.Controller
}

// @Title 勿扰模式接口
// @Description 勿扰模式接口
// @Success 200 {"code":200,"messgae":"ok","data":{"accessToken":"ok"}}
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Param personId query   string true       "身份id"
// @Param userId query   string true       "userid"
// @Param starttime  query   string true       "开始时间"
// @Param endtime  query   string true       "结束时间"
// @Failure 400 {"code":400,"message":"..."}
// @router /set [post]
func (c *QuietHoursController) Set() {
	id, _ := c.GetInt64("id")

	var reply QuietHours
	res := Response{}
	args := &QuietHours{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "QuietHours", "FindById", args, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = err.Error()
		c.Data["json"] = res
		c.ServeJSON()
	}
	timestamp := time.Now().UnixNano() / 1e6
	reply.Id = id

	reply.UpdateTime = timestamp

	err = client.Call(beego.AppConfig.String("EtcdURL"), "QuietHours", "UpdateById", reply, nil)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "信息修改失败！"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "信息修改成功！"
		c.Data["json"] = res
		c.ServeJSON()
	}
}

// @Title 勿扰模式接口
// @Description 勿扰模式接口
// @Success 200 {"code":200,"messgae":"ok","data":{"accessToken":"ok"}}
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Param personId query   string true       "身份id"
// @Param userId query   string true       "userid"
// @Param starttime  query   string true       "开始时间"
// @Param endtime  query   string true       "结束时间"
// @Failure 400 {"code":400,"message":"..."}
// @router /detail [post]
func (c *QuietHoursController) Detail() {
	id, _ := c.GetInt64("id")

	var reply QuietHours
	res := Response{}
	args := &QuietHours{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "QuietHours", "FindById", args, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = err.Error()
		c.Data["json"] = res
		c.ServeJSON()
	}
	timestamp := time.Now().UnixNano() / 1e6
	reply.Id = id

	reply.UpdateTime = timestamp

	err = client.Call(beego.AppConfig.String("EtcdURL"), "QuietHours", "UpdateById", reply, nil)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "信息修改失败！"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "信息修改成功！"
		c.Data["json"] = res
		c.ServeJSON()
	}
}
