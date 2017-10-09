package alipay

import (
	"github.com/smartwalle/alipay"
	"net/url"
)

var client = alipay.New(appID, partnerID, publicKey, privateKey, true)

func TradePagePay(Subject, OutTradeNo, TotalAmount, QRCodeWidth string) (*url.URL, error) {

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
