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
)

// USER API
type UserController struct {
	beego.Controller
}

// @Title 用户登录接口
// @Description 用户登录接口
// @Success 200 {"code":200,"message":"登录成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   username     query   string true       "用户名：百鸽账号、邮箱、手机号码 三种登录方式"
// @Param   password query   string true       "密码"
// @Failure 400 {"code":400,"message":"登录失败"}
// @router /login [post]
func (c *UserController) Login() {
	res := UserLoginResponse{}
	username := c.GetString("username")
	password := c.GetString("password")
	if username == "" || password == "" {
		res.Code = ResponseSystemErr
		res.Message = "登录失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	username_type, _ := utils.DetermineStringType(username) // 2 判断username 类型 属于 百鸽账号、邮箱、手机号码中的哪一种？
	args := action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: username_type, Val: username},
			action.CondValue{Type: "And", Key: "status", Val: "0"},
			action.CondValue{Type: "And", Key: "password", Val: password},
		},
	}
	var reply user.User
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", args, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "登录失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var updateArgs action.UpdateByIdCond
	updateArgs.Id = []int64{reply.Id}
	nawAccessTicket := reply.Username + strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	updateArgs.UpdateList = append(updateArgs.UpdateList, action.UpdateValue{
		Key: "update_time",
		Val: time.Now().UnixNano() / 1e6,
	}, action.UpdateValue{
		Key: "AccessTicket",
		Val: nawAccessTicket, // 更新船票 应该判断时效，再做更新
	})
	var updateReply action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", updateArgs, &updateReply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "登录失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.Code = ResponseNormal
	res.Message = "登录成功"
	res.Data.Username = reply.Username
	res.Data.AccessTicket = nawAccessTicket
	res.Data.ExpireIn = reply.ExpireIn
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 用户信息接口
// @Description 用户信息接口
// @Success 200 {"code":200,"message":"获取用户信息成功","data":{"id":"xxxx","username":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取用户信息失败"}
// @router /detail [get]
func (c *UserController) Detail() {
	res := UserDetailResponse{}
	access_token := c.GetString("access_token")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type: "And",
		Key:  "access_token",
		Val:  access_token,
	})
	args.Fileds = []string{"user_id"}
	var replyAccessToken user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken) // 检测 accessToken
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	if replyAccessToken.UserId == 0 {
		res.Code = ResponseSystemErr
		res.Message = "获取用户信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var userArgs user.User
	userArgs.Id = replyAccessToken.UserId
	var reply user.User
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", userArgs, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "获取用户信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.Code = ResponseNormal
	res.Message = "获取用户信息成功"
	res.Data.Id = reply.Id
	res.Data.Username = reply.Username
	res.Data.Email = reply.Email
	res.Data.Phone = reply.Phone
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 用户退出接口
// @Description 用户退出接口
// @Success 200 {"code":200,"message":"用户退出成功","data":{"access_ticket":"xxxx","expire_in":0}}// @Param   password query   string true       "密码"
// @Param   access_token query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户退出失败"}
// @router /logout [post]
func (c *UserController) Logout() {
	res := UserDetailResponse{}
	//access_token := c.GetString("access_token")
	//if access_token == "" {
	//	res.Code = ResponseSystemErr
	//	res.Message = "访问令牌无效"
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
	//	res.Message = "访问令牌失效"
	//	c.Data["json"] = res
	//	c.ServeJSON()
	//	return
	//}

	res.Code = ResponseNormal
	res.Message = "验证成功"
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 用户修改密码接口
// @Description 用户修改密码接口
// @Success 200 {"code":200,"message":"用户修改密码成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户修改密码失败"}
// @router /modifyPassword [post]
func (c *UserController) ModifyPassword() {
	res := ModifyPasswordResponse{}
	access_token := c.GetString("access_token")
	password := c.GetString("password")
	newPassword := c.GetString("newPassword")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	//检测 accessToken
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type: "And",
		Key:  "access_token",
		Val:  access_token,
	})
	var replyAccessToken user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var userArgs action.FindByCond
	userArgs.CondList = append(userArgs.CondList, action.CondValue{
		Type: "And",
		Key:  "id",
		Val:  replyAccessToken.UserId,
	}, action.CondValue{
		Type: "And",
		Key:  "password",
		Val:  password,
	})
	userArgs.Fileds = []string{"id"}
	var reply user.User
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", userArgs, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "用户密码错误！"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var updateArgs action.UpdateByIdCond
	updateArgs.Id = []int64{reply.Id}
	updateArgs.UpdateList = append(updateArgs.UpdateList, action.UpdateValue{
		Key: "update_time",
		Val: time.Now().UnixNano() / 1e6,
	}, action.UpdateValue{
		Key: "password",
		Val: newPassword,
	})
	var updateReply action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", updateArgs, &updateReply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "密码修改失败！"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.Code = ResponseNormal
	res.Message = "密码修改成功！"
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 用户信息修改接口
// @Description 用户信息修改接口
// @Success 200 {"code":200,"message":"用户修改密码成功","data":{"access_ticket":"xxxx","expire_in":0}}
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
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	//检测 accessToken
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type: "And",
		Key:  "access_token",
		Val:  access_token,
	})
	var replyAccessToken user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var updateArgs action.UpdateByIdCond
	updateArgs.Id = []int64{userId}
	updateArgs.UpdateList = append(updateArgs.UpdateList, action.UpdateValue{
		Key: "update_time",
		Val: time.Now().UnixNano() / 1e6,
	}, action.UpdateValue{
		Key: "phone",
		Val: phone,
	}, action.UpdateValue{
		Key: "email",
		Val: email,
	})
	var updateReply action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", updateArgs, &updateReply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "用户信息修改失败！"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.Code = ResponseNormal
	res.Message = "用户信息修改成功！"
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 验证用户名
// @Description 验证用户名
// @Success 200 {"code":200,"message":"用户名可用"}
// @Param   type     query   string true      "类型:0已注册手机号码不发送; !0不限制"
// @Param   phone     query   string true      "手机号码"
// @Failure 400 {"code":400,"message":"用户已被注册"}
// @router /sendMessageCode [post]
func (c *UserController) SendMessageCode() {
	res := SendMessageCodeResponse{}
	_type, _ := c.GetInt64("type", 0)
	phone := c.GetString("phone")

	// 1.验证
	if _type == 0 {
		// 1.验证是否可以注册
		var countReply action.Num
		err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "CountByCond", &action.CountByCond{
			CondList: []action.CondValue{
				action.CondValue{Type: "And", Key: "phone", Val: phone},
				action.CondValue{Type: "And", Key: "status", Val: 0},
			},
		}, &countReply)
		log.Println("err:", err, "countReply:", countReply)
		if err != nil {
			res.Code = ResponseSystemErr
			res.Message = "系统异常[手机号码验证失败]"
			c.Data["json"] = res
			c.ServeJSON()
			return
		}
		if countReply.Value > 0 {
			res.Code = ResponseLogicErr
			res.Message = "手机号码已被注册"
			c.Data["json"] = res
			c.ServeJSON()
			return
		}
	}

	// 2. 发送
	err := send.MessageCode("百鸽互联科技有限公司", "95888", phone)
	log.Println("err:", err)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "验证码发送失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	} else {
		res.Code = ResponseNormal
		res.Message = "验证码发送成功"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
}

// @Title 验证用户名
// @Description 验证用户名
// @Success 200 {"code":200,"message":"用户名可用"}
// @Param   key     query   string true      "键：phone|username"
// @Param   val     query   string true       "值"
// @Failure 400 {"code":400,"message":"用户已被注册"}
// @router /existKey [post]
func (c *UserController) ExistKey() {
	res := ExistKeyResponse{}
	key := c.GetString("key")
	val := c.GetString("val")

	// 1.验证是否可以注册
	var countReply action.Num
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "CountByCond", &action.CountByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: key, Val: val},
			action.CondValue{Type: "And", Key: "status", Val: 0},
		},
	}, &countReply)
	log.Println("err:", err, "countReply:", countReply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "系统异常[" + key + "验证失败]"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if countReply.Value > 0 {
		res.Code = ResponseLogicErr
		res.Message = key + "已被注册"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	res.Code = ResponseNormal
	res.Message = key + "可用"
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// @Title 用户注册
// @Description 用户注册
// @Success 200 {"code":200,"message":"用户注册成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   username     query   string true       "用户名"
// @Param   password     query   string true       "密码"
// @Param   verify_password query   string true       "确认密码"
// @Param   phone query   string true       "手机号码"
// @Param   verify_code query   string true       "验证码"
// @Failure 400 {"code":400,"message":"用户注册失败"}
// @router /register [post]
func (c *UserController) Register() {
	res := UserRegisterResponse{}
	username := c.GetString("username")
	password := c.GetString("password")
	verify_password := c.GetString("verify_password")
	phone := c.GetString("phone")
	verify_code := c.GetString("verify_code")

	// 1.验证验证码
	if verify_code != "95888" {
		res.Code = ResponseSystemErr
		res.Message = "验证码错误"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	log.Println("验证验证码成功")

	// 2.验证确认密码
	if password != verify_password {
		res.Code = ResponseSystemErr
		res.Message = "密码与确认密码不一致"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	log.Println("验证确认密码成功")

	// 3.验证是否可以注册
	var countReply action.Num
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "CountByCond", &action.CountByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "phone", Val: phone},
			action.CondValue{Type: "Or", Key: "username", Val: username},
			action.CondValue{Type: "And", Key: "status", Val: 0},
		},
	}, &countReply)
	log.Println("err:", err, "countReply:", countReply)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "系统异常[用户名或手机号码验证失败]"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	if countReply.Value > 0 {
		res.Code = ResponseLogicErr
		res.Message = "用户名或手机号码已被注册"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	log.Println("验证是否可以注册成功")
	// 4.添加注册信息
	o_time := time.Now().UnixNano() / 1e6
	args := user.User{
		CreateTime: o_time,
		UpdateTime: o_time,
		Username:   username,
		Password:   password,
		Phone:      phone,
	}
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "Add", &args, &args)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "系统异常[用户注册失败]"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 4.添加注册信息
	res.Code = ResponseNormal
	res.Message = "用户注册成功"
	c.Data["json"] = res
	c.ServeJSON()
	return
}
