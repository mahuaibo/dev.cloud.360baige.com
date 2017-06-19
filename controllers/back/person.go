package back

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/models/personnel"
)

type PersonController struct {
	beego.Controller
}

// @router /add [post]
func (c *PersonController) Add() {
	companyId, _ := c.GetInt64("companyId")
	Type, _ := c.GetInt8("type")
	status, _ := c.GetInt8("status")
	var reply Person
	args := &Person{
		Name:      c.GetString("name"),
		CompanyId: companyId,
		Type:      Type,
		Status:    status,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Person", "AddPerson", args, &reply)
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

// @router /detail [post]
func (c *PersonController) Detail() {
	var data map[string]interface{} = make(map[string]interface{})
	id, _ := c.GetInt64("id")
	var personReply Person
	personArgs := &Person{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Person", "Details", personArgs, &personReply)
	data["Id"] = personReply.Id
	data["CreateTime"] = personReply.CreateTime
	data["UpdateTime"] = personReply.UpdateTime
	data["Name"] = personReply.Name
	data["Type"] = personReply.Type
	data["Status"] = personReply.Status

	var reply []AssociatedAll
	args := AssociatedArgs{
		AssociatedId:  personReply.Id,
		AssociationId: personReply.Id,
	}
	//拼接关联人
	client.Call(beego.AppConfig.String("EtcdURL"), "PersonRelation", "GetAssociatedAll", args, &reply)
	jointDetailData(data, reply, "association")
	//拼接被关联人
	client.Call(beego.AppConfig.String("EtcdURL"), "PersonRelation", "GetBeAssociatedAll", args, &reply)
	jointDetailData(data, reply, "associated")
	//拼接结构
	client.Call(beego.AppConfig.String("EtcdURL"), "Structure", "GetStructure", args, &reply)
	jointDetailData(data, reply, "structure")
	var response Response // http 返回体
	if err != nil {
		response.Code = ResponseSystemErr
		response.Messgae = "获取失败"
		c.Data["json"] = response
		c.ServeJSON()
	}
	response.Code = ResponseNormal
	response.Messgae = "获取成功"
	response.Data = reply
	c.Data["json"] = response
	c.ServeJSON()
}

func jointDetailData(result map[string]interface{}, data []AssociatedAll, name string) {
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
	var reply Person
	args := &Person{
		Id:     id,
		Name:   c.GetString("name"),
		Type:   Type,
		Status: status,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Person", "ModifyPerson", args, &reply)
	var response Response // http 返回体
	if err != nil {
		response.Code = ResponseSystemErr
		response.Messgae = "修改失败"
		c.Data["json"] = response
		c.ServeJSON()
	}
	response.Code = ResponseNormal
	response.Messgae = "修改成功"
	response.Data = reply
	c.Data["json"] = response
	c.ServeJSON()
}

//// @router /delete [post]
//func (c *PersonController) Delete() {
//	ids := c.GetString("ids")
//	var reply Person
//	var err error
//	if strings.Contains(ids, ",") {
//		for _, val := range strings.Split(ids, ",") {
//			id, _ := strconv.ParseInt(val, 10, 64)
//			args := &Person{
//				Id: id,
//			}
//			err = client.Call(beego.AppConfig.String("EtcdURL"), "Person", "ModifyPersonStatus", args, &reply)
//		}
//	} else {
//		id, _ := strconv.ParseInt(ids, 10, 64)
//		args := &Person{
//			Id: id,
//		}
//		err = client.Call(beego.AppConfig.String("EtcdURL"), "Person", "ModifyPersonStatus", args, &reply)
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
	var reply PersonList
	args := &PersonPaging{
		Type: Type,
		Page: page,
		Rows: rows,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Person", "GetPersonList", args, &reply)
	var response Response // http 返回体
	if err != nil {
		response.Code = ResponseSystemErr
		response.Messgae = "获取失败"
		c.Data["json"] = response
		c.ServeJSON()
	}
	response.Code = ResponseNormal
	response.Messgae = "获取成功"
	response.Data = reply
	c.Data["json"] = response
	c.ServeJSON()
}

type Paginator struct {
	Type      int            //分页方式 按页码:1 按ID:2
	Total     int            //总页数
	PageSize  int            //每页数量
	Current   int            //当前页码
	MarkID    int64          //最大ID 或 最小ID
	Direction int            //翻页方向 上:-1 平:0 下:1
	List      []*interface{} //内容
}
