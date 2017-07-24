package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/models/company"
	//. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/http/window"
	. "dev.model.360baige.com/models/user"
	"time"
	"dev.model.360baige.com/action"
)

// COMPANY API
type CompanyController struct {
	beego.Controller
}

// @Title 企业信息接口
// @Description 企业信息接口
// @Success 200 {"code":200,"messgae":"获取企业信息成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取企业信息失败"}
// @router /detail [get]
func (c *CompanyController) Detail() {
	res := CompanyDetailResponse{}
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
	com_id := replyAccessToken.CompanyId
	if com_id == 0 {
		res.Code = ResponseSystemErr
		res.Messgae = "获取公司信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	var reply Company
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "FindById", &Company{
		Id: com_id,
	}, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取公司信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	res.Code = ResponseNormal
	res.Messgae = "获取公司信息成功"
	res.Data.Id = reply.Id
	res.Data.Logo = reply.Logo
	res.Data.Name = reply.Name
	res.Data.ShortName = reply.ShortName
	res.Data.ProvinceId = reply.ProvinceId
	res.Data.CityId = reply.CityId
	res.Data.DistrictId = reply.DistrictId
	res.Data.Address = reply.Address
	res.Data.PositionX = reply.PositionX
	res.Data.PositionY = reply.PositionY
	res.Data.Remark = reply.Remark
	res.Data.Brief = reply.Brief
	res.Data.Status = reply.Status
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 企业信息修改接口
// @Description 企业信息修改接口
// @Success 200 {"code":200,"messgae":"企业信息修改成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"企业信息修改失败"}
// @router /modify [post]
func (c *CompanyController) Modify() {
	res := CompanyModifyResponse{}
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
	com_id := replyAccessToken.CompanyId
	if com_id == 0 {
		res.Code = ResponseSystemErr
		res.Messgae = "获取公司信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	var reply Company
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "FindById", &Company{
		Id: com_id,
	}, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取公司信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if (reply.Status != 1) {
		res.Code = ResponseSystemErr
		res.Messgae = "公司状态不可修改"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	provinceId, _ := c.GetInt64("province_id")
	cityId, _ := c.GetInt64("city_id")
	districtId, _ := c.GetInt64("district_id")
	positionX, _ := c.GetFloat("position_x", 64)
	positionY, _ := c.GetFloat("position_y", 64)
	var updateArgs action.UpdateByIdCond
	updateArgs.UpdateList = append(updateArgs.UpdateList, action.UpdateValue{
		Key: "update_time",
		Val: time.Now().UnixNano() / 1e6,
	}, action.UpdateValue{
		Key: "logo",
		Val: c.GetString("logo"),
	}, action.UpdateValue{
		Key: "name",
		Val: c.GetString("name"),
	}, action.UpdateValue{
		Key: "short_name",
		Val: c.GetString("short_name"),
	}, action.UpdateValue{
		Key: "province_id",
		Val: provinceId,
	}, action.UpdateValue{
		Key: "city_id",
		Val: cityId,
	}, action.UpdateValue{
		Key: "district_id",
		Val: districtId,
	}, action.UpdateValue{
		Key: "address",
		Val: c.GetString("address"),
	}, action.UpdateValue{
		Key: "position_x",
		Val: positionX,
	}, action.UpdateValue{
		Key: "position_y",
		Val: positionY,
	}, action.UpdateValue{
		Key: "remark",
		Val: c.GetString("remark"),
	}, action.UpdateValue{
		Key: "brief",
		Val: c.GetString("brief"),
	})
	updateArgs.Id = []int64{com_id}
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "UpdateById", updateArgs, &action.Num{})
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "企业信息修改失败！"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.Code = ResponseNormal
	res.Messgae = "企业信息修改成功！"
	c.Data["json"] = res
	c.ServeJSON()
}
