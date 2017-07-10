package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"time"
	. "dev.model.360baige.com/models/application"
	. "dev.model.360baige.com/models/paginator"
	. "dev.model.360baige.com/models/response"
)

type ApplicationCommentController struct {
	beego.Controller
}
// @Title 获取应用菜单接口
// @Description 获取应用菜单接口
// @Success 200 {"code":200,"messgae":"ok","data":{"list":{ ... ... },"accessToken":"ok"}}
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Param personId query   string true       "身份id"
// @Param datatype query   string true       "类型(1,首页，2 基础 ，3已购 )"
// @Failure 400 {"code":400,"message":"..."}
// @router /get-commentlist [post]
func (c *ApplicationCommentController) GetCommentList() {
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
	err := client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationComment", "List", args, &reply)
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
// @Title 应用评价接口
// @Description 应用评价接口
// @Success 200 {"code":200,"messgae":"ok","data":{"accessToken":"ok"}}
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Param personId query   string true       "身份id"
// @Param userId query   string true       "userid"
// @Param appId query   string true       "应用id"
// @Param comment query   string true       "评论内容"
// @Failure 400 {"code":400,"message":"..."}
// @router /add [post]
func (c *ApplicationCommentController) Add() {
	timestamp := time.Now().UnixNano() / 1e6
	var (
		res   Response // http 返回体
		reply ApplicationComment
	)
	args := &ApplicationComment{
		CreateTime: timestamp,
		UpdateTime: timestamp,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationComment", "Add", args, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "新增失败"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
        res.Code = ResponseNormal
        res.Messgae = "新增成功"
        res.Data = reply
        c.Data["json"] = res
        c.ServeJSON()
	}
}
