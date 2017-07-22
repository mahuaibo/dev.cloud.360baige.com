package mobile

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/mobile"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/account"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/utils"
	"time"
)

// USER API
type AccountController struct {
	beego.Controller
}

// @Title 账务统计接口
// @Description 账务统计接口
// @Success 200 {"code":200,"messgae":"获取账务统计信息成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   date     query   string true       "账单日期：2017-07"
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @router /accountstatistics [get]
func (c *AccountController) AccountStatistics() {
	res := AccountStatisticsResponse{}
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
		Type:  "And",
		Key: "accessToken",
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
			res.Messgae = "获取账务统计信息失败"
			c.Data["json"] = res
			c.ServeJSON()
		} else {
			var reply Account
			//检测 accessToken
			var args2 action.FindByCond
			args2.CondList = append(args2.CondList, action.CondValue{
				Type:  "And",
				Key: "company_id",
				Val:  com_id,
			})
			args2.CondList = append(args2.CondList, action.CondValue{
				Type:  "And",
				Key: "user_id",
				Val:  user_id,
			})
			args2.CondList = append(args2.CondList, action.CondValue{
				Type:  "And",
				Key: "user_position_id",
				Val:  user_position_id,
			})
			args2.CondList = append(args2.CondList, action.CondValue{
				Type:  "And",
				Key: "user_position_type",
				Val:  user_position_type,
			})
			args2.Fileds = []string{"id", "balance"}
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindByCond", args2, &reply)
			if err != nil {
				res.Code = ResponseSystemErr
				res.Messgae = "获取账务统计信息失败"
				c.Data["json"] = res
				c.ServeJSON()
			} else {
				account_id := reply.Id
				var reply2 AccountItemStatisticsCond
				err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "AccountItemStatistics", &AccountItemStatisticsCond{
					AccountId: account_id,
				}, &reply2)
				if err != nil {
					res.Code = ResponseSystemErr
					res.Messgae = "获取账务统计信息失败"
					c.Data["json"] = res
					c.ServeJSON()
				} else {
					res.Data.TotalDischarge = reply2.Pay
					res.Data.TotalEntry = reply2.Income
					var reply3 AccountItemStatisticsCond
					var stime, etime int64
					current := c.GetString("date")
					if (current == "") {
						current = time.Now().Format("2006-01")
					}
					//获取指定时间的月初、下个月初时间戳
					stime = utils.GetMonthStartUnix(current + "-01")
					etime = utils.GetNextMonthStartUnix(current + "-01")
					err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "AccountItemStatistics", &AccountItemStatisticsCond{
						AccountId: account_id,
						StartTime: stime,
						EndTime:   etime,
					}, &reply3)
					res.Code = ResponseNormal
					res.Messgae = "获取账务统计信息成功"
					res.Data.Balance = reply.Balance
					res.Data.MonthPay = reply3.Pay
					res.Data.MonthIncome = reply3.Income
					c.Data["json"] = res
					c.ServeJSON()
				}
			}
		}
	}
}
