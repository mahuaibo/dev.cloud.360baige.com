package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.model.360baige.com/models/company"
	. "dev.model.360baige.com/http/window/center"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/utils"
)

// Company API
type CompanyController struct {
	beego.Controller
}

// @Title 企业信息接口
// @Description 企业信息接口
// @Success 200 {"code":200,"message":"获取企业信息成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取企业信息失败"}
// @router /detail [post]
func (c *CompanyController) Detail() {
	type data CompanyDetailResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")

	err := utils.Unable(map[string]string{"accessToken": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: Message(40000, err.Error())}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{CondList: []action.CondValue{action.CondValue{Type: "And", Key: "access_token", Val: accessToken }, action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTimestamp }, }, Fileds: []string{"id", "user_id", "company_id", "type"}, }, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: Message(50000)}
		c.ServeJSON()
		return
	}

	if replyUserPosition.UserId == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: Message(30000)}
		c.ServeJSON()
		return
	}

	var replyCompany company.Company
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "FindById", &company.Company{
		Id: replyUserPosition.CompanyId,
	}, &replyCompany)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常[公司信息失败]"}
		c.ServeJSON()
		return
	}
	logoUrl := utils.SignURLSample(replyCompany.Logo)
	c.Data["json"] = data{Code: Normal, Message: "获取公司信息成功",
		Data: CompanyDetail{
			Id:         replyCompany.Id,
			Logo:       logoUrl,
			Name:       replyCompany.Name,
			ProvinceId: replyCompany.ProvinceId,
			CityId:     replyCompany.CityId,
			DistrictId: replyCompany.DistrictId,
			Address:    replyCompany.Address,
			PositionX:  replyCompany.PositionX,
			PositionY:  replyCompany.PositionY,
			Remark:     replyCompany.Remark,
			Brief:      replyCompany.Brief,
			Status:     replyCompany.Status,
		},
	}
	c.ServeJSON()
	return
}

// @Title 企业信息修改接口
// @Description 企业信息修改接口
// @Success 200 {"code":200,"message":"企业信息修改成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"企业信息修改失败"}
// @router /modify [post]
func (c *CompanyController) Modify() {
	type data CompanyModifyResponse
	currentTimestamp := utils.CurrentTimestamp()
	name := c.GetString("name")
	short_name := c.GetString("short_name")
	provinceId, _ := c.GetInt64("province_id")
	cityId, _ := c.GetInt64("city_id")
	districtId, _ := c.GetInt64("district_id")
	address := c.GetString("address")
	positionX, _ := c.GetFloat("position_x", 64)
	positionY, _ := c.GetFloat("position_y", 64)
	remark := c.GetString("remark")
	brief := c.GetString("brief")
	accessToken := c.GetString("accessToken")

	err := utils.Unable(map[string]string{"accessToken": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: Message(40000, err.Error())}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{CondList: []action.CondValue{action.CondValue{Type: "And", Key: "access_token", Val: accessToken }, action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTimestamp }, }, Fileds: []string{"id", "user_id", "company_id", "type"}, }, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: Message(50000)}
		c.ServeJSON()
		return
	}

	if replyUserPosition.UserId == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: Message(30000)}
		c.ServeJSON()
		return
	}

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "UpdateById", action.UpdateByIdCond{
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "update_time", Val: currentTimestamp },
			//action.UpdateValue{Key: "logo", Val: logo},
			action.UpdateValue{Key: "name", Val: name},
			action.UpdateValue{Key: "short_name", Val: short_name},
			action.UpdateValue{Key: "province_id", Val: provinceId},
			action.UpdateValue{Key: "city_id", Val: cityId},
			action.UpdateValue{Key: "district_id", Val: districtId, },
			action.UpdateValue{Key: "address", Val: address },
			action.UpdateValue{Key: "position_x", Val: positionX },
			action.UpdateValue{Key: "position_y", Val: positionY },
			action.UpdateValue{Key: "remark", Val: remark },
			action.UpdateValue{Key: "brief", Val: brief },
		},
		Id: []int64{replyUserPosition.CompanyId},
	}, &replyNum)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: Message(50000)}
		c.ServeJSON()
		return
	}
	if replyNum.Value == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: Message(40000)}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: Message(20000)}
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
func (c *CompanyController) UploadLogo() {
	currentTimestamp := utils.CurrentTimestamp()
	requestType := c.Ctx.Request.Method
	if requestType == "POST" {
		type data UploadLogoResponse
		accessToken := c.GetString("accessToken")
		companyId, _ := c.GetInt64("id")
		_, handle, _ := c.Ctx.Request.FormFile("uploadFile")

		err := utils.Unable(map[string]string{"accessToken": "string:true"}, c.Ctx.Input)
		if err != nil {
			c.Data["json"] = data{Code: ErrorLogic, Message: Message(40000, err.Error())}
			c.ServeJSON()
			return
		}

		var replyUserPosition user.UserPosition
		err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{CondList: []action.CondValue{action.CondValue{Type: "And", Key: "access_token", Val: accessToken }, action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTimestamp }, }, Fileds: []string{"id", "user_id", "company_id", "type"}, }, &replyUserPosition)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: Message(50000)}
			c.ServeJSON()
			return
		}

		if replyUserPosition.UserId == 0 {
			c.Data["json"] = data{Code: ErrorLogic, Message: Message(30000)}
			c.ServeJSON()
			return
		}

		objectKey, err := utils.UploadImage(handle, "Company/logoImages/")
		if err != nil {
			c.Data["json"] = data{Code: ErrorLogic, Message: Message(50000)}
			c.ServeJSON()
			return
		}

		var replyCompany company.Company
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "FindByCond", action.FindByCond{
			CondList: []action.CondValue{action.CondValue{Type: "And", Key: "id", Val: companyId },
			},
			Fileds: []string{"id"},
		}, &replyCompany)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: Message(50000)}
			c.ServeJSON()
			return
		}

		var replyNum action.Num
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "UpdateById", action.UpdateByIdCond{
			Id: []int64{replyCompany.Id},
			UpdateList: []action.UpdateValue{
				action.UpdateValue{Key: "update_time", Val: currentTimestamp},
				action.UpdateValue{Key: "logo", Val: objectKey},
			},
		}, &replyNum)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: Message(50000)}
			c.ServeJSON()
			return
		}
		logoUrl := utils.SignURLSample(objectKey)
		c.Data["json"] = data{Code: Normal, Data: logoUrl, Message: Message(20000)}
		c.ServeJSON()
		return
	}
}
