package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/mobile/center"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/order"
	"dev.model.360baige.com/models/application"
	"dev.model.360baige.com/action"
	"encoding/json"
	"dev.cloud.360baige.com/utils"
	"dev.cloud.360baige.com/utils/pay/wechat"
	"strings"
	"dev.cloud.360baige.com/log"
	"dev.model.360baige.com/models/account"
)

type OrderController struct {
	beego.Controller
}

// @Title 订单列表接口
// @Description 订单列表接口
// @Success 200 {"code":200,"message":"获取订单列表成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   date     query   string true       "账单日期：2017-07"
// @Param   current     query   string true       "当前页"
// @Param   pageSize     query   string true       "每页数量"
// @Param   status     query   string true       "状态 -1销毁 0 待付款 1待发货 2 待收货 3待评价 4完成 5退货/售后"
// @Failure 400 {"code":400,"message":"获取订单列表信息失败"}
// @router /list [post]
func (c *OrderController) List() {
	type data OrderListResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	status, _ := c.GetInt("status", -100)
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

	var condValue action.CondValue
	if status != -100 {
		condValue = action.CondValue{Type: "And", Key: "status", Val: status }
	}

	var replyPageByCond action.PageByCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "PageByCond", action.PageByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId },
			action.CondValue{Type: "And", Key: "user_id", Val: replyUserPosition.UserId },
			action.CondValue{Type: "And", Key: "user_position_id", Val: replyUserPosition.Id},
			action.CondValue{Type: "And", Key: "user_position_type", Val: replyUserPosition.Type},
			action.CondValue{Type: "And", Key: "status__gt", Val: -1},
			condValue,
		},
		Cols:     []string{"id", "create_time", "image", "code", "price", "num", "brief", "total_price", "product_type", "product_id", "pay_type", "status"},
		OrderBy:  []string{"-id"},
		PageSize: pageSize,
		Current:  currentPage,
	}, &replyPageByCond)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	orderList := []order.Order{}
	err = json.Unmarshal([]byte(replyPageByCond.Json), &orderList)
	var resData []OrderValue
	for _, value := range orderList {
		resData = append(resData, OrderValue{
			Id:          value.Id,
			CreateTime:  value.CreateTime,
			Image:       value.Image,
			Code:        value.Code,
			Price:       value.Price,
			Num:         value.Num,
			Brief:       value.Brief,
			TotalPrice:  value.TotalPrice,
			ProductType: value.ProductType,
			ProductId:   value.ProductId,
			Status:      value.Status,
		})
	}
	c.Data["json"] = data{Code: Normal, Message: "获取信息成功", Data: OrderList{
		Total:       replyPageByCond.Total,
		Current:     currentPage,
		CurrentSize: replyPageByCond.CurrentSize,
		OrderBy:     replyPageByCond.OrderBy,
		PageSize:    pageSize,
		Status:      status,
		List:        resData,
	}}
	c.ServeJSON()
	return
}

// @Title 订单关闭
// @Description 订单关闭
// @Success 200 {"code":200,"message":"订单关闭成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"订单关闭失败"}
// @router /cancel [post]
func (c *OrderController) Cancel() {
	type data OrderCancelResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	orderId, _ := c.GetInt64("id")
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

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "DeleteById", &action.DeleteByIdCond{
		Value: []int64{orderId},
	}, &replyNum)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if replyNum.Value == 0 {
		c.Data["json"] = data{Code: ErrorSystem, Message: "订单关闭失败"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = data{Code: Normal, Message: "订单关闭成功"}
	c.ServeJSON()
	return
}

// @Title 订单详情接口
// @Description 订单详情接口
// @Success 200 {"code":200,"message":"获取订单详情成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"获取取订单信息失败"}
// @router /detail [post]
func (c *OrderController) Detail() {
	type data OrderDetailResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	orderId, _ := c.GetInt64("id")
	err := utils.Unable(map[string]string{"accessToken": "string:true", "orderId": "int:true"}, c.Ctx.Input)
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

	var replyOrder order.Order
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "FindById", order.Order{
		Id: orderId,
	}, &replyOrder)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "获取订单详情成功", Data: OrderDetail{
		Price:   replyOrder.Price,
		Num:     replyOrder.Num,
		Status:  replyOrder.Status,
		CodeUrl: replyOrder.CodeUrl,
	}}
	c.ServeJSON()
	return
}

// @Title 订单新增
// @Description 订单新增
// @Success 200 {"code":200,"message":"获取账务详情成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   code     query   string true       "code"
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @router /add [get,post]
func (c *OrderController) Add() {
	type data OrderAddResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	applicationTplId, _ := c.GetInt64("applicationTplId")
	num, _ := c.GetInt64("num", 0)
	payType, _ := c.GetInt("payType", 1)
	tradeType := c.GetString("tradeType", "NATIVE")
	remoteAddr := strings.Split(c.Ctx.Request.RemoteAddr, ":")

	err := utils.Unable(map[string]string{"accessToken": "string:true", "payType": "int:true", "applicationTplId": "int:true", "num": "int:true"}, c.Ctx.Input)
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

	var replyApplicationTpl application.ApplicationTpl
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindById", &application.ApplicationTpl{
		Id: applicationTplId,
	}, &replyApplicationTpl)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	orderCode := utils.Datetime(currentTimestamp, "20060102150405") + utils.RandomNum(6)
	var replyOrder order.Order
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "Add", &order.Order{
		CreateTime:       currentTimestamp,
		UpdateTime:       currentTimestamp,
		CompanyId:        replyUserPosition.CompanyId,
		UserId:           replyUserPosition.UserId,
		UserPositionType: replyUserPosition.Type,
		UserPositionId:   replyUserPosition.Id,
		Code:             orderCode,
		Price:            replyApplicationTpl.Price,
		Num:              num,
		TotalPrice:       replyApplicationTpl.Price * num,
		ProductType:      0,
		Image:            replyApplicationTpl.Image,
		ProductId:        replyApplicationTpl.Id,
		PayType:          payType,
		Brief:            replyApplicationTpl.Name,
		Status:           0,
	}, &replyOrder)
	log.Println("replyOrder:", replyOrder)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if replyOrder.Id == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}
	// NATIVE MWEB
	unifyOrderResponse, err := wechat.UnifiedOrder(remoteAddr[0], replyApplicationTpl.Name, orderCode, tradeType, replyApplicationTpl.Price*num)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "统一下单失败"}
		c.ServeJSON()
		return
	}
	url := unifyOrderResponse.Code_url
	if url == "" {
		url = unifyOrderResponse.Mweb_url
	}
	var replyNum *action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "UpdateById", &action.UpdateByIdCond{
		Id: []int64{replyOrder.Id},
		UpdateList: []action.UpdateValue{
			action.UpdateValue{"UpdateTime", currentTimestamp},
			action.UpdateValue{"TradeType", unifyOrderResponse.Trade_type},
			action.UpdateValue{"PrepayId", unifyOrderResponse.Prepay_id},
			action.UpdateValue{"CodeUrl", url},
			action.UpdateValue{"Openid", unifyOrderResponse.Openid},
		},
	}, &replyNum)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "订单新增成功", Data: OrderAdd{Id: replyOrder.Id, CodeUrl: unifyOrderResponse.Code_url}}
	c.ServeJSON()
	return
}

// @Title 订单支付结果
// @Description 订单支付结果
// @Success 200 {"code":200,"message":"获取订单详情成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"获取取订单信息失败"}
// @router /payResult [post]
func (c *OrderController) PayResult() {
	type data OrderPayResultResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	orderId, _ := c.GetInt64("id")

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

	var replyOrder order.Order
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "FindById", &order.Order{
		Id: orderId,
	}, &replyOrder)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	if replyOrder.Status >= 4 {
		c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: OrderPayResult{
			TradeState: "SUCCESS",
		}}
		c.ServeJSON()
		return
	}

	var tradeState = ""
	if replyOrder.Status == 0 {
		var replyApplication application.Application
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "FindByCond", &action.FindByCond{
			CondList: []action.CondValue{
				action.CondValue{"And", "company_id", replyUserPosition.CompanyId},
				action.CondValue{"And", "user_id", replyUserPosition.UserId},
				action.CondValue{"And", "user_position_type", replyUserPosition.Type},
				action.CondValue{"And", "user_position_id", replyUserPosition.Id},
				action.CondValue{"And", "application_tpl_id", replyOrder.ProductId},
			},
		}, &replyApplication)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}

		var replyApplicationTpl application.ApplicationTpl
		err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindById", &application.ApplicationTpl{
			Id: replyOrder.ProductId,
		}, &replyApplicationTpl)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}

		orderQuery, err := wechat.OrderQuery(replyOrder.Code)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: err.Error()}
			c.ServeJSON()
			return
		}

		tradeState = orderQuery.TradeState
		if orderQuery.ReturnCode == "SUCCESS" && orderQuery.ResultCode == "SUCCESS" && orderQuery.TradeState == "SUCCESS" {
			// have application
			if replyApplication.Id > 0 {
				// 修改时间
				endTime := replyApplication.EndTime
				if currentTimestamp > endTime {
					endTime = currentTimestamp
				}
				endTime = utils.ServiceTime(replyApplicationTpl.PayType, replyApplicationTpl.PayCycle, replyOrder.Num, endTime)
				var replyNum *action.Num
				err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "UpdateById", &action.UpdateByIdCond{
					Id: []int64{replyApplication.Id},
					UpdateList: []action.UpdateValue{
						action.UpdateValue{"update_time", currentTimestamp},
						action.UpdateValue{"StartTime", replyApplication.EndTime},
						action.UpdateValue{"end_time", endTime},
						action.UpdateValue{"status", 0},
					},
				}, &replyNum)

			} else {
				var replyApplication application.Application
				endTime := utils.ServiceTime(replyApplicationTpl.PayType, replyApplicationTpl.PayCycle, replyOrder.Num, currentTimestamp)
				// 新增时间
				err = client.Call(beego.AppConfig.String("EtcdURL"), "Application", "Add", &application.Application{
					CreateTime:       currentTimestamp,
					UpdateTime:       currentTimestamp,
					CompanyId:        replyUserPosition.CompanyId,
					UserId:           replyUserPosition.UserId,
					ApplicationTplId: replyOrder.ProductId,
					UserPositionType: replyUserPosition.Type,
					UserPositionId:   replyUserPosition.Id,
					Name:             replyOrder.Brief,
					Image:            replyOrder.Image,
					Status:           0,
					StartTime:        currentTimestamp,
					EndTime:          endTime,
				}, &replyApplication)
			}

			if err != nil {
				c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
				c.ServeJSON()
				return
			}

			var replyNum *action.Num
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "UpdateById", &action.UpdateByIdCond{
				Id: []int64{replyOrder.Id},
				UpdateList: []action.UpdateValue{
					action.UpdateValue{"update_time", currentTimestamp},
					action.UpdateValue{"Status", 4},
				},
			}, &replyNum)
			if err != nil {
				c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
				c.ServeJSON()
				return
			}
			if replyNum.Value == 0 {
				c.Data["json"] = data{Code: ErrorLogic, Message: "系统异常，请稍后重试"}
				c.ServeJSON()
				return
			}
			err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "UpdateById", &action.UpdateByIdCond{
				Id: []int64{replyApplicationTpl.Id},
				UpdateList: []action.UpdateValue{
					action.UpdateValue{"update_time", currentTimestamp},
					action.UpdateValue{"subscription", replyApplicationTpl.Subscription + 1},
				},
			}, &replyNum)

			// 充值
			AccountTransaction(&action.FindByCond{
				CondList: []action.CondValue{
					action.CondValue{"And", "CompanyId", user.UserPositionCompanyIdAudit},
					action.CondValue{"And", "UserId", user.UserIdAudit},
					action.CondValue{"And", "UserPositionType", user.UserPositionTypeAudit},
					action.CondValue{"And", "UserPositionId", user.UserPositionIdAudit},
					action.CondValue{"And", "Type", account.AccountTypeMoney},
				},
			}, &action.FindByCond{
				CondList: []action.CondValue{
					action.CondValue{"And", "CompanyId", replyUserPosition.CompanyId},
					action.CondValue{"And", "UserId", replyUserPosition.UserId},
					action.CondValue{"And", "UserPositionType", replyUserPosition.Type},
					action.CondValue{"And", "UserPositionId", replyUserPosition.Id},
					action.CondValue{"And", "Type", account.AccountTypeMoney},
				},
			}, replyOrder.TotalPrice, replyOrder.Code, replyOrder.Brief)

			// 消费
			AccountTransaction(&action.FindByCond{
				CondList: []action.CondValue{
					action.CondValue{"And", "CompanyId", replyUserPosition.CompanyId},
					action.CondValue{"And", "UserId", replyUserPosition.UserId},
					action.CondValue{"And", "UserPositionType", replyUserPosition.Type},
					action.CondValue{"And", "UserPositionId", replyUserPosition.Id},
					action.CondValue{"And", "Type", account.AccountTypeMoney},
				},
			}, &action.FindByCond{
				CondList: []action.CondValue{
					action.CondValue{"And", "CompanyId", user.UserPositionCompanyIdAudit},
					action.CondValue{"And", "UserId", user.UserIdAudit},
					action.CondValue{"And", "UserPositionType", user.UserPositionTypeAudit},
					action.CondValue{"And", "UserPositionId", user.UserPositionIdAudit},
					action.CondValue{"And", "Type", account.AccountTypeMoney},
				},
			}, replyOrder.TotalPrice, replyOrder.Code, replyOrder.Brief)

		} else if orderQuery.ReturnCode == "SUCCESS" && orderQuery.ResultCode == "SUCCESS" && orderQuery.TradeState == "SUCCESS" {

		}
	}

	c.Data["json"] = data{Code: Normal, Message: "获取订单信息", Data: OrderPayResult{
		TradeState: tradeState,
	}}
	c.ServeJSON()
	return
}

// @Title 获取二维码
// @Param   url     query   string true       "url"
// @Param   size     query   int true       "size 默认256"
// @router /qr [get,post]
func (c *OrderController) Qr() {
	url := c.GetString("url")
	size, _ := c.GetInt("size", 256)
	c.Ctx.Output.Body(utils.Qr(url, size))
}

// form - | to +
func AccountTransaction(fromFindByCond, toFindByCond *action.FindByCond, amount int64, orderCode, remark string) error {
	currentTimestamp := utils.CurrentTimestamp()

	// formAccount
	var replyFormAccount account.Account
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindByCond", &fromFindByCond, &replyFormAccount)

	// toAccount
	var replyToAccount account.Account
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindByCond", &toFindByCond, &replyToAccount)

	// transaction
	var replyTransaction account.Transaction
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Transaction", "Add", &account.Transaction{
		CreateTime:    currentTimestamp,
		UpdateTime:    currentTimestamp,
		FromAccountId: replyFormAccount.Id,
		ToAccountId:   replyToAccount.Id,
		Amount:        amount,
		OrderCode:     orderCode,
		Remark:        remark,
		Status:        account.TransactionStatusNormal,
	}, &replyTransaction)

	// formAccountItem
	var replyFormAccountItem account.AccountItem
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "Add", &account.AccountItem{
		CreateTime:    currentTimestamp,
		UpdateTime:    currentTimestamp,
		TransactionId: replyTransaction.Id,
		AccountId:     replyFormAccount.Id,
		Balance:       replyFormAccount.Balance - amount,
		Amount:        -amount,
		Remark:        remark,
		Status:        account.AccountItemStatusNormal,
	}, &replyFormAccountItem)
	// toAccountItem
	var replyToAccountItem account.AccountItem
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "Add", &account.AccountItem{
		CreateTime:    currentTimestamp,
		UpdateTime:    currentTimestamp,
		TransactionId: replyTransaction.Id,
		AccountId:     replyToAccount.Id,
		Balance:       replyToAccount.Balance + amount,
		Amount:        amount,
		Remark:        remark,
		Status:        account.AccountItemStatusNormal,
	}, &replyToAccountItem)

	// formAccountUpdate
	var replyFormNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "UpdateById", &action.UpdateByIdCond{
		Id: []int64{replyFormAccount.Id},
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "UpdateTime", Val: currentTimestamp},
			action.UpdateValue{Key: "Balance", Val: replyFormAccount.Balance - amount},
		},
	}, &replyFormNum)
	// toAccountUpdate
	var replyToNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "UpdateById", &action.UpdateByIdCond{
		Id: []int64{replyToAccount.Id},
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "UpdateTime", Val: currentTimestamp},
			action.UpdateValue{Key: "Balance", Val: replyToAccount.Balance + amount},
		},
	}, &replyToNum)

	return err
}
