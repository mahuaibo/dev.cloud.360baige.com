package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	//"dev.cloud.360baige.com/utils"
	. "dev.model.360baige.com/models/logger"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/http/window"
	"time"
	"dev.model.360baige.com/action"
	"encoding/json"
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
	res := LoggerAddResponse{}
	access_token := c.GetString("access_token")
	Type, _ := c.GetInt8("type")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	//检测 accessToken
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type:  "And",
		Key: "accessToken",
		Val:  access_token,
	})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	com_id := replyAccessToken.CompanyId
	user_id := replyAccessToken.UserId
	user_position_id := replyAccessToken.Id
	user_position_type := replyAccessToken.Type
	if com_id == 0 || user_id == 0 || user_position_id == 0 {
		res.Code = ResponseSystemErr
		res.Messgae = "获取应用信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var reply Logger
	loggerArgs := Logger{
		CreateTime:       time.Now().Unix(),
		Content:          c.GetString("content"),
		Remark:           c.GetString("remark"),
		Type:             Type,
		UserId:           user_id,
		CompanyId:        com_id,
		UserPositionId:   user_position_id,
		UserPositionType: user_position_type,
	}
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Logger", "Add", loggerArgs, &reply)
	if err == nil {
		res.Code = ResponseNormal
		res.Messgae = "新增成功"
		res.Data.Id = reply.Id
	} else {
		res.Code = ResponseSystemErr
		res.Messgae = "新增失败"
	}
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 列表接口
// @Description 列表接口
// @Success 200 {"code":200,"messgae":"获取列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   current     query   string true       "当前页"
// @Param   page_size     query   string true       "每页数量"
// @Failure 400 {"code":400,"message":"获取信息失败"}
// @router /list [get]
func (c *LoggerController) List() {
	res := LoggerListResponse{}
	access_token := c.GetString("access_token")
	currentPage, _ := c.GetInt64("current")
	pageSize, _ := c.GetInt64("page_size")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	//检测 accessToken
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type:  "And",
		Key: "accessToken",
		Val:  access_token,
	})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	//company_id、user_id、user_position_id、user_position_type
	com_id := replyAccessToken.CompanyId
	user_id := replyAccessToken.UserId
	user_position_id := replyAccessToken.Id
	user_position_type := replyAccessToken.Type
	if com_id == 0 || user_id == 0 || user_position_id == 0 {
		res.Code = ResponseSystemErr
		res.Messgae = "获取信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var loggerArgs action.PageByCond
	var reply action.PageByCond
	loggerArgs.CondList = append(loggerArgs.CondList, action.CondValue{
		Type: "And",
		Key:  "company_id",
		Val:  com_id,
	}, action.CondValue{
		Type: "And",
		Key:  "user_id",
		Val:  user_id,
	}, action.CondValue{
		Type: "And",
		Key:  "user_position_id",
		Val:  user_position_id,
	}, action.CondValue{
		Type: "And",
		Key:  "user_position_type",
		Val:  user_position_type,
	})
	loggerArgs.OrderBy = []string{"id"}
	loggerArgs.Cols = []string{"id", "create_time", "content", "remark", "type", }
	loggerArgs.PageSize = pageSize
	loggerArgs.Current = currentPage
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Logger", "PageByCond", loggerArgs, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var resData []LoggerValue
	replyList := []Logger{}
	err = json.Unmarshal([]byte(reply.Json), &replyList)
	for _, value := range replyList {
		var retype string
		if value.Type == 1 {
			retype = "增"
		} else if value.Type == 2 {
			retype = "删"
		} else if value.Type == 3 {
			retype = "改"
		} else {
			retype = "查"
		}
		resData = append(resData, LoggerValue{
			CreateTime: time.Unix(value.CreateTime / 1000, 0).Format("2006-01-02"),
			Content:    value.Content,
			Remark:     value.Remark,
			Type:       retype,
		})
	}

	res.Code = ResponseNormal
	res.Messgae = "获取成功"
	res.Data.Total = reply.Total
	res.Data.Current = currentPage
	res.Data.CurrentSize = reply.CurrentSize
	res.Data.OrderBy = reply.OrderBy
	res.Data.PageSize = pageSize
	res.Data.List = resData
	c.Data["json"] = res
	c.ServeJSON()
}
