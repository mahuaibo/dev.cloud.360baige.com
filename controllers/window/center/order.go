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

// @Title 账务详情接口
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

// @Title 账务详情接口
// @Description 账务详情接口
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
	num, _ := c.GetFloat("num", 0)
	payType, _ := c.GetInt8("payType", 1)

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
	brief := replyApplicationTpl.Name + "(" + utils.Amount(replyApplicationTpl.Price*num) + ")"
	orderCode := utils.Datetime(utils.CurrentTimestamp(), "20060102030405") + utils.RandomNum(6)

	var replyOrder order.Order
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "Add", &order.Order{
		CreateTime:       currentTimestamp,
		UpdateTime:       currentTimestamp,
		CompanyId:        replyUserPosition.CompanyId,
		UserId:           replyUserPosition.UserId,
		UserPositionType: replyUserPosition.Type,
		UserPositionId:   replyUserPosition.Id,
		Code:             orderCode,
		Price:            replyApplicationTpl.Price * num,
		Type:             0,
		PayType:          payType,
		Brief:            brief,
		Status:           0,
	}, &replyOrder)
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
	c.Data["json"] = data{Code: Normal, Message: "订单新增成功"}
	c.ServeJSON()
	return
}
