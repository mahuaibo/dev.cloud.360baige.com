package website

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.model.360baige.com/models/company"
	. "dev.model.360baige.com/http/window/website"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/utils"
	"dev.model.360baige.com/models/website"
)

// Company API
type MaterialController struct {
	beego.Controller
}

// @Title 企业信息修改接口
// @Description 企业信息修改接口
// @Success 200 {"code":200,"message":"企业信息修改成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"企业信息修改失败"}
// @router /add [post]
func (c *MaterialController) Add() {
	type data MaterialAddResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	menuId, _ := c.GetInt64("menuId")
	Type, _ := c.GetInt("type")
	name := c.GetString("name")
	photo := c.GetString("photo")
	url := c.GetString("url")
	desc := c.GetString("desc")
	status, _ := c.GetInt("status")

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

	var replyMaterial website.Material
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Material", "Add", website.Material{
		CreateTime:currentTimestamp,
		UpdateTime:currentTimestamp,
		MenuId:menuId,
		Type:Type,
		Name:name,
		Photo:photo,
		Url:url,
		Desc:desc,
		Status:status,
	}, &replyMaterial)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if replyMaterial.Id == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "资源新增失败"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "资源新增成功"}
	c.ServeJSON()
	return
}

// @Title 企业信息修改接口
// @Description 企业信息修改接口
// @Success 200 {"code":200,"message":"企业信息修改成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"企业信息修改失败"}
// @router /modify [post]
func (c *MaterialController) Modify() {
	type data MaterialModifyResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	id, _ := c.GetInt64("id")
	menuId, _ := c.GetInt64("menuId")
	Type := c.GetString("type")
	name, _ := c.GetInt64("name")
	photo := c.GetString("photo")
	url, _ := c.GetInt64("url")
	desc := c.GetString("desc")
	status := c.GetString("status")

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
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Material", "UpdateById", action.UpdateByIdCond{
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "update_time", Val: currentTimestamp },
			action.UpdateValue{Key: "menu_id", Val: menuId},
			action.UpdateValue{Key: "type", Val: Type},
			action.UpdateValue{Key: "name", Val: name},
			action.UpdateValue{Key: "photo", Val: photo},
			action.UpdateValue{Key: "url", Val: url, },
			action.UpdateValue{Key: "desc", Val: desc },
			action.UpdateValue{Key: "status", Val: status },
		},
		Id: []int64{id},
	}, &replyNum)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if replyNum.Value == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "资源修改失败"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "资源修改成功"}
	c.ServeJSON()
	return
}

// @Title 企业信息修改接口
// @Description 企业信息修改接口
// @Success 200 {"code":200,"message":"企业信息修改成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"企业信息修改失败"}
// @router /delete [post]
func (c *MaterialController) Delete() {
	type data MaterialModifyResponse
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
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Material", "UpdateById", action.UpdateByIdCond{
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "update_time", Val: currentTimestamp },
			action.UpdateValue{Key: "status", Val: -1 },
		},
		Id: []int64{id},
	}, &replyNum)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if replyNum.Value == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "资源删除失败"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "资源删除成功"}
	c.ServeJSON()
	return
}


// @Title 企业信息接口
// @Description 企业信息接口
// @Success 200 {"code":200,"message":"获取企业信息成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取企业信息失败"}
// @router /detail [post]
func (c *MaterialController) Detail() {
	type data MaterialDetailResponse
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

	var replyMaterial website.Material
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Material", "FindById", &company.Company{
		Id: id,
	}, &replyMaterial)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常[资源获取失败]"}
		c.ServeJSON()
		return
	}

	photoUrl := utils.SignURLSample(replyMaterial.Photo, 3600)
	c.Data["json"] = data{Code: Normal, Message: "资源获取成功",
		Data: MaterialDetail{
			Id:replyMaterial.Id,
			Photo:photoUrl,
			Name:replyMaterial.Name,
			Url:replyMaterial.Url,
			Status:replyMaterial.Status,
			Desc:replyMaterial.Desc,
		},
	}
	c.ServeJSON()
	return
}