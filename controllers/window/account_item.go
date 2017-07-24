package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/account"
	. "dev.model.360baige.com/models/company"
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

// @Title 账务列表接口
// @Description 账务列表接口
// @Success 200 {"code":200,"messgae":"获取账务列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   date     query   string true       "账单日期：2017-07"
// @Param   current     query   string true       "当前页"
// @Param   page_size     query   string true       "每页数量"
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @router /list [get]
func (c *AccountItemController) List() {
	res := AccountItemListResponse{}
	access_token := c.GetString("access_token")
	current := c.GetString("date")
	currentPage, _ := c.GetInt64("current")
	pageSize, _ := c.GetInt64("page_size")
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
		res.Messgae = "获取失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var accountReply Account
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
	accountArgs.Fileds = []string{"id"}
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindByCond", accountArgs, &accountReply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var accountItemArgs action.PageByCond
	var accountItemReply action.PageByCond
	if (current == "") {
		current = time.Now().Format("2006-01")
	}
	accountItemArgs.CondList = append(accountItemArgs.CondList, action.CondValue{
		Type:  "And",
		Key: "account_id",
		Val:  accountReply.Id,
	}, action.CondValue{
		Type:  "And",
		Key: "create_time__gte",
		Val:  utils.GetMonthStartUnix(current + "-01"),
	}, action.CondValue{
		Type:  "And",
		Key: "create_time__lt",
		Val:  utils.GetNextMonthStartUnix(current + "-01"),
	})
	accountItemArgs.Cols = []string{"id", "create_time", "amount", }
	accountItemArgs.OrderBy = []string{"id"}
	accountItemArgs.PageSize = pageSize
	accountItemArgs.Current = currentPage
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "PageByCond", accountItemArgs, &accountItemReply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	reply2List := []AccountItem{}
	err = json.Unmarshal([]byte(accountItemReply.Json), &reply2List)
	var aType string
	for _, value := range reply2List {
		if (value.Amount < 0) {
			aType = "收入"
			value.Amount, _ = strconv.ParseFloat(strings.Replace(strconv.FormatFloat(value.Amount, 'f', 5, 64), "-", "", 1), 64)
		} else {
			aType = "支出"
		}
		res.Data.List = append(res.Data.List, AccountItemListValue{
			Id:         value.Id,
			CreateTime: time.Unix(value.CreateTime / 1000, 0).Format("2006-01-02"),
			Amount:     value.Amount,
			AmountType: aType,
		})
	}

	res.Data.Total = accountItemReply.Total
	res.Data.Current = currentPage
	res.Data.CurrentSize = accountItemReply.CurrentSize
	res.Data.OrderBy = accountItemReply.OrderBy
	res.Data.PageSize = pageSize
	res.Code = ResponseNormal
	res.Messgae = "获取信息成功"
	c.Data["json"] = res
	c.ServeJSON()
}
// @Title 交易详情接口
// @Description 交易详情接口
// @Success 200 {"code":200,"messgae":"获取账务列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   start_date     query   string true       "开始日期：2017-01-01"
// @Param   end_date     query   string true       "结束日期：2017-01-01"
// @Param   current     query   string true       "当前页"
// @Param   page_size     query   string true       "每页数量"
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @router /tradinglist [get]
func (c *AccountItemController) TradingList() {
	res := AccountItemListResponse{}
	access_token := c.GetString("access_token")
	sdate := c.GetString("start_date")
	edate := c.GetString("end_date")
	currentPage, _ := c.GetInt64("current")
	pageSize, _ := c.GetInt64("page_size")
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
		res.Messgae = "获取信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var reply Account
	var args2 action.FindByCond
	args2.CondList = append(args2.CondList, action.CondValue{
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
	args2.Fileds = []string{"id"}
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindByCond", args2, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var accountItemArgs action.PageByCond
	var AccountItemReply action.PageByCond
	tm2, _ := time.ParseInLocation("2006-01-02", sdate, time.Local)
	accountItemArgs.CondList = append(accountItemArgs.CondList, action.CondValue{
		Type:  "And",
		Key: "account_id",
		Val:  reply.Id,
	}, action.CondValue{
		Type:  "And",
		Key: "create_time__gte",
		Val:  tm2.UnixNano() / 1e6,
	}, action.CondValue{
		Type:  "And",
		Key: "create_time__lt",
		Val:  utils.GetNextDayUnix(edate),
	})
	accountItemArgs.Cols = []string{"id", "create_time", "amount", }
	accountItemArgs.OrderBy = []string{"id"}
	accountItemArgs.PageSize = pageSize
	accountItemArgs.Current = currentPage
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "PageByCond", accountItemArgs, &AccountItemReply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	reply2List := []AccountItem{}
	err = json.Unmarshal([]byte(AccountItemReply.Json), &reply2List)
	//List 循环赋值
	var aType string
	for _, value := range reply2List {
		if (value.Amount < 0) {
			aType = "收入"
			value.Amount, _ = strconv.ParseFloat(strings.Replace(strconv.FormatFloat(value.Amount, 'f', 5, 64), "-", "", 1), 64)
		} else {
			aType = "支出"
		}
		re := time.Unix(value.CreateTime / 1000, 0).Format("2006-01-02")
		res.Data.List = append(res.Data.List, AccountItemListValue{
			Id:         value.Id,
			CreateTime: re,
			Amount:     value.Amount,
			AmountType: aType,
		})
	}

	res.Data.Total = AccountItemReply.Total
	res.Data.Current = currentPage
	res.Data.CurrentSize = AccountItemReply.CurrentSize
	res.Data.OrderBy = AccountItemReply.OrderBy
	res.Data.PageSize = pageSize
	res.Code = ResponseNormal
	res.Messgae = "获取信息成功"
	c.Data["json"] = res
	c.ServeJSON()
}
// @Title 账务详情接口
// @Description 账务详情接口
// @Success 200 {"code":200,"messgae":"获取账务详情成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"获取账务详情失败"}
// @router /detail [get]
func (c *AccountItemController) Detail() {
	res := AccountItemDetailResponse{}
	access_token := c.GetString("access_token")
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
	ai_id, _ := c.GetInt64("id")
	if ai_id == 0 {
		res.Code = ResponseSystemErr
		res.Messgae = "获取信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var reply AccountItem
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "FindById", &AccountItem{
		Id: ai_id,
	}, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var aType string
	re := time.Unix(reply.CreateTime / 1000, 0).Format("2006-01-02")
	if reply.Amount < 0 {
		aType = "收入"
		va, _ := strconv.ParseFloat(strings.Replace(strconv.FormatFloat(reply.Amount, 'f', 5, 64), "-", "", 1), 64)
		reply.Amount = va
	} else {
		aType = "支出"
	}
	if reply.Balance < 0 {
		vb, _ := strconv.ParseFloat(strings.Replace(strconv.FormatFloat(reply.Balance, 'f', 5, 64), "-", "", 1), 64)
		reply.Balance = vb
	}
	res.Code = ResponseSystemErr
	res.Messgae = "获取账户信息成功"
	res.Data.CreateTime = re
	res.Data.Amount = reply.Amount
	res.Data.AmountType = aType
	res.Data.Balance = reply.Balance
	res.Data.Remark = reply.Remark

	if reply.TransactionId > 0 {
		var reply2 Transaction
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Transaction", "FindById", &Transaction{
			Id: reply.TransactionId,
		}, &reply2)
		if err == nil {
			res.Data.OrderCode = reply2.OrderCode
			if reply2.ToAccountId > 0 {
				var reply3 Account
				err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindById", &Account{
					Id: reply2.ToAccountId,
				}, &reply3)
				if err == nil && reply3.CompanyId > 0 {
					var reply4 Company
					err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "FindById", &Company{
						Id: reply3.CompanyId,
					}, &reply4)
					if err == nil {
						res.Data.ToAccount = reply4.Name
					}
				}
			}
		}
	}
	c.Data["json"] = res
	c.ServeJSON()
}


