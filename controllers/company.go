package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.cloud.360baige.com/models/response"
	"dev.model.360baige.com/models"
	"dev.cloud.360baige.com/models/constant"
	"time"
)

type CompanyController struct {
	beego.Controller
}

// @Title 企业新增
// @Description 企业信息
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   userId     query   string true       "用户ID"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /detail [get]
func (c *CompanyController) Add() {
	var reply models.Company
	args := &models.Company{
		CreateTime:time.Now().UnixNano() / 1e6,
		UpdateTime:time.Now().UnixNano() / 1e6,
		Name: c.GetString("name"),
	}
	err := client.Call("http://127.0.0.1:2379", "Company", "AddCompany", args, &reply)
	if err == nil {
		// TODO 注册成功添加角色
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
	c.ServeJSON()
}

// @Title 企业信息
// @Description 企业信息
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   userId     query   string true       "用户ID"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /detail [get]
func (c *CompanyController) Detail() {
	id, _ := c.GetInt64("Id")
	var reply models.Company
	response := response.Response{}
	args := &models.Company{
		Id: id,
	}
	err := client.Call("http://127.0.0.1:2379", "Company", "FindCompanyById", args, &reply)

	if err != nil {
		response.Code = constant.ResponseSystemErr
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}
	response.Code = constant.ResponseNormal
	response.Messgae = "获取企业信息成功"
	response.Data = reply
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title 企业信息修改
// @Description 企业信息修改
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   userId     query   string true       "用户ID"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /modify [post]
func (c *CompanyController) Modify() {
	id, _ := c.GetInt64("Id")
	name := c.GetString("Name")
	shortName := c.GetString("ShortName")
	address := c.GetString("Address")

	var reply models.Company
	response := response.Response{}
	args := &models.Company{
		Id: id,
	}
	err := client.Call("http://127.0.0.1:2379", "Company", "FindCompanyById", args, &reply)

	if err != nil {
		response.Code = constant.ResponseSystemErr
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}
	timestamp := time.Now().UnixNano() / 1e6
	reply.Id = id
	if name == "" {
		reply.Name = name
	}
	reply.ShortName = shortName
	reply.UpdateTime = timestamp
	reply.Address = address

	err = client.Call("http://127.0.0.1:2379", "Company", "UpdateCompanyById", reply, nil)

	if err != nil {
		response.Code = constant.ResponseSystemErr
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}

	response.Code = constant.ResponseNormal
	response.Messgae = "用户信息修改成功！"
	c.Data["json"] = response
	c.ServeJSON()
}
