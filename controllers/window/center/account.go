package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/center"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/account"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/utils"
)

// Account API
type AccountController struct {
	beego.Controller
}

// @Title 账户统计接口
// @Description 账户统计接口
// @Success 200 {"code":200,"message":"获取账务统计信息成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @router /statistics [post]
func (c *AccountController) Statistics() {
	type data AccountStatisticsResponse
	accessToken := c.GetString("accessToken")
	currentTimestamp := utils.CurrentTimestamp()
	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "access_token", Val: accessToken },
			action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTimestamp },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "验证访问令牌失效"}
		c.ServeJSON()
		return
	}

	if replyUserPosition.Id == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌失效"}
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
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取账务统计信息失败1"}
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
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取账务统计信息失败2"}
		c.ServeJSON()
		return
	}

	var inAccount, outAccount float64
	for _, accountItem := range replyAccountItemList {
		if accountItem.Amount > 0 {
			inAccount += accountItem.Amount
		} else {
			outAccount += accountItem.Amount
		}
	}

	c.Data["json"] = data{Code: Normal, Message: "获取账务统计信息成功", Data: AccountStatistics{
		Balance:    replyAccount.Balance,
		InAccount:  inAccount,
		OutAccount: outAccount,
	}}
	c.ServeJSON()
	return
}
