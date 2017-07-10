package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"time"
	. "dev.model.360baige.com/models/application"
	. "dev.model.360baige.com/models/paginator"
	. "dev.model.360baige.com/models/response"
)

type ApplicationController struct {
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
// @router /getapplist [post]
func (c *ApplicationController) GetAppList() {
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
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Application", "List", args, &reply)
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
// @Title 应用菜单修改接口
// @Description 应用菜单修改接口
// @Success 200 {"code":200,"messgae":"ok","accessToken":"ok"}}
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Param personId query   string true       "身份id"
// @Param appId query   string true       "应用id"
// @Param apppos query   string true       "应用顺序位置"
// @Failure 400 {"code":400,"message":"..."}
// @router /modifyapp [post]
func (c *ApplicationController) ModifyApp() {
	id, _ := c.GetInt64("id")

	var reply Application
	res := Response{}
	args := &Application{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindById", args, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = err.Error()
		c.Data["json"] = res
		c.ServeJSON()
	}
	timestamp := time.Now().UnixNano() / 1e6
	reply.Id = id

	reply.UpdateTime = timestamp

	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "UpdateById", reply, nil)

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

// @Title 应用搜索接口
// @Description 应用搜索接口
// @Success 200 {"code":200,"messgae":"ok","data":{"list":{ ... ... },"accessToken":"ok"}}
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Param personId query   string true       "身份id"
// @Param name query   string true       "应用名称"
// @Failure 400 {"code":400,"message":"..."}
// @router /search-applist [post]
func (c *ApplicationController) SearchAppList() {
	id, _ := c.GetInt64("id")

	var reply Application
	res := Response{}
	args := &Application{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindById", args, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = err.Error()
		c.Data["json"] = res
		c.ServeJSON()
	}
	timestamp := time.Now().UnixNano() / 1e6
	reply.Id = id

	reply.UpdateTime = timestamp

	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "UpdateById", reply, nil)

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

// @Title 应用相关接口
// @Description 应用相关接口
// @Success 200 {"code":200,"messgae":"ok","data":{"list":{ ... ... },"accessToken":"ok"}}
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Param personId query   string true       "身份id"
// @Param appId query   string true       "应用id"
// @Param datatype query   string true       "类型(1,相关，2 热门 3 同作者 )"
// @Failure 400 {"code":400,"message":"..."}
// @router /get-appcorrelationtlist [post]
func (c *ApplicationController) GetAppcorrelationtList() {
	id, _ := c.GetInt64("id")

	var reply Application
	res := Response{}
	args := &Application{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindById", args, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = err.Error()
		c.Data["json"] = res
		c.ServeJSON()
	}
	timestamp := time.Now().UnixNano() / 1e6
	reply.Id = id

	reply.UpdateTime = timestamp

	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "UpdateById", reply, nil)

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
func (c *ApplicationController) Add() {
	timestamp := time.Now().UnixNano() / 1e6
	var (
		res   Response // http 返回体
		reply Application
	)
	args := &Application{
		CreateTime: timestamp,
		UpdateTime: timestamp,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Application", "Add", args, &reply)
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
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Param personId query   string true       "身份id"
// @Failure 400 {"code":400,"message":"..."}
// @router /detail [post]
func (c *ApplicationController) Detail() {
	id, _ := c.GetInt64("id")
	res := Response{}
	var reply Application
	args := &Application{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindById", args, &reply)

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
func (c *ApplicationController) Modify() {
	id, _ := c.GetInt64("id")

	var reply Application
	res := Response{}
	args := &Application{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindById", args, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = err.Error()
		c.Data["json"] = res
		c.ServeJSON()
	}
	timestamp := time.Now().UnixNano() / 1e6
	reply.Id = id

	reply.UpdateTime = timestamp

	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "UpdateById", reply, nil)

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
