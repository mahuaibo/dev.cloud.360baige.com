package alipay

import (
	"net/http"
)

type AlipayClient struct {
	client *http.Client
	//URL               string
	//APPID             string
	//APP_PRIVATE_KEY   string
	//FORMAT            string
	//CHARSET           string
	//ALIPAY_PUBLIC_KEY string
	//SIGN_TYPE         string

	//partnerId         string
	//publicKey         string
	//privateKey        string
	//AliPayPublicKey   string
	//SignType          string
}

func DefaultAlipayClient() *AlipayClient {
	return &AlipayClient{
		client: http.DefaultClient,
		//URL:               URL,
		//APPID:             APPID,
		//APP_PRIVATE_KEY:   APP_PRIVATE_KEY,
		//FORMAT:            FORMAT,
		//CHARSET:           CHARSET,
		//ALIPAY_PUBLIC_KEY: ALIPAY_PUBLIC_KEY,
		//SIGN_TYPE:         SIGN_TYPE,
	}
}
