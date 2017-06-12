package controllers

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
)

type UserController struct {
	beego.Controller
}

// @Title 登录接口
// @Description 登录接口
// @Success 200 {"code":200,"messgae":"ok","data":{"username":"zhangs","accessToken":"ok"}}
// @Param   username     query   string true       "用户名"
// @Param   username     query   string true       "密码"
// @Failure 400 {"code":400,"message":"..."}
// @router /login [post]
func (c *UserController) Login() {
	// TODO
	error := client.Call("http://127.0.0.1:2379","User","Login","","")
	if error == nil {

	} else {

	}
}

// @Title 退出接口
// @Description 退出接口
// @Success 200 {"code":200,"messgae":"ok"}
// @Param   username     query   string true       "用户名"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /logout [post]
func (c *UserController) Logout() {
	// TODO
}

// @Title 用户信息详情
// @Description 用户信息详情
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   userId     query   string true       "用户ID"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /detail [post]
func (c *UserController) Detail() {
	// TODO
}

// @Title 用户信息修改
// @Description 用户信息修改
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   userId     query   string true       "用户ID"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /modify [post]
func (c *UserController) Modify() {
	// TODO
}

// @Title 用户密码修改
// @Description 用户密码修改
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   userId     query   string true       "用户ID"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /modifyPassword [post]
func (c *UserController) ModifyPassword() {
	// TODO
}