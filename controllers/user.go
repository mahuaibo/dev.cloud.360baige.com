package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.model.360baige.com/models"
	"fmt"
)

type UserController struct {
	beego.Controller
}

// @Title 注册账号接口
// @Description 注册账号接口
// @Param	page		path 	int	true		"分页页码"
// @Param	rows		path 	int	true		"展示条数"
// @Success 200 {object} models.logger
// @Failure 403 :uid is empty
// @router /register [post]
func (c *UserController) Register() {
	var reply models.User
	args := &models.User{
		Username: c.GetString("Username"),
		Password: c.GetString("Password"),
	}
	err := client.Call("http://127.0.0.1:2379", "Admin", "Register", args, &reply)
	fmt.Println(reply, err)
	if err == nil {
		// TODO 注册成功添加角色
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
	c.ServeJSON()
}

// @Title 登录接口
// @Description 登录接口
// @Success 200 {"code":200,"messgae":"ok","data":{"username":"zhangs","accessToken":"ok"}}
// @Param   Username     query   string true       "用户名"
// @Param   Password     query   string true       "密码"
// @Failure 400 {"code":400,"message":"..."}
// @router /login [post]
func (c *UserController) Login() {
	var reply models.User
	args := &models.User{
		Username: c.GetString("Username"),
		Password: c.GetString("Password"),
	}
	err := client.Call("http://127.0.0.1:2379", "User", "Login", args, &reply)
	fmt.Println(reply, err)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
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
		Username: c.GetString("Username"),
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


// @Title 用户注销接口
// @Description 用户注销接口
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   Username     query   string true       "用户ID"
// @Param   AccessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /cancel [post]
func (c *UserController) Cancel() {
	var reply models.User
	args := &models.User{
		Username: c.GetString("Username"),
		AccessToken: c.GetString("AccessToken"),
	}
	err := plugin.Call("http://127.0.0.1:2379", "Admin", "Cancel", args, &reply)
	fmt.Println(reply, err)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
	c.ServeJSON()
}

// @Title 用户删除接口
// @Description 用户删除接口
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   Id     query   string true       "用户ID"
// @Failure 400 {"code":400,"message":"..."}
// @router /delete/:id [get]
func (c *UserController) Delete() {
	id, _ := c.GetInt64("Id")
	var reply models.User
	args := &models.User{
		Id:id,
	}
	err := plugin.Call("http://127.0.0.1:2379", "Admin", "Delete", args, &reply)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
	c.ServeJSON()
}

// @Title 账号激活接口
// @Description 账号激活接口
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   Username     query   string true       "用户名称"
// @Failure 400 {"code":400,"message":"..."}
// @router /activation/:username [get]
func (c *UserController) Activation() {
	var reply models.User
	args := &models.User{
		Username:c.GetString("Username"),
	}
	err := plugin.Call("http://127.0.0.1:2379", "Admin", "Activation", args, &reply)
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
// @Param   Code     query   string true       "验证码"
// @Param   Username     query   string true       "用户名称"
// @Param   Password     query   string true       "用户密码"
// @Failure 400 {"code":400,"message":"..."}
// @router /findPassword [post]
func (c *UserController) FindPassword() {
	var reply models.User
	args := &models.User{
		Username:c.GetString("Username"),
		Password:c.GetString("password"),
		Code:c.GetString("Code"),
		CodeTime:time.Now().UnixNano() / 1e6,
	}
	err := plugin.Call("http://127.0.0.1:2379", "Admin", "FindPassword", args, &reply)
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
	args := &models.User{
		Id: id,
	}
	err := client.Call("http://127.0.0.1:2379", "User", "Detail", args, &reply)
	fmt.Println(reply, err)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
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
	var reply models.User
	args := &models.User{
		Id: id,
		Phone:c.GetString("Phone"),
		Email:c.GetString("Email"),
	}
	err := client.Call("http://127.0.0.1:2379", "User", "Modify", args, &reply)
	fmt.Println(reply, err)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
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
	if password == newPassword {
		c.Data["json"] = "新密码不能与旧密码重复！"
	}
	var reply, loginReply models.Admin
	args := &models.User{
		Id:loginReply.Id,
		Username: username,
		Password: password,
		Password: newPassword,
	}
	e := plugin.Call("http://127.0.0.1:2379", "Admin", "Login", args, &loginReply)
	if e != nil {
		c.Data["json"] = "密码错误！"
		c.ServeJSON()
	}
	err := plugin.Call("http://127.0.0.1:2379", "Admin", "ModifyPassword", args, &reply)
	fmt.Println(reply, err)
	if err == nil {
		c.Data["json"] = reply
	} else {
		c.Data["json"] = err
	}
	c.ServeJSON()
}