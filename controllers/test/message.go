package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/models/paginator"
	//. "dev.model.360baige.com/models/message"
	_"fmt"
)

type MessageController struct {
	beego.Controller
}

// @Title 最新推广消息接口
// @Description 最新推广消息接口
// @Success 200 {"code":200,"messgae":"ok","data":{"list":{ ... ... },"accessToken":"ok"}}
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Param personId query   string true       "身份id"
// @Param userId query   string true       "userid"
// @Param datatype query   string true       "类型((1,首页推广，2 首页消息列表 ，3消息首页消息列表，4系统消息)"
// @Failure 400 {"code":400,"message":"..."}
// @router /get-message-list [post]
func (c *MessageController) GetMessageList () {
	var reply Paginator
	res := Response{}
	pageSize, _ := c.GetInt("pageSize")
	current, _ := c.GetInt("current")
	markID, _ := c.GetInt64("markid")
	direction, _ := c.GetInt("direction")
	filters := c.GetString("filters")
	args := &Paginator{
		PageSize:  pageSize,
		Current:   current,
		MarkID:    markID,
		Direction: direction,
		Filters:   filters,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Message", "List", args, &reply)
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