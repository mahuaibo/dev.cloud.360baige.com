package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"time"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/models/paginator"
)

type UserPositionController struct {
	beego.Controller
}
// @Title 切换身份接口
// @Description 切换身份接口
// @Success 200 {"code":200,"messgae":"ok","data":{"accessToken":"ok"}}
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Param personId query   string true       "身份id"
// @Param userId query   string true       "userid"
// @Param position  query   string true       "身份类型"
// @Failure 400 {"code":400,"message":"..."}
// @router /positionlist [post]
func (c *UserPositionController) PositionList() {
	var reply Paginator
	res := Response{}
	args := &Paginator{
		PageSize:  0,
		Current:   0,
		MarkID:    0,
		Direction: 0,
		Filters:   "",
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "List", args, &reply)
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
// @Title 切换身份接口
// @Description 切换身份接口
// @Success 200 {"code":200,"messgae":"ok","data":{"accessToken":"ok"}}
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Param personId query   string true       "身份id"
// @Param userId query   string true       "userid"
// @Param position  query   string true       "身份类型"
// @Failure 400 {"code":400,"message":"..."}
// @router /setposition [post]
func (c *UserPositionController) SetPosition() {
	id, _ := c.GetInt64("id")

	var reply UserPosition
	res := Response{}
	args := &UserPosition{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindById", args, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = err.Error()
		c.Data["json"] = res
		c.ServeJSON()
	}
	timestamp := time.Now().UnixNano() / 1e6
	reply.Id = id

	reply.UpdateTime = timestamp

	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "UpdateById", reply, nil)

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
// @Title 新增
// @Description 新增
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /add [get]
func (c *UserPositionController) Add() {
	timestamp := time.Now().UnixNano() / 1e6
	var (
		res   Response // http 返回体
		reply UserPosition
	)
	args := &UserPosition{
		CreateTime: timestamp,
		UpdateTime: timestamp,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "Add", args, &reply)
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
func (c *UserPositionController) Detail() {
	id, _ := c.GetInt64("id")
	res := Response{}
	var reply UserPosition
	args := &UserPosition{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindById", args, &reply)

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
func (c *UserPositionController) Modify() {
	id, _ := c.GetInt64("id")

	var reply UserPosition
	res := Response{}
	args := &UserPosition{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindById", args, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = err.Error()
		c.Data["json"] = res
		c.ServeJSON()
	}
	timestamp := time.Now().UnixNano() / 1e6
	reply.Id = id

	reply.UpdateTime = timestamp

	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "UpdateById", reply, nil)

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
