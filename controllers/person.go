package controllers

import (
	"github.com/astaxie/beego"
	"dev.model.360baige.com/models"
	"dev.cloud.360baige.com/rpc/client"
	"dev.cloud.360baige.com/models/response"
	"dev.cloud.360baige.com/models/constant"
)

type PersonController struct {
	beego.Controller
}

// @router /add [post]
func (c *PersonController) Add() {
	companyId, _ := c.GetInt64("companyId")
	Type, _ := c.GetInt8("type")
	status, _ := c.GetInt8("status")
	var reply models.Person
	args := &models.Person{
		Name:c.GetString("name"),
		CompanyId:companyId,
		Type:Type,
		Status:status,
	}
	err := client.Call("http://127.0.0.1:2379", "Person", "AddPerson", args, &reply)
	var response response.Response // http 返回体
	if err != nil {
		response.Code = constant.ResponseSystemErr
		response.Messgae = "新增失败"
		c.Data["json"] = response
		c.ServeJSON()
	}
	response.Code = constant.ResponseNormal
	response.Messgae = "新增成功"
	response.Data = reply
	c.Data["json"] = response
	c.ServeJSON()
}

// @router /detail [post]
func (c *PersonController) Detail() {
	var data map[string]interface{} = make(map[string]interface{})
	id, _ := c.GetInt64("id")
	var personReply models.Person
	personArgs := &models.Person{
		Id: id,
	}
	err := client.Call("http://127.0.0.1:2379", "Person", "Details", personArgs, &personReply)
	data["Id"] = personReply.Id
	data["CreateTime"] = personReply.CreateTime
	data["UpdateTime"] = personReply.UpdateTime
	data["Name"] = personReply.Name
	data["Type"] = personReply.Type
	data["Status"] = personReply.Status

	var reply []models.AssociatedAll
	args := models.AssociatedArgs{
		AssociatedId:personReply.Id,
		AssociationId:personReply.Id,
	}
	//拼接关联人
	client.Call("http://127.0.0.1:2379", "PersonRelation", "GetAssociatedAll", args, &reply)
	jointDetailData(data, reply, "association")
	//拼接被关联人
	client.Call("http://127.0.0.1:2379", "PersonRelation", "GetBeAssociatedAll", args, &reply)
	jointDetailData(data, reply, "associated")
	//拼接结构
	client.Call("http://127.0.0.1:2379", "Structure", "GetStructure", args, &reply)
	jointDetailData(data, reply, "structure")
	var response response.Response // http 返回体
	if err != nil {
		response.Code = constant.ResponseSystemErr
		response.Messgae = "获取失败"
		c.Data["json"] = response
		c.ServeJSON()
	}
	response.Code = constant.ResponseNormal
	response.Messgae = "获取成功"
	response.Data = reply
	c.Data["json"] = response
	c.ServeJSON()
}

func jointDetailData(result map[string]interface{}, data []models.AssociatedAll, name string) {
	if len(data) > 0 {
		for key, val := range data {
			data[key] = val
		}
		result[name] = data
	}
}

// @router /modify [post]
func (c *PersonController) Modify() {
	id, _ := c.GetInt64("id")
	Type, _ := c.GetInt8("type")
	status, _ := c.GetInt8("status")
	var reply models.Person
	args := &models.Person{
		Id: id,
		Name:c.GetString("name"),
		Type:Type,
		Status:status,
	}
	err := client.Call("http://127.0.0.1:2379", "Person", "ModifyPerson", args, &reply)
	var response response.Response // http 返回体
	if err != nil {
		response.Code = constant.ResponseSystemErr
		response.Messgae = "修改失败"
		c.Data["json"] = response
		c.ServeJSON()
	}
	response.Code = constant.ResponseNormal
	response.Messgae = "修改成功"
	response.Data = reply
	c.Data["json"] = response
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
	Type, _ := c.GetInt8("type")
	var reply models.PersonList
	args := &models.PersonPaging{
		Type:Type,
		Page:page,
		Rows:rows,
	}
	err := client.Call("http://127.0.0.1:2379", "Person", "GetPersonList", args, &reply)
	var response response.Response // http 返回体
	if err != nil {
		response.Code = constant.ResponseSystemErr
		response.Messgae = "获取失败"
		c.Data["json"] = response
		c.ServeJSON()
	}
	response.Code = constant.ResponseNormal
	response.Messgae = "获取成功"
	response.Data = reply
	c.Data["json"] = response
	c.ServeJSON()
}
