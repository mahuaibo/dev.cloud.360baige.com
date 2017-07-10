package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/response"
)

// USER API
type UserController struct {
	beego.Controller
}

// @Title 登录接口
// @Description 登录接口
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
	username_type := 1
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
		res.Data = reply
		c.Data["json"] = res
		c.ServeJSON()
	}

}

// @Title 登录接口
// @Description 登录接口
// @Success 200 {"code":200,"messgae":"登录成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   username     query   string true       "用户名：百鸽账号、邮箱、手机号码 三种登录方式"
// @Param   password query   string true       "密码"
// @Failure 400 {"code":400,"message":"登录失败"}
// @router /login [post]
func (c *UserController) GetUserPosition() {
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
	username_type := 1
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
		res.Data = reply
		c.Data["json"] = res
		c.ServeJSON()
	}

}
