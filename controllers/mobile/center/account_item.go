package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/mobile"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/account"
	. "dev.model.360baige.com/models/company"
	"dev.cloud.360baige.com/utils"
	"time"
	"fmt"
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
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
	}
	// 检测 accessToken
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
		// company_id、user_id、user_position_id、user_position_type
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
			args2.Fileds = []string{"id"}
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindByCond", args2, &reply)
			fmt.Println("acccount-id:", reply.Id)
			if err != nil {
				res.Code = ResponseSystemErr
				res.Messgae = "获取账户信息失败"
				c.Data["json"] = res
				c.ServeJSON()
			} else {
				account_id := reply.Id
				var reply2 action.PageByCond
				var cond1 []action.CondValue
				cond1 = append(cond1,action.CondValue{
					Type:  "And",
					Key: "account_id",
					Val:  account_id,
				})
				var stime, etime int64
				current := c.GetString("date")
				if (current == "") {
					current = time.Now().Format("2006-01")
				}
				// 获取指定时间的月初、下个月初时间戳
				stime = utils.GetMonthStartUnix(current + "-01")
				etime = utils.GetNextMonthStartUnix(current + "-01")
				cond1 = append(cond1, action.CondValue{
					Type:  "And",
					Key: "create_time__gte",
					Val:  stime,
				})
				cond1 = append(cond1,action.CondValue{
					Type:  "And",
					Key: "create_time__lt",
					Val:  etime,
				})
				currentPage, _ := c.GetInt64("current")
				pageSize, _ := c.GetInt64("page_size")
				err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "PageByCond", &action.PageByCond{
					CondList:     cond1,
					Cols:     []string{"id", "create_time", "amount", },
					OrderBy:  []string{"id"},
					PageSize: pageSize,
					Current:  currentPage,
				}, &reply2)
				if err != nil {
					res.Code = ResponseSystemErr
					res.Messgae = "获取账务信息失败"
					c.Data["json"] = res
					c.ServeJSON()
				} else {
					reply2List := []AccountItem{}
					err = json.Unmarshal([]byte(reply2.Json), &reply2List)
					// List 循环赋值
					var aType string
					for _, value := range reply2List{
						if (value.Amount < 0) {
							aType = "收入"
							va, _ := strconv.ParseFloat(strings.Replace(strconv.FormatFloat(value.Amount, 'f', 5, 64), "-", "", 1), 64)
							value.Amount = va
						} else {
							aType = "支出"
						}
						re := time.Unix(value.CreateTime/1000, 0).Format("2006-01-02")
						res.Data.List = append(res.Data.List, AccountItemListValue{
							Id:         value.Id,
							CreateTime: re,
							Amount:     value.Amount,
							AmountType: aType,
						})
					}
					res.Data.Total = reply2.Total
					res.Data.Current = currentPage
					res.Data.CurrentSize = reply2.CurrentSize
					res.Data.OrderBy = reply2.OrderBy
					res.Data.PageSize = pageSize
					res.Code = ResponseNormal
					res.Messgae = "获取账务统计信息成功"
					c.Data["json"] = res
					c.ServeJSON()
				}
			}
		}
	}
}
// @Title 账务列表接口
// @Description 账务列表接口
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
			res.Messgae = "获取信息失败"
			c.Data["json"] = res
			c.ServeJSON()
		} else {
			var reply Account
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
			args2.Fileds = []string{"id"}
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindByCond", args2, &reply)
			fmt.Println("acccount-id:", reply.Id)
			if err != nil {
				res.Code = ResponseSystemErr
				res.Messgae = "获取账户信息失败"
				c.Data["json"] = res
				c.ServeJSON()
			} else {
				account_id := reply.Id
				var reply2 action.PageByCond
				var cond1 []action.CondValue
				cond1 = append(cond1, action.CondValue{
					Type:  "And",
					Key: "account_id",
					Val:  account_id,
				})
				var stime, etime int64
				sdate := c.GetString("start_date")
				edate := c.GetString("end_date")
				tm2, _ := time.ParseInLocation("2006-01-02", sdate, time.Local)
				stime = tm2.UnixNano() / 1e6
				etime = utils.GetNextDayUnix(edate)
				cond1 = append(cond1, action.CondValue{
					Type:  "And",
					Key: "create_time__gte",
					Val:  stime,
				})
				cond1 = append(cond1, action.CondValue{
					Type:  "And",
					Key: "create_time__lt",
					Val:  etime,
				})
				currentPage, _ := c.GetInt64("current")
				pageSize, _ := c.GetInt64("page_size")
				err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "PageByCond", &action.PageByCond{
					CondList:     cond1,
					Cols:     []string{"id", "create_time", "amount", },
					OrderBy:  []string{"id"},
					PageSize: pageSize,
					Current:  currentPage,
				}, &reply2)
				if err != nil {
					res.Code = ResponseSystemErr
					res.Messgae = "获取账务信息失败"
					c.Data["json"] = res
					c.ServeJSON()
				} else {
					reply2List := []AccountItem{}
					err = json.Unmarshal([]byte(reply2.Json), &reply2List)
					//List 循环赋值
					var aType string
					for _, value := range reply2List{
						if (value.Amount < 0) {
							aType = "收入"
							va, _ := strconv.ParseFloat(strings.Replace(strconv.FormatFloat(value.Amount, 'f', 5, 64), "-", "", 1), 64)
							value.Amount = va
						} else {
							aType = "支出"
						}
						re := time.Unix(value.CreateTime/1000, 0).Format("2006-01-02")
						res.Data.List = append(res.Data.List, AccountItemListValue{
							Id:         value.Id,
							CreateTime: re,
							Amount:     value.Amount,
							AmountType: aType,
						})
					}
					res.Data.Total = reply2.Total
					res.Data.Current = currentPage
					res.Data.CurrentSize = reply2.CurrentSize
					res.Data.OrderBy = reply2.OrderBy
					res.Data.PageSize = pageSize
					res.Code = ResponseNormal
					res.Messgae = "获取账务统计信息成功"
					c.Data["json"] = res
					c.ServeJSON()
				}
			}
		}
	}
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
		ai_id, _ := c.GetInt64("id")
		if ai_id == 0  {
			res.Code = ResponseSystemErr
			res.Messgae = "获取信息失败"
			c.Data["json"] = res
			c.ServeJSON()
		} else {

			var reply AccountItem
			err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "FindById", &AccountItem{
				Id: ai_id,
			}, &reply)
			if err != nil {
				res.Code = ResponseSystemErr
				res.Messgae = "获取信息失败"
				c.Data["json"] = res
				c.ServeJSON()
			} else {
				var aType string
				re := time.Unix(reply.CreateTime/1000, 0).Format("2006-01-02")
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
		}
	}
}


