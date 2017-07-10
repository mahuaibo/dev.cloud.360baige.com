package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/response"
)

// USER API
type UserPositionController struct {
	beego.Controller
}

// @Title 获取用户身份接口
// @Description 获取用户身份接口
// @Success 200 {"code":200,"messgae":"获取用户身份成功","data":{"accessToken":"ok","expire_in":0}}
// @Param access_ticket     query   string true       "访问票据"
// @Failure 400 {"code":400,"message":"获取用户身份失败"}
// @router /list [get]
func (c *UserPositionController) List() {
	res := Response{}
	access_ticket := c.GetString("access_ticket")

	var replyUser User
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByAccessTicket", &User{
		AccessTicket: access_ticket,
	}, &replyUser)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "访问票据无效"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		//判断时效是否超时 TODO
	}

	var replyUserPosition []UserPosition

	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "ListByUserId", &UserPosition{
		UserId: replyUser.Id,
	}, &replyUserPosition)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取用户身份失败"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "获取用户身份成功"
		//replyUserPositionData := UserPosition{
		//	Id:replyUserPosition.Id,
		//	CreateTime:replyUserPosition.CreateTime,
		//	UpdateTime:replyUserPosition.UpdateTime,
		//	Type:replyUserPosition.Type,
		//	CompanyId:replyUserPosition.CompanyId,
		//	UserId:replyUserPosition.UserId,
		//	PersonId:replyUserPosition.PersonId,
		//	AccessToken:replyUserPosition.AccessToken,
		//	ExpireIn:replyUserPosition.ExpireIn,
		//	Status:replyUserPosition.Status,
		//}
		res.Data = replyUserPosition
		c.Data["json"] = res
		c.ServeJSON()
	}
}

// @Title 获取一个登录用户身份权限接口
// @Description 获取一个登录用户身份权限接口
// @Success 200 {"code":200,"messgae":"获取一个登录用户身份权限成功","data":{"accessToken":"ok","expire_in":0}}
// @Param user_position_id     query   string true       "用户身份Id"
// @Failure 400 {"code":400,"message":"获取一个登录用户身份权限失败"}
// @router /detail [get]
func (c *UserPositionController) Detail() {
	res := Response{}
	user_position_id, _ := c.GetInt64("user_position_id", 0)

	var replyUserPosition UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindById", &UserPosition{
		Id: user_position_id,
	}, &replyUserPosition)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取用户身份失败"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "获取用户身份成功"
		res.Data = replyUserPosition
		c.Data["json"] = res
		c.ServeJSON()
	}
}
