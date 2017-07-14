package window

import (
	"github.com/astaxie/beego"
	//"dev.cloud.360baige.com/rpc/client"
	////"dev.cloud.360baige.com/utils"
	//. "dev.model.360baige.com/models/logger"
	//. "dev.model.360baige.com/models/response"
	//. "dev.model.360baige.com/http/window"
	////"fmt"
	//"time"
)

type LoggerController struct {
	beego.Controller
}

// @Title 新增add
// @Description post user by uid
// @Param	content		path 	string	true		"内容"
// @Param	remark		path 	string	true		"描述"
// @Param	ip		path 	string	true		"IP地址	"
// @Param	type		path 	int	true		"类别（增、删、改、查）"
// @Param	ownerId		path 	int64	true		"操作者ID"
// @Param	companyId		path 	int64	true		"公司ID"
// @Success 200 {object} models.logger
// @Failure 403 :uid is empty
// @router /add [post]
func (c *LoggerController) Add() {
	//res := LoggerResponse{}
	//Type, _ := c.GetInt("type")
	//ownerId, _ := c.GetInt64(" ownerId")
	//companyId, _ := c.GetInt64("companyId")
	//var reply Logger
	//args := Logger{
	//	CreateTime: time.Now().Unix(),
	//	Content:c.GetString("content"),
	//	Remark:c.GetString("remark"),
	//	Type:Type,
	//	OwnerId:ownerId,
	//	CompanyId:companyId,
	//}
	//err = client.Call(beego.AppConfig.String("EtcdURL"), "Logger", "Add", args, &reply)
	//if err == nil {
	//	res.Code = ResponseNormal
	//	res.Messgae = "新增成功"
	//	res.Data.ExpireIn = reply.ExpireIn
	//} else {
	//	res.Code = ResponseSystemErr
	//	res.Messgae = "新增失败"
	//}
	//c.Data["json"] = res
	//c.ServeJSON()
}

// @Title 列表list
// @Description post user by uid
// @Param	page		path 	int	true		"分页页码"
// @Param	rows		path 	int	true		"展示条数"
// @Param	companyId		path 	int64	true		"公司ID"
// @Success 200 {object} models.logger
// @Failure 403 :uid is empty
// @router /getList [post]
func (c *LoggerController) GetList() {
	//res := LoggerResponse{}
	//page, _ := c.GetInt("page")
	//rows, _ := c.GetInt("rows")
	//companyId, _ := c.GetInt64("companyId")
	//var reply models.LoggerList
	//args := models.LoggerPaging{
	//	Page:page,
	//	Rows:rows,
	//	CompanyId:companyId,
	//}
	//err = client.Call(beego.AppConfig.String("EtcdURL"), "Logger", "List", args, &reply)
	//if err == nil {
	//	c.Data["json"] = reply
	//} else {
	//	c.Data["json"] = err
	//}
	//c.ServeJSON()
}