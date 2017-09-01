package wechat

import (
	"sort"
	"fmt"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"github.com/silenceper/wechat"
	"dev.cloud.360baige.com/utils"
)

var (
	unifiedorder = "https://api.mch.weixin.qq.com/pay/unifiedorder"
)

type UnifyOrderRequest struct {
	Appid            string `xml:"appid"`
	Body             string `xml:"body"`
	Mch_id           string `xml:"mch_id"`
	Nonce_str        string `xml:"nonce_str"`
	Notify_url       string `xml:"notify_url"`
	Trade_type       string `xml:"trade_type"`
	Spbill_create_ip string `xml:"spbill_create_ip"`
	Total_fee        int    `xml:"total_fee"`
	Out_trade_no     string `xml:"out_trade_no"`
	Sign             string `xml:"sign"`
}

type UnifyOrderResponse struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
	Appid       string `xml:"appid"`
	Mch_id      string `xml:"mch_id"`
	Nonce_str   string `xml:"nonce_str"`
	Sign        string `xml:"sign"`
	Result_code string `xml:"result_code"`
	Prepay_id   string `xml:"prepay_id"`
	Trade_type  string `xml:"trade_type"`
	Code_url    string `xml:"code_url"`
	Openid      string `xml:"openid"`
}

func UnifiedOrder(ip, body, out_trade_no string, total_fee int64) error {
	//remoteAddr := strings.Split(c.Ctx.Request.RemoteAddr, ":")
	//log.Println("remoteAddr:", remoteAddr[0])
	params := map[string]interface{}{
		"appid":            "wxc2cbb6b5a46fc13e",
		"mch_id":           1457803302,
		"nonce_str":        utils.RandomString(20),
		"trade_type":       "NATIVE",
		"body":             body,
		"notify_url":       "http://wxpay.figool.cn/account",
		"spbill_create_ip": ip,
		"total_fee":        total_fee,
		"out_trade_no":     out_trade_no,
	}
	unifyOrder := UnifyOrderRequest{
		Appid:            fmt.Sprintf("%v", params["appid"]),
		Body:             fmt.Sprintf("%v", params["body"]),
		Mch_id:           fmt.Sprintf("%v", params["mch_id"]),
		Nonce_str:        fmt.Sprintf("%v", params["nonce_str"]),
		Notify_url:       fmt.Sprintf("%v", params["notify_url"]),
		Trade_type:       fmt.Sprintf("%v", params["trade_type"]),
		Spbill_create_ip: fmt.Sprintf("%v", params["spbill_create_ip"]),
		Total_fee:        params["total_fee"].(int),
		Out_trade_no:     fmt.Sprintf("%v", params["out_trade_no"]),
		Sign:             Sign(params, "17DA9CAF1E16CF508609FEB6944CE97A"), //39f22f62edf5163a8efc107d63e81c9c 17DA9CAF1E16CF508609FEB6944CE97A
	}
	//urlStr := "https://api.mch.weixin.qq.com/pay/unifiedorder"
	fmt.Println(unifyOrder)
	return nil
}

func Url(params map[string]interface{}, prefix string) string {
	var signStrings string
	sep := ""
	for k, v := range params {
		value := fmt.Sprintf("%v", v)
		if value != "" {
			signStrings = signStrings + sep + k + "=" + value
			sep = "&"
		}
	}
	return prefix + signStrings
}

func Sign(params map[string]interface{}, key ...string) string {
	//STEP 1, 对key进行升序排序.
	sorted_keys := make([]string, 0)
	for k, _ := range params {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)
	fmt.Println("sorted_keys:", sorted_keys)
	// STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	sep := ""
	for i, k := range sorted_keys {
		if i > 0 {
			sep = "&"
		}
		value := fmt.Sprintf("%v", params[k])
		if value != "" {
			signStrings = signStrings + sep + k + "=" + value
		}
	}
	fmt.Println("signStrings:", signStrings)
	//STEP3, 在键值对的最后加上key=API_KEY
	if key != nil {
		signStrings += "&key=" + key[0]
	}
	fmt.Println("signStrings:", signStrings)
	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}
