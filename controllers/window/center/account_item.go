package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/center"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/account"
	"dev.model.360baige.com/models/company"
	"dev.cloud.360baige.com/utils"
	"time"
	"strconv"
	"strings"
	"dev.model.360baige.com/action"
	"encoding/json"
)

type AccountItemController struct {
	beego.Controller
}

// @Title 账户列表接口
// @Description 账户列表接口
// @Success 200 {"code":200,"message":"获取账务列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   date     query   string true       "账单日期：2017-07"
// @Param   current     query   string true       "当前页"
// @Param   page_size     query   string true       "每页数量"
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @router /list [get]
func (c *AccountItemController) List() {
	type data AccountItemListResponse
	access_token := c.GetString("access_token")
	currentPage, _ := c.GetInt64("current")
	pageSize, _ := c.GetInt64("page_size")
	if access_token == "" {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "accessToken", Val: access_token },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)

	if err != nil {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}
	if replyUserPosition.UserId == 0 {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "访问令牌失效"}
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
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取失败"}
		c.ServeJSON()
		return
	}

	var replyPageByCond action.PageByCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "PageByCond", action.PageByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "account_id", Val: replyAccount.Id },
		},
		Cols:     []string{"id", "create_time", "amount" },
		OrderBy:  []string{"id"},
		PageSize: pageSize,
		Current:  currentPage,
	}, &replyPageByCond)
	if err != nil {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取失败"}
		c.ServeJSON()
		return
	}

	var ailv []AccountItemListValue
	reply2List := []account.AccountItem{}
	err = json.Unmarshal([]byte(replyPageByCond.Json), &reply2List)
	var aType string
	for _, value := range reply2List {
		if (value.Amount < 0) {
			aType = "收入"
			value.Amount, _ = strconv.ParseFloat(strings.Replace(strconv.FormatFloat(value.Amount, 'f', 5, 64), "-", "", 1), 64)
		} else {
			aType = "支出"
		}
		ailv = append(ailv, AccountItemListValue{
			Id:         value.Id,
			CreateTime: time.Unix(value.CreateTime/1000, 0).Format("2006-01-02"),
			Amount:     value.Amount,
			AmountType: aType,
		})
	}
	c.Data["json"] = data{Code: ResponseNormal, Message: "获取信息成功", Data: AccountItemList{
		Total:       replyPageByCond.Total,
		Current:     currentPage,
		CurrentSize: replyPageByCond.CurrentSize,
		OrderBy:     replyPageByCond.OrderBy,
		PageSize:    pageSize,
		List:        ailv,
	}}
	c.ServeJSON()
}

// @Title 交易详情列表接口
// @Description 交易详情列表接口
// @Success 200 {"code":200,"message":"获取账务列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   start_date     query   string true       "开始日期：2017-01-01"
// @Param   end_date     query   string true       "结束日期：2017-01-01"
// @Param   current     query   string true       "当前页"
// @Param   page_size     query   string true       "每页数量"
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @router /tradinglist [get]
func (c *AccountItemController) TradingList() {
	type data AccountItemListResponse
	access_token := c.GetString("access_token")
	sdate := c.GetString("start_date")
	edate := c.GetString("end_date")
	currentPage, _ := c.GetInt64("current")
	pageSize, _ := c.GetInt64("page_size")
	if access_token == "" {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}
	//检测 accessToken
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type: "And",
		Key:  "accessToken",
		Val:  access_token,
	})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取失败"}
		c.ServeJSON()
		return
	}

	com_id := replyAccessToken.CompanyId
	user_id := replyAccessToken.UserId
	user_position_id := replyAccessToken.Id
	user_position_type := replyAccessToken.Type
	if com_id == 0 || user_id == 0 || user_position_id == 0 {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取失败"}
		c.ServeJSON()
		return
	}

	var reply account.Account
	var AccountArgs action.FindByCond
	AccountArgs.CondList = append(AccountArgs.CondList, action.CondValue{
		Type: "And",
		Key:  "company_id",
		Val:  com_id,
	}, action.CondValue{
		Type: "And",
		Key:  "user_id",
		Val:  user_id,
	}, action.CondValue{
		Type: "And",
		Key:  "user_position_id",
		Val:  user_position_id,
	}, action.CondValue{
		Type: "And",
		Key:  "user_position_type",
		Val:  user_position_type,
	})
	AccountArgs.Fileds = []string{"id"}
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindByCond", AccountArgs, &reply)
	if err != nil {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取失败"}
		c.ServeJSON()
		return
	}

	var accountItemArgs action.PageByCond
	var replyPageByCond action.PageByCond
	accountItemArgs.CondList = append(accountItemArgs.CondList, action.CondValue{
		Type: "And",
		Key:  "account_id",
		Val:  reply.Id,
	})
	if sdate != "" && edate != "" {
		tm2, _ := time.ParseInLocation("2006-01-02", sdate, time.Local)
		accountItemArgs.CondList = append(accountItemArgs.CondList, action.CondValue{
			Type: "And",
			Key:  "create_time__gte",
			Val:  tm2.UnixNano() / 1e6,
		}, action.CondValue{
			Type: "And",
			Key:  "create_time__lt",
			Val:  utils.GetNextDayUnix(edate),
		})
	}
	accountItemArgs.Cols = []string{"id", "create_time", "amount", }
	accountItemArgs.OrderBy = []string{"id"}
	accountItemArgs.PageSize = pageSize
	accountItemArgs.Current = currentPage
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "PageByCond", accountItemArgs, &replyPageByCond)
	if err != nil {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取失败"}
		c.ServeJSON()
		return
	}
	var ailv []AccountItemListValue
	accountItemList := []account.AccountItem{}
	err = json.Unmarshal([]byte(replyPageByCond.Json), &accountItemList)
	//List 循环赋值
	var aType string
	for _, value := range accountItemList {
		if (value.Amount < 0) {
			aType = "收入"
			value.Amount, _ = strconv.ParseFloat(strings.Replace(strconv.FormatFloat(value.Amount, 'f', 5, 64), "-", "", 1), 64)
		} else {
			aType = "支出"
		}
		re := time.Unix(value.CreateTime/1000, 0).Format("2006-01-02")
		ailv = append(ailv, AccountItemListValue{
			Id:         value.Id,
			CreateTime: re,
			Amount:     value.Amount,
			AmountType: aType,
		})
	}

	c.Data["json"] = data{Code: ResponseNormal, Message: "获取信息成功", Data: AccountItemList{
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
// @Success 200 {"code":200,"message":"获取账务详情成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"获取账务详情失败"}
// @router /detail [get]
func (c *AccountItemController) Detail() {
	type data AccountItemDetailResponse
	access_token := c.GetString("access_token")
	ai_id, _ := c.GetInt64("id")
	if access_token == "" {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取失败"}
		c.ServeJSON()
		return
	}
	//检测 accessToken
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type: "And",
		Key:  "accessToken",
		Val:  access_token,
	})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取失败"}
		c.ServeJSON()
		return
	}
	if ai_id == 0 {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取失败"}
		c.ServeJSON()
		return
	}

	var accountItemArgs account.AccountItem
	accountItemArgs.Id = ai_id
	var accountItemReply account.AccountItem
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "FindById", accountItemArgs, &accountItemReply)
	if err != nil {
		c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取失败"}
		c.ServeJSON()
		return
	}

	var aType string
	re := time.Unix(accountItemReply.CreateTime/1000, 0).Format("2006-01-02")
	if accountItemReply.Amount < 0 {
		aType = "收入"
		accountItemReply.Amount, _ = strconv.ParseFloat(strings.Replace(strconv.FormatFloat(accountItemReply.Amount, 'f', 5, 64), "-", "", 1), 64)
	} else {
		aType = "支出"
	}
	if accountItemReply.Balance < 0 {
		accountItemReply.Balance, _ = strconv.ParseFloat(strings.Replace(strconv.FormatFloat(accountItemReply.Balance, 'f', 5, 64), "-", "", 1), 64)
	}
	var ToAccount, OrderCode string
	if accountItemReply.TransactionId > 0 {
		var transactionArgs account.Transaction
		transactionArgs.Id = accountItemReply.TransactionId
		var transactionReply account.Transaction
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Transaction", "FindById", transactionArgs, &transactionReply)
		if err != nil {
			c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取失败"}
			c.ServeJSON()
			return
		}
		OrderCode = transactionReply.OrderCode
		if transactionReply.ToAccountId > 0 {
			c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取失败"}
			c.ServeJSON()
			return
		}
		var accountArgs account.Account
		accountArgs.Id = transactionReply.ToAccountId
		var accountReply account.Account
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindById", accountArgs, &accountReply)
		if err == nil && accountReply.CompanyId > 0 {
			c.Data["json"] = data{Code: ResponseSystemErr, Message: "获取失败"}
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

	c.Data["json"] = data{Code: ResponseNormal, Message: "获取信息成功", Data: AccountItemDetail{
		CreateTime: re,
		Amount:     accountItemReply.Amount,
		AmountType: aType,
		Balance:    accountItemReply.Balance,
		Remark:     accountItemReply.Remark,
		OrderCode:  OrderCode,
		ToAccount:  ToAccount,
	}}
	c.ServeJSON()
}
