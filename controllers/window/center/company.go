package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.model.360baige.com/models/company"
	. "dev.model.360baige.com/http/window/center"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/action"
	"time"
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
	accessToken := c.GetString("accessToken")
	if accessToken == "" {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌不能为空"}
		c.ServeJSON()
		return
	}

	//检测 accessToken
	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: [] action.CondValue{
			action.CondValue{Type: "And", Key: "accessToken", Val: accessToken },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)

	if err != nil {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "系统异常[访问令牌失效]"}
		c.ServeJSON()
		return
	}
	if replyUserPosition.Id == 0 {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	var replyCompany company.Company
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "FindById", &company.Company{
		Id: replyUserPosition.CompanyId,
	}, &replyCompany)

	if err != nil {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "系统异常[公司信息失败]"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: ResponseNormal, Message: "获取公司信息成功",
		Data: CompanyDetail{
			Logo:       replyCompany.Logo,
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
	logo := c.GetString("logo")
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
	if accessToken == "" {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌无效"}
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
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}
	if replyUserPosition.Id == 0 {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取公司信息失败"}
		c.ServeJSON()
		return
	}

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "UpdateById", action.UpdateByIdCond{
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "update_time", Val: time.Now().UnixNano() / 1e6 },
			action.UpdateValue{Key: "logo", Val: logo},
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
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "企业信息修改失败"}
		c.ServeJSON()
		return
	}
	if replyNum.Value == 0 {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "企业信息修改失败"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: ResponseNormal, Message: "企业信息修改成功"}
	c.ServeJSON()
	return
}
