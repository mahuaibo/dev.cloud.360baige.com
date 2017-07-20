package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.cloud.360baige.com/utils"
	. "dev.model.360baige.com/models/user"
	//. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/http/window"
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
	}

	username_type, _ := utils.DetermineStringType(username) // 2 判断username 类型 属于 百鸽账号、邮箱、手机号码中的哪一种？

	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type:"And",
		Key:username_type,
		Val:username,
	})
	args.CondList = append(args.CondList, action.CondValue{
		Type:"And",
		Key:"status",
		Val:"0",
	})
	var reply User
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", args, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "登录失败"
		c.Data["json"] = res
		c.ServeJSON()
	}

	var updateArgs action.UpdateByIdCond
	updateArgs.Id = []int64{reply.Id}
	updateArgs.UpdateList = append(updateArgs.UpdateList, action.UpdateValue{
		Key:"update_time",
		Val:time.Now().UnixNano() / 1e6,
	})
	updateArgs.UpdateList = append(updateArgs.UpdateList, action.UpdateValue{
		Key:"AccessTicket",
		Val:reply.Username + strconv.FormatInt(time.Now().UnixNano() / 1e6, 10), // 更新船票 应该判断时效，再做更新
	})
	var updateReply action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", updateArgs, &updateReply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "登录失败"
		c.Data["json"] = res
		c.ServeJSON()
	}

	res.Code = ResponseNormal
	res.Messgae = "登录成功"
	res.Data.Username = reply.Username
	res.Data.AccessTicket = reply.AccessTicket
	res.Data.ExpireIn = reply.ExpireIn
	c.Data["json"] = res
	c.ServeJSON()
}

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
	fmt.Println("err>>>", err)
	fmt.Println("replyAccessToken>>>", replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
	}

	if replyAccessToken.UserId == 0 {
		res.Code = ResponseSystemErr
		res.Messgae = "获取用户信息失败"
		c.Data["json"] = res
		c.ServeJSON()
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
	access_token := c.GetString("access_token")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
	}

	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type:"And",
		Key:"access_token",
		Val:access_token,
	})
	var replyAccessToken UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken) // 检测 accessToken
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
	}

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
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
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
	}

	if replyAccessToken.UserId == 0 {
		res.Code = ResponseSystemErr
		res.Messgae = "获取用户信息失败"
		c.Data["json"] = res
		c.ServeJSON()
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
	}

	var updateArgs action.UpdateByIdCond
	updateArgs.Id = []int64{reply.Id}
	updateArgs.UpdateList = append(updateArgs.UpdateList, action.UpdateValue{
		Key:"update_time",
		Val:time.Now().UnixNano() / 1e6,
	})
	updateArgs.UpdateList = append(updateArgs.UpdateList, action.UpdateValue{
		Key:"password",
		Val:password,
	})
	var updateReply action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", updateArgs, &updateReply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "用户信息修改失败！"
		c.Data["json"] = res
		c.ServeJSON()
	}

	res.Code = ResponseNormal
	res.Messgae = "用户信息修改成功！"
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

}