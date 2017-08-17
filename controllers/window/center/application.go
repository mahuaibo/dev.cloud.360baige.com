package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/center"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/application"
	. "dev.model.360baige.com/models/company"
	"time"
	"dev.model.360baige.com/action"
	"encoding/json"
	"fmt"
)

// APPLICATION API
type ApplicationController struct {
	beego.Controller
}

// @Title 应用列表接口
// @Description 应用列表接口
// @Success 200 {"code":200,"message":"获取应用列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   current     query   string true       "当前页"
// @Param   page_size     query   string true       "每页数量"
// @Param   name     query   string true       "搜索名称"
// @Failure 400 {"code":400,"message":"获取应用信息失败"}
// @router /list [get]
func (c *ApplicationController) List() {
	res := ApplicationListResponse{}
	access_token := c.GetString("access_token")
	appname := c.GetString("name")
	pageSize, _ := c.GetInt64("page_size")
	currentPage, _ := c.GetInt64("current")

	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	//检测 accessToken
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type:  "And",
		Key: "access_token",
		Val:  access_token,
	})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var positionReply UserPosition
	fmt.Println("args", args)
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &positionReply)
	fmt.Println("positionReply", positionReply)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	//company_id、user_id、user_position_id、user_position_type
	com_id := positionReply.CompanyId
	user_id := positionReply.UserId
	user_position_id := positionReply.Id
	user_position_type := positionReply.Type
	if com_id == 0 || user_id == 0 || user_position_id == 0 {
		res.Code = ResponseSystemErr
		res.Message = "获取信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var applicationArge action.PageByCond
	applicationArge.CondList = append(applicationArge.CondList, action.CondValue{
		Type:  "And",
		Key: "company_id",
		Val:  com_id,
	}, action.CondValue{
		Type:  "And",
		Key: "user_id",
		Val:  user_id,
	}, action.CondValue{
		Type:  "And",
		Key: "user_position_id",
		Val:  user_position_id,
	}, action.CondValue{
		Type:  "And",
		Key: "user_position_type",
		Val:  user_position_type,
	})
	if appname != "" {
		applicationArge.CondList = append(applicationArge.CondList, action.CondValue{
			Type:  "And",
			Key: "name__icontains",
			Val:  appname,
		})
	}
	applicationArge.Cols = []string{"id", "create_time", "name", "image", "status", "application_tpl_id" }
	applicationArge.OrderBy = []string{"id"}
	applicationArge.PageSize = pageSize
	applicationArge.Current = currentPage
	var reply action.PageByCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "PageByCond", applicationArge, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "获取应用信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	replyList := []Application{}
	err = json.Unmarshal([]byte(reply.Json), &replyList)
	//获取应用其他信息tpl
	var idarg []int64
	idmap := make(map[int64]int64)
	for _, value := range replyList {
		idmap[value.ApplicationTplId] = value.ApplicationTplId
	}
	for _, value := range idmap {
		idarg = append(idarg, value)
	}

	var applicationTplArgs action.ListByCond
	applicationTplArgs.CondList = append(applicationTplArgs.CondList, action.CondValue{
		Type:  "And",
		Key: "id__in",
		Val:  idarg,
	})
	applicationTplArgs.Cols = []string{"id", "name", "image", "status", "site"}
	applicationTplArgs.OrderBy = []string{"id"}
	applicationTplArgs.PageSize = -1
	var applicationTplReply []ApplicationTpl
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "ListByCond", applicationTplArgs, &applicationTplReply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "获取应用失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	applicationTplByIds := make(map[int64]ApplicationTpl)
	for _, value := range applicationTplReply {
		applicationTplByIds[value.Id] = value
	}
	for _, value := range replyList {
		var rename, reimage string
		if value.Name == "" && applicationTplByIds[value.ApplicationTplId].Name != "" {
			rename = applicationTplByIds[value.ApplicationTplId].Name
		} else {
			rename = value.Name
		}
		if value.Image == "" && applicationTplByIds[value.ApplicationTplId].Image != "" {
			reimage = applicationTplByIds[value.ApplicationTplId].Image
		} else {
			reimage = value.Image
		}
		res.Data.List = append(res.Data.List, ApplicationValue{
			Id:         value.Id,
			CreateTime: time.Unix(value.CreateTime / 1000, 0).Format("2006-01-02"),
			Name:       rename,
			Image:      reimage,
			Status:     value.Status,
			Site:       applicationTplByIds[value.ApplicationTplId].Site,
		})
	}

	res.Code = ResponseNormal
	res.Message = "获取应用成功"
	res.Data.Total = reply.Total
	res.Data.Current = currentPage
	res.Data.CurrentSize = reply.CurrentSize
	res.Data.OrderBy = reply.OrderBy
	res.Data.PageSize = pageSize
	res.Data.Name = appname
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 应用详情接口
// @Description 应用详情接口
// @Success 200 {"code":200,"message":"获取应用详情成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"获取应用详情失败"}
// @router /detail [get]
func (c *ApplicationController) Detail() {
	res := ApplicationDetailResponse{}
	access_token := c.GetString("access_token")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Message = "访问令牌无效"
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
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	//company_id、user_id、user_position_id、user_position_type
	ap_id, _ := c.GetInt64("id")

	if ap_id == 0 {
		res.Code = ResponseSystemErr
		res.Message = "获取信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	var reply Application
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindById", &Application{
		Id: ap_id,
	}, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "获取应用信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if reply.ApplicationTplId == 0 {
		res.Code = ResponseSystemErr
		res.Message = "获取应用信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	//获取应用其他信息tpl
	var replyApplicationTpl ApplicationTpl
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindById", &ApplicationTpl{
		Id: reply.ApplicationTplId,
	}, &replyApplicationTpl)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "获取应用信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	re := time.Unix(reply.CreateTime / 1000, 0).Format("2006-01-02")
	var rename, reimage string
	if reply.Name == "" && replyApplicationTpl.Name != "" {
		rename = replyApplicationTpl.Name
	} else {
		rename = reply.Name
	}
	if reply.Image == "" && replyApplicationTpl.Image != "" {
		reimage = replyApplicationTpl.Image
	} else {
		reimage = reply.Image
	}
	//开发者
	var replyUser User
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", &User{
		Id: replyApplicationTpl.UserId,
	}, &replyUser)
	var username, cname string
	if err == nil {
		username = replyUser.Username
	}
	//开发公司
	var replyCompany Company
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "FindById", &Company{
		Id: replyApplicationTpl.CompanyId,
	}, &replyCompany)
	if err == nil {
		cname = replyCompany.Name
	}

	res.Code = ResponseNormal
	res.Message = "获取应用信息成功"
	res.Data.CreateTime = re
	res.Data.Name = rename
	res.Data.Image = reimage
	res.Data.Desc = replyApplicationTpl.Desc
	res.Data.Price = replyApplicationTpl.Price
	res.Data.PayType = GetPayTypeName(replyApplicationTpl.PayType)
	res.Data.PayCycle = GetPayCycleName(replyApplicationTpl.PayCycle)
	res.Data.CompanyName = cname
	res.Data.UserName = username
	c.Data["json"] = res
	c.ServeJSON()
}
func GetPayTypeName(ptype int8) string {
	// 0:限免 1:永久免费 2:1次性收费 3:周期收费tpl
	var rPtype string
	switch ptype {
	case 0:
		rPtype = "限免"
	case 1:
		rPtype = "永久免费"
	case 2:
		rPtype = "1次性收费"
	case 3:
		rPtype = "周期收费"
	}
	return rPtype
}

func GetPayCycleName(ptype int8) string {
	// 0无1月2季3半年4年tpl
	var rPtype string
	switch ptype {
	case 0:
		rPtype = "无"
	case 1:
		rPtype = "月"
	case 2:
		rPtype = "季"
	case 3:
		rPtype = "半年"
	case 4:
		rPtype = "年"
	}
	return rPtype
}

// @Title 应用修改状态接口
// @Description 应用修改状态接口
// @Success 200 {"code":200,"message":"获取应用修改状态成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Param   status     query   string true       " 0 启用 1 停用 2 退订"
// @Failure 400 {"code":400,"message":"获取应用修改状态失败"}
// @router /modifystatus [get]
func (c *ApplicationController) ModifyStatus() {
	res := ModifyApplicationStatusResponse{}
	access_token := c.GetString("access_token")
	ap_id, _ := c.GetInt64("id")
	status, _ := c.GetInt8("status")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Message = "访问令牌无效"
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
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var reply Application
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindById", &Application{
		Id: ap_id,
	}, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "获取应用信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var updateArgs []action.UpdateValue
	updateArgs = append(updateArgs, action.UpdateValue{
		Key: "update_time",
		Val:  time.Now().UnixNano() / 1e6,
	}, action.UpdateValue{
		Key: "status",
		Val:  status,
	})
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "UpdateById", &action.UpdateByIdCond{
		Id: []int64{ap_id},
		UpdateList:updateArgs,
	}, nil)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "应用信息修改失败！"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	res.Code = ResponseNormal
	res.Message = "应用信息修改成功！"
	c.Data["json"] = res
	c.ServeJSON()
}
