package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/models/account"
	"dev.cloud.360baige.com/utils"
	"time"
	"fmt"
)

type AccountItemController struct {
	beego.Controller
}

// @Title 账务列表接口
// @Description 账务列表接口
// @Success 200 {"code":200,"messgae":"获取账务列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   date     query   string true       "账单日期：2017-07"
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @router /list [get]
func (c *AccountItemController) List() {
	res := AccountItemListResponse{}
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
			var reply Account
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindByUserPos", &Account{
				CompanyId:        com_id,
				UserId:           user_id,
				UserPositionId:   user_position_id,
				UserPositionType: user_position_type,
			}, &reply)
			fmt.Println("acccount-id:", reply.Id)
			if err != nil {
				res.Code = ResponseSystemErr
				res.Messgae = "获取账户信息失败"
				c.Data["json"] = res
				c.ServeJSON()
			} else {
				account_id := reply.Id
				var reply2 AccountItemListPaginator
				var cond1 []CondValue
				cond1 = append(cond1, CondValue{
					Type:  "And",
					Exprs: "account_id",
					Args:  account_id,
				})
				var stime, etime int64
				current := c.GetString("date")
				if (current == "") {
					current = time.Now().Format("2006-01")
				}
				//获取指定时间的月初、下个月初时间戳
				stime = utils.GetMonthStartUnix(current + "-01")
				etime = utils.GetNextMonthStartUnix(current + "-01")
				cond1 = append(cond1, CondValue{
					Type:  "And",
					Exprs: "create_time__gte",
					Args:  stime,
				})
				cond1 = append(cond1, CondValue{
					Type:  "And",
					Exprs: "create_time__lt",
					Args:  etime,
				})
				err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "PageBy", &AccountItemListPaginator{
					Cond:     cond1,
					Cols:     []string{"id", "create_time", "amount", },
					OrderBy:  []string{"id"},
					PageSize: 0,
				}, &reply2)
				if err != nil {
					res.Code = ResponseSystemErr
					res.Messgae = "获取账务信息失败"
					c.Data["json"] = res
					c.ServeJSON()
				} else {
					//res.Data.List = reply2.List
					//List 循环赋值
					var aType string
					for _, value := range reply2.List {
						if (value.Amount < 0) {
							aType = "收入"
						} else {
							aType = "支出"
						}
						res.Data.List = append(res.Data.List, AccountItemListValue{
							Id:         value.Id,
							CreateTime: value.CreateTime,
							Amount:     value.Amount,
							AmountType: aType,
						})
					}
					res.Data.Total = reply2.Total
					res.Data.Current = reply2.Current
					res.Data.CurrentSize = reply2.CurrentSize
					res.Data.OrderBy = reply2.OrderBy
					res.Data.PageSize = reply2.PageSize
					res.Code = ResponseNormal
					res.Messgae = "获取账务统计信息成功"
					c.Data["json"] = res
					c.ServeJSON()
				}
			}
		}
	}
}
