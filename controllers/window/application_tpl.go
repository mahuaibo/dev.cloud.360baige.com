package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/application"
	. "dev.model.360baige.com/models/company"
	"time"
	"dev.model.360baige.com/action"
	"encoding/json"
	"fmt"
)

// APPLICATIONTPL API
type ApplicationTplController struct {
	beego.Controller
}

// @Title 应用列表接口
// @Description 应用列表接口
// @Success 200 {"code":200,"messgae":"获取应用列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   current     query   string true       "当前页"
// @Param   page_size     query   string true       "每页数量"
// @Param   name     query   string true       "搜索名称"
// @Failure 400 {"code":400,"message":"获取应用信息失败"}
// @router /list [get]
func (c *ApplicationTplController) List() {
	res := ApplicationTplListResponse{}
	access_token := c.GetString("access_token")
	currentPage, _ := c.GetInt64("current")
	pageSize, _ := c.GetInt64("page_size")
	appname := c.GetString("name")
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
		Type: "And",
		Key:  "accessToken",
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

	var appTplArgs action.PageByCond
	if appname != "" {
		appTplArgs.CondList = append(appTplArgs.CondList, action.CondValue{
			Type: "And",
			Key:  "name__icontains",
			Val:  appname,
		})
	}
	if Type != 0 {
		appTplArgs.CondList = append(appTplArgs.CondList, action.CondValue{
			Type: "And",
			Key:  "type",
			Val:  Type,
		})
	}
	appTplArgs.CondList = append(appTplArgs.CondList, action.CondValue{
		Type: "And",
		Key:  "status__gt",
		Val:  -1,
	})
	appTplArgs.Cols = []string{"id", "name", "image", "status", "desc"}
	appTplArgs.OrderBy = []string{"id"}
	appTplArgs.PageSize = -1 // 目前显示全部
	appTplArgs.Current = currentPage
	var reply action.PageByCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "PageByCond", appTplArgs, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取应用信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var appArgs action.ListByCond
	appArgs.CondList = append(appArgs.CondList, action.CondValue{
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
	}, action.CondValue{
		Type: "And",
		Key:  "status__in",
		Val:  []string{"0", "1"},
	})
	appArgs.Cols = []string{"id", "application_tpl_id" }
	appArgs.OrderBy = []string{"id"}
	appArgs.PageSize = -1
	var replyApplication []Application
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "ListByCond", appArgs, &replyApplication)

	idmap := make(map[int64]int64)
	for _, value := range replyApplication {
		idmap[value.ApplicationTplId] = value.ApplicationTplId
	}

	var resData []ApplicationTplValue
	replyList := []ApplicationTpl{}
	err = json.Unmarshal([]byte(reply.Json), &replyList)
	//循环赋值
	for _, value := range replyList {
		var restatus int8
		if idmap[value.Id] > 0 {
			restatus = 1
		} else {
			restatus = 0
		}
		resData = append(resData, ApplicationTplValue{
			Id:                 value.Id,
			Name:               value.Name,
			Image:              value.Image,
			SubscriptionStatus: restatus,
			Desc:               value.Desc,
		})
	}

	res.Code = ResponseNormal
	res.Messgae = "获取应用成功"
	res.Data.Total = reply.Total
	res.Data.Current = currentPage
	res.Data.CurrentSize = reply.CurrentSize
	res.Data.OrderBy = reply.OrderBy
	res.Data.PageSize = pageSize
	res.Data.Name = appname
	res.Data.List = resData
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 应用详情接口
// @Description 应用详情接口
// @Success 200 {"code":200,"messgae":"获取应用详情成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"获取应用详情失败"}
// @router /detail [get]
func (c *ApplicationTplController) Detail() {
	res := ApplicationDetailResponse{}
	access_token := c.GetString("access_token")
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
		Type: "And",
		Key:  "accessToken",
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

	var appTplArgs ApplicationTpl // 获取应用信息tpl
	appTplArgs.Id, _ = c.GetInt64("id")
	if appTplArgs.Id == 0 {
		res.Code = ResponseSystemErr
		res.Messgae = "获取信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	var replyApplicationTpl ApplicationTpl
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindById", appTplArgs, &replyApplicationTpl)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取应用信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var userArgs User // 开发者
	userArgs.Id = replyApplicationTpl.UserId
	var replyUser User
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", userArgs, &replyUser)
	if err == nil {
		res.Data.UserName = replyUser.Username
	}

	var companyArgs Company // 开发公司
	companyArgs.Id = replyApplicationTpl.CompanyId
	var replyCompany Company
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "FindById", companyArgs, &replyCompany)
	if err == nil {
		res.Data.CompanyName = replyCompany.Name
	}

	res.Code = ResponseNormal
	res.Messgae = "获取应用成功"
	res.Data.CreateTime = time.Unix(replyApplicationTpl.CreateTime / 1000, 0).Format("2006-01-02")
	res.Data.Name = replyApplicationTpl.Name
	res.Data.Image = replyApplicationTpl.Image
	res.Data.Desc = replyApplicationTpl.Desc
	res.Data.Price = replyApplicationTpl.Price
	res.Data.PayType = GetPayTypeName(replyApplicationTpl.PayType)
	res.Data.PayCycle = GetPayCycleName(replyApplicationTpl.PayCycle)
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 应用订阅接口
// @Description 应用订阅接口
// @Success 200 {"code":200,"messgae":"获取应用详情成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"获取应用详情失败"}
// @router /subscription [get]
func (c *ApplicationTplController) Subscription() {
	res := ModifyApplicationTplStatusResponse{}
	access_token := c.GetString("access_token")
	ap_id, _ := c.GetInt64("id")
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
		Type: "And",
		Key:  "accessToken",
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
	if ap_id == 0 || com_id == 0 || user_id == 0 || user_position_id == 0 {
		res.Code = ResponseSystemErr
		res.Messgae = "获取应用信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	//判断此应用是否订阅过
	var appArgs action.FindByCond
	appArgs.CondList = append(appArgs.CondList, action.CondValue{
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
	}, action.CondValue{
		Type: "And",
		Key:  "application_tpl_id",
		Val:  ap_id,
	})
	appArgs.Fileds = []string{"id", "application_tpl_id" }
	var replyApplication Application
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindByCond", appArgs, &replyApplication)
	if err == nil {
		res.Code = ResponseSystemErr
		res.Messgae = "此应用已经订阅过"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var appTplArgs ApplicationTpl
	appTplArgs.Id = ap_id
	var reply ApplicationTpl
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindById", appTplArgs, &reply)
	fmt.Println("reply>>>", reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取应用信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if reply.Status == 0 {
		res.Code = ResponseSystemErr
		res.Messgae = "此应用已经下架"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var addReply Application
	addReply.CreateTime = time.Now().UnixNano() / 1e6
	addReply.UpdateTime = time.Now().UnixNano() / 1e6
	addReply.CompanyId = com_id
	addReply.UserId = user_id
	addReply.UserPositionId = user_position_id
	addReply.UserPositionType = user_position_type
	addReply.ApplicationTplId = reply.Id // 应用ID
	addReply.Name = reply.Name           // 名称
	addReply.Image = reply.Image         // 图片链接
	addReply.Status = 1                  // 状态
	addReply.StartTime = 1               // 开始时间
	addReply.EndTime = 1                 // 结束时间

	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "Add", addReply, &addReply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "应用订阅失败！"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.Code = ResponseNormal
	res.Messgae = "应用订阅成功！"
	res.Data.ApplicationTplId = reply.Id
	res.Data.AppId = addReply.Id
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 应用修改状态接口
// @Description 应用修改状态接口
// @Success 200 {"code":200,"messgae":"获取应用修改状态成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Param   status     query   string true      " 0 上架 1 下架 "
// @Failure 400 {"code":400,"message":"获取应用修改状态失败"}
// @router /modifystatus [get]
func (c *ApplicationTplController) ModifyStatus() {
	res := ModifyApplicationStatusResponse{}
	access_token := c.GetString("access_token")
	status, _ := c.GetInt8("status")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	//检测 accessToken
	var replyAccessToken UserPosition
	var err error
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByAccessToken", &UserPosition{
		AccessToken: access_token,
	}, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	ap_id, _ := c.GetInt64("id")
	var reply ApplicationTpl
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindById", &ApplicationTpl{
		Id: ap_id,
	}, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取应用信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	reply.UpdateTime = time.Now().UnixNano() / 1e6
	reply.Status = status
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "UpdateById", reply, nil)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "应用信息修改失败！"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.Code = ResponseNormal
	res.Messgae = "应用信息修改成功！"
	c.Data["json"] = res
	c.ServeJSON()
}
