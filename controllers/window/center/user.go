package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.cloud.360baige.com/utils"
	"dev.cloud.360baige.com/log"
	"dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/http/window/center"
	"dev.model.360baige.com/action"
	"sms.sdk.360baige.com/send"
	"strconv"
	"time"
	"fmt"
)

// USER API
type UserController struct {
	beego.Controller
}

// @Title 用户登录接口
// @Description 用户登录接口
// @Success 200 {"code":200,"message":"登录成功"}
// @Param   username     query   string true       "用户名：百鸽账号、邮箱、手机号码 三种登录方式"
// @Param   password query   string true       "密码"
// @Failure 400 {"code":400,"message":"登录失败"}
// @router /login [post]
func (c *UserController) Login() {
	type data UserLoginResponse
	username := c.GetString("username")
	password := c.GetString("password")
	if username == "" || password == "" {
		c.Data["json"] = data{Code: ErrorSystem, Message: "用户名或密码不能为空"}
		c.ServeJSON()
		return
	}

	usernameType, _ := utils.DetermineStringType(username) // 2 判断username 类型 属于 百鸽账号、邮箱、手机号码中的哪一种？
	fmt.Print("usernameType", usernameType)
	var replyUser user.User
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: usernameType, Val: username},
			action.CondValue{Type: "And", Key: "status", Val: "0"},
			action.CondValue{Type: "And", Key: "password", Val: password},
		},
	}, &replyUser)
	fmt.Print("replyUser", replyUser)
	if err != nil || replyUser.Id == 0 {
		c.Data["json"] = data{Code: ErrorSystem, Message: "登录失败"}
		c.ServeJSON()
		return
	}

	currentTime := time.Now().UnixNano() / 1e6
	if currentTime > replyUser.ExpireIn {
		createAccessTicket := utils.CreateAccessValue(replyUser.Username + "#" + strconv.FormatInt(currentTime, 10))
		var updateReply action.Num
		expireIn := currentTime + 60*1000
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", action.UpdateByIdCond{
			Id: []int64{replyUser.Id},
			UpdateList: []action.UpdateValue{
				action.UpdateValue{Key: "expire_in", Val: expireIn },
				action.UpdateValue{Key: "access_ticket", Val: createAccessTicket },
			},
		}, &updateReply)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "登录失败"}
			c.ServeJSON()
			return
		} else {
			replyUser.AccessTicket = createAccessTicket
			replyUser.ExpireIn = expireIn
		}
	}

	c.Data["json"] = data{Code: Normal, Message: "登录成功", Data: UserLogin{
		Head:         replyUser.Head,
		AccessTicket: replyUser.AccessTicket,
		ExpireIn:     replyUser.ExpireIn,
	}}
	c.ServeJSON()
	return
}

// @Title 用户信息接口
// @Description 用户信息接口
// @Success 200 {"code":200,"message":"获取用户信息成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取用户信息失败"}
// @router /detail [post]
func (c *UserController) Detail() {
	type data UserDetailResponse
	accessToken := c.GetString("accessToken")
	if accessToken == "" {
		c.Data["json"] = data{Code: Normal, Message: "访问令牌不能为空"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "access_token", Val: accessToken },
		},
		Fileds: []string{"user_id"},
	}, &replyUserPosition)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	if replyUserPosition.UserId == 0 {
		c.Data["json"] = data{Code: Normal, Message: "获取用户信息失败"}
		c.ServeJSON()
		return
	}

	var reply user.User
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", user.User{
		Id: replyUserPosition.UserId,
	}, &reply)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取用户信息失败"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "获取用户信息成功", Data: UserDetail{
		Id:       reply.Id,
		Username: reply.Username,
		Email:    reply.Email,
		Phone:    reply.Phone,
	}}
	c.ServeJSON()
}

// @Title 用户退出接口
// @Description 用户退出接口
// @Success 200 {"code":200,"message":"用户退出成功"}
// @Param   password query   string true       "密码"
// @Param   accessToken query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户退出失败"}
// @router /logout [post]
func (c *UserController) Logout() {
	type data UserDetailResponse
	c.Data["json"] = data{Code: Normal, Message: "退出成功"}
	c.ServeJSON()
}

// @Title 用户修改密码接口
// @Description 用户修改密码接口
// @Success 200 {"code":200,"message":"用户修改密码成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户修改密码失败"}
// @router /modifyPassword [post]
func (c *UserController) ModifyPassword() {
	type data ModifyPasswordResponse
	accessToken := c.GetString("accessToken")
	password := c.GetString("password")
	newPassword := c.GetString("newPassword")
	if accessToken == "" {
		c.Data["json"] = data{Code: Normal, Message: "访问令牌不能为空"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "access_token", Val: accessToken},
		},
	}, &replyUserPosition)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	var replyUser user.User
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "id", Val: replyUserPosition.UserId },
			action.CondValue{Type: "And", Key: "password", Val: password },
		},
		Fileds: []string{"id"},
	}, &replyUser)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "用户密码错误"}
		c.ServeJSON()
		return
	}

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", action.UpdateByIdCond{
		Id: []int64{replyUser.Id},
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "update_time", Val: time.Now().UnixNano() / 1e6 },
			action.UpdateValue{Key: "password", Val: newPassword},
		},
	}, &replyNum)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "密码修改失败"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "密码修改成功"}
	c.ServeJSON()
	return
}

// @Title 用户信息修改接口
// @Description 用户信息修改接口
// @Success 200 {"code":200,"message":"用户修改密码成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户修改密码失败"}
// @router /modify [get]
func (c *UserController) Modify() {
	type data ModifyPasswordResponse
	accessToken := c.GetString("accessToken")
	userId, _ := c.GetInt64("id")
	phone := c.GetString("phone")
	email := c.GetString("email")
	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "access_token", Val: accessToken },
		},
	}, &replyUserPosition)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", action.UpdateByIdCond{
		Id: []int64{userId},
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "update_time", Val: time.Now().UnixNano() / 1e6 },
			action.UpdateValue{Key: "phone", Val: phone },
			action.UpdateValue{Key: "email", Val: email },
		},
	}, &replyNum)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "用户信息修改失败"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "用户信息修改成功"}
	c.ServeJSON()
	return
}

// @Title 验证用户名
// @Description 验证用户名
// @Success 200 {"code":200,"message":"用户名可用"}
// @Param   verify     query   string true      "类型:0已注册手机号码不发送; !0不限制"
// @Param   phone     query   string true      "手机号码"
// @Failure 400 {"code":400,"message":"用户已被注册"}
// @router /sendMessageCode [post]
func (c *UserController) SendMessageCode() {
	type data SendMessageCodeResponse
	verify, _ := c.GetInt64("verify", 0)
	phone := c.GetString("phone")

	if verify == 0 {
		var replyNum action.Num
		err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "CountByCond", &action.CountByCond{
			CondList: []action.CondValue{
				action.CondValue{Type: "And", Key: "phone", Val: phone},
				action.CondValue{Type: "And", Key: "status", Val: 0},
			},
		}, &replyNum)

		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常[手机号码验证失败]"}
			c.ServeJSON()
			return
		}
		if replyNum.Value > 0 {
			c.Data["json"] = data{Code: ErrorLogic, Message: "手机号码已被注册"}
			c.ServeJSON()
			return
		}
	}

	err := send.MessageCode("百鸽互联科技有限公司", "95888", phone)
	log.Println("err:", err)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "验证码发送失败"}
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = data{Code: Normal, Message: "验证码发送成功"}
		c.ServeJSON()
		return
	}
}

// @Title 验证用户名
// @Description 验证用户名
// @Success 200 {"code":200,"message":"用户名可用"}
// @Param   key     query   string true      "键：phone|username 默认phone"
// @Param   val     query   string true       "值"
// @Failure 400 {"code":400,"message":"用户已被注册"}
// @router /existKey [post]
func (c *UserController) ExistKey() {
	type data ExistKeyResponse
	key := c.GetString("key", "phone")
	val := c.GetString("val")

	var replyNum action.Num
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "CountByCond", &action.CountByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: key, Val: val},
			action.CondValue{Type: "And", Key: "status", Val: 0},
		},
	}, &replyNum)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常[" + key + "验证失败]"}
		c.ServeJSON()
		return
	}
	if replyNum.Value > 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: key + "已被注册"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = data{Code: Normal, Message: key + "可用"}
	c.ServeJSON()
	return
}

// @Title 用户注册
// @Description 用户注册
// @Success 200 {"code":200,"message":"用户注册成功"}
// @Param   username     query   string true       "用户名"
// @Param   password     query   string true       "密码"
// @Param   phone query   string true       "手机号码"
// @Param   verifyCode query   string true       "验证码"
// @Failure 400 {"code":400,"message":"用户注册失败"}
// @router /register [post]
func (c *UserController) Register() {
	type data UserRegisterResponse
	username := c.GetString("username")
	password := c.GetString("password")
	phone := c.GetString("phone")
	verifyCode := c.GetString("verifyCode")
	currentTimestamp := utils.CurrentTimestamp()

	if verifyCode != "95888" {
		c.Data["json"] = data{Code: ErrorLogic, Message: "验证码错误"}
		c.ServeJSON()
		return
	}

	var replyNum action.Num
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "CountByCond", &action.CountByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "phone", Val: phone},
			action.CondValue{Type: "Or", Key: "username", Val: username},
			action.CondValue{Type: "And", Key: "status", Val: 0},
		},
	}, &replyNum)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常[用户名或手机号码验证失败]"}
		c.ServeJSON()
		return
	}
	if replyNum.Value > 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "用户名或手机号码已被注册"}
		c.ServeJSON()
		return
	}
	var replyUser user.User
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "Add", &user.User{
		CreateTime: currentTimestamp,
		UpdateTime: currentTimestamp,
		Username:   username,
		Password:   password,
		Phone:      phone,
	}, &replyUser)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常[用户注册失败]"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "用户注册成功"}
	c.ServeJSON()
	return
}
