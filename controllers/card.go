package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"time"
	. "dev.model.360baige.com/models/card"
	. "dev.model.360baige.com/models/response"
)

type CardController struct {
	beego.Controller
}

// @Title 新增
// @Description 新增
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /add [get]
func (c *CardController) Add() {
	timestamp := time.Now().UnixNano() / 1e6
	var (
		res   Response // http 返回体
		reply Card
	)
	args := &Card{
		CreateTime: timestamp,
		UpdateTime: timestamp,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Card", "Add", args, &reply)
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

// @Title 信息
// @Description 信息
// @Success 200 {"code":200,"messgae":"信息查询成功", "data":{ ... ... }}
// @Param   id     query   string true       "ID"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /detail [get]
func (c *CardController) Detail() {
	id, _ := c.GetInt64("id")
	res := Response{}
	var reply Card
	args := &Card{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Card", "FindById", args, &reply)

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

// @Title 信息修改
// @Description 信息修改
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   id     query   string true       "ID"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /modify [post]
func (c *CardController) Modify() {
	id, _ := c.GetInt64("id")

	var reply Card
	res := Response{}
	args := &Card{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Card", "FindById", args, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = err.Error()
		c.Data["json"] = res
		c.ServeJSON()
	}
	timestamp := time.Now().UnixNano() / 1e6
	reply.Id = id

	reply.UpdateTime = timestamp

	err = client.Call(beego.AppConfig.String("EtcdURL"), "Card", "UpdateById", reply, nil)

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