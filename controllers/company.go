package controllers

import (
	"github.com/astaxie/beego"
	"test/models"
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
	args := &models.Company{
		Id: id,
	}
	err := client.Call("http://127.0.0.1:2379", "Company", "CompanyDetail", args, &reply)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
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

}