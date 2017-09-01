package wechat

import (
	"sort"
	"fmt"
	"strings"
	"crypto/md5"
	"encoding/hex"
	"dev.cloud.360baige.com/utils"
	"encoding/xml"
	"net/http"
	"bytes"
	"io/ioutil"
)

var (
	appid        = "wxc2cbb6b5a46fc13e"
	mch_id       = 1457803302
	notify_url   = "http://wxpay.figool.cn/account"
	unifiedorder = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	orderquery   = "https://api.mch.weixin.qq.com/pay/orderquery"
	key          = "17DA9CAF1E16CF508609FEB6944CE97A"
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

func UnifiedOrder(ip, body, out_trade_no, total_fee string) (UnifyOrderResponse, error) {
	xmlResp := UnifyOrderResponse{}
	params := map[string]interface{}{
		"appid":            appid,
		"mch_id":           mch_id,
		"nonce_str":        utils.RandomString(20),
		"trade_type":       "NATIVE",
		"body":             body,
		"notify_url":       notify_url,
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
		Sign:             Sign(params, key),
	}
	bytes_req, err := xml.Marshal(unifyOrder)
	if err != nil {
		return xmlResp, err
	}
	str_req := string(bytes_req)
	str_req = strings.Replace(str_req, "UnifyOrderRequest", "xml", -1)
	bytes_req = []byte(str_req)
	req, err := http.NewRequest("POST", unifiedorder, bytes.NewReader(bytes_req))
	if err != nil {
		return xmlResp, err
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return xmlResp, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return xmlResp, err
	}
	err = xml.Unmarshal(respBytes, &xmlResp)
	if err != nil {
		return xmlResp, err
	}
	return xmlResp, nil
}

type OrderQueryRequest struct {
	Appid        string `xml:"appid"`        //
	Mch_id       string `xml:"mch_id"`       //
	Nonce_str    string `xml:"nonce_str"`    //
	Out_trade_no string `xml:"out_trade_no"` //
	Sign         string `xml:"sign"`         //
}

type OrderQueryResponse struct {
	ReturnCode     string `xml:"return_code"`
	ReturnMsg      string `xml:"return_msg"`
	Appid          string `xml:"appid"`
	MchId          string `xml:"mch_id"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	ResultCode     string `xml:"result_code"`
	ErrCode        string `xml:"err_code"`
	ErrCodeDes     string `xml:"err_code_des"`
	TradeType      string `xml:"trade_type"`
	TradeState     string `xml:"trade_state"`
	BankType       string `xml:"bank_type"`
	TotalFee       string `xml:"total_fee"`
	CashFee        string `xml:"cash_fee"`
	TransactionId  string `xml:"transaction_id"`
	OutTradeNo     string `xml:"out_trade_no"`
	Attach         string `xml:"attach"`
	TimeEnd        string `xml:"time_end"`
	TradeStateDesc string `xml:"trade_state_desc"`
}

func OrderQuery(out_trade_no string) (OrderQueryResponse, error) {
	xmlResp := OrderQueryResponse{}
	params := map[string]interface{}{
		"appid":        appid,
		"mch_id":       mch_id,
		"nonce_str":    utils.RandomString(20),
		"out_trade_no": out_trade_no,
	}
	orderQuery := OrderQueryRequest{
		Appid:        fmt.Sprintf("%v", params["appid"]),
		Mch_id:       fmt.Sprintf("%v", params["mch_id"]),
		Nonce_str:    fmt.Sprintf("%v", params["nonce_str"]),
		Out_trade_no: fmt.Sprintf("%v", params["out_trade_no"]),
		Sign:         Sign(params, key),
	}
	bytes_req, err := xml.Marshal(orderQuery)
	if err != nil {
		return xmlResp, err
	}
	str_req := string(bytes_req)
	str_req = strings.Replace(str_req, "OrderQueryRequest", "xml", -1)
	bytes_req = []byte(str_req)
	req, err := http.NewRequest("POST", unifiedorder, bytes.NewReader(bytes_req))
	if err != nil {
		return xmlResp, err
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return xmlResp, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return xmlResp, err
	}
	err = xml.Unmarshal(respBytes, &xmlResp)
	if err != nil {
		return xmlResp, err
	}
	return xmlResp, nil
}
