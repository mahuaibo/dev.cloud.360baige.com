package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	//. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/models/paginator"
)

type UserPositionController struct {
	beego.Controller
}

// @Title 获取用户身份接口
// @Description 获取用户身份接口
// @Success 200 {"code":200,"messgae":"获取用户身份成功","data":{"accessToken":"ok","expire_in":0}}
// @Param access_ticket     query   string true       "访问票据"
// @Param expire_in     query   string true       "访问时效"
// @Failure 400 {"code":400,"message":"获取用户身份失败"}
// @router /list [get]
func (c *UserPositionController) List() {
	res := Response{}
	//access_ticket := c.GetString("access_ticket")
	//expire_in := c.GetString("expire_in")

	//1 校验 access_ticket 获取 user_id

	//2 通过 user_id 获取 user_position 登录权限
	var reply Paginator
	args := &Paginator{
		PageSize:  0,
		Current:   0,
		MarkID:    0,
		Direction: 0,
		Filters:   "",
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "List", args, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "信息查询失败"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "信息查询成功"
		res.Data = reply
		c.Data["json"] = res
		c.ServeJSON()
	}
}
