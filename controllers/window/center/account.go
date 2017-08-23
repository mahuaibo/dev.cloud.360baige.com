package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/center"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/account"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/utils"
	"time"
)

// Account API
type AccountController struct {
	beego.Controller
}

// @Title 账户统计接口
// @Description 账户统计接口
// @Success 200 {"code":200,"message":"获取账务统计信息成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   date     query   string true       "账单日期：2017-07"
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @router /accountStatistics [get]
func (c *AccountController) AccountStatistics() {
	type data AccountStatisticsResponse
	accessToken := c.GetString("accessToken")
	current := c.GetString("date", time.Now().Format("2006-01"))
	if accessToken == "" {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "access_token", Val: accessToken },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)

	if err != nil {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "验证访问令牌失效"}
		c.ServeJSON()
		return
	}

	if replyUserPosition.Id == 0 {
		c.Data["json"] = data{Code: ResponseLogicErr, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	var replyAccount account.Account
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId },
			action.CondValue{Type: "And", Key: "user_id", Val: replyUserPosition.UserId },
			action.CondValue{Type: "And", Key: "user_position_id", Val: replyUserPosition.Id },
			action.CondValue{Type: "And", Key: "user_position_type", Val: replyUserPosition.Type },
		},
		Fileds: []string{"id", "balance"},
	}, &replyAccount)

	if err != nil {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取账务统计信息失败"}
		c.ServeJSON()
		return
	}

	var reply2 AccountItemStatisticsCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "AccountItemStatistics", AccountItemStatisticsCond{
		AccountId: replyAccount.Id,
	}, &reply2)

	if err != nil {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取账务统计信息失败"}
		c.ServeJSON()
		return
	}

	var AccountItemReply AccountItemStatisticsCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "AccountItemStatistics", AccountItemStatisticsCond{
		AccountId: reply2.AccountId,
		StartTime: utils.GetMonthStartUnix(current + "-01"),
		EndTime:   utils.GetNextMonthStartUnix(current + "-01"),
	}, &AccountItemReply)

	c.Data["json"] = data{Code: ResponseNormal, Message: "获取账务统计信息成功", Data: AccountStatistics{
		Balance:        replyAccount.Balance,
		MonthPay:       AccountItemReply.Pay,
		MonthIncome:    AccountItemReply.Income,
		TotalDischarge: reply2.Pay,
		TotalEntry:     reply2.Income,
	}}
	c.ServeJSON()
	return

}
