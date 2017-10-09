package alipay

import (
	"net/url"
	"time"
	"sort"
	"encoding/json"
)

func CommonParams() url.Values {
	var p url.Values = make(url.Values)
	p.Add("app_id", APPID)
	//p.Add("method", "")
	p.Add("format", FORMAT)
	p.Add("charset", CHARSET)
	p.Add("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	p.Add("version", VERSION)
	p.Add("sign_type", SIGN_TYPE)
	//p.Add("app_auth_token", "")
	//p.Add("return_url", "")
	//p.Add("notify_url", "")
	//p.Add("biz_content", "")
	//p.Add("sign", "")
	return p
}

func URLValues(action string, ps interface{}) url.Values {
	p := CommonParams()
	var keys = make([]string, 0, 0)
	for key, _ := range p {
		keys = append(keys, key)
	}
	p.Add("method", action)
	var bytes, err = json.Marshal(ps)
	if err == nil {
		p.Add("biz_content", string(bytes))
	}

	if action == AlipayTradePagePayAction {
		p.Add("return_url", ReturnUrl)
		p.Add("notify_url", NotifyUrl)
	} else if action == AlipayTradeRefundAction {
		p.Add("app_auth_token", "")
	} else if action == AlipayTradeQueryAction {
		p.Add("app_auth_token", "")
	} else if action == AlipayTradeCloseAction {
		p.Add("app_auth_token", "")
		p.Add("notify_url", NotifyUrl)
	} else if action == AlipayTradeFastpayRefundQueryAction {

	} else if action == AlipayDataDataserviceBillDownloadurlQueryAction {

	}

	sort.Strings(keys)
	if p.Get("sign_type") == SIGN_TYPE {
		p.Add("sign", sign_rsa2(keys, p, []byte(APP_PRIVATE_KEY)))
	} else {
		p.Add("sign", sign_rsa(keys, p, []byte(APP_PRIVATE_KEY)))
	}
	return p
}
