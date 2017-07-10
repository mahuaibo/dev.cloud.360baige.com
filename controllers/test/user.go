package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"time"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/response"
)

type UserController struct {
	beego.Controller
}

// @Title 登录接口
// @Description 登录接口
// @Success 200 {"code":200,"messgae":"ok","data":{"data":{... ...},"accessToken":"ok"}}
// @Param   username     query   string true       "用户名"
// @Param   password query   string true       "密码"
// @Param   datatype query   string true       "1 app账号 2 qq 3 新浪 4 支付宝 5 微信 "
// @Failure 400 {"code":400,"message":"..."}
// @router /login [post]
func (c *UserController) Login() {
	id, _ := c.GetInt64("id")
	res := Response{}
	var reply User
	args := &User{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", args, &reply)

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

// @Title 登录其他方式接口
// @Description 登录其他方式接口
// @Success 200 {"code":200,"messgae":"ok","data":{"data":{... ...},"accessToken":"ok"}}
// @Param   username     query   string true       "用户名"
// @Param   password query   string true       "密码"
// @Param   datatype query   string true       "2 qq 3 新浪 4 支付宝 5 微信 "
// @Failure 400 {"code":400,"message":"..."}
// @router /other-login [post]
func (c *UserController) OtherLogin() {
	id, _ := c.GetInt64("id")
	res := Response{}
	var reply User
	args := &User{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", args, &reply)

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
// @Title 退出接口
// @Description 退出接口
// @Success 200 {"code":200,"messgae":"ok"}
// @Param   username     query   string true       "用户名"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /logout [post]
func (c *UserController) LogOut() {
	id, _ := c.GetInt64("id")
	res := Response{}
	var reply User
	args := &User{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", args, &reply)

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
// @Title 验证码接口
// @Description 验证码接口
// @Success 200 {"code":200,"messgae":"ok"}
// @Param   username     query   string true       "用户名"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /send-message-code [post]
func (c *UserController) SendMessageCode() {
	id, _ := c.GetInt64("id")
	res := Response{}
	var reply User
	args := &User{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", args, &reply)

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
// @Title 注册接口
// @Description 注册接口
// @Success 200 {"code":200,"messgae":"ok"}
// @Param   username     query   string true       "用户名"
// @Param   password query   string true       "密码"
// @Param   code query   string true       "验证码"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /register [post]
func (c *UserController) Register() {
	timestamp := time.Now().UnixNano() / 1e6
	var (
		res   Response // http 返回体
		reply User
	)
	args := &User{
		CreateTime: timestamp,
		UpdateTime: timestamp,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "Add", args, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "新增失败"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "新增成功"
		res.Data = reply
		c.Data["json"] = res
		c.ServeJSON()
	}
}

// @Title 修改密码接口
// @Description 修改密码接口
// @Success 200 {"code":200,"messgae":"ok"}
// @Param   username     query   string true       "用户名"
// @Param   password query   string true       "密码"
// @Param   code query   string true       "验证码"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /modify-pwd [post]
func (c *UserController) ModifyPWD() {
	id, _ := c.GetInt64("id")

	var reply User
	res := Response{}
	args := &User{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", args, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = err.Error()
		c.Data["json"] = res
		c.ServeJSON()
	}
	timestamp := time.Now().UnixNano() / 1e6
	reply.Id = id

	reply.UpdateTime = timestamp

	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", reply, nil)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "信息修改失败！"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "信息修改成功！"
		c.Data["json"] = res
		c.ServeJSON()
	}
}

// @Title 个人信息接口
// @Description 个人信息接口
// @Success 200 {"code":200,"messgae":"ok","data":{"data":{ ... ... },"accessToken":"ok"}}
// @Param accessToken     query   string true       "访问令牌"
// @Param companyId     query   string true       "学校id"
// @Param userId query   string true       "userid"
// @Failure 400 {"code":400,"message":"..."}
// @router /personal-information [post]
func (c *UserController) PersonalInformation() {
	id, _ := c.GetInt64("id")
	res := Response{}
	var reply User
	args := &User{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", args, &reply)

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

// @Title 新增
// @Description 新增
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /add [post]
func (c *UserController) Add() {
	timestamp := time.Now().UnixNano() / 1e6
	var (
		res   Response // http 返回体
		reply User
	)
	args := &User{
		CreateTime: timestamp,
		UpdateTime: timestamp,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "Add", args, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "新增失败"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "新增成功"
		res.Data = reply
		c.Data["json"] = res
		c.ServeJSON()
	}
}

// @Title 信息
// @Description 信息
// @Success 200 {"code":200,"messgae":"信息查询成功", "data":{ ... ... }}
// @Param   id     query   string true       "ID"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /detail [get]
func (c *UserController) Detail() {
	id, _ := c.GetInt64("id")
	res := Response{}
	var reply User
	args := &User{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", args, &reply)

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

// @Title 信息修改
// @Description 信息修改
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   id     query   string true       "ID"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /modify [post]
func (c *UserController) Modify() {
	id, _ := c.GetInt64("id")

	var reply User
	res := Response{}
	args := &User{
		Id: id,
	}
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", args, &reply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = err.Error()
		c.Data["json"] = res
		c.ServeJSON()
	}
	timestamp := time.Now().UnixNano() / 1e6
	reply.Id = id

	reply.UpdateTime = timestamp

	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", reply, nil)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "信息修改失败！"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "信息修改成功！"
		c.Data["json"] = res
		c.ServeJSON()
	}
}
