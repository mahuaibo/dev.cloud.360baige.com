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
	"fmt"
	"dev.cloud.360baige.com/utils"
)

// APPLICATIONTPL API
type ApplicationTplController struct {
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
func (c *ApplicationTplController) List() {
	type data ApplicationTplListResponse
	accessToken := c.GetString("accessToken")
	currentPage, _ := c.GetInt64("current", 1)
	pageSize, _ := c.GetInt64("pageSize", 20)
	appName := c.GetString("name", "")
	Type, _ := c.GetInt8("type")

	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "accessToken", Val: accessToken },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	if replyUserPosition.UserId == 0 {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取信息失败"}
		c.ServeJSON()
		return
	}

	var condValue action.CondValue
	if Type != 0 {
		condValue = action.CondValue{Type: "And", Key: "type", Val: Type }
	}
	var reply action.PageByCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "PageByCond", action.PageByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "name__icontains", Val: appName },
			action.CondValue{Type: "And", Key: "status__gt", Val: -1 },
			condValue,
		},
		Cols:     []string{"id", "name", "image", "status", "desc", "price", "pay_cycle"},
		OrderBy:  []string{"id"},
		PageSize: -1,
		Current:  currentPage,
	}, &reply)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取应用信息失败"}
		c.ServeJSON()
		return
	}

	var replyApplication []application.Application
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "ListByCond", action.ListByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId },
			action.CondValue{Type: "And", Key: "user_id", Val: replyUserPosition.UserId },
			action.CondValue{Type: "And", Key: "user_position_id", Val: replyUserPosition.Id },
			action.CondValue{Type: "And", Key: "user_position_type", Val: replyUserPosition.Type },
			action.CondValue{Type: "And", Key: "status__in", Val: []string{"0", "1"} },
		},
		Cols:     []string{"id", "application_tpl_id" },
		OrderBy:  []string{"id"},
		PageSize: -1,
	}, &replyApplication)

	applicationList := make(map[int64]int64)
	for _, value := range replyApplication {
		applicationList[value.ApplicationTplId] = value.ApplicationTplId
	}

	var resData []ApplicationTplValue
	replyList := []application.ApplicationTpl{}
	err = json.Unmarshal([]byte(reply.Json), &replyList)

	fmt.Println("replyList", replyList)
	for _, value := range replyList {
		var restatus int8
		if applicationList[value.Id] > 0 {
			restatus = 1
		} else {
			restatus = 0
		}
		resData = append(resData, ApplicationTplValue{
			Id:                 value.Id,
			Name:               value.Name,
			Image:              value.Image,
			Subscription:       value.Subscription,
			SubscriptionStatus: restatus,
			Desc:               value.Desc,
			Price:              value.Price,
			PayCycle:           value.PayCycle,
		})
	}

	c.Data["json"] = data{Code: Normal, Message: "获取应用成功", Data: ApplicationTplList{
		Total:       reply.Total,
		Current:     currentPage,
		CurrentSize: reply.CurrentSize,
		OrderBy:     reply.OrderBy,
		PageSize:    pageSize,
		Name:        appName,
		List:        resData,
	}}
	c.ServeJSON()
	return
}

// @Title 应用详情接口
// @Description 应用详情接口
// @Success 200 {"code":200,"message":"获取应用详情成功"}
// @Failure 400 {"code":400,"message":"获取应用详情失败"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   applicationTplId     query   string true       "applicationTplId"
// @router /detail [post]
func (c *ApplicationTplController) Detail() {
	type data ApplicationTalDetailResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	applicationTplId, _ := c.GetInt64("applicationTplId")

	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{CondList: []action.CondValue{action.CondValue{Type: "And", Key: "access_token", Val: accessToken }, action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTimestamp }, }, Fileds: []string{"id", "user_id", "company_id", "type"}, }, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	if replyUserPosition.UserId == 0 {
		c.Data["json"] = data{Code: Normal, Message: "获取用户信息失败"}
		c.ServeJSON()
		return
	}

	var replyApplicationTpl application.ApplicationTpl
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindById", &application.ApplicationTpl{
		Id: applicationTplId,
	}, &replyApplicationTpl)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取应用信息失败"}
		c.ServeJSON()
		return
	}
	var replyApplication application.Application
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindByCond", &action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "ApplicationTplId", Val: applicationTplId},
			action.CondValue{Type: "And", Key: "CompanyId", Val: replyUserPosition.CompanyId},
			action.CondValue{Type: "And", Key: "UserId", Val: replyUserPosition.UserId},
			action.CondValue{Type: "And", Key: "UserPositionId", Val: replyUserPosition.Id},
			action.CondValue{Type: "And", Key: "UserPositionType", Val: replyUserPosition.Type},
		},
	}, &replyApplication)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取应用信息失败"}
		c.ServeJSON()
		return
	}
	var subscriptionStatus int8 = 0
	if replyApplication.Id != 0 {
		subscriptionStatus = replyApplication.Status
	}

	c.Data["json"] = data{Code: Normal, Message: "获取应用成功", Data: ApplicationTalDetail{
		Id:                 replyApplicationTpl.Id,
		Name:               replyApplicationTpl.Name,
		Image:              replyApplicationTpl.Image,
		Desc:               replyApplicationTpl.Desc,
		PriceDesc:          "该应用功能￥72.00/月，您可以根据自己的需求选择是否订购使用。",
		Price:              replyApplicationTpl.Price,
		PayType:            replyApplicationTpl.PayType,
		PayCycle:           GetPayCycleName(replyApplicationTpl.PayCycle),
		SubscriptionStatus: subscriptionStatus,
		StartTime:          utils.Datetime(replyApplication.StartTime, "2016-01-02 03:04:05"),
		EndTime:            utils.Datetime(replyApplication.EndTime, "2016-01-02 03:04:05"),
	}}
	c.ServeJSON()
	return
}

// @Title 应用订阅接口
// @Description 应用订阅接口
// @Success 200 {"code":200,"message":"获取应用详情成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"获取应用详情失败"}
// @router /subscription [get]
func (c *ApplicationTplController) Subscription() {
	type data ModifyApplicationTplStatusResponse
	accessToken := c.GetString("accessToken")
	applicationTplId, _ := c.GetInt64("id")
	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}
	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "accessToken", Val: accessToken },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	if replyUserPosition.UserId == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "获取应用信息失败"}
		c.ServeJSON()
		return
	}

	var replyApplication application.Application
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId },
			action.CondValue{Type: "And", Key: "user_id", Val: replyUserPosition.UserId },
			action.CondValue{Type: "And", Key: "user_position_id", Val: replyUserPosition.Id },
			action.CondValue{Type: "And", Key: "user_position_type", Val: replyUserPosition.Type },
			action.CondValue{Type: "And", Key: "application_tpl_id", Val: applicationTplId },
		},
		Fileds: []string{"id", "application_tpl_id" },
	}, &replyApplication)
	if err == nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "此应用已经订阅过"}
		c.ServeJSON()
		return
	}

	var reply application.ApplicationTpl
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindById", application.ApplicationTpl{
		Id: applicationTplId,
	}, &reply)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取应用信息失败"}
		c.ServeJSON()
		return
	}
	if reply.Status == 0 {
		c.Data["json"] = data{Code: ErrorSystem, Message: "此应用已经下架"}
		c.ServeJSON()
		return
	}

	var replyApplication2 application.Application
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "Add", application.Application{
		CreateTime:       time.Now().UnixNano() / 1e6,
		UpdateTime:       time.Now().UnixNano() / 1e6,
		CompanyId:        replyUserPosition.CompanyId,
		UserId:           replyUserPosition.UserId,
		UserPositionId:   replyUserPosition.Id,
		UserPositionType: replyUserPosition.Type,
		ApplicationTplId: reply.Id,
		Name:             reply.Name,
		Image:            reply.Image,
		Status:           1,
		StartTime:        1,
		EndTime:          1,
	}, &replyApplication2)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "应用订阅失败"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "应用订阅成功", Data: ApplicationTplStatus{
		ApplicationTplId: reply.Id,
		AppId:            replyApplication2.Id,
	}}
	c.ServeJSON()
	return
}

// @Title 应用修改状态接口
// @Description 应用修改状态接口
// @Success 200 {"code":200,"message":"获取应用修改状态成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Param   status     query   string true      " 0 上架 1 下架 "
// @Failure 400 {"code":400,"message":"获取应用修改状态失败"}
// @router /modifystatus [get]
func (c *ApplicationTplController) ModifyStatus() {
	type data ModifyApplicationStatusResponse
	accessToken := c.GetString("accessToken")
	status, _ := c.GetInt8("status")
	applicationTplId, _ := c.GetInt64("id")
	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	var err error
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByAccessToken", &user.UserPosition{
		AccessToken: accessToken,
	}, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyApplicationTpl application.ApplicationTpl
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindById", &application.ApplicationTpl{
		Id: applicationTplId,
	}, &replyApplicationTpl)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取应用信息失败"}
		c.ServeJSON()
		return
	}

	replyApplicationTpl.UpdateTime = time.Now().UnixNano() / 1e6
	replyApplicationTpl.Status = status
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "UpdateById", replyApplicationTpl, nil)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "应用信息修改失败"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "应用信息修改成功"}
	c.ServeJSON()
}
