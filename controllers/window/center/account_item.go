package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/center"
	"dev.model.360baige.com/models/account"
	"dev.model.360baige.com/models/company"
	"time"
	"dev.model.360baige.com/action"
	"encoding/json"
	"dev.cloud.360baige.com/utils"
	"dev.cloud.360baige.com/log"
)

type AccountItemController struct {
	beego.Controller
}

// @Title 账户列表接口
// @Description 账户列表接口
// @Success 200 {"code":200,"message":"获取账务列表成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   cycleType     query   string true       "周期类型 1:近1月 2:近3月 3:近6月 4:近1年"
// @Param   current     query   string true       "当前页"
// @Param   pageSize     query   string true       "每页数量"
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @router /list [post]
func (c *AccountItemController) List() {
	type data AccountItemListResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	cycleType := c.GetString("cycleType", "4")
	currentPage, _ := c.GetInt64("current", 1)
	pageSize, _ := c.GetInt64("pageSize", 50)
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

	t := time.Now()
	year, month, _ := t.Date()
	thisMonthFirstDay := time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	if cycleType == "1" {

	} else if cycleType == "2" {
		thisMonthFirstDay = thisMonthFirstDay.AddDate(0, -2, 0)

	} else if cycleType == "3" {
		thisMonthFirstDay = thisMonthFirstDay.AddDate(0, -5, 0)

	} else if cycleType == "4" {
		thisMonthFirstDay = thisMonthFirstDay.AddDate(0, -11, 0)
	}
	create_time := thisMonthFirstDay.UnixNano() / 1e6

	var replyAccount account.Account
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId },
			action.CondValue{Type: "And", Key: "user_id", Val: replyUserPosition.UserId },
			action.CondValue{Type: "And", Key: "user_position_id", Val: replyUserPosition.Id },
			action.CondValue{Type: "And", Key: "user_position_type", Val: replyUserPosition.Type },
		},
		Fileds: []string{"id"},
	}, &replyAccount)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	var replyPageByCond action.PageByCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "PageByCond", action.PageByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "account_id", Val: replyAccount.Id },
			action.CondValue{Type: "And", Key: "create_time__gt", Val: create_time },
			action.CondValue{Type: "And", Key: "status__gt", Val: -1 },
		},
		Cols:     []string{"id", "create_time", "amount", "remark" },
		OrderBy:  []string{"-id"},
		PageSize: pageSize,
		Current:  currentPage,
	}, &replyPageByCond)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	var replyAccountItemList []account.AccountItem
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "ListByCond", action.ListByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "account_id", Val: replyAccount.Id },
			action.CondValue{Type: "And", Key: "create_time__gt", Val: create_time },
			action.CondValue{Type: "And", Key: "status__gt", Val: -1 },
		},
		Cols: []string{"amount", "balance"},
	}, &replyAccountItemList)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	var inAccount, outAccount int64 = 0, 0
	for _, accountItem := range replyAccountItemList {
		if accountItem.Amount > 0 {
			inAccount += accountItem.Amount
		} else {
			outAccount += accountItem.Amount
		}
	}

	var ailv []AccountItemListValue
	reply2List := []account.AccountItem{}
	err = json.Unmarshal([]byte(replyPageByCond.Json), &reply2List)
	var aType string
	for _, value := range reply2List {
		ailv = append(ailv, AccountItemListValue{
			Id:         value.Id,
			CreateTime: utils.Datetime(value.CreateTime, "2006-01-02 15:04:05"),
			Amount:     value.Amount,
			AmountType: aType,
			Remark:     value.Remark,
		})
	}
	c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: AccountItemList{
		Total:       replyPageByCond.Total,
		Current:     currentPage,
		CurrentSize: replyPageByCond.CurrentSize,
		OrderBy:     replyPageByCond.OrderBy,
		PageSize:    pageSize,
		List:        ailv,
		InAccount:   inAccount,
		OutAccount:  outAccount,
	}}
	c.ServeJSON()
}

// @Title 交易详情列表接口
// @Description 交易详情列表接口
// @Success 200 {"code":200,"message":"获取账务列表成功"}
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   startDate     query   string true       "开始日期：2017-01-01"
// @Param   endDate     query   string true       "结束日期：2017-01-01"
// @Param   current     query   string true       "当前页"
// @Param   pageSize     query   string true       "每页数量"
// @router /tradingList [post]
func (c *AccountItemController) TradingList() {
	type data AccountItemListResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	startDate := c.GetString("startDate")
	endDate := c.GetString("endDate")
	currentPage, _ := c.GetInt64("current", 1)
	pageSize, _ := c.GetInt64("pageSize", 50)

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
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId },
			action.CondValue{Type: "And", Key: "user_id", Val: replyUserPosition.UserId },
			action.CondValue{Type: "And", Key: "user_position_id", Val: replyUserPosition.Id },
			action.CondValue{Type: "And", Key: "user_position_type", Val: replyUserPosition.Type },
		},
		Fileds: []string{"id"},
	}, &replyAccount)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	startTime, _ := time.Parse("2006-01-02", startDate)
	endTime, _ := time.Parse("2006-01-02", endDate)
	var replyPageByCond action.PageByCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "PageByCond", action.PageByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "account_id", Val: replyAccount.Id },
			action.CondValue{Type: "And", Key: "create_time__gte", Val: startTime.UnixNano() / 1e6 },
			action.CondValue{Type: "And", Key: "create_time__lt", Val: endTime.UnixNano() / 1e6  },
		},
		Cols:     []string{"id", "create_time", "amount", "remark", "balance"},
		OrderBy:  []string{"-id"},
		PageSize: pageSize,
		Current:  currentPage,
	}, &replyPageByCond)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	var ailv []AccountItemListValue
	accountItemList := []account.AccountItem{}
	err = json.Unmarshal([]byte(replyPageByCond.Json), &accountItemList)
	//List 循环赋值
	var aType string
	for _, value := range accountItemList {
		ailv = append(ailv, AccountItemListValue{
			Id:         value.Id,
			CreateTime: utils.Datetime(value.CreateTime, "2006-01-02 15:04:05"),
			Amount:     value.Amount,
			Balance:    value.Balance,
			AmountType: aType,
			Remark:     value.Remark,
		})
	}

	c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: AccountItemList{
		Total:       replyPageByCond.Total,
		Current:     currentPage,
		CurrentSize: replyPageByCond.CurrentSize,
		OrderBy:     replyPageByCond.OrderBy,
		PageSize:    pageSize,
		List:        ailv,
	}}
	c.ServeJSON()
}

// @Title 账务详情接口
// @Description 账务详情接口
// @Success 200 {"code":200,"message":"获取账务详情成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"获取账务详情失败"}
// @router /detail [post]
func (c *AccountItemController) Detail() {
	type data AccountItemDetailResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	accountItemId, _ := c.GetInt64("id")

	err := utils.Unable(map[string]string{"accessToken": "string:true", "id": "int:true"}, c.Ctx.Input)
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
	log.Println("replyUserPosition:", replyUserPosition)
	var replyAccountItem account.AccountItem
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "FindById", &account.AccountItem{
		Id: accountItemId,
	}, &replyAccountItem)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试282"}
		c.ServeJSON()
		return
	}

	var ToAccount, OrderCode string
	if replyAccountItem.TransactionId > 0 {
		var transactionReply account.Transaction
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Transaction", "FindById", &account.Transaction{
			Id: replyAccountItem.TransactionId,
		}, &transactionReply)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试294"}
			c.ServeJSON()
			return
		}
		OrderCode = transactionReply.OrderCode
		if transactionReply.ToAccountId == 0 {
			c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
			c.ServeJSON()
			return
		}
		var accountReply account.Account
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindById", &account.Account{Id: transactionReply.ToAccountId}, &accountReply)
		if err == nil && accountReply.CompanyId > 0 {
			var accountReply account.Account
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindById", account.Account{
				Id: transactionReply.ToAccountId,
			}, &accountReply)
			if err != nil {
				c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试312"}
				c.ServeJSON()
				return
			}
			var companyArgs company.Company
			companyArgs.Id = accountReply.CompanyId
			var companyReply company.Company
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "FindById", companyArgs, &companyReply)
			if err == nil {
				ToAccount = companyReply.Name
			}
		}

		c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: AccountItemDetail{
			CreateTime: utils.Datetime(replyAccountItem.CreateTime, "2006-01-02 15:04:05"),
			Amount:     replyAccountItem.Amount,
			Balance:    replyAccountItem.Balance,
			Remark:     replyAccountItem.Remark,
			OrderCode:  OrderCode,
			ToAccount:  ToAccount,
		}}
		c.ServeJSON()
	}
}
