package schoolfee

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.cloud.360baige.com/utils/pay/wechat"
	"fmt"
	"strings"
	//"dev.cloud.360baige.com/utils"
	. "dev.model.360baige.com/http/mobile/schoolfee"
	"dev.cloud.360baige.com/utils"
	"dev.model.360baige.com/models/schoolfee"
	"dev.model.360baige.com/models/order"
	"dev.model.360baige.com/action"
)

// Order API
type OrderController struct {
	beego.Controller
}

// @Title 查询缴费历史接口
// @Description Limit Project List 查询缴费历史接口
// @Success 200 {"code":200,"message":"查询缴费历史成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   searchKey     query   string true       "查询值"
// @Failure 400 {"code":400,"message":"查询缴费历史失败"}
// @router /add [*]
func (c *OrderController) Add() {
	type data OrderAddResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
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

	recordIds := c.GetString("recordIds", "")
	if recordIds == "" {
		// 自由支付
		projectId, _ := c.GetInt64("projectId", 0)
		//var replyRecord schoolfee.Record
		//err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "Add", &schoolfee.Record{
		//	CreateTime: currentTimestamp,
		//	UpdateTime: currentTimestamp,
		//	CompanyId:  replyUserPosition.CompanyId,
		//	ProjectId:  projectId,
		//	Name:       "name",
		//	ClassName:  "无",
		//	IdCard:     "",
		//	Num:        "",
		//	Phone:      "",
		//	Price:      1,
		//	IsFee:      0,
		//	FeeTime:    0,
		//	Desc:       "",
		//	Status:     0,
		//}, &replyRecord)

		// 下单
		var replyOrder order.Order
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "Add", &order.Order{
			CreateTime:       currentTimestamp,
			UpdateTime:       currentTimestamp,
			CompanyId:        replyUserPosition.CompanyId,
			UserId:           replyUserPosition.UserId,
			UserPositionType: replyUserPosition.Type,
			UserPositionId:   replyUserPosition.Id,
			Code:             "",
			SubCode:          "",
			Price:            1,
			Num:              1,
			TotalPrice:       1,
			ProductType:      1,
			ProductId:        projectId,
			Image:            "schoolfee",
			PayType:          1,
			Brief:            "",
			TradeType:        "MWEB",
			Status:           0,
		}, &replyOrder)
		//replyOrder.Code

		//Id               int64    `db:"id" json:"id"`                              // 主键
		//Image            string   `db:"image" json:"image"`                        // 订单图片
		//PayType          int     `db:"pay_type" json:"payType"`                    // 支付方式
		//Brief            string   `db:"brief" json:"brief"`                        // 详情
		//Status           int     `db:"status" json:"status"`                       // 订单状态
		//TradeType        string   `db:"trade_type" json:"tradeType"`               // 交易类型
		//PrepayId         string   `db:"prepay_id" json:"prepayId"`                 // 预支付交易会话标识
		//CodeUrl          string   `db:"code_url" json:"codeUrl"`                   // 二维码链接
		//Openid           string   `db:"openid" json:"openid"`                      // 用户标识

	} else {
		// 限制支付
		var replyListOfRecord []schoolfee.Record
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Record", "ListByCond", action.ListByCond{
			CondList: []action.CondValue{
				action.CondValue{"And", "company_id", replyUserPosition.CompanyId},
				action.CondValue{"And", "is_fee", 0},
				action.CondValue{"And", "id__in", recordIds},
			},
		}, &replyListOfRecord)
		fmt.Println("replyListOfRecord:", replyListOfRecord)
		// 下单

		remoteAddr := c.Ctx.Request.Header.Get("X-Real-IP")
		if remoteAddr == "" {
			remoteAddr = strings.Split(c.Ctx.Request.RemoteAddr, ":")[0]
		}
		ou := utils.RandomString(10)
		var totolFee int64 = 0
		body := "测试商品"
		unifyOrderResponse, err := wechat.UnifiedOrder(remoteAddr, body, ou, "MWEB", totolFee)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "统一下单失败"}
			c.ServeJSON()
			return
		}
		fmt.Println("unifyOrderResponse:", unifyOrderResponse)

	}

	//projectId := c.GetString("projectId")
	//name := c.GetString("name")
	//phone := c.GetString("phone")
	//id_card := c.GetString("id_card")
	//currentTimestamp := utils.CurrentTimestamp()
	//accessToken := c.GetString("accessToken")
	//applicationTplId, _ := c.GetInt64("applicationTplId")
	num, _ := c.GetInt64("num", 0)
	fmt.Println("商品数量", num)
	payType, _ := c.GetInt("payType", 1)
	fmt.Println("支付方式：1微信 2支付宝", payType)
	//tradeType := c.GetString("tradeType", "MWEB")
	orderId, _ := c.GetInt64("orderId", 0)
	fmt.Println("订单编号", orderId)

	//err := utils.Unable(map[string]string{"accessToken": "string:true", "payType": "int:true", "applicationTplId": "int:true", "num": "int:true"}, c.Ctx.Input)
	//if err != nil {
	//	c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
	//	c.ServeJSON()
	//	return
	//}

	//c.Ctx.ResponseWriter.Write([]byte("<a href='" + unifyOrderResponse.Mweb_url + "'>调转</a>"))
	//fmt.Println("unifyOrderResponse:", unifyOrderResponse)
	//c.Data["json"] = unifyOrderResponse.Mweb_url
	//c.ServeJSON()
}
