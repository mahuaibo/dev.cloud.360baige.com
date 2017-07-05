package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/models/response"
	//. "dev.model.360baige.com/models/feedback"
	_"fmt"
)

type FeedbackController struct {
	beego.Controller
}

// @Title 意见反馈接口
// @Description 意见反馈接口
// @Success 200 {"code":200,"messgae":"ok","data":{"accessToken":"ok"}}
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Param personId query   string true       "身份id"
// @Param userId query   string true       "userid"
// @Param comment query   string true       "反馈内容"
// @Failure 400 {"code":400,"message":"..."}
// @router /add [post]
func (c *FeedbackController) Add() {
	//accessToken := c.GetString("accessToken")
	var (
		res   Response // http 返回体
		reply Response
		args  Response // 需要更改
	)

	err := client.Call(beego.AppConfig.String("EtcdURL"), "Feedback", "Add", args, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "新增失败"
		c.Data["json"] = res
		c.ServeJSON()
	}
	res.Code = ResponseNormal
	res.Messgae = "新增成功"
	res.Data = reply
	c.Data["json"] = res
	c.ServeJSON()
}
