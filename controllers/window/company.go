package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/models/company"
	//. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/http/window"
	. "dev.model.360baige.com/models/user"
	"time"
	"fmt"
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
		var replyAccessToken UserPosition
		var err error
		err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByAccessToken", &UserPosition{
			AccessToken: access_token,
		}, &replyAccessToken)
		fmt.Println(err)
		fmt.Println(replyAccessToken)
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
					reply.UpdateTime = timestamp
					reply.Logo = logo
					reply.Name = name
					reply.ShortName = shortName
					reply.ProvinceId = provinceId
					reply.CityId = cityId
					reply.DistrictId = districtId
					reply.Address = address
					reply.PositionX = positionX
					reply.PositionY = positionY
					reply.Remark = remark
					reply.Brief = brief

					err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "UpdateById", reply, nil)

					if err != nil {
						res.Code = ResponseSystemErr
						res.Messgae = "企业信息修改失败！"
						c.Data["json"] = res
						c.ServeJSON()
					}

					res.Code = ResponseNormal
					res.Messgae = "企业信息修改成功！"
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

			}

		}

	}

}
