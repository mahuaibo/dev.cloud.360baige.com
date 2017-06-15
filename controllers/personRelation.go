package controllers

import (
	"github.com/astaxie/beego"
	"fmt"
	"dev.model.360baige.com/models"
	"dev.cloud.360baige.com/rpc/client"
)

type PersonRelationController struct {
	beego.Controller
}

// @router /add [post]
func (c *PersonRelationController) Add() {
	associationId, _ := c.GetInt64("associationId")
	associatedId, _ := c.GetInt64("associatedId")
	ownerId, _ := c.GetInt64("ownerId")
	ownerType, _ := c.GetInt("ownerType")
	Type, _ := c.GetInt("type")
	var reply models.PersonRelation
	args := &models.PersonRelation{
		AssociationId:associationId,
		AssociatedId:associatedId,
		OwnerId:ownerId,
		OwnerType:ownerType,
		Type:Type,
		Name:c.GetString("name"),
		Detail:c.GetString("detail"),
	}
	err := client.Call("http://127.0.0.1:2379", "PersonRelation", "GetAssociated", args, &reply)
	if err == nil {
		args := &models.PersonRelation{
			Id:reply.Id,
			Status:1,
		}
		err = client.Call("http://127.0.0.1:2379", "PersonRelation", "ModifyPersonRelation", args, &reply)
	} else {
		err = client.Call("http://127.0.0.1:2379", "PersonRelation", "AddPersonRelation", args, &reply)
	}
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
	c.ServeJSON()
}


// @router /modify [post]
func (c *PersonRelationController) Modify() {
	id, _ := c.GetInt64("id")
	status, _ := c.GetInt("status")
	var reply models.PersonRelation
	args := &models.PersonRelation{
		Id:id,
		Status:status,
	}
	err := client.Call("http://127.0.0.1:2379", "PersonRelation", "ModifyPersonRelation", args, &reply)
	fmt.Println(reply, err)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
	c.ServeJSON()
}


// @router /associatedList [post]
func (c *PersonRelationController) AssociatedList() {
	associationId, _ := c.GetInt64("id")
	var reply models.PersonRelation
	args := &models.PersonRelation{
		AssociationId:associationId,
	}
	err := client.Call("http://127.0.0.1:2379", "PersonRelation", "GetAssociatedAll", args, &reply)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = "关注列表为空"
	}
	c.ServeJSON()
}


