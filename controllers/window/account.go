package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/account"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/utils"
	"time"
	"fmt"
)

// USER API
type AccountController struct {
	beego.Controller
}

// @Title 账户统计接口
// @Description 账户统计接口
// @Success 200 {"code":200,"messgae":"获取账务统计信息成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   date     query   string true       "账单日期：2017-07"
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @router /accountstatistics [get]
func (c *AccountController) AccountStatistics() {
	res := AccountStatisticsResponse{}
	access_token := c.GetString("access_token")
	current := c.GetString("date")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
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
		return
	}
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
		return
	}
	var reply Account
	//检测 accessToken
	var accountArgs action.FindByCond
	accountArgs.CondList = append(accountArgs.CondList, action.CondValue{
		Type:  "And",
		Key: "company_id",
		Val:  com_id,
	}, action.CondValue{
		Type:  "And",
		Key: "user_id",
		Val:  user_id,
	}, action.CondValue{
		Type:  "And",
		Key: "user_position_id",
		Val:  user_position_id,
	}, action.CondValue{
		Type:  "And",
		Key: "user_position_type",
		Val:  user_position_type,
	})
	accountArgs.Fileds = []string{"id", "balance"}
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindByCond", accountArgs, &reply)
	fmt.Println("reply>>>>", reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取账务统计信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	account_id := reply.Id
	var AccountItemStatisticsArgs AccountItemStatisticsCond
	AccountItemStatisticsArgs.AccountId = account_id
	var reply2 AccountItemStatisticsCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "AccountItemStatistics", AccountItemStatisticsArgs, &reply2)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取账务统计信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	res.Data.TotalDischarge = reply2.Pay
	res.Data.TotalEntry = reply2.Income

	if (current == "") {
		current = time.Now().Format("2006-01")
	}
	//获取指定时间的月初、下个月初时间戳
	stime := utils.GetMonthStartUnix(current + "-01")
	etime := utils.GetNextMonthStartUnix(current + "-01")
	var AccountItemArgs AccountItemStatisticsCond
	AccountItemArgs.AccountId = account_id
	AccountItemArgs.StartTime = stime
	AccountItemArgs.EndTime = etime
	var AccountItemReply AccountItemStatisticsCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "AccountItemStatistics", AccountItemArgs, &AccountItemReply)

	res.Code = ResponseNormal
	res.Messgae = "获取账务统计信息成功"
	res.Data.Balance = reply.Balance
	res.Data.MonthPay = AccountItemReply.Pay
	res.Data.MonthIncome = AccountItemReply.Income
	c.Data["json"] = res
	c.ServeJSON()

}
