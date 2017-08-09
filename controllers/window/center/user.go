package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.cloud.360baige.com/utils"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/http/window/center"
	"time"
	"strconv"
	"dev.model.360baige.com/action"
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
	res := UserLoginResponse{}
	//res := Response{}
	username := c.GetString("username")
	password := c.GetString("password")
	if (username == "" || password == "") {
		res.Code = ResponseSystemErr
		res.Messgae = "登录失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	username_type, _ := utils.DetermineStringType(username) // 2 判断username 类型 属于 百鸽账号、邮箱、手机号码中的哪一种？

	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type:"And",
		Key:username_type,
		Val:username,
	}, action.CondValue{
		Type:"And",
		Key:"status",
		Val:"0",
	}, action.CondValue{
		Type:"And",
		Key:"password",
		Val:password,
	})
	fmt.Print("args>>>", args)
	var reply User
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", args, &reply)
	fmt.Print("reply>>>", reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "登录失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var updateArgs action.UpdateByIdCond
	updateArgs.Id = []int64{reply.Id}
	nawAccessTicket := reply.Username + strconv.FormatInt(time.Now().UnixNano() / 1e6, 10)
	updateArgs.UpdateList = append(updateArgs.UpdateList, action.UpdateValue{
		Key:"update_time",
		Val:time.Now().UnixNano() / 1e6,
	}, action.UpdateValue{
		Key:"AccessTicket",
		Val:nawAccessTicket, // 更新船票 应该判断时效，再做更新
	})
	var updateReply action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", updateArgs, &updateReply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "登录失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.Code = ResponseNormal
	res.Messgae = "登录成功"
	res.Data.Id = reply.Id
	res.Data.Username = reply.Username
	res.Data.AccessTicket = nawAccessTicket
	res.Data.ExpireIn = reply.ExpireIn
	c.Data["json"] = res
	c.ServeJSON()
}

//// @Title 用户注册接口
//// @Description 用户注册接口
//// @Success 200 {"code":200,"messgae":"注册成功","data":{}
//// @Param   username     query   string true       "用户名：百鸽账号、邮箱、手机号码 三种登录方式"
//// @Param   password query   string true       "密码"
//// @Failure 400 {"code":400,"message":"注册失败"}
//// @router /register [post]
//func (c *UserController) Register() {
//	res := UserLoginResponse{}
//	//res := Response{}
//	username := c.GetString("username")
//	password := c.GetString("password")
//	if (username == "" || password == "") {
//		res.Code = ResponseSystemErr
//		res.Messgae = "注册失败"
//		c.Data["json"] = res
//		c.ServeJSON()
//	}
//
//	var args User
//	args.UpdateTime, args.UpdateTime = time.Now().UnixNano() / 1e6
//	args.Username = username
//	args.Password = password
//	args.AccessTicket = username + strconv.FormatInt(time.Now().UnixNano() / 1e6, 10)
//	var reply User
//	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "Add", args, &reply)
//	if err != nil {
//		res.Code = ResponseSystemErr
//		res.Messgae = "注册失败"
//		c.Data["json"] = res
//		c.ServeJSON()
//	}
//
//	var positionArgs UserPosition
//	positionArgs.UpdateTime, positionArgs.UpdateTime = time.Now().UnixNano() / 1e6
//	positionArgs.UserId = reply.Id
//	args.AccessTicket = username + strconv.FormatInt(time.Now().UnixNano() / 1e6, 10)
//
//	var reply UserPosition
//	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "Add", updateArgs, &updateReply)
//	if err != nil {
//		res.Code = ResponseSystemErr
//		res.Messgae = "注册失败"
//		c.Data["json"] = res
//		c.ServeJSON()
//	}
//
//
//
//	res.Code = ResponseNormal
//	res.Messgae = "注册成功"
//	res.Data.Id = reply.Id
//	res.Data.Username = reply.Username
//	res.Data.AccessTicket = nawAccessTicket
//	res.Data.ExpireIn = reply.ExpireIn
//	c.Data["json"] = res
//	c.ServeJSON()
//}

// @Title 用户信息接口
// @Description 用户信息接口
// @Success 200 {"code":200,"messgae":"获取用户信息成功","data":{"id":"xxxx","username":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取用户信息失败"}
// @router /detail [get]
func (c *UserController) Detail() {
	res := UserDetailResponse{}
	access_token := c.GetString("access_token")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type:"And",
		Key:"access_token",
		Val:access_token,
	})
	args.Fileds = []string{"user_id"}
	var replyAccessToken UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken) // 检测 accessToken
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	if replyAccessToken.UserId == 0 {
		res.Code = ResponseSystemErr
		res.Messgae = "获取用户信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var userArgs User
	userArgs.Id = replyAccessToken.UserId
	var reply User
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", userArgs, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取用户信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.Code = ResponseNormal
	res.Messgae = "获取用户信息成功"
	res.Data.Id = reply.Id
	res.Data.Username = reply.Username
	res.Data.Email = reply.Email
	res.Data.Phone = reply.Phone
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 用户退出接口
// @Description 用户退出接口
// @Success 200 {"code":200,"messgae":"用户退出成功","data":{"access_ticket":"xxxx","expire_in":0}}// @Param   password query   string true       "密码"
// @Param   access_token query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户退出失败"}
// @router /logout [post]
func (c *UserController) Logout() {
	res := UserDetailResponse{}
	//access_token := c.GetString("access_token")
	//if access_token == "" {
	//	res.Code = ResponseSystemErr
	//	res.Messgae = "访问令牌无效"
	//	c.Data["json"] = res
	//	c.ServeJSON()
	//	return
	//}
	//
	//var args action.FindByCond
	//args.CondList = append(args.CondList, action.CondValue{
	//	Type:"And",
	//	Key:"access_ticket",
	//	Val:access_token,
	//})
	//var reply User
	//err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", args, &reply) // 检测 accessToken
	//
	//if err != nil {
	//	res.Code = ResponseLogicErr
	//	res.Messgae = "访问令牌失效"
	//	c.Data["json"] = res
	//	c.ServeJSON()
	//	return
	//}

	res.Code = ResponseNormal
	res.Messgae = "验证成功"
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 用户修改密码接口
// @Description 用户修改密码接口
// @Success 200 {"code":200,"messgae":"用户修改密码成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户修改密码失败"}
// @router /modifypassword [post]
func (c *UserController) ModifyPassword() {
	res := ModifyPasswordResponse{}
	access_token := c.GetString("access_token")
	password := c.GetString("password")
	newPassword := c.GetString("newPassword")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	//检测 accessToken
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type:"And",
		Key:"access_token",
		Val:access_token,
	})
	var replyAccessToken UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var userArgs action.FindByCond
	userArgs.CondList = append(userArgs.CondList, action.CondValue{
		Type:"And",
		Key:"id",
		Val:replyAccessToken.UserId,
	}, action.CondValue{
		Type:"And",
		Key:"password",
		Val:password,
	})
	userArgs.Fileds = []string{"id"}
	var reply User
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", userArgs, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "用户密码错误！"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var updateArgs action.UpdateByIdCond
	updateArgs.Id = []int64{reply.Id}
	updateArgs.UpdateList = append(updateArgs.UpdateList, action.UpdateValue{
		Key:"update_time",
		Val:time.Now().UnixNano() / 1e6,
	}, action.UpdateValue{
		Key:"password",
		Val:newPassword,
	})
	var updateReply action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", updateArgs, &updateReply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "密码修改失败！"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.Code = ResponseNormal
	res.Messgae = "密码修改成功！"
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 用户信息修改接口
// @Description 用户信息修改接口
// @Success 200 {"code":200,"messgae":"用户修改密码成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户修改密码失败"}
// @router /modify [get]
func (c *UserController) Modify() {
	res := ModifyPasswordResponse{}
	access_token := c.GetString("access_token")
	userId, _ := c.GetInt64("id")
	phone := c.GetString("phone")
	email := c.GetString("email")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	//检测 accessToken
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type:"And",
		Key:"access_token",
		Val:access_token,
	})
	var replyAccessToken UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var updateArgs action.UpdateByIdCond
	updateArgs.Id = []int64{userId}
	updateArgs.UpdateList = append(updateArgs.UpdateList, action.UpdateValue{
		Key:"update_time",
		Val:time.Now().UnixNano() / 1e6,
	}, action.UpdateValue{
		Key:"phone",
		Val:phone,
	}, action.UpdateValue{
		Key:"email",
		Val:email,
	})
	var updateReply action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", updateArgs, &updateReply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "用户信息修改失败！"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.Code = ResponseNormal
	res.Messgae = "用户信息修改成功！"
	c.Data["json"] = res
	c.ServeJSON()
}