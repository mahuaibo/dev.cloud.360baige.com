package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.cloud.360baige.com/models/response"
	"dev.cloud.360baige.com/models/constant"
	"dev.model.360baige.com/models"
	"time"
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
	// 获取参数 和 准备参数
	username := c.GetString("Username")      // 用户名
	password := c.GetString("Password")      //密码
	var reply models.User                    // rpc 返回参数
	response := response.Response{}          //http 返回体
	timestamp := time.Now().UnixNano() / 1e6 //操作时间戳
	args := &models.User{//新增用户参数
		CreateTime:  timestamp,
		UpdateTime:  timestamp,
		Username:    username,
		Password:    password,
		Status:      constant.UserStatusValid,
		AccessToken: "",
		ExpireIn:    0,

	}
	//
	err := client.Call("http://127.0.0.1:2379", "User", "AddUser", args, &reply)

	if err != nil {
		response.Code = constant.ResponseSystemErr
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}

	if reply.Id == 0 {
		response.Code = constant.ResponseLogicErr
		response.Messgae = "用户新增失败"
		c.Data["json"] = response
		c.ServeJSON()
	}

	resUser := models.User{
		Id:          reply.Id,
		Username:    reply.Username,
		AccessToken: reply.AccessToken,
		ExpireIn:    reply.ExpireIn,
	}
	response.Code = constant.ResponseNormal
	response.Messgae = "登录成功"
	response.Data = resUser
	c.Data["json"] = response
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
	username := c.GetString("Username")
	password := c.GetString("Password")
	response := response.Response{}
	var reply models.User
	args := &models.User{
		Username: username,
		Password:password,
	}
	err := client.Call("http://127.0.0.1:2379", "User", "FindByUsername", args, &reply)
	fmt.Println("reply", reply)
	if err != nil {
		// 返回错误信息
		response.Code = constant.ResponseSystemErr
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}
	if reply.Password != password {
		// 返回错误信息
		response.Code = constant.ResponseLogicErr
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
	response.Code = constant.ResponseNormal
	response.Messgae = "登录成功"
	response.Data = resUser
	fmt.Println("response", response)
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title 退出接口
// @Description 退出接口
// @Success 200 {"code":200,"messgae":"ok"}
// @Param   Id     query   string true       "用户名"
// @Param   AccessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /logout [post]
func (c *UserController) Logout() {
	id, _ := c.GetInt64("Id")
	accessToken := c.GetString("AccessToken")
	var reply models.User
	response := response.Response{}
	args := &models.User{
		Id: id,
	}
	err := client.Call("http://127.0.0.1:2379", "User", "FindUserById", args, &reply)
	if err != nil {
		// 返回错误信息
		response.Code = constant.ResponseSystemErr
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}

	if reply.AccessToken != accessToken {
		response.Code = constant.ResponseLogicErr
		response.Messgae = "访问令牌无效"
		c.Data["json"] = response
		c.ServeJSON()
	}

	response.Code = constant.ResponseNormal
	response.Messgae = "退出成功"
	c.Data["json"] = response
	c.ServeJSON()
}

// @Title 用户注销接口
// @Description 用户注销接口
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   Id     query   string true       "用户ID"
// @Param   AccessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /cancel [post]
func (c *UserController) Cancel() {
	id, _ := c.GetInt64("Id")
	accessToken := c.GetString("AccessToken")
	var reply models.User
	response := response.Response{}
	args := &models.User{
		Id: id,
		AccessToken:accessToken,
	}
	err := client.Call("http://127.0.0.1:2379", "User", "FindUserById", args, &reply)
	if err != nil {
		// 返回错误信息
		response.Code = constant.ResponseSystemErr
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}

	c.ServeJSON()
}

// @Title 账号激活接口
// @Description 账号激活接口
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   Id     query   string true       "用户ID"
// @Param   AccessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /activation/:username [get]
func (c *UserController) Activation() {
	id, _ := c.GetInt64("Id")
	accessToken := c.GetString("AccessToken")
	var reply models.User
	response := response.Response{}
	args := &models.User{
		Id: id,
	}
	err := client.Call("http://127.0.0.1:2379", "User", "FindUserById", args, &reply)
	if err != nil {
		// 返回错误信息
		response.Code = constant.ResponseSystemErr
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}

	if reply.AccessToken != accessToken {
		response.Code = constant.ResponseLogicErr
		response.Messgae = "访问令牌无效"
		c.Data["json"] = response
		c.ServeJSON()
	}
	timestamp := time.Now().UnixNano() / 1e6
	reply.Status = constant.UserStatusValid
	reply.UpdateTime = timestamp
	err = client.Call("http://127.0.0.1:2379", "User", "UpdateUserById", reply, nil)

	if err != nil {
		// 返回错误信息
		response.Code = constant.ResponseSystemErr
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}

	response.Code = constant.ResponseNormal
	response.Messgae = "用户注销成功"
	c.Data["json"] = response
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
	response := response.Response{}
	args := &models.User{
		Id: id,
	}
	err := client.Call("http://127.0.0.1:2379", "User", "FindUserById", args, &reply)

	if err != nil {
		response.Code = constant.ResponseSystemErr
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}

	response.Code = constant.ResponseNormal
	response.Messgae = "获取用户信息成功"
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
	response := response.Response{}
	args := &models.User{
		Id: id,
	}
	err := client.Call("http://127.0.0.1:2379", "User", "FindUserById", args, &reply)

	if err != nil {
		response.Code = constant.ResponseSystemErr
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}
	reply.Id = id
	reply.Phone = phone
	reply.Email = email

	err = client.Call("http://127.0.0.1:2379", "User", "UpdateUserById", reply, nil)

	if err != nil {
		response.Code = constant.ResponseSystemErr
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}

	response.Code = constant.ResponseNormal
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
	id, _ := c.GetInt64("Id")
	username := c.GetString("Username")
	password := c.GetString("Password")
	newPassword := c.GetString("NewPassword")
	response := response.Response{}

	if password == newPassword {
		response.Code = constant.ResponseSystemErr
		response.Messgae = "新密码不能与原密码相同"
		c.Data["json"] = response
		c.ServeJSON()
	}

	args := &models.User{
		Id: id,
	}
	var loginReply models.User

	err := client.Call("http://127.0.0.1:2379", "User", "FindUserById", args, &loginReply)

	if err != nil {
		response.Code = constant.ResponseSystemErr
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}

	if loginReply.Username != username || loginReply.Password != password {
		response.Code = constant.ResponseSystemErr
		response.Messgae = "原密码错误"
		c.Data["json"] = response
		c.ServeJSON()
	}

	loginReply.Password = newPassword
	timestamp := time.Now().UnixNano() / 1e6
	loginReply.UpdateTime = timestamp

	err = client.Call("http://127.0.0.1:2379", "User", "UpdateUserById", args, nil)

	if err != nil {
		response.Code = constant.ResponseSystemErr
		response.Messgae = err.Error()
		c.Data["json"] = response
		c.ServeJSON()
	}

	response.Code = constant.ResponseNormal
	response.Messgae = "密码修改成功"
	c.Data["json"] = response
	c.ServeJSON()
}
