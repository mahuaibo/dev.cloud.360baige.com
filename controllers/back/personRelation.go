package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/models/personnel"
)

type PersonRelationController struct {
	beego.Controller
}

// @router /add [post]
func (c *PersonRelationController) Add() {
	companyId, _ := c.GetInt64("companyId")
	associationId, _ := c.GetInt64("associationId")
	associatedId, _ := c.GetInt64("associatedId")
	status, _ := c.GetInt8("status")
	Type, _ := c.GetInt8("type")
	var reply PersonRelation
	args := &PersonRelation{
		CompanyId:     companyId,
		AssociationId: associationId,
		AssociatedId:  associatedId,
		Type:          Type,
		Status:        status,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "PersonRelation", "GetAssociated", args, &reply)
	if err == nil {
		args := &PersonRelation{
			Id:     reply.Id,
			Status: 1,
		}
		err = client.Call(beego.AppConfig.String("EtcdURL"), "PersonRelation", "ModifyPersonRelation", args, &reply)
	} else {
		err = client.Call(beego.AppConfig.String("EtcdURL"), "PersonRelation", "AddPersonRelation", args, &reply)
	}
	var response Response // http 返回体
	if err != nil {
		response.Code = ResponseSystemErr
		response.Messgae = "新增失败"
		c.Data["json"] = response
		c.ServeJSON()
	}
	response.Code = ResponseNormal
	response.Messgae = "新增成功"
	response.Data = reply
	c.Data["json"] = response
	c.ServeJSON()
}

// @router /modify [post]
func (c *PersonRelationController) Modify() {
	id, _ := c.GetInt64("id")
	Type, _ := c.GetInt8("type")
	status, _ := c.GetInt8("status")
	var reply PersonRelation
	args := &PersonRelation{
		Id:     id,
		Type:   Type,
		Status: status,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "PersonRelation", "Modify", args, &reply)
	var response Response // http 返回体
	if err != nil {
		response.Code = ResponseSystemErr
		response.Messgae = "修改失败！"
		c.Data["json"] = response
		c.ServeJSON()
	}
	response.Code = ResponseNormal
	response.Messgae = "修改成功"
	response.Data = reply
	c.Data["json"] = response
	c.ServeJSON()
}

// @router /associatedList [post]
func (c *PersonRelationController) AssociatedList() {
	associationId, _ := c.GetInt64("id")
	var reply PersonRelation
	args := &PersonRelation{
		AssociationId: associationId,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "PersonRelation", "GetAssociatedAll", args, &reply)
	var response Response // http 返回体
	if err != nil {
		response.Code = ResponseSystemErr
		response.Messgae = "关注列表为空！"
		c.Data["json"] = response
		c.ServeJSON()
	}
	response.Code = ResponseNormal
	response.Messgae = "关注列表获取成功"
	response.Data = reply
	c.Data["json"] = response
	c.ServeJSON()
}