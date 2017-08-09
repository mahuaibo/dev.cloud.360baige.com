package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/center"
	. "dev.model.360baige.com/models/user"
	//. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/models/order"
	"time"
	"dev.model.360baige.com/action"
	"encoding/json"
)

type OrderController struct {
	beego.Controller
}

// @Title 订单列表接口
// @Description 订单列表接口
// @Success 200 {"code":200,"messgae":"获取订单列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   date     query   string true       "账单日期：2017-07"
// @Param   current     query   string true       "当前页"
// @Param   page_size     query   string true       "每页数量"
// @Param   status     query   string true       "订单状态：-2 全部 0:撤回 1：待审核 2：已通过 3：未通过 4：发货中 5：完成"
// @Failure 400 {"code":400,"message":"获取订单列表信息失败"}
// @router /list [get]
func (c *OrderController) List() {
	res := OrderListResponse{}
	access_token := c.GetString("access_token")
	status, _ := c.GetInt8("status")
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

	var orderArgs action.PageByCond
	var orderReply action.PageByCond
	orderArgs.CondList = append(orderArgs.CondList, action.CondValue{
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
	if status != -1 {
		orderArgs.CondList = append(orderArgs.CondList, action.CondValue{
			Type:  "And",
			Key: "status",
			Val:  status,
		})
	}
	orderArgs.Cols = []string{"id", "create_time", "code", "price", "type", "pay_type", "brief", "status"}
	orderArgs.OrderBy = []string{"id"}
	orderArgs.PageSize = pageSize
	orderArgs.Current = currentPage
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "PageByCond", orderArgs, &orderReply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取订单信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	replyList := []Order{}
	err = json.Unmarshal([]byte(orderReply.Json), &replyList)
	//List 循环赋值
	for _, value := range replyList {
		var rPayType, rStatus string
		if value.PayType == 1 {
			rPayType = "在线支付"
		} else {
			rPayType = "线下支付"
		}
		rStatus = GetStatus(value.Status)
		res.Data.List = append(res.Data.List, OrderValue{
			Id:         value.Id,
			CreateTime: time.Unix(value.CreateTime / 1000, 0).Format("2006-01-02"),
			Code:       value.Code,
			Price:      value.Price,
			Type:       value.Type,
			PayType:    rPayType,
			Brief:      value.Brief,
			Status:     rStatus,
		})
	}

	res.Data.Total = orderReply.Total
	res.Data.Current = currentPage
	res.Data.CurrentSize = orderReply.CurrentSize
	res.Data.OrderBy = orderReply.OrderBy
	res.Data.PageSize = pageSize
	res.Data.Status = status
	res.Code = ResponseNormal
	res.Messgae = "获取信息成功"
	c.Data["json"] = res
	c.ServeJSON()
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
// @Success 200 {"code":200,"messgae":"获取订单详情成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"获取取订单信息失败"}
// @router /detail [get]
func (c *OrderController) Detail() {
	res := OrderDetailResponse{}
	access_token := c.GetString("access_token")
	ai_id, _ := c.GetInt64("id")
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
	if ai_id == 0 {
		res.Code = ResponseSystemErr
		res.Messgae = "获取订单信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var orderArgs Order
	orderArgs.Id = ai_id
	var orderReply Order
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "FindById", orderArgs, &orderReply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取订单信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.Code = ResponseNormal
	res.Messgae = "获取订单详情成功"
	res.Data.CreateTime = time.Unix(orderReply.CreateTime / 1000, 0).Format("2006-01-02")
	res.Data.Code = orderReply.Code
	res.Data.Price = orderReply.Price
	res.Data.Type = orderReply.Type
	if orderReply.PayType == 1 {
		res.Data.PayType = "在线支付"
	} else {
		res.Data.PayType = "线下支付"
	}
	res.Data.Brief = orderReply.Brief
	res.Data.Status = GetStatus(orderReply.Status)
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 账务详情接口
// @Description 账务详情接口
// @Success 200 {"code":200,"messgae":"获取账务详情成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   code     query   string true       "code"
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @router /detailbycode [get]
func (c *OrderController) DetailByCode() {
	res := OrderDetailResponse{}
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

	code := c.GetString("code")
	if code == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "获取信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var reply Order
	var orderArgs action.FindByCond
	orderArgs.CondList = append(orderArgs.CondList, action.CondValue{
		Type:  "And",
		Key: "code",
		Val:  code,
	})
	orderArgs.Fileds = []string{"id", "create_time", "code", "price", "type", "pay_type", "brief", "status"}
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "FindByCond", orderArgs, &reply)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.Code = ResponseSystemErr
	res.Messgae = "获取账户信息成功"
	res.Data.CreateTime = time.Unix(reply.CreateTime / 1000, 0).Format("2006-01-02")
	res.Data.Code = reply.Code
	res.Data.Price = reply.Price
	res.Data.Type = reply.Type
	if reply.PayType == 1 {
		res.Data.PayType = "在线支付"
	} else {
		res.Data.PayType = "线下支付"
	}
	res.Data.Brief = reply.Brief
	res.Data.Status = GetStatus(reply.Status)
	c.Data["json"] = res
	c.ServeJSON()
}
