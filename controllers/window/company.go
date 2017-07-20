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
	}
	com_id := replyAccessToken.CompanyId
	if com_id == 0 {
		res.Code = ResponseSystemErr
		res.Messgae = "获取公司信息失败"
		c.Data["json"] = res
		c.ServeJSON()
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
	}
	res.Code = ResponseNormal
	res.Messgae = "获取公司信息成功"
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
	} else {
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
		} else {
			com_id := replyAccessToken.CompanyId
			if com_id == 0 {
				res.Code = ResponseSystemErr
				res.Messgae = "获取公司信息失败"
				c.Data["json"] = res
				c.ServeJSON()
			} else {
				var reply Company
				err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "FindById", &Company{
					Id: com_id,
				}, &reply)

				if err != nil {
					res.Code = ResponseSystemErr
					res.Messgae = "获取公司信息失败"
					c.Data["json"] = res
					c.ServeJSON()
				}

				if (reply.Status != 1) {
					res.Code = ResponseSystemErr
					res.Messgae = "公司状态不可修改"
					c.Data["json"] = res
					c.ServeJSON()
				} else {
					logo := c.GetString("logo")
					name := c.GetString("name")
					shortName := c.GetString("short_name")
					provinceId, _ := c.GetInt64("province_id")
					cityId, _ := c.GetInt64("city_id")
					districtId, _ := c.GetInt64("district_id")
					address := c.GetString("address")
					positionX, _ := c.GetFloat("position_x", 64)
					positionY, _ := c.GetFloat("position_y", 64)
					remark := c.GetString("remark")
					brief := c.GetString("brief")

					timestamp := time.Now().UnixNano() / 1e6
					var updateArgs []action.UpdateValue
					updateArgs = append(updateArgs, action.UpdateValue{
						Key: "update_time",
						Val: timestamp,
					})
					updateArgs = append(updateArgs, action.UpdateValue{
						Key: "logo",
						Val: logo,
					})
					updateArgs = append(updateArgs, action.UpdateValue{
						Key: "name",
						Val: name,
					})
					updateArgs = append(updateArgs, action.UpdateValue{
						Key: "short_name",
						Val: shortName,
					})
					updateArgs = append(updateArgs, action.UpdateValue{
						Key: "province_id",
						Val: provinceId,
					})
					updateArgs = append(updateArgs, action.UpdateValue{
						Key: "city_id",
						Val: cityId,
					})
					updateArgs = append(updateArgs, action.UpdateValue{
						Key: "district_id",
						Val: districtId,
					})
					updateArgs = append(updateArgs, action.UpdateValue{
						Key: "address",
						Val: address,
					})
					updateArgs = append(updateArgs, action.UpdateValue{
						Key: "position_x",
						Val: positionX,
					})
					updateArgs = append(updateArgs, action.UpdateValue{
						Key: "position_y",
						Val: positionY,
					})
					updateArgs = append(updateArgs, action.UpdateValue{
						Key: "remark",
						Val: remark,
					})
					updateArgs = append(updateArgs, action.UpdateValue{
						Key: "brief",
						Val: brief,
					})

					err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "UpdateById", &action.UpdateByIdCond{
						Id:         []int64{com_id},
						UpdateList: updateArgs,
					}, nil)

					if err != nil {
						res.Code = ResponseSystemErr
						res.Messgae = "企业信息修改失败！"
						c.Data["json"] = res
						c.ServeJSON()
					}

					res.Code = ResponseNormal
					res.Messgae = "企业信息修改成功！"
					res.Data.Logo = logo
					res.Data.Name = name
					res.Data.ShortName = shortName
					res.Data.ProvinceId = provinceId
					res.Data.CityId = cityId
					res.Data.DistrictId = districtId
					res.Data.Address = address
					res.Data.PositionX = positionX
					res.Data.PositionY = positionY
					res.Data.Remark = remark
					res.Data.Brief = brief
					res.Data.Status = reply.Status
					c.Data["json"] = res
					c.ServeJSON()
				}

			}

		}

	}

}
