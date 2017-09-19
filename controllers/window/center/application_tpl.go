package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/center"
	"dev.model.360baige.com/models/application"
	"dev.model.360baige.com/action"
	"encoding/json"
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
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	currentPage, _ := c.GetInt64("current", 1)
	pageSize, _ := c.GetInt64("pageSize", 50)
	appName := c.GetString("name", "")
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

	var condValue action.CondValue
	if Type != 0 {
		condValue = action.CondValue{Type: "And", Key: "type", Val: Type }
	}
	var reply action.PageByCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "PageByCond", action.PageByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "name__icontains", Val: appName },
			action.CondValue{Type: "And", Key: "user_position_type", Val: replyUserPosition.Type },
			action.CondValue{Type: "And", Key: "status__gt", Val: -1 },
			condValue,
		},
		Cols:     []string{"id", "name", "image", "status", "desc", "price", "pay_cycle", "subscription"},
		OrderBy:  []string{"id"},
		PageSize: -1,
		Current:  currentPage,
	}, &reply)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
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
	for _, value := range replyList {
		var restatus int
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

	c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: ApplicationTplList{
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
// @router /detail [get,post]
func (c *ApplicationTplController) Detail() {
	type data ApplicationTalDetailResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	applicationTplId, _ := c.GetInt64("applicationTplId")

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

	var replyApplicationTpl application.ApplicationTpl
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindById", &application.ApplicationTpl{
		Id: applicationTplId,
	}, &replyApplicationTpl)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	var replyApplication application.Application
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindByCond", &action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "CompanyId", Val: replyUserPosition.CompanyId},
			action.CondValue{Type: "And", Key: "UserId", Val: replyUserPosition.UserId},
			action.CondValue{Type: "And", Key: "UserPositionType", Val: replyUserPosition.Type},
			action.CondValue{Type: "And", Key: "UserPositionId", Val: replyUserPosition.Id},
			action.CondValue{Type: "And", Key: "ApplicationTplId", Val: applicationTplId},
		},
	}, &replyApplication)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	var subscriptionStatus int = 0
	if replyApplication.Id != 0 {
		subscriptionStatus = replyApplication.Status
	}

	c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: ApplicationTalDetail{
		Id:                 replyApplicationTpl.Id,
		Name:               replyApplicationTpl.Name,
		Image:              replyApplicationTpl.Image,
		Desc:               replyApplicationTpl.Desc,
		Price:              replyApplicationTpl.Price,
		PayType:            replyApplicationTpl.PayType,
		PayCycle:           GetPayCycleName(replyApplicationTpl.PayCycle),
		SubscriptionStatus: subscriptionStatus,
		EndTime:            utils.Datetime(replyApplication.EndTime, "2006-01-02 15:04:05"),
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
// @router /subscribe [post]
func (c *ApplicationTplController) Subscribe() {
	type data SubscribeResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	applicationTplId, _ := c.GetInt64("id")

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

	var replyApplication application.Application
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId },
			action.CondValue{Type: "And", Key: "user_id", Val: replyUserPosition.UserId },
			action.CondValue{Type: "And", Key: "user_position_id", Val: replyUserPosition.Id },
			action.CondValue{Type: "And", Key: "user_position_type", Val: replyUserPosition.Type },
			action.CondValue{Type: "And", Key: "application_tpl_id", Val: applicationTplId},
		},
		Fileds: []string{"id", "application_tpl_id", "status" },
	}, &replyApplication)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	var reply application.ApplicationTpl
	var replyApplication2 application.Application
	if replyApplication.Status == -1 {
		var replyNum action.Num
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "UpdateById", action.UpdateByIdCond{
			Id: []int64{replyApplication.Id},
			UpdateList: []action.UpdateValue{
				action.UpdateValue{Key: "update_time", Val: currentTimestamp},
				action.UpdateValue{Key: "status", Val: 0},
			},
		}, &replyNum)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}
		reply.Id = replyApplication.Id
		replyApplication2.Id = applicationTplId
	} else {
		err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindById", application.ApplicationTpl{
			Id: applicationTplId,
		}, &reply)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "Add", application.Application{
			CreateTime:       currentTimestamp,
			UpdateTime:       currentTimestamp,
			CompanyId:        replyUserPosition.CompanyId,
			UserId:           replyUserPosition.UserId,
			UserPositionId:   replyUserPosition.Id,
			UserPositionType: replyUserPosition.Type,
			ApplicationTplId: reply.Id,
			Name:             reply.Name,
			Image:            reply.Image,
			Status:           0,
			StartTime:        currentTimestamp,
			EndTime:          currentTimestamp,
		}, &replyApplication2)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}
	}
	c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: ApplicationTplStatus{
		ApplicationTplId: reply.Id,
		AppId:            replyApplication2.Id,
	}}
	c.ServeJSON()
	return
}

// @Title 应用退订接口
// @Description 应用订阅接口
// @Success 200 {"code":200,"message":"获取应用详情成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"获取应用详情失败"}
// @router /unSubscribe [post]
func (c *ApplicationTplController) UnSubscribe() {
	type data UnSubscribeResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	applicationTplId, _ := c.GetInt64("id")

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
	var replyApplication application.Application
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId },
			action.CondValue{Type: "And", Key: "user_id", Val: replyUserPosition.UserId },
			action.CondValue{Type: "And", Key: "user_position_id", Val: replyUserPosition.Id },
			action.CondValue{Type: "And", Key: "user_position_type", Val: replyUserPosition.Type },
			action.CondValue{Type: "And", Key: "application_tpl_id", Val: applicationTplId},
		},
		Fileds: []string{"id", "application_tpl_id", "status" },
	}, &replyApplication)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "UpdateById", action.UpdateByIdCond{
		Id: []int64{replyApplication.Id},
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "update_time", Val: currentTimestamp},
			action.UpdateValue{Key: "status", Val: -1},
		},
	}, &replyNum)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "SUCCESS"}
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
//func (c *ApplicationTplController) ModifyStatus() {
//	type data ModifyApplicationStatusResponse
//	currentTimestamp := utils.CurrentTimestamp()
//	accessToken := c.GetString("accessToken")
//	status, _ := c.GetInt("status")
//	applicationTplId, _ := c.GetInt64("id")
//
//	err := utils.Unable(map[string]string{"accessToken": "string:true"}, c.Ctx.Input)
//	if err != nil {
//		c.Data["json"] = data{Code: ErrorLogic, Message: Message(40000, err.Error())}
//		c.ServeJSON()
//		return
//	}
//
//	replyUserPosition, errCode := utils.UserPosition(accessToken, currentTimestamp)
//	if errCode != 0 {
//		c.Data["json"] = data{Code: ErrorPower, Message: Message(errCode)}
//		c.ServeJSON()
//		return
//	}
//
//	var replyApplicationTpl application.ApplicationTpl
//	//args *action.FindByCond, reply *application.ApplicationTpl
//	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindByCond", &action.FindByCond{
//		CondList: []action.CondValue{
//			action.CondValue{"And", "Id", applicationTplId},
//		},
//	}, &replyApplicationTpl)
//	if err != nil {
//		c.Data["json"] = data{Code: ErrorSystem, Message: Message(50000)}
//		c.ServeJSON()
//		return
//	}
//
//	replyApplicationTpl.UpdateTime = currentTimestamp
//	replyApplicationTpl.Status = status
//	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "UpdateById", replyApplicationTpl, nil)
//	if err != nil {
//		c.Data["json"] = data{Code: ErrorSystem, Message: Message(50000)}
//		c.ServeJSON()
//		return
//	}
//
//	c.Data["json"] = data{Code: Normal, Message: Message(20000)}
//	c.ServeJSON()
//}
