package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/models/paginator"
	//. "dev.model.360baige.com/models/resource"
	_"fmt"
)

type ResourceController struct {
	beego.Controller
}


// @Title 启动页接口
// @Description 启动页接口
// @Success 200 {"code":200,"messgae":"ok","data":{"list":{ ... ... },"accessToken":"ok"}}
// @Param accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /get-startpage-img [post]
func (c *ResourceController) StartPageList() {
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
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Resource", "List", args, &reply)
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
// @Title 首页图片接口
// @Description 首页图片接口
// @Success 200 {"code":200,"messgae":"ok","data":{"list":{ ... ... },"accessToken":"ok"}}
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Failure 400 {"code":400,"message":"..."}
// @router /get-homepage-img [post]
func (c *ResourceController) HomepageImgList() {
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
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Resource", "List", args, &reply)
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