package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/center"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/order"
	"dev.model.360baige.com/models/application"
	"dev.model.360baige.com/action"
	"encoding/json"
	"dev.cloud.360baige.com/utils"
	"dev.cloud.360baige.com/utils/pay/wechat"
	"strings"
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
			action.CondValue{Type: "And", Key: "status__gt", Val: -1},
			condValue,
		},
		Cols:     []string{"id", "create_time", "image", "code", "price", "num", "brief", "total_price", "product_type", "product_id", "pay_type", "status"},
		OrderBy:  []string{"-id"},
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
		resData = append(resData, OrderValue{
			Id:          value.Id,
			CreateTime:  utils.Datetime(value.CreateTime, "2006-01-02 15:04:05"),
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
		c.Data["json"] = data{Code: ErrorLogic, Message: "订单关闭失败"}
		c.ServeJSON()
		return
	}

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "DeleteById", &action.DeleteByIdCond{
		Value: []int64{orderId},
	}, &replyNum)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "订单关闭失败"}
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

	c.Data["json"] = data{Code: Normal, Message: "获取订单详情成功", Data: OrderDetail{
		Price:       replyOrder.Price,
		Num:         replyOrder.Num,
		Status:      replyOrder.Status,
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
//func (c *OrderController) DetailByCode() {
//	type data OrderDetailResponse
//	accessToken := c.GetString("accessToken")
//	code := c.GetString("code")
//	if accessToken == "" {
//		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌无效"}
//		c.ServeJSON()
//		return
//	}
//
//	var replyUserPosition user.UserPosition
//	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
//		CondList: []action.CondValue{
//			action.CondValue{Type: "And", Key: "accessToken", Val: accessToken },
//		},
//		Fileds: []string{"id", "user_id", "company_id", "type"},
//	}, &replyUserPosition)
//
//	if err != nil {
//		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌失效"}
//		c.ServeJSON()
//		return
//	}
//	if code == "" {
//		c.Data["json"] = data{Code: ErrorLogic, Message: "获取信息失败"}
//		c.ServeJSON()
//		return
//	}
//
//	var replyOrder order.Order
//	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "FindByCond", action.FindByCond{
//		CondList: []action.CondValue{
//			action.CondValue{Type: "And", Key: "code", Val: code },
//		},
//		Fileds: []string{"id", "create_time", "code", "price", "type", "pay_type", "brief", "status"},
//	}, &replyOrder)
//
//	if err != nil {
//		c.Data["json"] = data{Code: ErrorSystem, Message: "获取信息失败"}
//		c.ServeJSON()
//		return
//	}
//
//	c.Data["json"] = data{Code: Normal, Message: "获取信息成功", Data: OrderDetail{
//		CreateTime:  time.Unix(replyOrder.CreateTime / 1000, 0).Format("2006-01-02"),
//		Code:        replyOrder.Code,
//		Price:       replyOrder.Price,
//		ProductType: replyOrder.ProductType,
//		PayType:     replyOrder.PayType,
//		Brief:       replyOrder.Brief,
//		Status:      replyOrder.Status,
//	}}
//	c.ServeJSON()
//	return
//}

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
	orderCode := utils.Datetime(utils.CurrentTimestamp(), "20060102150405") + utils.RandomNum(6)
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
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}
	if replyOrder.Id == 0 {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}
	// NATIVE MWEB
	unifyOrderResponse, err := wechat.UnifiedOrder(remoteAddr[0], replyApplicationTpl.Name, orderCode, "MWEB", replyApplicationTpl.Price*num)
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

	var tradeState = ""
	if replyOrder.Status == 0 {
		orderQuery, err := wechat.OrderQuery(replyOrder.Code)
		log.Println("orderQuery:", orderQuery, err)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: err.Error()}
			c.ServeJSON()
			return
		}
		tradeState = orderQuery.TradeState
		if orderQuery.ReturnCode == "SUCCESS" && orderQuery.ResultCode == "SUCCESS" && orderQuery.TradeState == "SUCCESS" {
			var replyNum *action.Num
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "UpdateById", &action.UpdateByIdCond{
				Id: []int64{replyOrder.Id},
				UpdateList: []action.UpdateValue{
					action.UpdateValue{"UpdateTime", currentTimestamp},
					action.UpdateValue{"Status", 4},
				},
			}, &replyNum)
			if err != nil {
				c.Data["json"] = data{Code: Normal, Message: "修改订单信息失败"}
				c.ServeJSON()
				return
			}
			if replyNum.Value > 1 {
				// 添加application
				//Id               int64    `db:"id" json:"id"`                               // 主键
				//CreateTime       int64    `db:"create_time" json:"createTime"`              // 创建时间
				//UpdateTime       int64    `db:"update_time" json:"updateTime"`              // 更新时间
				//CompanyId        int64    `db:"company_id" json:"companyId"`                // 所属公司ID
				//UserId           int64    `db:"user_id" json:"userId"`                      // 购买者ID
				//ApplicationTplId int64    `db:"application_tpl_id" json:"applicationTplId"` // 应用ID
				//UserPositionType int8    `db:"user_position_type" json:"userPositionType"`  // 身份类型
				//UserPositionId   int64    `db:"user_position_id" json:"userPositionId"`     // 身份ID
				//Name             string    `db:"name" json:"name"`                          // 名称
				//Image            string    `db:"image" json:"image"`                        // 图片链接
				//Status           int8    `db:"status" json:"status"`                        // 状态 0停用  1启用  2退订
				//StartTime        int64    `db:"start_time" json:"startTime"`                // 开始时间
				//EndTime          int64    `db:"end_time" json:"endTime"`                    // 结束时间
				var replyApplication application.Application
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
				}, &replyApplication)
			}

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
