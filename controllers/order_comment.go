package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"time"
	. "dev.model.360baige.com/models/order"
	. "dev.model.360baige.com/models/response"
)

type OrderCommentController struct {
	beego.Controller
}
// @Title 订单评价接口
// @Description 订单评价接口
// @Success 200 {"code":200,"messgae":"ok","data":{accessToken":"ok"}}
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Param personId query   string true       "身份id"
// @Param userId query   string true       "userid"
// @Param childId query   string true       "childid"
// @Param orderId query   string true       "orderid"
// @Param   params     query   string true       "各种参数"
// @Failure 400 {"code":400,"message":"..."}
// @router /add [post]
func (c *OrderCommentController) Add() {
	timestamp := time.Now().UnixNano() / 1e6
	var (
		res   Response // http 返回体
		reply OrderComment
	)
	args := &OrderComment{
		CreateTime: timestamp,
		UpdateTime: timestamp,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "OrderComment", "Add", args, &reply)
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
