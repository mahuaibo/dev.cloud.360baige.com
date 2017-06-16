package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.model.360baige.com/models/http"
	"dev.cloud.360baige.com/models/constant"
	. "dev.model.360baige.com/models/personnel"
)

type PersonStructureController struct {
	beego.Controller
}

// @router /add [post]
func (c *PersonStructureController) Add() {
	PersonId, _ := c.GetInt64("PersonId")
	StructureId, _ := c.GetInt64("StructureId")
	Type, _ := c.GetInt8("Type")
	Status, _ := c.GetInt8("Status")
	var reply PersonStructure
	args := &PersonStructure{
		PersonId:    PersonId,
		StructureId: StructureId,
		Type:        Type,
		Status:      Status,
	}
	err := client.Call("http://127.0.0.1:2379", "PersonStructure", "AddPersonStructure", args, &reply)
	var response http.Response // http 返回体
	if err != nil {
		response.Code = constant.ResponseSystemErr
		response.Messgae = "新增失败！"
		c.Data["json"] = response
		c.ServeJSON()
	}
	response.Code = constant.ResponseNormal
	response.Messgae = "新增成功"
	response.Data = reply
	c.Data["json"] = response
	c.ServeJSON()
}

//// @router /delete [post]
//func (c *PersonStructureController) Delete() {
//	ids := c.GetString("ids")
//	var reply PersonStructure
//	var err error
//	if strings.Contains(ids, ",") {
//		for _, val := range strings.Split(ids, ",") {
//			id, _ := strconv.ParseInt(val, 10, 64)
//			args := &PersonStructure{
//				Id:id,
//			}
//			err = client.Call("http://127.0.0.1:2379", "PersonStructure", "Delete", args, &reply)
//		}
//	} else {
//		id, _ := strconv.ParseInt(ids, 10, 64)
//		args := &PersonStructure{
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
