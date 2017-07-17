package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.cloud.360baige.com/utils"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/http/window"
	"fmt"
	"time"
	"strconv"
	"dev.model.360baige.com/models/batch"
	"dev.action.360baige.com/actions/user"
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

	//1 检测数据完整性
	if (username == "" || password == "") {
		res.Code = ResponseSystemErr
		res.Messgae = "登录失败"
		c.Data["json"] = res
		c.ServeJSON()
	}

	//2 判断username 类型 属于 百鸽账号、邮箱、手机号码中的哪一种？
	username_type, _ := utils.DetermineStringType(username)
	fmt.Println("Type1:", username_type)
	var reply User
	var err error

	switch username_type {
	case 1:
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByUsername", &User{
			Username: username,
		}, &reply)
	case 2:
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByEmail", &User{
			Email: username,
		}, &reply)
	case 3:
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByPhone", &User{
			Phone: username,
		}, &reply)
	}

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "登录失败"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "登录成功"
		timestamp := time.Now().UnixNano() / 1e6
		newAccessTicket := reply.Username + strconv.FormatInt(timestamp, 10)
		reply.UpdateTime = timestamp
		reply.AccessTicket = newAccessTicket //更新船票 应该判断时效，再做更新
		err2 := client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", reply, nil)
		if err2 != nil {
			res.Code = ResponseSystemErr
			res.Messgae = "登录失败"
			c.Data["json"] = res
			c.ServeJSON()
		} else {
			res.Data.Username = reply.Username
			res.Data.AccessTicket = reply.AccessTicket
			res.Data.ExpireIn = reply.ExpireIn
			c.Data["json"] = res
			c.ServeJSON()
		}

	}

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
	//检测 accessToken
	var replyAccessToken UserPosition
	var err error
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByAccessToken", &UserPosition{
		AccessToken: access_token,
	}, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		user_id := replyAccessToken.UserId
		if user_id == 0 {
			res.Code = ResponseSystemErr
			res.Messgae = "获取用户信息失败"
			c.Data["json"] = res
			c.ServeJSON()
		} else {
			var reply User
			err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", &User{
				Id: user_id,
			}, &reply)

			if err != nil {
				res.Code = ResponseSystemErr
				res.Messgae = "获取用户信息失败"
				c.Data["json"] = res
				c.ServeJSON()
			}
			res.Code = ResponseNormal
			res.Messgae = "获取用户信息成功"
			res.Data.Username = reply.Username
			res.Data.Email = reply.Email
			res.Data.Phone = reply.Phone
			c.Data["json"] = res
			c.ServeJSON()
		}

	}

}

// @Title 用户退出接口
// @Description 用户退出接口
// @Success 200 {"code":200,"messgae":"用户退出成功","data":{"access_ticket":"xxxx","expire_in":0}}// @Param   password query   string true       "密码"
// @Param   access_token query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户退出失败"}
// @router /logout [post]
func (c *UserController) Logout() {

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
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		//检测 accessToken
		var replyAccessToken UserPosition
		var err error
		err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByAccessToken", &UserPosition{
			AccessToken: access_token,
		}, &replyAccessToken)
		fmt.Println(err)
		fmt.Println(replyAccessToken)
		if err != nil {
			res.Code = ResponseLogicErr
			res.Messgae = "访问令牌失效"
			c.Data["json"] = res
			c.ServeJSON()
		} else {
			u_id := replyAccessToken.UserId
			if u_id == 0 {
				res.Code = ResponseSystemErr
				res.Messgae = "获取用户信息失败"
				c.Data["json"] = res
				c.ServeJSON()
			} else {
				var reply User
				err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", &User{
					Id: u_id,
				}, &reply)

				if err != nil {
					res.Code = ResponseSystemErr
					res.Messgae = "获取用户信息失败"
					c.Data["json"] = res
					c.ServeJSON()
				} else {
					password := c.GetString("password")
					timestamp := time.Now().UnixNano() / 1e6
					reply.UpdateTime = timestamp
					reply.Password = password
					err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", reply, nil)
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
			}
		}
	}
}

// @Title 用户信息修改接口
// @Description 用户信息修改接口
// @Success 200 {"code":200,"messgae":"用户修改密码成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户修改密码失败"}
// @router /modify [get]
func (c *UserController) Modify() {
	res := ModifyPasswordResponse{}
	var args user.UpdateByCond
	//args.Id = 1
	//user := User{
	//	CreateTime: 1499760680777,
	//	UpdateTime: 1499760680777,
	//	Username: "bbbbbbbbbbb",
	//	Password: "18910110013",
	//	Email: "ashdiasjf@aseq.com",
	//	Phone: "18919000919",
	//}
	//args = append(args, user)
	//var reply User
	var reply batch.BackNumm
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateByCond", args, &reply)
	fmt.Println("err>>>>>", err)
	fmt.Println("reply>>>>>", reply)
	res.Code = ResponseNormal
	res.Messgae = ""
	c.Data["json"] = res
	c.ServeJSON()
}