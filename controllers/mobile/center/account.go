package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/mobile/center"
	"dev.model.360baige.com/models/account"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/utils"
)

// Account API
type AccountController struct {
	beego.Controller
}

// @Title 账户统计
// @Description 账户统计
// @Success 200 {"code":200,"message":"获取账务统计信息成功"}
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @Param   accessToken     query   string true       "访问令牌"
// @router /statistics [post]
func (c *AccountController) Statistics() {
	type data AccountStatisticsResponse
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

	var replyAccount account.Account
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindByCond", &action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId },
			action.CondValue{Type: "And", Key: "user_id", Val: replyUserPosition.UserId },
			action.CondValue{Type: "And", Key: "user_position_id", Val: replyUserPosition.Id },
			action.CondValue{Type: "And", Key: "user_position_type", Val: replyUserPosition.Type },
		},
		Fileds: []string{"id", "balance"},
	}, &replyAccount)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	var replyAccountItemList []account.AccountItem
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "ListByCond", &action.ListByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "account_id", Val: replyAccount.Id },
			action.CondValue{Type: "And", Key: "status__gt", Val: -1 },
		},
		Cols: []string{"amount", "balance"},
	}, &replyAccountItemList)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	var inAccount, outAccount int64
	for _, accountItem := range replyAccountItemList {
		if accountItem.Amount > 0 {
			inAccount += accountItem.Amount
		} else {
			outAccount += accountItem.Amount
		}
	}

	c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: AccountStatistics{
		Balance:    replyAccount.Balance,
		InAccount:  inAccount,
		OutAccount: outAccount,
	}}
	c.ServeJSON()
	return
}
