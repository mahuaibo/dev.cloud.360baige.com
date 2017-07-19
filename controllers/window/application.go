package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window"
	. "dev.model.360baige.com/models/user"
	//. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/models/application"
	. "dev.model.360baige.com/models/company"
	"time"
)

// APPLICATION API
type ApplicationController struct {
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
func (c *ApplicationController) List() {
	res := ApplicationListResponse{}
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
			var reply ApplicationListPaginator
			var cond1 []CondValue
			cond1 = append(cond1, CondValue{
				Type:  "And",
				Exprs: "company_id",
				Args:  com_id,
			})
			cond1 = append(cond1, CondValue{
				Type:  "And",
				Exprs: "user_id",
				Args:  user_id,
			})
			cond1 = append(cond1, CondValue{
				Type:  "And",
				Exprs: "user_position_id",
				Args:  user_position_id,
			})
			cond1 = append(cond1, CondValue{
				Type:  "And",
				Exprs: "user_position_type",
				Args:  user_position_type,
			})
			appname := c.GetString("name")
			if appname != "" {
				cond1 = append(cond1, CondValue{
					Type:  "And",
					Exprs: "name__icontains",
					Args:  appname,
				})

			}
			currentPage, _ := c.GetInt64("current")
			pageSize, _ := c.GetInt64("page_size")
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "PageBy", &ApplicationListPaginator{
				Cond:     cond1,
				Cols:     []string{"id", "create_time", "name", "image", "status", "application_tpl_id" },
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
				//获取应用其他信息tpl
				var idarg []int64
				idmap := make(map[int64]int64)
				for _, value := range reply.List {
					idmap[value.ApplicationTplId] = value.ApplicationTplId
				}
				for _, value := range idmap {
					idarg = append(idarg, value)
				}

				var cond2 []CondValue
				cond2 = append(cond2, CondValue{
					Type:  "And",
					Exprs: "id__in",
					Args:  idarg,
				})
				var replyApplicationTpl *ApplicationTplListPaginator
				err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "ListAll", &ApplicationTplListPaginator{
					Cond:     cond2,
					Cols:     []string{"id", "name", "image", "status", "site"},
					OrderBy:  []string{"id"},
					PageSize: -1,
				}, &replyApplicationTpl)

				var resData []ApplicationValue
				//循环赋值
				if err != nil {
					res.Code = ResponseSystemErr
					res.Messgae = "获取应用失败"
					c.Data["json"] = res
					c.ServeJSON()
				}
				applicationTplByIds := make(map[int64]ApplicationTpl)
				for _, value := range replyApplicationTpl.List {
					applicationTplByIds[value.Id] = value
				}
				for _, value := range reply.List {
					re := time.Unix(value.CreateTime/1000, 0).Format("2006-01-02")
					var rename, reimage, restatus string
					if value.Name == "" {
						if applicationTplByIds[value.ApplicationTplId].Name != "" {
							rename = applicationTplByIds[value.ApplicationTplId].Name
						}
					} else {
						rename = value.Name
					}
					if value.Image == "" {
						if applicationTplByIds[value.ApplicationTplId].Image != "" {
							reimage = applicationTplByIds[value.ApplicationTplId].Image
						}
					} else {
						reimage = value.Image
					}
					if value.Status == 0 {
						restatus = "启用"
					} else if value.Status == 1 {
						restatus = "停用"
					} else {
						restatus = "退订"
					}
					resData = append(resData, ApplicationValue{
						Id:         value.Id,
						CreateTime: re,
						Name:       rename,
						Image:      reimage,
						Status:     restatus,
						Site:       applicationTplByIds[value.ApplicationTplId].Site,
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
func (c *ApplicationController) Detail() {
	res := ApplicationDetailResponse{}
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
		//company_id、user_id、user_position_id、user_position_type
		ap_id, _ := c.GetInt64("id")

		if ap_id == 0 {
			res.Code = ResponseSystemErr
			res.Messgae = "获取信息失败"
			c.Data["json"] = res
			c.ServeJSON()
		} else {
			var reply Application
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindById", &Application{
				Id: ap_id,
			}, &reply)

			if err != nil {
				res.Code = ResponseSystemErr
				res.Messgae = "获取应用信息失败"
				c.Data["json"] = res
				c.ServeJSON()
			} else {
				if reply.ApplicationTplId == 0 {
					res.Code = ResponseSystemErr
					res.Messgae = "获取应用信息失败"
					c.Data["json"] = res
					c.ServeJSON()
				}
				//获取应用其他信息tpl
				var replyApplicationTpl ApplicationTpl
				err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindById", &ApplicationTpl{
					Id: reply.ApplicationTplId,
				}, &replyApplicationTpl)
				if err != nil {
					res.Code = ResponseSystemErr
					res.Messgae = "获取应用信息失败"
					c.Data["json"] = res
					c.ServeJSON()
				}
				re := time.Unix(reply.CreateTime/1000, 0).Format("2006-01-02")
				var rename, reimage string
				if reply.Name == "" {
					if replyApplicationTpl.Name != "" {
						rename = replyApplicationTpl.Name
					}
				} else {
					rename = reply.Name
				}
				if reply.Image == "" {
					if replyApplicationTpl.Image != "" {
						reimage = replyApplicationTpl.Image
					}
				} else {
					reimage = reply.Image
				}
				//开发者
				var replyUser User
				err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByIdNoStatus", &User{
					Id: replyApplicationTpl.UserId,
				}, &replyUser)
				var username, cname string
				if err == nil {
					username = replyUser.Username
				}
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
				res.Data.Name = rename
				res.Data.Image = reimage
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
}
func GetPayTypeName(ptype int8) string {
	// 0:限免 1:永久免费 2:1次性收费 3:周期收费tpl
	var rPtype string
	switch ptype {
	case 0:
		rPtype = "限免"
	case 1:
		rPtype = "永久免费"
	case 2:
		rPtype = "1次性收费"
	case 3:
		rPtype = "周期收费"
	}
	return rPtype
}

func GetPayCycleName(ptype int8) string {
	// 0无1月2季3半年4年tpl
	var rPtype string
	switch ptype {
	case 0:
		rPtype = "无"
	case 1:
		rPtype = "月"
	case 2:
		rPtype = "季"
	case 3:
		rPtype = "半年"
	case 4:
		rPtype = "年"
	}
	return rPtype
}

// @Title 应用修改状态接口
// @Description 应用修改状态接口
// @Success 200 {"code":200,"messgae":"获取应用修改状态成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Param   status     query   string true       " 0 上架 1 下架 "
// @Failure 400 {"code":400,"message":"获取应用修改状态失败"}
// @router /modifystatus [get]
func (c *ApplicationController) ModifyStatus() {
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
		var reply Application
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindById", &Application{
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
			reply.UpdateTime = timestamp
			reply.Status = status
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "UpdateById", reply, nil)
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
