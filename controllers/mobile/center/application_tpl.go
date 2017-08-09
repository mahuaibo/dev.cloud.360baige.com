package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/mobile/center"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/application"
	. "dev.model.360baige.com/models/company"
	"time"
	"fmt"
	"dev.model.360baige.com/action"
	"encoding/json"
)

// APPLICATIONTPL API
type ApplicationTplController struct {
	beego.Controller
}

// @Title 应用列表接口
// @Description 应用列表接口
// @Success 200 {"code":200,"messgae":"获取应用列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   current     query   string true       "当前页"
// @Param   page_size     query   string true       "每页数量"
// @Param   name     query   string true       "搜索名称"
// @Failure 400 {"code":400,"message":"获取应用信息失败"}
// @router /list [get]
func (c *ApplicationTplController) List() {
	res := ApplicationTplListResponse{}
	access_token := c.GetString("access_token")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
	}
	//检测 accessToken
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type: "And",
		Key:  "accessToken",
		Val:  access_token,
	})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		//company_id、user_id、user_position_id、user_position_type
		com_id := replyAccessToken.CompanyId
		user_id := replyAccessToken.UserId
		user_position_id := replyAccessToken.Id
		user_position_type := replyAccessToken.Type
		if com_id == 0 || user_id == 0 || user_position_id == 0 {
			res.Code = ResponseSystemErr
			res.Messgae = "获取信息失败"
			c.Data["json"] = res
			c.ServeJSON()
		} else {

			var reply action.PageByCond
			var cond1 []action.CondValue
			appname := c.GetString("name")
			if appname != "" {
				cond1 = append(cond1, action.CondValue{
					Type: "And",
					Key:  "name__icontains",
					Val:  appname,
				})

			}
			cond1 = append(cond1, action.CondValue{
				Type: "And",
				Key:  "status__gt",
				Val:  -1,
			})
			currentPage, _ := c.GetInt64("current")
			pageSize, _ := c.GetInt64("page_size")
			err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "PageByCond", &action.PageByCond{
				CondList: cond1,
				Cols:     []string{"id", "name", "image", "status", "desc"},
				OrderBy:  []string{"id"},
				PageSize: pageSize,
				Current:  currentPage,
			}, &reply)
			if err != nil {
				res.Code = ResponseSystemErr
				res.Messgae = "获取应用信息失败"
				c.Data["json"] = res
				c.ServeJSON()
			} else {
				var replyApplication []Application
				var cond2 []action.CondValue
				cond2 = append(cond2, action.CondValue{
					Type: "And",
					Key:  "company_id",
					Val:  com_id,
				})
				cond2 = append(cond2, action.CondValue{
					Type: "And",
					Key:  "user_id",
					Val:  user_id,
				})
				cond2 = append(cond2, action.CondValue{
					Type: "And",
					Key:  "user_position_id",
					Val:  user_position_id,
				})
				cond2 = append(cond2, action.CondValue{
					Type: "And",
					Key:  "user_position_type",
					Val:  user_position_type,
				})
				cond2 = append(cond2, action.CondValue{
					Type: "And",
					Key:  "status__gt",
					Val:  0,
				})
				err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "ListByCond", &action.ListByCond{
					CondList: cond2,
					Cols:     []string{"id", "application_tpl_id" },
					OrderBy:  []string{"id"},
					PageSize: -1,
				}, &replyApplication)
				idmap := make(map[int64]int64)
				for _, value := range replyApplication {
					idmap[value.ApplicationTplId] = value.ApplicationTplId
				}

				var resData []ApplicationTplValue
				replyList := []ApplicationTpl{}
				err = json.Unmarshal([]byte(reply.Json), &replyList)
				//循环赋值
				for _, value := range replyList {
					var restatus int8
					if idmap[value.Id] > 0 {
						restatus = 1
					} else {
						restatus = 0
					}
					resData = append(resData, ApplicationTplValue{
						Id:                 value.Id,
						Name:               value.Name,
						Image:              value.Image,
						SubscriptionStatus: restatus,
						Desc:               value.Desc,
					})

				}
				res.Code = ResponseNormal
				res.Messgae = "获取应用成功"
				res.Data.Total = reply.Total
				res.Data.Current = currentPage
				res.Data.CurrentSize = reply.CurrentSize
				res.Data.OrderBy = reply.OrderBy
				res.Data.PageSize = pageSize
				res.Data.Name = appname
				res.Data.List = resData
				c.Data["json"] = res
				c.ServeJSON()
			}

		}
	}
}

// @Title 应用详情接口
// @Description 应用详情接口
// @Success 200 {"code":200,"messgae":"获取应用详情成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"获取应用详情失败"}
// @router /detail [get]
func (c *ApplicationTplController) Detail() {
	res := ApplicationDetailResponse{}
	access_token := c.GetString("access_token")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
	}
	//检测 accessToken
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type: "And",
		Key:  "accessToken",
		Val:  access_token,
	})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		ap_id, _ := c.GetInt64("id")
		if ap_id == 0 {
			res.Code = ResponseSystemErr
			res.Messgae = "获取信息失败"
			c.Data["json"] = res
			c.ServeJSON()
		} else {
			//获取应用信息tpl
			var replyApplicationTpl ApplicationTpl
			err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindById", &ApplicationTpl{
				Id: ap_id,
			}, &replyApplicationTpl)
			if err != nil {
				res.Code = ResponseSystemErr
				res.Messgae = "获取应用信息失败"
				c.Data["json"] = res
				c.ServeJSON()
			}
			re := time.Unix(replyApplicationTpl.CreateTime/1000, 0).Format("2006-01-02")
			//开发者
			var replyUser User
			err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindById", &User{
				Id: replyApplicationTpl.UserId,
			}, &replyUser)
			var username, cname string
			if err == nil {
				username = replyUser.Username
			}
			fmt.Println(replyApplicationTpl.UserId)
			fmt.Println(replyUser)
			//开发公司
			var replyCompany Company
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "FindById", &Company{
				Id: replyApplicationTpl.CompanyId,
			}, &replyCompany)
			if err == nil {
				cname = replyCompany.Name
			}
			res.Code = ResponseNormal
			res.Messgae = "获取应用成功"
			res.Data.CreateTime = re
			res.Data.Name = replyApplicationTpl.Name
			res.Data.Image = replyApplicationTpl.Image
			res.Data.Desc = replyApplicationTpl.Desc
			res.Data.Price = replyApplicationTpl.Price
			res.Data.PayType = GetPayTypeName(replyApplicationTpl.PayType)
			res.Data.PayCycle = GetPayCycleName(replyApplicationTpl.PayCycle)
			res.Data.CompanyName = cname
			res.Data.UserName = username
			c.Data["json"] = res
			c.ServeJSON()

		}

	}
}

// @Title 应用详情接口
// @Description 应用详情接口
// @Success 200 {"code":200,"messgae":"获取应用详情成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"获取应用详情失败"}
// @router /subscription [get]
func (c *ApplicationTplController) Subscription() {
	res := ModifyApplicationTplStatusResponse{}
	access_token := c.GetString("access_token")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
	}
	//检测 accessToken
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type: "And",
		Key:  "accessToken",
		Val:  access_token,
	})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		ap_id, _ := c.GetInt64("id")
		//company_id、user_id、user_position_id、user_position_type
		com_id := replyAccessToken.CompanyId
		user_id := replyAccessToken.UserId
		user_position_id := replyAccessToken.Id
		user_position_type := replyAccessToken.Type
		if ap_id == 0 || com_id == 0 || user_id == 0 || user_position_id == 0 {
			res.Code = ResponseSystemErr
			res.Messgae = "获取应用信息失败"
			c.Data["json"] = res
			c.ServeJSON()
		} else {
			//判断此应用是否订阅过
			var replyApplication Application
			var cond2 []action.CondValue
			cond2 = append(cond2, action.CondValue{
				Type: "And",
				Key:  "company_id",
				Val:  com_id,
			})
			cond2 = append(cond2, action.CondValue{
				Type: "And",
				Key:  "user_id",
				Val:  user_id,
			})
			cond2 = append(cond2, action.CondValue{
				Type: "And",
				Key:  "user_position_id",
				Val:  user_position_id,
			})
			cond2 = append(cond2, action.CondValue{
				Type: "And",
				Key:  "user_position_type",
				Val:  user_position_type,
			})
			cond2 = append(cond2, action.CondValue{
				Type: "And",
				Key:  "application_tpl_id",
				Val:  ap_id,
			})
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindByCond", &action.FindByCond{
				CondList: cond2,
				Fileds:   []string{"id", "application_tpl_id" },
			}, &replyApplication)
			if err == nil {
				res.Code = ResponseSystemErr
				res.Messgae = "此应用已经订阅过"
				c.Data["json"] = res
				c.ServeJSON()
			} else {
				var reply ApplicationTpl
				err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindById", &ApplicationTpl{
					Id: ap_id,
				}, &reply)
				if err != nil {
					res.Code = ResponseSystemErr
					res.Messgae = "获取应用信息失败"
					c.Data["json"] = res
					c.ServeJSON()
				} else {
					if reply.Status > 0 {
						res.Code = ResponseSystemErr
						res.Messgae = "此应用已经下架"
						c.Data["json"] = res
						c.ServeJSON()
					} else {
						var addReply Application
						addReply.CreateTime = time.Now().UnixNano() / 1e6
						addReply.UpdateTime = time.Now().UnixNano() / 1e6
						addReply.CompanyId = com_id
						addReply.UserId = user_id
						addReply.UserPositionId = user_position_id
						addReply.UserPositionType = user_position_type
						addReply.ApplicationTplId = reply.Id // 应用ID
						addReply.Name = reply.Name           // 名称
						addReply.Image = reply.Image         // 图片链接
						addReply.Status = 1                  // 状态
						addReply.StartTime = 1               // 开始时间
						addReply.EndTime = 1                 // 结束时间

						err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "Add", addReply, &addReply)
						if err != nil {
							res.Code = ResponseSystemErr
							res.Messgae = "应用订阅失败！"
							c.Data["json"] = res
							c.ServeJSON()
						}
						res.Code = ResponseNormal
						res.Messgae = "应用订阅成功！"
						res.Data.ApplicationTplId = reply.Id
						res.Data.AppId = addReply.Id
						c.Data["json"] = res
						c.ServeJSON()
					}

				}
			}
		}

	}

}

// @Title 应用修改状态接口
// @Description 应用修改状态接口
// @Success 200 {"code":200,"messgae":"获取应用修改状态成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Param   status     query   string true      " 0 上架 1 下架 "
// @Failure 400 {"code":400,"message":"获取应用修改状态失败"}
// @router /modifystatus [get]
func (c *ApplicationTplController) ModifyStatus() {
	res := ModifyApplicationStatusResponse{}
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
		ap_id, _ := c.GetInt64("id")
		var reply ApplicationTpl
		err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindById", &ApplicationTpl{
			Id: ap_id,
		}, &reply)
		if err != nil {
			res.Code = ResponseSystemErr
			res.Messgae = "获取应用信息失败"
			c.Data["json"] = res
			c.ServeJSON()
		} else {
			status, _ := c.GetInt8("status")
			timestamp := time.Now().UnixNano() / 1e6
			var updateArgs []action.UpdateValue
			updateArgs = append(updateArgs, action.UpdateValue{
				Key: "update_time",
				Val:  timestamp,
			})
			updateArgs = append(updateArgs, action.UpdateValue{
				Key: "status",
				Val:  status,
			})
			err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "UpdateById", &action.UpdateByIdCond{
				Id: []int64{ap_id},
				UpdateList:updateArgs,
			}, nil)
			if err != nil {
				res.Code = ResponseSystemErr
				res.Messgae = "应用信息修改失败！"
				c.Data["json"] = res
				c.ServeJSON()
			}
			res.Code = ResponseNormal
			res.Messgae = "应用信息修改成功！"
			c.Data["json"] = res
			c.ServeJSON()
		}
	}

}
