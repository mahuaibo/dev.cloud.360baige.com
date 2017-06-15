package controllers

import (
	"github.com/astaxie/beego"
	"dev.model.360baige.com/models"
	"dev.cloud.360baige.com/rpc/client"
	"fmt"
	"strings"
	"strconv"
)

type PersonStructureController struct {
	beego.Controller
}


// @router /add [post]
func (c *PersonStructureController) Add() {
	OwnerId, _ := c.GetInt64("OwnerId")
	PersonId, _ := c.GetInt64("PersonId")
	StructureId, _ := c.GetInt64("StructureId")
	OwnerType, _ := c.GetInt("OwnerType")
	Type, _ := c.GetInt("Type")
	Status, _ := c.GetInt("Status")
	var reply models.PersonStructure
	args := &models.PersonStructure{
		OwnerId:OwnerId,
		PersonId:PersonId,
		StructureId:StructureId,
		OwnerType:OwnerType,
		Type:Type,
		Status:Status,
		Name:c.GetString("Name"),
		Detail:c.GetString("Detail"),
	}
	err := client.Call("http://127.0.0.1:2379", "PersonStructure", "AddPersonStructure", args, &reply)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
	c.ServeJSON()
}


//// @router /delete [post]
//func (c *PersonStructureController) Delete() {
//	ids := c.GetString("ids")
//	var reply models.PersonStructure
//	var err error
//	if strings.Contains(ids, ",") {
//		for _, val := range strings.Split(ids, ",") {
//			id, _ := strconv.ParseInt(val, 10, 64)
//			args := &models.PersonStructure{
//				Id:id,
//			}
//			err = client.Call("http://127.0.0.1:2379", "PersonStructure", "Delete", args, &reply)
//		}
//	} else {
//		id, _ := strconv.ParseInt(ids, 10, 64)
//		args := &models.PersonStructure{
//			Id:id,
//		}
//		err = client.Call("http://127.0.0.1:2379", "PersonStructure", "Delete", args, &reply)
//	}
//	fmt.Println(reply, err)
//	if err == nil {
//		c.Data["json"] = reply
//	} else {
//		c.Data["json"] = err
//	}
//	c.ServeJSON()
//}
