package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.cloud.360baige.com/utils"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/response"
	"fmt"
)

// USER API
type UserController struct {
	beego.Controller
}

// @Title 用户登录接口
// @Description 用户登录接口
// @Success 200 {"code":200,"messgae":"登录成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   username     query   string true       "用户名：百鸽账号、邮箱、手机号码 三种登录方式"
// @Param   password query   string true       "密码"
// @Failure 400 {"code":400,"message":"登录失败"}
// @router /login [post]
func (c *UserController) Login() {
	res := Response{}
	username := c.GetString("username")
	password := c.GetString("password")

	//1 检测数据完整性
	if (username == "" || password == "") {
		res.Code = ResponseSystemErr
		res.Messgae = "登录失败"
		c.Data["json"] = res
		c.ServeJSON()
	}

	//2 判断username 类型 属于 百鸽账号、邮箱、手机号码中的哪一种？
	username_type,_ := utils.DetermineStringType(username)
	fmt.Println("Type1:",username_type)
	var reply User
	var err error

	switch username_type {
	case 1:
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByUsername", &User{
			Username: username,
		}, &reply)
	case 2:
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByEmail", &User{
			Email: username,
		}, &reply)
	case 3:
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByPhone", &User{
			Phone: username,
		}, &reply)
	}

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "登录失败"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "登录成功"
		replyData := User{
			Username:     reply.Username,
			AccessTicket: reply.AccessTicket,
			ExpireIn:     reply.ExpireIn,
		}
		res.Data = replyData
		c.Data["json"] = res
		c.ServeJSON()
	}

}

// @Title 用户信息接口
// @Description 用户信息接口
// @Success 200 {"code":200,"messgae":"获取用户信息成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   user_id     query   string true       "用户ID"
// @Failure 400 {"code":400,"message":"获取用户信息失败"}
// @router /detail [get]
func (c *UserController) Detail() {
	res := Response{}
	user_id, _ := c.GetInt64("user_id", 0)
	if user_id == 0{
		res.Code = ResponseSystemErr
		res.Messgae = "获取用户信息失败"
		c.Data["json"] = res
		c.ServeJSON()
	}
	var reply User
	var err error

	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", &User{
		Id: user_id,
	}, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取用户信息失败"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "获取用户信息成功"
		replyData := User{
			Id:         reply.Id,
			CreateTime: reply.CreateTime,
			UpdateTime: reply.UpdateTime,
			Username:   reply.Username,
			Email:      reply.Email,
			Phone:      reply.Phone,
			Status:     reply.Status,
		}
		res.Data = replyData
		c.Data["json"] = res
		c.ServeJSON()
	}
}

// @Title 用户退出接口
// @Description 用户退出接口
// @Success 200 {"code":200,"messgae":"用户退出成功","data":{"access_ticket":"xxxx","expire_in":0}}// @Param   password query   string true       "密码"
// @Param   access_token query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户退出失败"}
// @router /login [post]
func (c *UserController) Logout() {

}

// @Title 用户修改密码接口
// @Description 用户修改密码接口
// @Success 200 {"code":200,"messgae":"用户修改密码成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户修改密码失败"}
// @router /login [post]
func (c *UserController) ModifyPassword() {

}


// @Title 用户信息修改接口
// @Description 用户信息修改接口
// @Success 200 {"code":200,"messgae":"用户修改密码成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户修改密码失败"}
// @router /login [post]
func (c *UserController) Modify() {

}
