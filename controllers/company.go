package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"time"
	. "dev.model.360baige.com/models/company"
	. "dev.model.360baige.com/models/response"
)

type CompanyController struct {
	beego.Controller
}
// @Title 学校地区接口
// @Description 学校地区接口
// @Success 200 {"code":200,"messgae":"ok","data":{"list":{... ...},"accessToken":"ok"}}
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Param personId query   string true       "身份id"
// @Param userId query   string true       "userid"
// @Failure 400 {"code":400,"message":"..."}
// @router /getschoolarea [post]
func (c *PersonController) GetSchoolArea() {
	var (
		res   Response // http 返回体
		reply Company
		args Company
	)
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Company", "List", args, &reply)
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
// @Title 企业新增
// @Description 企业信息
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   name     query   string true       "公司名称"
// @Param   shortName     query   string true       "公司简称"
// @Param   type     query   string true       "公司类型"
// @Param   logo     query   string true       "公司LOGO"
// @Param   brief     query   string true       "公司简介"
// @Failure 400 {"code":400,"message":"..."}
// @router /add [get]
func (c *CompanyController) Add() {
	name := c.GetString("name")
	shortName := c.GetString("shortName")
	logo := c.GetString("logo")
	brief := c.GetString("brief")
	companyType, _ := c.GetInt8("type")
	//accessToken := c.GetString("accessToken")
	timestamp := time.Now().UnixNano() / 1e6
	var (
		res   Response // http 返回体
		reply Company
	)
	args := &Company{
		CreateTime: timestamp,
		UpdateTime: timestamp,
		Name:       name,
		ShortName:  shortName,
		Type:       companyType,
		Logo:       logo,
		Brief:      brief,
		Status:     CompanyStatusNormal,
	}

	err := client.Call(beego.AppConfig.String("EtcdURL"), "Company", "Add", args, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "企业新增失败"
		c.Data["json"] = res
		c.ServeJSON()
	}
	res.Code = ResponseNormal
	res.Messgae = "企业新增成功"
	res.Data = reply
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 企业信息
// @Description 企业信息
// @Success 200 {"code":200,"messgae":"企业信息查询成功", "data":{ ... ... }}
// @Param   companyId     query   string true       "企业ID"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /detail [get]
func (c *CompanyController) Detail() {
	companyId, _ := c.GetInt64("companyId")
	res := Response{}
	var reply Company
	args := &Company{
		Id: companyId,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Company", "FindById", args, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "企业信息查询失败"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "企业信息查询成功"
		res.Data = reply
		c.Data["json"] = res
		c.ServeJSON()
	}
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

	var reply Company
	res := Response{}
	args := &Company{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Company", "FindById", args, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = err.Error()
		c.Data["json"] = res
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

	err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "UpdateById", reply, nil)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "企业信息修改失败！"
		c.Data["json"] = res
		c.ServeJSON()
	}

	res.Code = ResponseNormal
	res.Messgae = "企业信息修改成功！"
	c.Data["json"] = res
	c.ServeJSON()
}
