package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.model.360baige.com/models"
	"fmt"
	"time"
)

type UserController struct {
	beego.Controller
}

type Response struct {
	Code    string        `json:"code"`
	Messgae string        `json:"messgae"`
	Data    interface{}   `json:"data,omitempty"`
}

// @Title 登录接口
// @Description 登录接口
// @Success 200 {"code":200,"messgae":"ok","data":{"username":"zhangs","accessToken":"ok"}}
// @Param   Username     query   string true       "用户名"
// @Param   Password     query   string true       "密码"
// @Failure 400 {"code":400,"message":"..."}
// @router /login [post]
func (c *UserController) Login() {
	username := c.GetString("Username")
	password := c.GetString("Password")
	response := Response{}
	var reply models.User
	args := &models.User{
		Username: username,
	}
	err := client.Call("http://127.0.0.1:2379", "User", "GetOneByUsername", args, &reply)
	if err != nil {
		// 返回错误信息
		response.Code = "500"
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}
	if reply.Password != password {
		// 返回错误信息
		response.Code = "400"
		response.Messgae = "用户名或密码错误"
		c.Data["json"] = response
		c.ServeJSON()
	}

	resUser := models.User{
		Id:          reply.Id,
		Username:    reply.Username,
		AccessToken: reply.AccessToken,
		ExpireIn:    reply.ExpireIn,
	}
	response.Code = "200"
	response.Messgae = "登录成功"
	response.Data = resUser
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title 退出接口
// @Description 退出接口
// @Success 200 {"code":200,"messgae":"ok"}
// @Param   Username     query   string true       "用户名"
// @Param   AccessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /logout [post]
func (c *UserController) Logout() {
	var reply models.User
	args := &models.User{
		Username:    c.GetString("Username"),
		AccessToken: c.GetString("AccessToken"),
	}
	err := client.Call("http://127.0.0.1:2379", "User", "Logout", args, &reply)
	fmt.Println(reply, err)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
	c.ServeJSON()
}

// @Title 用户信息详情
// @Description 用户信息详情
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   Id     query   string true       "用户ID"
// @Failure 400 {"code":400,"message":"..."}
// @router /detail [post]
func (c *UserController) Detail() {
	id, _ := c.GetInt64("Id")
	var reply models.User
	response := Response{}
	args := &models.User{
		Id: id,
	}
	err := client.Call("http://127.0.0.1:2379", "User", "Detail", args, &reply)

	if err != nil {
		response.Code = "500"
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}

	response.Code = "200"
	response.Messgae = "获取用户信息成功"
	reply.AccessToken = nil
	reply.Password = nil
	response.Data = reply
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title 用户信息修改
// @Description 用户信息修改
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   Id     query   int64 true       "用户ID"
// @Param   Phone     query   string false       "手机号码"
// @Param   Email     query   string false       "邮箱"
// @Failure 400 {"code":400,"message":"..."}
// @router /modify [post]
func (c *UserController) Modify() {

	id, _ := c.GetInt64("Id")
	phone := c.GetString("Phone")
	email := c.GetString("Email")

	var reply models.User
	response := Response{}
	args := &models.User{
		Id: id,
	}
	err := client.Call("http://127.0.0.1:2379", "User", "Detail", args, &reply)

	if err != nil {
		response.Code = "500"
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}
	reply.Id = id
	reply.Phone = phone
	reply.Email = email

	err = client.Call("http://127.0.0.1:2379", "User", "Modify", reply, nil)

	if err != nil {
		response.Code = "500"
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}

	response.Code = "200"
	response.Messgae = "用户信息修改成功！"
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title 用户密码修改
// @Description 用户密码修改
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   Username     query   string true       "用户名称"
// @Param   Password     query   string true       "密码"
// @Param   NewPassword     query   string true    "新密码"
// @Failure 400 {"code":400,"message":"..."}
// @router /modifyPassword [post]
func (c *UserController) ModifyPassword() {
	username := c.Ctx.Input.Param("Username")
	password := c.Ctx.Input.Param("Password")
	newPassword := c.Ctx.Input.Param("NewPassword")
	response := Response{}

	if password == newPassword {
		response.Code = "500"
		response.Messgae = "新密码不能与原密码相同"
		c.Data["json"] = response
		c.ServeJSON()
	}

	args := &models.User{
		Username: username,
	}
	var loginReply models.User

	err := client.Call("http://127.0.0.1:2379", "User", "GetOneByUsername", args, &loginReply)

	if err != nil {
		response.Code = "500"
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}

	if loginReply.Password != password {
		response.Code = "500"
		response.Messgae = "原密码错误"
		c.Data["json"] = response
		c.ServeJSON()
	}

	loginReply.Password = newPassword
	timestamp := time.Now().Unix()
	loginReply.UpdateTime = timestamp

	err = client.Call("http://127.0.0.1:2379", "User", "ModifyPassword", args, nil)

	if err != nil {
		response.Code = "500"
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}

	response.Code = "200"
	response.Messgae = "密码修改成功"
	c.Data["json"] = response
	c.ServeJSON()
}
