package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.cloud.360baige.com/utils"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/http/mobile/center"
	"fmt"
	"time"
	"strconv"
	"dev.model.360baige.com/action"
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
	//res := Response{}
	username := c.GetString("username")
	password := c.GetString("password")

	//1 检测数据完整性
	if (username == "" || password == "") {
		res.Code = ResponseSystemErr
		res.Message = "登录失败"
		c.Data["json"] = res
		c.ServeJSON()
	}

	//2 判断username 类型 属于 百鸽账号、邮箱、手机号码中的哪一种？
	username_type, _ := utils.DetermineStringType(username)
	fmt.Println("Type1:", username_type)
	var reply User
	var err error
	var args action.FindByCond
	args.Fileds = []string{"id", "username", "access_ticket", "expire_in"}
	switch username_type {
	case "1":
		args.CondList = append(args.CondList, action.CondValue{
			Type: "And",
			Key:  "username",
			Val:  username,
		})
	case "2":
		args.CondList = append(args.CondList, action.CondValue{
			Type: "And",
			Key:  "email",
			Val:  username,
		})
	case "3":
		args.CondList = append(args.CondList, action.CondValue{
			Type: "And",
			Key:  "phone",
			Val:  username,
		})
	}
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", args, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "登录失败"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Message = "登录成功"
		timestamp := time.Now().UnixNano() / 1e6
		newAccessTicket := reply.Username + strconv.FormatInt(timestamp, 10)
		var updateArgs []action.UpdateValue
		updateArgs = append(updateArgs, action.UpdateValue{
			Key: "update_time",
			Val: timestamp,
		})
		updateArgs = append(updateArgs, action.UpdateValue{
			Key: "access_ticket",
			Val: newAccessTicket,
		}) // 更新船票 应该判断时效，再做更新
		// 更新ExpireIn
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", &action.UpdateByIdCond{
			Id:         []int64{reply.Id},
			UpdateList: updateArgs,
		}, nil)
		if err != nil {
			res.Code = ResponseSystemErr
			res.Message = "登录失败"
			c.Data["json"] = res
			c.ServeJSON()
		} else {
			res.Data.Username = reply.Username
			res.Data.AccessTicket = newAccessTicket
			res.Data.ExpireIn = reply.ExpireIn
			c.Data["json"] = res
			c.ServeJSON()
		}

	}

}

// @Title 用户退出接口
// @Description 用户退出接口
// @Success 200 {"code":200,"message":"用户退出成功","data":{"access_ticket":"xxxx","expire_in":0}}// @Param   password query   string true       "密码"
// @Param   access_token query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"用户退出失败"}
// @router /logout [post]
func (c *UserController) Logout() {

}

// @Title 用户修改密码接口
// @Description 用户修改密码接口
// @Success 200 {"code":200,"message":"用户修改密码成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token query   string true       "访问令牌"
// @Param   password query   string true       "password"
// @Failure 400 {"code":400,"message":"用户修改密码失败"}
// @router /modifypassword [post]
func (c *UserController) ModifyPassword() {
	res := ModifyPasswordResponse{}
	access_token := c.GetString("access_token")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		//检测 accessToken
		var args action.FindByCond
		args.CondList = append(args.CondList, action.CondValue{
			Type:  "And",
			Key: "accessToken",
			Val:  access_token,
		})
		args.Fileds = []string{"id", "user_id", "company_id", "type"}
		var replyAccessToken UserPosition
		err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
		if err != nil {
			res.Code = ResponseLogicErr
			res.Message = "访问令牌失效"
			c.Data["json"] = res
			c.ServeJSON()
		} else {
			u_id := replyAccessToken.UserId
			if u_id == 0 {
				res.Code = ResponseSystemErr
				res.Message = "获取用户信息失败"
				c.Data["json"] = res
				c.ServeJSON()
			} else {
				var reply User
				err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", &User{
					Id: u_id,
				}, &reply)

				if err != nil {
					res.Code = ResponseSystemErr
					res.Message = "获取用户信息失败"
					c.Data["json"] = res
					c.ServeJSON()
				} else {
					password := c.GetString("password")
					timestamp := time.Now().UnixNano() / 1e6
					var updateArgs []action.UpdateValue
					updateArgs = append(updateArgs, action.UpdateValue{
						Key: "update_time",
						Val: timestamp,
					})
					updateArgs = append(updateArgs, action.UpdateValue{
						Key: "password",
						Val: password,
					})
					err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "UpdateById", &action.UpdateByIdCond{
						Id:         []int64{reply.Id},
						UpdateList: updateArgs,
					}, nil)
					if err != nil {
						res.Code = ResponseSystemErr
						res.Message = "用户信息修改失败！"
						c.Data["json"] = res
						c.ServeJSON()
					}
					res.Code = ResponseNormal
					res.Message = "用户信息修改成功！"
					c.Data["json"] = res
					c.ServeJSON()
				}
			}
		}
	}
}

