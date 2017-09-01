package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/center"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/account"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/utils"
	"dev.cloud.360baige.com/utils/pay/wechat"
	"strings"
	"dev.cloud.360baige.com/log"
	"encoding/xml"
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
)

// Account API
type AccountController struct {
	beego.Controller
}

// @router /unifiedorder [get]
func (c *AccountController) UnifiedOrder() {
	remoteAddr := strings.Split(c.Ctx.Request.RemoteAddr, ":")
	log.Println("remoteAddr:", remoteAddr[0])
	params := map[string]interface{}{
		"appid":            "wxc2cbb6b5a46fc13e",
		"body":             "服务费1元",
		"mch_id":           1457803302,
		"nonce_str":        utils.RandomString(20),
		"notify_url":       "http://wxpay.figool.cn/account",
		"trade_type":       "NATIVE",
		"spbill_create_ip": remoteAddr[0],
		"total_fee":        1,
		"out_trade_no":     "20170829041525144258",
	}
	unifyOrder := wechat.UnifyOrderRequest{
		Appid:            fmt.Sprintf("%v", params["appid"]),
		Body:             fmt.Sprintf("%v", params["body"]),
		Mch_id:           fmt.Sprintf("%v", params["mch_id"]),
		Nonce_str:        fmt.Sprintf("%v", params["nonce_str"]),
		Notify_url:       fmt.Sprintf("%v", params["notify_url"]),
		Trade_type:       fmt.Sprintf("%v", params["trade_type"]),
		Spbill_create_ip: fmt.Sprintf("%v", params["spbill_create_ip"]),
		Total_fee:        params["total_fee"].(int),
		Out_trade_no:     fmt.Sprintf("%v", params["out_trade_no"]),
		Sign:             "2323",//Sign(params, "17DA9CAF1E16CF508609FEB6944CE97A"), //39f22f62edf5163a8efc107d63e81c9c 17DA9CAF1E16CF508609FEB6944CE97A
	}
	urlStr := "https://api.mch.weixin.qq.com/pay/unifiedorder"
	//url := Url(params, prefix)
	bytes_req, _ := xml.Marshal(unifyOrder)
	str_req := string(bytes_req)
	str_req = strings.Replace(str_req, "UnifyOrder", "xml", -1)
	bytes_req = []byte(str_req)
	req, _ := http.NewRequest("POST", urlStr, bytes.NewReader(bytes_req))
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")
	client := http.Client{}
	resp, _ := client.Do(req)
	xmlResp := wechat.UnifyOrderResponse{}
	respBytes, _ := ioutil.ReadAll(resp.Body)
	_ = xml.Unmarshal(respBytes, &xmlResp)
	//log.Println("xmlResp:", xmlResp)
	log.Println("xmlResp.Sign:", xmlResp.Sign)
	log.Println("xmlResp.Nonce_str:", xmlResp.Nonce_str)
	log.Println("xmlResp.Mch_id:", xmlResp.Mch_id)
	log.Println("xmlResp.Return_code:", xmlResp.Return_code)
	log.Println("xmlResp.Return_msg:", xmlResp.Return_msg)
	log.Println("xmlResp.Result_code:", xmlResp.Result_code)
	log.Println("xmlResp.Appid:", xmlResp.Appid)
	log.Println("xmlResp.Prepay_id:", xmlResp.Prepay_id)
	log.Println("xmlResp.Trade_type:", xmlResp.Trade_type)
	log.Println("xmlResp.Code_url:", xmlResp.Code_url)
	log.Println("xmlResp.Openid:", xmlResp.Openid)
}

// @router /qr [get]
func (c *AccountController) Qr() {
	//params := map[string]interface{}{
	//	"sign":       "",
	//	"appid":      "wxc2cbb6b5a46fc13e",
	//	"mch_id":     "1457803302",
	//	"product_id": "product_id",
	//	"time_stamp": strconv.FormatInt(utils.CurrentTimestamp()/1000, 10),
	//	"nonce_str":  utils.RandomString(20),
	//}
	//params["sign"] = Sign(params)
	//url := Url(params, "weixin://wxpay/bizpayurl?")
	//log.Println("url:", url)
	c.Ctx.Output.Body(utils.Qr("weixin://wxpay/bizpayurl?pr=zpQVnCK", 250))
}

// @Title 账户统计接口
// @Description 账户统计接口
// @Success 200 {"code":200,"message":"获取账务统计信息成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @router /statistics [post]
func (c *AccountController) Statistics() {
	type data AccountStatisticsResponse
	accessToken := c.GetString("accessToken")
	currentTimestamp := utils.CurrentTimestamp()
	if accessToken == "" {
		c.Data["json"] = data{Code: ErrorSystem, Message: "访问令牌无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "access_token", Val: accessToken },
			action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTimestamp },
		},
		Fileds: []string{"id", "user_id", "company_id", "type"},
	}, &replyUserPosition)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "验证访问令牌失效"}
		c.ServeJSON()
		return
	}

	if replyUserPosition.Id == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌失效"}
		c.ServeJSON()
		return
	}

	var replyAccount account.Account
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Account", "FindByCond", &action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "company_id", Val: replyUserPosition.CompanyId },
			action.CondValue{Type: "And", Key: "user_id", Val: replyUserPosition.UserId },
			action.CondValue{Type: "And", Key: "user_position_id", Val: replyUserPosition.Id },
			action.CondValue{Type: "And", Key: "user_position_type", Val: replyUserPosition.Type },
		},
		Fileds: []string{"id", "balance"},
	}, &replyAccount)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取账务统计信息失败1"}
		c.ServeJSON()
		return
	}

	var replyAccountItemList []account.AccountItem
	err = client.Call(beego.AppConfig.String("EtcdURL"), "AccountItem", "ListByCond", &action.ListByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "account_id", Val: replyAccount.Id },
			action.CondValue{Type: "And", Key: "status__gt", Val: -1 },
		},
		Cols: []string{"amount", "balance"},
	}, &replyAccountItemList)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取账务统计信息失败2"}
		c.ServeJSON()
		return
	}

	var inAccount, outAccount float64
	for _, accountItem := range replyAccountItemList {
		if accountItem.Amount > 0 {
			inAccount += accountItem.Amount
		} else {
			outAccount += accountItem.Amount
		}
	}

	c.Data["json"] = data{Code: Normal, Message: "获取账务统计信息成功", Data: AccountStatistics{
		Balance:    utils.Amount(replyAccount.Balance),
		InAccount:  utils.Amount(inAccount),
		OutAccount: utils.Amount(outAccount),
	}}
	c.ServeJSON()
	return
}
