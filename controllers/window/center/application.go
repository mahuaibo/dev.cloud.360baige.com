package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/center"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/application"
	"time"
	"dev.model.360baige.com/action"
	"encoding/json"
	"dev.cloud.360baige.com/utils"
)

// APPLICATION API
type ApplicationController struct {
	beego.Controller
}

// @Title 应用列表接口
// @Description 应用列表接口
// @Success 200 {"code":200,"message":"获取应用列表成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   current     query   string true       "当前页"
// @Param   pageSize     query   string true       "每页数量"
// @Param   name     query   string true       "搜索名称"
// @Failure 400 {"code":400,"message":"获取应用信息失败"}
// @router /list [post]
func (c *ApplicationController) List() {
	type data ApplicationListResponse
	accessToken := c.GetString("accessToken")
	appName := c.GetString("name")
	pageSize, _ := c.GetInt64("pageSize", 50)
	current, _ := c.GetInt64("current", 1)

	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}
	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "access_token", Val: accessToken },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}
	if replyUserPosition.UserId == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "获取信息失败"}
		c.ServeJSON()
		return
	}

	var reply action.PageByCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "PageByCond", action.PageByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId },
			action.CondValue{Type: "And", Key: "user_id", Val: replyUserPosition.UserId },
			action.CondValue{Type: "And", Key: "user_position_id", Val: replyUserPosition.Id },
			action.CondValue{Type: "And", Key: "user_position_type", Val: replyUserPosition.Type },
			action.CondValue{Type: "And", Key: "name__icontains", Val: appName},
		},
		Cols:     []string{"id", "end_time", "name", "image", "status", "application_tpl_id" },
		OrderBy:  []string{"id"},
		PageSize: pageSize,
		Current:  current,
	}, &reply)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取应用信息失败"}
		c.ServeJSON()
		return
	}

	replyList := []application.Application{}
	err = json.Unmarshal([]byte(reply.Json), &replyList)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "获取应用信息失败"}
		c.ServeJSON()
		return
	}

	var ids []int64
	for _, value := range replyList {
		ids = append(ids, value.ApplicationTplId)
	}
	var applicationTplReply []application.ApplicationTpl
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "ListByCond", action.ListByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "id__in", Val: ids },
		},
		Cols:     []string{"id", "name", "image", "status", "site"},
		OrderBy:  []string{"id"},
		PageSize: -1,
	}, &applicationTplReply)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取应用失败"}
		c.ServeJSON()
		return
	}

	applicationTpls := make(map[int64]application.ApplicationTpl)
	for _, value := range applicationTplReply {
		applicationTpls[value.Id] = value
	}
	var al []ApplicationValue
	for _, value := range replyList {
		var rename, reimage string
		if value.Name == "" && applicationTpls[value.ApplicationTplId].Name != "" {
			rename = applicationTpls[value.ApplicationTplId].Name
		} else {
			rename = value.Name
		}
		if value.Image == "" && applicationTpls[value.ApplicationTplId].Image != "" {
			reimage = applicationTpls[value.ApplicationTplId].Image
		} else {
			reimage = value.Image
		}

		al = append(al, ApplicationValue{
			Id:      value.Id,
			EndTime: utils.Datetime(value.EndTime, "2006-01-02 15:04:05"),
			Name:    rename,
			Image:   reimage,
			Status:  value.Status,
			Site:    applicationTpls[value.ApplicationTplId].Site,
		})
	}

	c.Data["json"] = data{Code: Normal, Message: "获取应用成功", Data: ApplicationList{
		Total:       reply.Total,
		Current:     current,
		CurrentSize: reply.CurrentSize,
		OrderBy:     reply.OrderBy,
		PageSize:    pageSize,
		Name:        appName,
		List:        al,
	}}
	c.ServeJSON()
	return
}

// @Title 应用修改状态接口
// @Description 应用修改状态接口
// @Success 200 {"code":200,"message":"获取应用修改状态成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Param   status     query   string true       " 0 启用 1 停用 2 退订"
// @Failure 400 {"code":400,"message":"获取应用修改状态失败"}
// @router /modifyStatus [post]
func (c *ApplicationController) ModifyStatus() {
	type data ModifyApplicationStatusResponse
	accessToken := c.GetString("accessToken")
	ap_id, _ := c.GetInt64("id")
	status, _ := c.GetInt("status")
	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌不能为空"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "access_token", Val: accessToken },
		},
		Fileds: []string{"user_id"},
	}, &replyUserPosition)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	var reply application.Application
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindById", &application.Application{
		Id: ap_id,
	}, &reply)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取应用信息失败"}
		c.ServeJSON()
		return
	}

	var updateArgs []action.UpdateValue
	updateArgs = append(updateArgs, action.UpdateValue{
		Key: "update_time",
		Val: time.Now().UnixNano() / 1e6,
	}, action.UpdateValue{
		Key: "status",
		Val: status,
	})
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "UpdateById", &action.UpdateByIdCond{
		Id:         []int64{ap_id},
		UpdateList: updateArgs,
	}, nil)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "应用信息修改失败"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = data{Code: Normal, Message: "应用信息修改成功"}
	c.ServeJSON()
	return
}
