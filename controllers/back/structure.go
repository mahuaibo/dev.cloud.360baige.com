package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/models/personnel"
)

type StructureController struct {
	beego.Controller
}

// @router /add [post]
func (c *StructureController) Add() {
	Type, _ := c.GetInt8("type")
	parentId, _ := c.GetInt64("parentId")
	status, _ := c.GetInt8("status")
	var reply Structure
	args := &Structure{
		ParentId: parentId,
		Type:     Type,
		Name:     c.GetString("name"),
		Status:   status,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Structure", "AddStructure", args, &reply)
	var response Response // http 返回体
	if err != nil {
		response.Code = ResponseSystemErr
		response.Messgae = "新增失败！"
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
func (c *StructureController) Modify() {
	id, _ := c.GetInt64("id")
	parentId, _ := c.GetInt64("parentId")
	status, _ := c.GetInt8("status")
	var reply Structure
	args := &Structure{
		Id:       id,
		ParentId: parentId,
		Name:     c.GetString("name"),
		Status:   status,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Structure", "ModifyStructure", args, &reply)
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

//
//// @router /delete [post]
//func (c *StructureController) Delete() {
//	ids := c.GetString("ids")
//	var reply Structure
//	var err error
//	if strings.Contains(ids, ",") {
//		for _, val := range strings.Split(ids, ",") {
//			id, _ := strconv.ParseInt(val, 10, 64)
//			args := &Structure{
//				Id:id,
//			}
//			err = client.Call(beego.AppConfig.String("EtcdURL"), "Structure", "Delete", args, &reply)
//		}
//	} else {
//		id, _ := strconv.ParseInt(ids, 10, 64)
//		args := &Structure{
//			Id:id,
//		}
//		err = client.Call(beego.AppConfig.String("EtcdURL"), "Structure", "Delete", args, &reply)
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
	var reply Structure
	args := &Structure{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Structure", "StructureDetails", args, &reply)
	var response Response // http 返回体
	if err != nil {
		response.Code = ResponseSystemErr
		response.Messgae = "获取失败！"
		c.Data["json"] = response
		c.ServeJSON()
	}
	response.Code = ResponseNormal
	response.Messgae = "获取成功"
	response.Data = reply
	c.Data["json"] = response
	c.ServeJSON()
}

// @router /structureList [post]
func (c *StructureController) StructureList() {
	Type, _ := c.GetInt8("type")
	var reply StructureList
	args := &Structure{
		ParentId: 0,
		Type:     Type,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Structure", "GetStructureList", args, &reply)
	var response Response // http 返回体
	if err != nil {
		response.Code = ResponseSystemErr
		response.Messgae = "查询失败！"
		c.Data["json"] = response
		c.ServeJSON()
	}
	response.Code = ResponseNormal
	response.Messgae = "查询成功"
	response.Data = reply
	c.Data["json"] = response
	c.ServeJSON()
}
