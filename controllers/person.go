package controllers

import (
	"github.com/astaxie/beego"
	"dev.model.360baige.com/models"
	"dev.cloud.360baige.com/rpc/client"
	"fmt"
	"strings"
	"strconv"
)

type PersonController struct {
	beego.Controller
}

// @router /add [post]
func (c *PersonController) Add() {
	ownerId, _ := c.GetInt64("ownerId")
	ownerType, _ := c.GetInt("ownerType")
	Type, _ := c.GetInt("type")
	var reply models.Person
	args := &models.Person{
		OwnerId: ownerId,
		OwnerType:ownerType,
		Type:Type,
		Name:c.GetString("name"),
		Detail:c.GetString("detail"),
	}
	err := client.Call("http://127.0.0.1:2379", "Person", "AddPerson", args, &reply)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
	c.ServeJSON()
}

//// @router /detail [post]
//func (c *PersonController) Detail() {
//	var data map[string]interface{} = make(map[string]interface{})
//	id, _ := c.GetInt64("id")
//	var personReply models.Person
//	personArgs := &models.Person{
//		Id: id,
//	}
//	err := client.Call("http://127.0.0.1:2379", "Person", "Details", personArgs, &personReply)
//	data["Id"] = personReply.Id
//	data["CreateTime"] = personReply.CreateTime
//	data["UpdateTime"] = personReply.UpdateTime
//	data["OwnerId"] = personReply.OwnerId
//	data["Name"] = personReply.Name
//	data["Detail"] = personReply.Detail
//	data["Type"] = personReply.Type
//	data["Status"] = personReply.Status
//
//	var reply []models.AssociatedAll
//	args := models.AssociatedArgs{
//		AssociatedId:personReply.Id,
//		AssociationId:personReply.Id,
//	}
//	//拼接关联人
//	client.Call("http://127.0.0.1:2379", "PersonRelation", "GetAssociatedAll", args, &reply)
//	jointDetailData(data, reply, "association")
//	//拼接被关联人
//	client.Call("http://127.0.0.1:2379", "PersonRelation", "GetBeAssociatedAll", args, &reply)
//	jointDetailData(data, reply, "associated")
//	//拼接结构
//	client.Call("http://127.0.0.1:2379", "Structure", "GetStructure", args, &reply)
//	jointDetailData(data, reply, "structure")
//	fmt.Println(data, err)
//	if err == nil {
//		c.Data["json"] = data
//	}else{
//		c.Data["json"] = err
// 	}
//	c.ServeJSON()
//}

//func jointDetailData(result map[string]interface{}, data []models.AssociatedAll, name string) {
//	if len(data) > 0 {
//		for key, val := range data {
//			data[key] = val
//		}
//		result[name] = data
//	}
//}

// @router /modify [post]
func (c *PersonController) Modify() {
	id, _ := c.GetInt64("id")
	ownerId, _ := c.GetInt64("ownerId")
	Type, _ := c.GetInt("type")
	status, _ := c.GetInt("status")
	var reply models.Person
	args := &models.Person{
		Id: id,
		OwnerId: ownerId,
		Type:Type,
		Status:status,
		Name:c.GetString("name"),
		Detail:c.GetString("detail"),
	}
	err := client.Call("http://127.0.0.1:2379", "Person", "ModifyPerson", args, &reply)
	fmt.Println(reply, err)
	if err != nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
	c.ServeJSON()
}

//// @router /delete [post]
//func (c *PersonController) Delete() {
//	ids := c.GetString("ids")
//	var reply models.Person
//	var err error
//	if strings.Contains(ids, ",") {
//		for _, val := range strings.Split(ids, ",") {
//			id, _ := strconv.ParseInt(val, 10, 64)
//			args := &models.Person{
//				Id: id,
//			}
//			err = client.Call("http://127.0.0.1:2379", "Person", "ModifyPersonStatus", args, &reply)
//		}
//	} else {
//		id, _ := strconv.ParseInt(ids, 10, 64)
//		args := &models.Person{
//			Id: id,
//		}
//		err = client.Call("http://127.0.0.1:2379", "Person", "ModifyPersonStatus", args, &reply)
//	}
//	fmt.Println(reply, err)
//	if err == nil {
//		c.Data["json"] = reply
//	} else {
//		c.Data["json"] = err
//	}
//	c.ServeJSON()
//}


// @router /personList [post]
func (c *PersonController) PersonList() {
	page, _ := c.GetInt("page")
	rows, _ := c.GetInt("rows")
	Type, _ := c.GetInt("type")
	var reply models.PersonList
	args := &models.PersonPaging{
		Type:Type,
		Page:page,
		Rows:rows,
	}
	err := client.Call("http://127.0.0.1:2379", "Person", "GetPersonList", args, &reply)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
	c.ServeJSON()
}
