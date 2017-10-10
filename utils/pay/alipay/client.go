package alipay

import (
	"github.com/smartwalle/alipay"
	"net/url"
	"fmt"
)

var client = alipay.New(appID, partnerID, publicKey, privateKey, true)

func TradePagePay(Subject, OutTradeNo, TotalAmount, QRCodeWidth string) (*url.URL, error) {
	fmt.Println("TradePagePay=OutTradeNo", OutTradeNo)
	res, err := client.TradePagePay(alipay.AliPayTradePagePay{
		Subject:     Subject,
		OutTradeNo:  OutTradeNo,
		TotalAmount: TotalAmount,
		ProductCode: "FAST_INSTANT_TRADE_PAY",
		QRPayMode:   "4",
		QRCodeWidth: QRCodeWidth,
	})
	return res, err
}

func TradeClose(OutTradeNo string) (*alipay.AliPayTradeCloseResponse, error) {
	fmt.Println("TradeClose=OutTradeNo", OutTradeNo)
	res, err := client.TradeClose(alipay.AliPayTradeClose{
		OutTradeNo: OutTradeNo,
	})
	return res, err
}

func TradeQuery(OutTradeNo string) (*alipay.AliPayTradeQueryResponse, error) {
	fmt.Println("TradeQuery=OutTradeNo", OutTradeNo)
	res, err := client.TradeQuery(alipay.AliPayTradeQuery{
		OutTradeNo: OutTradeNo,
	})
	return res, err
}
