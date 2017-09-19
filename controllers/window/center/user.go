package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.cloud.360baige.com/utils"
	"dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/http/window/center"
	"dev.model.360baige.com/action"
	"sms.sdk.360baige.com/send"
	"strconv"
	"dev.model.360baige.com/models/account"
	"dev.cloud.360baige.com/utils/login/wechat"
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
	currentTimestamp := utils.CurrentTimestamp()
	username := c.GetString("username")
	password := c.GetString("password")
	err := utils.Unable(map[string]string{"username": "string:true", "password": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	usernameType, _ := utils.DetermineStringType(username)
	var replyUser user.User
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: usernameType, Val: username},
			action.CondValue{Type: "And", Key: "status", Val: "0"},
			action.CondValue{Type: "And", Key: "password", Val: password},
		},
	}, &replyUser)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if replyUser.Id == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
		c.ServeJSON()
		return
	}

	if currentTimestamp > replyUser.ExpireIn {
		createAccessTicket := utils.CreateAccessValue(replyUser.Username + "#" + strconv.FormatInt(currentTimestamp, 10))
		var updateReply action.Num
		expireIn := currentTimestamp + user.UserExpireIn
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", action.UpdateByIdCond{
			Id: []int64{replyUser.Id},
			UpdateList: []action.UpdateValue{
				action.UpdateValue{Key: "expire_in", Val: expireIn },
				action.UpdateValue{Key: "access_ticket", Val: createAccessTicket },
			},
		}, &updateReply)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		} else {
			replyUser.AccessTicket = createAccessTicket
			replyUser.ExpireIn = expireIn
		}
	}

	c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: UserLogin{
		Username:     replyUser.Username,
		Head:         utils.SignURLSample(replyUser.Head),
		AccessTicket: replyUser.AccessTicket,
		ExpireIn:     replyUser.ExpireIn,
	}}
	c.ServeJSON()
	return
}

// @Title 微信登录接口
// @Description 微信登录接口
// @Success 200 {"code":200,"message":"登录成功"}
// @Param   code     query   string true       ""
// @Failure 400 {"code":400,"message":"登录失败"}
// @router /weChatLogin [post]
func (c *UserController) WeChatLogin() {
	type data UserLoginResponse
	code := c.GetString("code")
	err := utils.Unable(map[string]string{"code": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	userInfo, err := wechat.GetUserInfo(code)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if userInfo.Openid == "" {
		c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
		c.ServeJSON()
		return
	}

	var replyUser user.User
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "wx_open_id", Val: userInfo.Openid},
			action.CondValue{Type: "And", Key: "status", Val: user.UserStatusNormal},
		},
	}, &replyUser)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	if replyUser.Id != 0 && replyUser.Status == 0 {
		c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: UserLogin{
			Username:     replyUser.Username,
			Head:         utils.SignURLSample(replyUser.Head),
			AccessTicket: replyUser.AccessTicket,
			ExpireIn:     replyUser.ExpireIn,
		}}
		c.ServeJSON()
		return
	}
	type bindData UserBindResponse
	c.Data["json"] = bindData{Code: ErrorUnBind, Message: "", Data: UserBind{
		OpenId: userInfo.Openid,
	}}
	c.ServeJSON()
	return
}

// @Title 微信账号绑定
// @Description 微信账号绑定
// @Success 200 {"code":200,"message":"登录成功"}
// @Param   openId     query   string true       ""
// @Failure 400 {"code":400,"message":"登录失败"}
// @router /bindAccount [post]
func (c *UserController) WeChatBindAccount() {
	type data UserLoginResponse
	currentTimestamp := utils.CurrentTimestamp()
	Type, _ := c.GetInt64("type")
	username := c.GetString("username")
	password := c.GetString("password")
	phone := c.GetString("phone")
	verifyCode := c.GetString("verifyCode")
	openId := c.GetString("openId")
	err := utils.Unable(map[string]string{"username": "string:true", "password": "string:true", "phone": "string:true", "verifyCode": "string:true", "openId": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	var replyUser user.User
	if Type == 1 {
		usernameType, _ := utils.DetermineStringType(username)
		err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", action.FindByCond{
			CondList: []action.CondValue{
				action.CondValue{Type: "And", Key: usernameType, Val: username},
				action.CondValue{Type: "And", Key: "status", Val: "0"},
				action.CondValue{Type: "And", Key: "password", Val: password},
			},
		}, &replyUser)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}
		if replyUser.Id == 0 {
			c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
			c.ServeJSON()
			return
		}

		var replyWxUser user.User
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", action.FindByCond{
			CondList: []action.CondValue{
				action.CondValue{Type: "And", Key: "wx_open_id", Val: openId},
			},
		}, &replyWxUser)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}
		if replyWxUser.Id != 0 {
			c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
			c.ServeJSON()
			return
		}

		var replyNum action.Num
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", action.UpdateByIdCond{
			Id: []int64{replyUser.Id},
			UpdateList: []action.UpdateValue{
				action.UpdateValue{Key: "update_time", Val: currentTimestamp},
				action.UpdateValue{Key: "wx_open_id", Val: openId},
			},
		}, &replyNum)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}
		if replyNum.Value == 0 {
			c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
			c.ServeJSON()
			return
		}

		if currentTimestamp > replyUser.ExpireIn {
			createAccessTicket := utils.CreateAccessValue(replyUser.Username + "#" + strconv.FormatInt(currentTimestamp, 10))
			var updateReply action.Num
			expireIn := currentTimestamp + 60*1000
			err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", action.UpdateByIdCond{
				Id: []int64{replyUser.Id},
				UpdateList: []action.UpdateValue{
					action.UpdateValue{Key: "expire_in", Val: expireIn },
					action.UpdateValue{Key: "access_ticket", Val: createAccessTicket },
				},
			}, &updateReply)
			if err != nil {
				c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
				c.ServeJSON()
				return
			}
			replyUser.AccessTicket = createAccessTicket
			replyUser.ExpireIn = expireIn
		}
	} else {
		// 验证码
		if verifyCode != "95888" {
			c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
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
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}
		if replyNum.Value > 0 {
			c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
			c.ServeJSON()
			return
		}
		// 注册
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "Add", &user.User{
			CreateTime:   currentTimestamp,
			UpdateTime:   currentTimestamp,
			Username:     username,
			Password:     password,
			Phone:        phone,
			WxOpenId:     openId,
			AccessTicket: utils.CreateAccessValue(username),
		}, &replyUser)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}
		if replyUser.Id == 0 {
			c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
			c.ServeJSON()
			return
		}

		// 初始化身份
		var replyUserPosition user.UserPosition
		err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "Add", &user.UserPosition{
			CreateTime:  currentTimestamp,
			UpdateTime:  currentTimestamp,
			UserId:      replyUser.Id,
			CompanyId:   user.UserPositionCompanyIdInit,
			Type:        user.UserPositionTypeVisitor,
			AccessToken: utils.CreateAccessValue(username),
		}, &replyUserPosition)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}
		if replyUserPosition.Id == 0 {
			c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
			c.ServeJSON()
			return
		}
		// 初始化现金账户
		var replyAccount account.Account
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "Add", &account.Account{
			CreateTime:       currentTimestamp,
			UpdateTime:       currentTimestamp,
			CompanyId:        user.UserPositionCompanyIdInit,
			UserId:           replyUser.Id,
			UserPositionType: user.UserPositionTypeVisitor,
			UserPositionId:   replyUserPosition.Id,
			Type:             account.AccountTypeMoney,
			Balance:          account.AccountBalanceInit,
			Status:           account.AccountStatusNormal,
		}, &replyAccount)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}
		if replyAccount.Id == 0 {
			c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
			c.ServeJSON()
			return
		}
		replyUser.Username = username
		replyUser.AccessTicket = utils.CreateAccessValue(username)
	}

	c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: UserLogin{
		Username:     replyUser.Username,
		Head:         utils.SignURLSample(replyUser.Head),
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
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	err := utils.Unable(map[string]string{"accessToken": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	replyUserPosition, err := utils.UserPosition(accessToken, currentTimestamp)
	if err != nil {
		c.Data["json"] = data{Code: ErrorPower, Message: err.Error()}
		c.ServeJSON()
		return
	}

	var reply user.User
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", user.User{
		Id: replyUserPosition.UserId,
	}, &reply)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	headUrl := utils.SignURLSample(reply.Head)
	c.Data["json"] = data{Code: Normal, Message: "获取用户信息成功", Data: UserDetail{
		Id:       reply.Id,
		Username: reply.Username,
		Email:    reply.Email,
		Phone:    reply.Phone,
		Head:     headUrl,
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
// @Success 200 {"code":200,"message":"用户修改密码成功"}
// @Param   accessToken query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户修改密码失败"}
// @router /modifyPassword [post]
func (c *UserController) ModifyPassword() {
	type data ModifyPasswordResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	password := c.GetString("password")
	newPassword := c.GetString("newPassword")
	err := utils.Unable(map[string]string{"accessToken": "string:true", "password": "string:true", "newPassword": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	replyUserPosition, err := utils.UserPosition(accessToken, currentTimestamp)
	if err != nil {
		c.Data["json"] = data{Code: ErrorPower, Message: err.Error()}
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
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", action.UpdateByIdCond{
		Id: []int64{replyUser.Id},
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "update_time", Val: currentTimestamp},
			action.UpdateValue{Key: "password", Val: newPassword},
		},
	}, &replyNum)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "密码修改成功"}
	c.ServeJSON()
	return
}

// @Title 用户信息修改接口
// @Description 用户信息修改接口
// @Success 200 {"code":200,"message":"用户修改密码成功"}}
// @Param   accessToken query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户修改密码失败"}
// @router /modify [get]
func (c *UserController) Modify() {
	type data ModifyPasswordResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	phone := c.GetString("phone")
	email := c.GetString("email")

	err := utils.Unable(map[string]string{"accessToken": "string:true", "id": "int:true", "phone": "string:true", "email": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	replyUserPosition, err := utils.UserPosition(accessToken, currentTimestamp)
	if err != nil {
		c.Data["json"] = data{Code: ErrorPower, Message: err.Error()}
		c.ServeJSON()
		return
	}

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", action.UpdateByIdCond{
		Id: []int64{replyUserPosition.UserId},
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "update_time", Val: currentTimestamp},
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
	err := utils.Unable(map[string]string{"phone": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

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

	err = send.MessageCode("百鸽互联科技有限公司", "95888", phone)
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
	err := utils.Unable(map[string]string{"val": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "CountByCond", &action.CountByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: key, Val: val},
			action.CondValue{Type: "And", Key: "status", Val: user.UserStatusNormal},
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
	c.Data["json"] = data{Code: Normal, Message: "SUCCESS"}
	c.ServeJSON()
	return
}

// @Title 用户注册
// @Description 用户注册
// @Success 200 {"code":200,"message":"用户注册成功"}
// @Failure 400 {"code":400,"message":"用户注册失败"}
// @Param   username     query   string true       "用户名"
// @Param   password     query   string true       "密码"
// @Param   phone query   string true       "手机号码"
// @Param   verifyCode query   string true       "验证码"
// @router /register [post]
func (c *UserController) Register() {
	type data UserRegisterResponse
	currentTimestamp := utils.CurrentTimestamp()
	username := c.GetString("username")
	password := c.GetString("password")
	phone := c.GetString("phone")
	verifyCode := c.GetString("verifyCode")

	err := utils.Unable(map[string]string{"username": "string:true", "password": "string:true", "phone": "int:true", "verifyCode": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	// 验证码
	if verifyCode != "95888" {
		c.Data["json"] = data{Code: ErrorLogic, Message: "验证码错误"}
		c.ServeJSON()
		return
	}
	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "CountByCond", &action.CountByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "phone", Val: phone},
			action.CondValue{Type: "Or", Key: "username", Val: username},
			action.CondValue{Type: "And", Key: "status", Val: 0},
		},
	}, &replyNum)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if replyNum.Value > 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
		c.ServeJSON()
		return
	}

	// 注册
	var replyUser user.User
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "Add", &user.User{
		CreateTime:   currentTimestamp,
		UpdateTime:   currentTimestamp,
		Username:     username,
		Password:     password,
		Phone:        phone,
		Head:         user.HEAD,
		AccessTicket: utils.CreateAccessValue(username),
	}, &replyUser)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if replyUser.Id == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
		c.ServeJSON()
		return
	}

	// 初始化身份
	var replyUserPosition user.UserPosition
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "Add", &user.UserPosition{
		CreateTime:  currentTimestamp,
		UpdateTime:  currentTimestamp,
		UserId:      replyUser.Id,
		CompanyId:   user.UserPositionCompanyIdInit,
		Type:        user.UserPositionTypeVisitor,
		AccessToken: utils.CreateAccessValue(username),
	}, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if replyUserPosition.Id == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
		c.ServeJSON()
		return
	}

	// 初始化现金账户
	var replyAccount account.Account
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "Add", &account.Account{
		CreateTime:       currentTimestamp,
		UpdateTime:       currentTimestamp,
		CompanyId:        user.UserPositionCompanyIdInit,
		UserId:           replyUser.Id,
		UserPositionType: user.UserPositionTypeVisitor,
		UserPositionId:   replyUserPosition.Id,
		Type:             account.AccountTypeMoney,
		Balance:          account.AccountBalanceInit,
		Status:           account.AccountStatusNormal,
	}, &replyAccount)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if replyAccount.Id == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "SUCCESS"}
	c.ServeJSON()
	return
}

// @Title 用户注册
// @Description 用户注册
// @Success 200 {"code":200,"message":"用户头像上传失败"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   id     query   string true       "用户ID"
// @Param   uploadFile query   string true       "file"
// @Failure 400 {"code":400,"message":"用户头像上传成功"}
// @router /uploadHead [options,post]
func (c *UserController) UploadHead() {
	currentTimestamp := utils.CurrentTimestamp()
	requestType := c.Ctx.Request.Method
	if requestType == "POST" {
		type data UploadHeadResponse
		accessToken := c.GetString("accessToken")
		_, handle, _ := c.Ctx.Request.FormFile("uploadFile")

		err := utils.Unable(map[string]string{"accessToken": "string:true"}, c.Ctx.Input)
		if err != nil {
			c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
			c.ServeJSON()
			return
		}

		replyUserPosition, err := utils.UserPosition(accessToken, currentTimestamp)
		if err != nil {
			c.Data["json"] = data{Code: ErrorPower, Message: err.Error()}
			c.ServeJSON()
			return
		}

		objectKey, err := utils.UploadImage(handle, "User/HeadImages/")
		if err != nil {
			c.Data["json"] = data{Code: ErrorLogic, Message: "用户头像上传失败"}
			c.ServeJSON()
			return
		}

		var replyUser user.User
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", action.FindByCond{
			CondList: []action.CondValue{action.CondValue{Type: "And", Key: "id", Val: replyUserPosition.UserId },
			},
			Fileds: []string{"id"},
		}, &replyUser)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "用户信息错误"}
			c.ServeJSON()
			return
		}

		var replyNum action.Num
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", action.UpdateByIdCond{
			Id: []int64{replyUser.Id},
			UpdateList: []action.UpdateValue{
				action.UpdateValue{Key: "update_time", Val: currentTimestamp },
				action.UpdateValue{Key: "head", Val: objectKey},
			},
		}, &replyNum)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "用户头像上传失败"}
			c.ServeJSON()
			return
		}
		headUrl := utils.SignURLSample(objectKey)
		c.Data["json"] = data{Code: Normal, Data: headUrl, Message: "用户头像上传成功"}
		c.ServeJSON()
		return
	}
}
