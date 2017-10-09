package alipay
//
//import (
//	"fmt"
//	"testing"
//)
//
//func Test_TradePagePay(t *testing.T) {
//	alipay := DefaultAlipayClient()
//	response, err := alipay.TradePagePay(&AlipayTradePagePayRequest{
//		out_trade_no:            "20150320010101001",
//		product_code:            "FAST_INSTANT_TRADE_PAY",
//		total_amount:            "88.88",
//		subject:                 "Iphone6 16G",
//		body:                    "",
//		goods_detail:            "",
//		passback_params:         "",
//		extend_params:           "",
//		goods_type:              "",
//		timeout_express:         "",
//		enable_pay_channels:     "",
//		disable_pay_channels:    "",
//		auth_token:              "",
//		qr_pay_mode:             "",
//		qrcode_width:            "",
//		sys_service_provider_id: "",
//		hb_fq_num:               "",
//		hb_fq_seller_percent:    "",
//	})
//	if err == nil {
//		fmt.Println(response)
//	}
//}
//
//func Test_TradeRefund(t *testing.T) {
//	alipay := DefaultAlipayClient()
//	response, err := alipay.TradeRefund(&AlipayTradeRefundRequest{
//		out_trade_no:   "",
//		trade_no:       "",
//		refund_amount:  "",
//		refund_reason:  "",
//		out_request_no: "",
//		operator_id:    "",
//		store_id:       "",
//		terminal_id:    "",
//	})
//	if err == nil {
//		fmt.Println(response)
//	}
//}
//
//func Test_TradeQuery(t *testing.T) {
//	alipay := DefaultAlipayClient()
//	response, err := alipay.TradeQuery(&AlipayTradeQueryRequest{
//		out_trade_no: "",
//		trade_no:     "",
//	})
//	if err == nil {
//		fmt.Println(response)
//	}
//}
//
//func Test_TradeClose(t *testing.T) {
//	alipay := DefaultAlipayClient()
//	response, err := alipay.TradeClose(&AlipayTradeCloseRequest{
//		out_trade_no: "",
//		trade_no:     "",
//		operator_id:  "",
//	})
//	if err == nil {
//		fmt.Println(response)
//	}
//}
//
//func Test_TradeFastpayRefundQuery(t *testing.T) {
//	alipay := DefaultAlipayClient()
//	response, err := alipay.TradeFastpayRefundQuery(&AlipayTradeFastpayRefundQueryRequest{
//		out_trade_no:   "",
//		trade_no:       "",
//		out_request_no: "",
//	})
//	if err == nil {
//		fmt.Println(response)
//	}
//}
//
//func Test_DataDataserviceBillDownloadurlQuery(t *testing.T) {
//	alipay := DefaultAlipayClient()
//	response, err := alipay.DataDataserviceBillDownloadurlQuery(&AlipayDataDataserviceBillDownloadurlQueryRequest{
//		bill_type: "trade",
//		bill_date: "2016-04-05",
//	})
//	if err == nil {
//		fmt.Println(response)
//	}
//}
