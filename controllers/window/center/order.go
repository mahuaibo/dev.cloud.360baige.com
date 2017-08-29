package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/center"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/order"
	"time"
	"dev.model.360baige.com/action"
	"encoding/json"
)

type OrderController struct {
	beego.Controller
}

// @Title 订单列表接口
// @Description 订单列表接口
// @Success 200 {"code":200,"message":"获取订单列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   date     query   string true       "账单日期：2017-07"
// @Param   current     query   string true       "当前页"
// @Param   pageSize     query   string true       "每页数量"
// @Param   status     query   string true       "订单状态：-2 全部 0:撤回 1：待审核 2：已通过 3：未通过 4：发货中 5：完成"
// @Failure 400 {"code":400,"message":"获取订单列表信息失败"}
// @router /list [post]
func (c *OrderController) List() {
	type data OrderListResponse
	accessToken := c.GetString("accessToken")
	status, _ := c.GetInt8("status")
	currentPage, _ := c.GetInt64("current")
	pageSize, _ := c.GetInt64("pageSize")
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

func GetStatus(status int8) string {
	var rStatus string
	switch  status {
	case 0:
		rStatus = "撤回"
	case 1:
		rStatus = "待审核"
	case 2:
		rStatus = "已通过"
	case 3:
		rStatus = "未通过"
	case 4:
		rStatus = "发货中"
	default:
		rStatus = "完成"
	}
	return rStatus
}

// @Title 订单详情接口
// @Description 订单详情接口
// @Success 200 {"code":200,"message":"获取订单详情成功","data":{"access_ticket":"xxxx","expire_in":0}}
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
