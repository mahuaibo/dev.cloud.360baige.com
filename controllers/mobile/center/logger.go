package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.model.360baige.com/models/logger"
	. "dev.model.360baige.com/http/mobile/center"
	"dev.model.360baige.com/action"
	"encoding/json"
	"dev.cloud.360baige.com/utils"
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
	type data LoggerAddResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	content := c.GetString("content")
	remark := c.GetString("remark")
	Type, _ := c.GetInt("type")

	err := utils.Unable(map[string]string{"accessToken": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	replyUserPosition, err := utils.UserPosition(accessToken, currentTimestamp)
	if err != nil {
		c.Data["json"] = data{Code: ErrorPower, Message: err.Error()}
		c.ServeJSON()
		return
	}

	var replyLogger logger.Logger
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Logger", "Add", logger.Logger{
		CreateTime:       currentTimestamp,
		Content:          content,
		Remark:           remark,
		Type:             Type,
		UserId:           replyUserPosition.UserId,
		CompanyId:        replyUserPosition.CompanyId,
		UserPositionId:   replyUserPosition.Id,
		UserPositionType: replyUserPosition.Type,
	}, &replyLogger)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: LoggerAdd{
		Id: replyLogger.Id,
	}}
	c.ServeJSON()
	return
}

// @Title 列表接口
// @Description 列表接口
// @Success 200 {"code":200,"message":"获取列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   current     query   string true       "当前页"
// @Param   pageSize     query   string true       "每页数量"
// @Failure 400 {"code":400,"message":"获取信息失败"}
// @router /list [post]
func (c *LoggerController) List() {
	type data LoggerListResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	currentPage, _ := c.GetInt64("current")
	pageSize, _ := c.GetInt64("pageSize")

	err := utils.Unable(map[string]string{"accessToken": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	replyUserPosition, err := utils.UserPosition(accessToken, currentTimestamp)
	if err != nil {
		c.Data["json"] = data{Code: ErrorPower, Message: err.Error()}
		c.ServeJSON()
		return
	}

	var replyPageByCond action.PageByCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Logger", "PageByCond", action.PageByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId},
			action.CondValue{Type: "And", Key: "user_id", Val: replyUserPosition.UserId},
			action.CondValue{Type: "And", Key: "user_position_id", Val: replyUserPosition.Id},
			action.CondValue{Type: "And", Key: "user_position_type", Val: replyUserPosition.Type},
		},
		OrderBy:  []string{"id"},
		Cols:     []string{"id", "create_time", "content", "remark", "type"},
		PageSize: pageSize,
		Current:  currentPage,
	}, &replyPageByCond)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	var resData []LoggerValue
	loggerList := []logger.Logger{}
	err = json.Unmarshal([]byte(replyPageByCond.Json), &loggerList)
	for _, value := range loggerList {
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
			CreateTime: utils.Datetime(value.CreateTime, "2006-01-02 15:04:05"),
			Content:    value.Content,
			Remark:     value.Remark,
			Type:       retype,
		})
	}

	c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: LoggerList{
		Total:       replyPageByCond.Total,
		Current:     currentPage,
		CurrentSize: replyPageByCond.CurrentSize,
		OrderBy:     replyPageByCond.OrderBy,
		PageSize:    pageSize,
		List:        resData,
	}}
	c.ServeJSON()
	return
}
