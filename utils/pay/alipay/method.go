package alipay

import (
	"net/http"
	"strings"
	"io"
	"io/ioutil"
	"fmt"
	"encoding/json"
)

// 1.统一收单下单并支付页面接口
func (this *AlipayClient) TradePagePay(request AlipayTradePagePayRequest) (*AlipayTradePagePayResponse, error) {
	fmt.Println("1:", request)
	var buf io.Reader
	if &request != nil {
		buf = strings.NewReader(request.Params().Encode())
	}
	req, err := http.NewRequest(POST_METHOD, URL, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	rep, err := this.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rep.Body.Close()

	data, err := ioutil.ReadAll(rep.Body)
	fmt.Println(data)
	if err != nil {
		return nil, err
	}
	var response *AlipayTradePagePayResponse
	err = json.Unmarshal(data, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// 2.统一收单交易退款接口
func (this *AlipayClient) TradeRefund(param *AlipayTradeRefundRequest) (*AlipayTradeRefundResponse, error) {
	return nil, nil
}

// 3.统一收单交易退款查询接口
func (this *AlipayClient) TradeFastpayRefundQuery(param *AlipayTradeFastpayRefundQueryRequest) (*AlipayTradeFastpayRefundQueryResponse, error) {
	return nil, nil
}

// 4.统一收单线下交易查询接口
func (this *AlipayClient) TradeQuery(param *AlipayTradeQueryRequest) (*AlipayTradeQueryResponse, error) {
	return nil, nil
}

// 5.统一收单交易关闭接口
func (this *AlipayClient) TradeClose(param *AlipayTradeCloseRequest) (*AlipayTradeCloseResponse, error) {
	return nil, nil
}

// 6.查询对账单下载地址
func (this *AlipayClient) DataDataserviceBillDownloadurlQuery(param *AlipayDataDataserviceBillDownloadurlQueryRequest) (*AlipayDataDataserviceBillDownloadurlQueryResponse, error) {
	return nil, nil
}
