package controllers

import (
	"github.com/astaxie/beego"
	"dev.model.360baige.com/models"
	"dev.cloud.360baige.com/rpc/client"
	"fmt"
	"strings"
	"strconv"
)

type StructureController struct {
	beego.Controller
}

// @router /add [post]
func (c *StructureController) Add() {
	ownerId, _ := c.GetInt64("ownerId")
	ownerType, _ := c.GetInt("ownerType")
	Type, _ := c.GetInt("type")
	parentId, _ := c.GetInt64("parentId")
	status, _ := c.GetInt("status")
	var reply models.Structure
	args := &models.Structure{
		OwnerId:ownerId,
		OwnerType:ownerType,
		ParentId:parentId,
		Type:Type,
		Name: c.GetString("name"),
		Detail:c.GetString("detail"),
		Status:status,
	}
	err := client.Call("http://127.0.0.1:2379", "Structure", "AddStructure", args, &reply)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
	c.ServeJSON()
}

// @router /modify [post]
func (c *StructureController) Modify() {
	id, _ := c.GetInt64("id")
	parentId, _ := c.GetInt64("parentId")
	status, _ := c.GetInt("status")
	var reply models.Structure
	args := &models.Structure{
		Id:id,
		ParentId:parentId,
		Name:c.GetString("name"),
		Detail:c.GetString("detail"),
		Status:status,
	}
	err := client.Call("http://127.0.0.1:2379", "Structure", "ModifyStructure", args, &reply)
	fmt.Println(reply, err)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
	c.ServeJSON()
}
//
//// @router /delete [post]
//func (c *StructureController) Delete() {
//	ids := c.GetString("ids")
//	var reply models.Structure
//	var err error
//	if strings.Contains(ids, ",") {
//		for _, val := range strings.Split(ids, ",") {
//			id, _ := strconv.ParseInt(val, 10, 64)
//			args := &models.Structure{
//				Id:id,
//			}
//			err = client.Call("http://127.0.0.1:2379", "Structure", "Delete", args, &reply)
//		}
//	} else {
//		id, _ := strconv.ParseInt(ids, 10, 64)
//		args := &models.Structure{
//			Id:id,
//		}
//		err = client.Call("http://127.0.0.1:2379", "Structure", "Delete", args, &reply)
//	}
//	fmt.Println(reply, err)
//	if err == nil {
//		c.Data["json"] = reply
//	} else {
//		c.Data["json"] = err
//	}
//	c.ServeJSON()
//}


// @router /detail [post]
func (c *StructureController) Detail() {
	id, _ := c.GetInt64("id")
	var reply models.Structure
	args := &models.Structure{
		Id:id,
	}
	err := client.Call("http://127.0.0.1:2379", "Structure", "StructureDetails", args, &reply)
	fmt.Println(reply, err)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
	c.ServeJSON()
}


// @router /structureList [post]
func (c *StructureController) StructureList() {
	ownerId, _ := c.GetInt64("ownerId")
	Type, _ := c.GetInt("type")
	var reply models.StructureList
	args := &models.Structure{
		OwnerId:ownerId,
		ParentId:0,
		Type:Type,
	}
	err := client.Call("http://127.0.0.1:2379", "Structure", "GetStructureList", args, &reply)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
	c.ServeJSON()
}