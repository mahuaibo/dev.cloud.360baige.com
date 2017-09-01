package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/center"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/order"
	"dev.model.360baige.com/models/application"
	"time"
	"dev.model.360baige.com/action"
	"encoding/json"
	"dev.cloud.360baige.com/utils"
	"dev.cloud.360baige.com/utils/pay/wechat"
	"strings"
	"strconv"
	"dev.cloud.360baige.com/log"
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
	accessToken := c.GetString("accessToken")
	status, _ := c.GetInt8("status")
	currentPage, _ := c.GetInt64("current", 1)
	pageSize, _ := c.GetInt64("pageSize", 50)
	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}
	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "accessToken", Val: accessToken },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	if replyUserPosition.UserId == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "获取信息失败"}
		c.ServeJSON()
		return
	}
	var condValue action.CondValue
	if status != -1 {
		condValue = action.CondValue{Type: "And", Key: "status", Val: status }
	}

	var replyPageByCond action.PageByCond
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "PageByCond", action.PageByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId },
			action.CondValue{Type: "And", Key: "user_id", Val: replyUserPosition.UserId },
			action.CondValue{Type: "And", Key: "user_position_id", Val: replyUserPosition.Id},
			action.CondValue{Type: "And", Key: "user_position_type", Val: replyUserPosition.Type},
			condValue,
		},
		Cols:     []string{"id", "create_time", "code", "price", "type", "pay_type", "brief", "status"},
		OrderBy:  []string{"id"},
		PageSize: pageSize,
		Current:  currentPage,
	}, &replyPageByCond)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取订单信息失败"}
		c.ServeJSON()
		return
	}

	orderList := []order.Order{}
	err = json.Unmarshal([]byte(replyPageByCond.Json), &orderList)
	var resData []OrderValue
	for _, value := range orderList {
		var rPayType = "线下支付"
		if value.PayType == 1 {
			rPayType = "在线支付"
		}
		resData = append(resData, OrderValue{
			Id:         value.Id,
			CreateTime: time.Unix(value.CreateTime/1000, 0).Format("2006-01-02"),
			Code:       value.Code,
			Price:      value.Price,
			Type:       value.Type,
			PayType:    rPayType,
			Brief:      value.Brief,
			Status:     GetStatus(value.Status),
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

// @Title 订单详情接口
// @Description 订单详情接口
// @Success 200 {"code":200,"message":"获取订单详情成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"获取取订单信息失败"}
// @router /detail [post]
func (c *OrderController) Detail() {
	type data OrderDetailResponse
	accessToken := c.GetString("accessToken")
	orderId, _ := c.GetInt64("id")
	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "accessToken", Val: accessToken },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}
	if orderId == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "获取订单信息失败"}
		c.ServeJSON()
		return
	}

	var replyOrder order.Order
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "FindById", order.Order{
		Id: orderId,
	}, &replyOrder)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取订单信息失败"}
		c.ServeJSON()
		return
	}

	rPayType := "线下支付"
	if replyOrder.PayType == 1 {
		rPayType = "在线支付"
	}
	c.Data["json"] = data{Code: Normal, Message: "获取订单详情成功", Data: OrderDetail{
		CreateTime: time.Unix(replyOrder.CreateTime/1000, 0).Format("2006-01-02"),
		Code:       replyOrder.Code,
		Price:      replyOrder.Price,
		Type:       replyOrder.Type,
		PayType:    rPayType,
		Brief:      replyOrder.Brief,
		Status:     GetStatus(replyOrder.Status),
	}}
	c.ServeJSON()
	return
}

// @Title 订单详情接口
// @Description 账务详情接口
// @Success 200 {"code":200,"message":"获取账务详情成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   code     query   string true       "code"
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @router /detailByCode [post]
func (c *OrderController) DetailByCode() {
	type data OrderDetailResponse
	accessToken := c.GetString("accessToken")
	code := c.GetString("code")
	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "accessToken", Val: accessToken },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}
	if code == "" {
		c.Data["json"] = data{Code: ErrorLogic, Message: "获取信息失败"}
		c.ServeJSON()
		return
	}

	var replyOrder order.Order
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "code", Val: code },
		},
		Fileds: []string{"id", "create_time", "code", "price", "type", "pay_type", "brief", "status"},
	}, &replyOrder)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取信息失败"}
		c.ServeJSON()
		return
	}

	rPayType := "线下支付"
	if replyOrder.PayType == 1 {
		rPayType = "在线支付"
	}
	c.Data["json"] = data{Code: Normal, Message: "获取信息成功", Data: OrderDetail{
		CreateTime: time.Unix(replyOrder.CreateTime/1000, 0).Format("2006-01-02"),
		Code:       replyOrder.Code,
		Price:      replyOrder.Price,
		Type:       replyOrder.Type,
		PayType:    rPayType,
		Brief:      replyOrder.Brief,
		Status:     GetStatus(replyOrder.Status),
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
// @router /add [post]
func (c *OrderController) Add() {
	type data OrderAddResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	applicationTplId, _ := c.GetInt64("applicationTplId")
	num, _ := c.GetInt64("num", 0)
	payType, _ := c.GetInt8("payType", 1)
	remoteAddr := strings.Split(c.Ctx.Request.RemoteAddr, ":")

	err := utils.Unable(map[string]string{"accessToken": "string:true", "payType": "int:true", "applicationTplId": "int:true", "num": "int:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: Message(40000, err.Error())}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "accessToken", Val: accessToken },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: Message(50000)}
		c.ServeJSON()
		return
	}
	if replyUserPosition.UserId == 0 {
		c.Data["json"] = data{Code: ErrorSystem, Message: Message(30000)}
		c.ServeJSON()
		return
	}

	var replyApplicationTpl application.ApplicationTpl
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindById", &application.ApplicationTpl{
		Id: applicationTplId,
	}, &replyApplicationTpl)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: Message(50001)}
		c.ServeJSON()
		return
	}
	brief := replyApplicationTpl.Name + "(￥:" + strconv.FormatInt(replyApplicationTpl.Price*num, 10) + ")"
	orderCode := utils.Datetime(utils.CurrentTimestamp(), "20060102030405") + utils.RandomNum(6)
	total_fee := replyApplicationTpl.Price * num
	log.Println("brief:", brief)
	log.Println("orderCode:", orderCode)
	log.Println("total_fee:", total_fee)
	var replyOrder order.Order
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "Add", &order.Order{
		CreateTime:       currentTimestamp,
		UpdateTime:       currentTimestamp,
		CompanyId:        replyUserPosition.CompanyId,
		UserId:           replyUserPosition.UserId,
		UserPositionType: replyUserPosition.Type,
		UserPositionId:   replyUserPosition.Id,
		Code:             orderCode,
		Price:            total_fee,
		Type:             0,
		PayType:          payType,
		Brief:            brief,
		Status:           0,
	}, &replyOrder)
	log.Println("replyOrder:", replyOrder)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}
	if replyOrder.Id == 0 {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}
	unifyOrderResponse, err := wechat.UnifiedOrder(remoteAddr[0], brief, orderCode, total_fee)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "统一下单失败"}
		c.ServeJSON()
		return
	}
	var replyNum *action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "UpdateById", &action.UpdateByIdCond{
		Id: []int64{replyOrder.Id},
		UpdateList: []action.UpdateValue{
			action.UpdateValue{"UpdateTime", currentTimestamp},
			action.UpdateValue{"TradeType", unifyOrderResponse.Trade_type},
			action.UpdateValue{"PrepayId", unifyOrderResponse.Prepay_id},
			action.UpdateValue{"CodeUrl", unifyOrderResponse.Code_url},
			action.UpdateValue{"Openid", unifyOrderResponse.Openid},
		},
	}, &replyNum)
	if err != nil {
		c.Data["json"] = data{Code: Normal, Message: "修改订单信息失败"}
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
// @router /payResult [get,post]
func (c *OrderController) PayResult() {
	type data OrderPayResultResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	orderId, _ := c.GetInt64("id")
	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "accessToken", Val: accessToken },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}
	if orderId == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "获取订单信息失败"}
		c.ServeJSON()
		return
	}
	var replyOrder order.Order
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "FindById", &order.Order{
		Id: orderId,
	}, &replyOrder)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取订单信息失败"}
		c.ServeJSON()
		return
	}
	if replyOrder.Status == 0 {
		orderQuery, err := wechat.OrderQuery(replyOrder.Code)
		log.Println("orderQuery:", orderQuery, err)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: err.Error()}
			c.ServeJSON()
			return
		}
		if orderQuery.ReturnCode == "SUCCESS" && orderQuery.ResultCode == "SUCCESS" && orderQuery.TradeState == "SUCCESS" {
			var replyNum *action.Num
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "UpdateById", &order.Order{
				Id:         replyOrder.Id,
				UpdateTime: currentTimestamp,
				Status:     4,
			}, &replyNum)
			if err != nil {
				c.Data["json"] = data{Code: Normal, Message: "修改订单信息失败"}
				c.ServeJSON()
				return
			}
			// 充值 + 消费 + 交易记录

		}
	}

	c.Data["json"] = data{Code: Normal, Message: "获取订单信息"}
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
