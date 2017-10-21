package safe

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.cloud.360baige.com/utils"
	"dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/http/window/safe"
	"dev.model.360baige.com/action"
	"sms.sdk.360baige.com/send"
	mail "mail.sdk.360baige.com/client"
)

// USER API
type UserController struct {
	beego.Controller
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

// @Title 用户修改密码接口
// @Description 用户修改密码接口
// @Success 200 {"code":200,"message":"用户修改密码成功"}
// @Param   accessToken query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户修改密码失败"}
// @router /resetPassword [post]
func (c *UserController) ResetPassword() {
	type data ModifyPasswordResponse
	currentTimestamp := utils.CurrentTimestamp()
	modifyType, _ := c.GetInt("modifyType", 0)
	accessToken := c.GetString("accessToken")
	userId, _ := c.GetInt64("userId")
	password := c.GetString("password")
	err := utils.Unable(map[string]string{"modifyType": "int:true", "accessToken": "string:false", "password": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	var replyUserPosition *user.UserPosition
	var replyUser user.User
	if modifyType == 1 {
		replyUserPosition, err = utils.UserPosition(accessToken, currentTimestamp)
		if err != nil {
			c.Data["json"] = data{Code: ErrorPower, Message: err.Error()}
			c.ServeJSON()
			return
		}
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", user.User{
			Id: replyUserPosition.UserId,
		}, &replyUser)
	} else if modifyType == 0 {
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", user.User{
			Id: userId,
		}, &replyUser)
	}
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
			action.UpdateValue{Key: "password", Val: password},
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

// @Title 发送验证码
// @Description 发送验证码
// @Success 200 {"code":200,"message":"验证码发送成功"}
// @Param   verify     query   string true      "类型:0已注册手机号码不发送; !0不限制"
// @Param   phone     query   string true      "手机号码"
// @Failure 400 {"code":400,"message":"验证码发送失败"}
// @router /sendMessageCode [get,post]
func (c *UserController) SendMessageCode() {
	type data SendMessageCodeResponse
	userId, _ := c.GetInt64("userId")
	whether := c.GetString("whether", "0")
	key := c.GetString("key")
	val := c.GetString("val")
	err := utils.Unable(map[string]string{"key": "string:true", "val": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}
	var replyUser user.User
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", &action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: key, Val: val},
			action.CondValue{Type: "And", Key: "status", Val: 0},
		},
	}, &replyUser)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常[" + key + "验证失败]"}
		c.ServeJSON()
		return
	}

	if whether == "0" {
		if replyUser.Id <= 0 || replyUser.Id != userId {
			c.Data["json"] = data{Code: ErrorLogic, Message: key + "验证失败"}
			c.ServeJSON()
			return
		}
	} else if whether == "1" {
		if replyUser.Id > 0 {
			c.Data["json"] = data{Code: ErrorLogic, Message: key + "已绑定!请更换手机号码"}
			c.ServeJSON()
			return
		}
	}

	if key == "email" {
		user := "taihe@360baige.com" // 发送人
		password := "TaiHe0129" // 发送人密码
		host := "smtp.exmail.qq.com:25"
		to := val // 接收人
		subject := "使用Golang发送邮件" // 主题
		body := `<html><body><h3>"百鸽互联科技有限公司<br/>验证码为95888"</h3></body></html>` // 邮件内容
		err = mail.SendToMail(user, password, host, to, subject, body, "html")
	} else if key == "phone" {
		err = send.MessageCode("百鸽互联科技有限公司", "95888", val)
	}

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "验证码发送失败"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = data{Code: Normal, Message: "验证码发送成功"}
	c.ServeJSON()
	return
}

// @Title 验证字段
// @Description 验证字段
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

	var replyUser user.User
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: key, Val: val },
			action.CondValue{Type: "And", Key: "status", Val: user.UserStatusNormal },
		},
		Fileds: []string{"id"},
	}, &replyUser)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常[" + key + "验证失败]"}
		c.ServeJSON()
		return
	}
	if replyUser.Id <= 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "账号未注册"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data:ExistKeyUserData{
		UserId:replyUser.Id,
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

	c.Data["json"] = data{Code: Normal, Message: "获取用户信息成功", Data: UserDetail{
		Email:    reply.Email,
		Phone:    reply.Phone,
		Qq:       reply.QqOpenId,
		WeChat:   reply.WxOpenId,
	}}
	c.ServeJSON()
}


// @Title 密保信息绑定
// @Description 用户绑定
// @Success 200 {"code":200,"message":"获取用户信息成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取用户信息失败"}
// @router /bindOrUnbind [post]
func (c *UserController) BindOrUnbind() {
	type data BindOrUnbindResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	password := c.GetString("password")
	bindType := c.GetString("bindType", "1")
	key := c.GetString("key")
	val := c.GetString("val")
	err := utils.Unable(map[string]string{"accessToken": "string:true", "key": "string:true", "val": "string:false"}, c.Ctx.Input)
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
			action.CondValue{Type: "And", Key: "status", Val: user.UserStatusNormal },
		},
	}, &replyUser)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	if bindType == "2" && replyUser.Password != password {
		c.Data["json"] = data{Code: ErrorSystem, Message: "密码错误!"}
		c.ServeJSON()
		return
	}

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", action.UpdateByIdCond{
		Id: []int64{replyUser.Id},
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "update_time", Val: currentTimestamp},
			action.UpdateValue{Key: key, Val: val},
		},
	}, &replyNum)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: key + "绑定成功"}
	c.ServeJSON()
	return
}