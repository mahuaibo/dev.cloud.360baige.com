package website

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/website"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/utils"
	"dev.model.360baige.com/models/website"
)

// Company API
type MenuController struct {
	beego.Controller
}

// @Title 新增菜单接口
// @Description 新增菜单接口
// @Success 200 {"code":200,"message":"菜单新增成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"菜单新增失败"}
// @router /add [post]
func (c *MenuController) Add() {
	type data MenuAddResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	icon := c.GetString("icon")
	name := c.GetString("name")
	sn := c.GetString("sn")
	status, _ := c.GetInt("status")

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

	var replyMenu website.Menu
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Menu", "Add", &website.Menu{
		CreateTime:   currentTimestamp,
		UpdateTime:   currentTimestamp,
		Icon:     icon,
		CompanyId:     replyUserPosition.CompanyId,
		Name:        name,
		Sn:     sn,
		Status: status,
	}, &replyMenu)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if replyMenu.Id == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "菜单新增失败"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "菜单新增成功"}
	c.ServeJSON()
	return
}


// @Title 企业信息修改接口
// @Description 企业信息修改接口
// @Success 200 {"code":200,"message":"企业信息修改成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"企业信息修改失败"}
// @router /modify [post]
func (c *MenuController) Modify() {
	type data MenuModifyResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	id, _ := c.GetInt64("id")
	icon := c.GetString("icon")
	name := c.GetString("name")
	sn := c.GetString("sn")
	status := c.GetString("status")

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

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Menu", "UpdateById", action.UpdateByIdCond{
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "update_time", Val: currentTimestamp },
			action.UpdateValue{Key: "company_id", Val: replyUserPosition.CompanyId },
			action.UpdateValue{Key: "icon", Val: icon},
			action.UpdateValue{Key: "name", Val: name},
			action.UpdateValue{Key: "sn", Val: sn},
			action.UpdateValue{Key: "status", Val: status},
		},
		Id: []int64{id},
	}, &replyNum)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if replyNum.Value == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "菜单修改失败"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "菜单修改成功"}
	c.ServeJSON()
	return
}

// @Title 企业信息修改接口
// @Description 企业信息修改接口
// @Success 200 {"code":200,"message":"企业信息修改成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"企业信息修改失败"}
// @router /delete [post]
func (c *MenuController) Delete() {
	type data MenuDeleteResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	id, _ := c.GetInt64("id")

	err := utils.Unable(map[string]string{"accessToken": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	_, err = utils.UserPosition(accessToken, currentTimestamp)
	if err != nil {
		c.Data["json"] = data{Code: ErrorPower, Message: err.Error()}
		c.ServeJSON()
		return
	}

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Menu", "UpdateById", action.UpdateByIdCond{
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "update_time", Val: currentTimestamp },
			action.UpdateValue{Key: "status", Val: -1},
		},
		Id: []int64{id},
	}, &replyNum)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if replyNum.Value == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "菜单删除失败"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "菜单删除成功"}
	c.ServeJSON()
	return
}

// @Title 菜单信息接口
// @Description 菜单信息接口
// @Success 200 {"code":200,"message":"获取菜单信息成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取菜单信息失败"}
// @router /detail [post]
func (c *MenuController) Detail() {
	type data MenuDetailResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	menuId, _ := c.GetInt64("menuId")

	err := utils.Unable(map[string]string{"accessToken": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	_, err = utils.UserPosition(accessToken, currentTimestamp)
	if err != nil {
		c.Data["json"] = data{Code: ErrorPower, Message: err.Error()}
		c.ServeJSON()
		return
	}

	var replyMenu website.Menu
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Menu", "FindById", &website.Menu{
		Id: menuId,
	}, &replyMenu)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常[菜单信息获取失败]"}
		c.ServeJSON()
		return
	}

	iconUrl := utils.SignURLSample(replyMenu.Icon, 3600)
	c.Data["json"] = data{Code: Normal, Message: "获取菜单详情成功",
		Data: MenuDetail{
			Id:         replyMenu.Id,
			Icon:       iconUrl,
			Name:       replyMenu.Name,
			Status:     replyMenu.Status,
		},
	}
	c.ServeJSON()
	return
}

// @Title 用户注册
// @Description 用户注册
// @Success 200 {"code":200,"message":"企业logo上传失败"}
// @Param   accessToken     query   string true       "方文令牌"
// @Param   id     query   string true       "企业ID"
// @Param   uploadFile query   string true       "file"
// @Failure 400 {"code":400,"message":"企业logo上传成功"}
// @router /uploadLogo [options,post]
func (c *MenuController) UploadLogo() {
	currentTimestamp := utils.CurrentTimestamp()
	requestType := c.Ctx.Request.Method
	if requestType == "POST" {
		type data UploadIconResponse
		accessToken := c.GetString("accessToken")
		_, handle, _ := c.Ctx.Request.FormFile("uploadFile")

		err := utils.Unable(map[string]string{"accessToken": "string:true"}, c.Ctx.Input)
		if err != nil {
			c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
			c.ServeJSON()
			return
		}

		_, err = utils.UserPosition(accessToken, currentTimestamp)
		if err != nil {
			c.Data["json"] = data{Code: ErrorPower, Message: err.Error()}
			c.ServeJSON()
			return
		}

		objectKey, err := utils.UploadImage(handle, "Company/logoImages/")
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}

		c.Data["json"] = data{Code: Normal, Data: IconData{
			ObjectKet:objectKey,
			Icon:utils.SignURLSample(objectKey, 3600),
		}, Message: "菜单Icon上传成功"}
		c.ServeJSON()
		return
	}
}
