package alipay

const (
	AlipayTradePagePayAction                        = "alipay.trade.page.pay"
	AlipayTradeRefundAction                         = "alipay.trade.refund"
	AlipayTradeFastpayRefundQueryAction             = "alipay.trade.fastpay.refund.query"
	AlipayTradeQueryAction                          = "alipay.trade.query"
	AlipayTradeCloseAction                          = "alipay.trade.close"
	AlipayDataDataserviceBillDownloadurlQueryAction = "alipay.data.dataservice.bill.downloadurl.query"
)

const (
	GET_METHOD  = "GET"
	POST_METHOD = "POST"
)

const (
	URL               = "https://openapi.alipay.com/gateway.do" //支付宝网关（固定）
	APPID             = "2017071307740722"                      //即创建应用后生成
	FORMAT            = "json"                                  //参数返回格式，只支持json	json（固定）
	CHARSET           = "UTF-8"                                 //编码集，支持GBK/UTF-8	开发者根据实际工程编码配置
	ALIPAY_PUBLIC_KEY = "tGPCd9uZD8Qo6C5ehbK52Q=="              //支付宝公钥，由支付宝生成
	SIGN_TYPE         = "RSA2"                                  //商户生成签名字符串所使用的签名算法类型，目前支持RSA2和RSA，推荐使用RSA2
	VERSION           = "1.0"
)

const (
	APP_PRIVATE_KEY = "MIIEogIBAAKCAQEAytpPmnL39SJzDC0GkkUiDHz9jmvQR2n0HuRxqmw7QVJkHJlNKPy3sUDK3vIHglgoMGUDU6ODA+H/eNIlUeu1dSWC0qwF6fwDlFJHOddX46mIGFBVy6HGYXWGlsFHV+D+s3jyL7xbsEHG/HP+v5dYz9wK2H7CLgAQ6HvMxz+XpkKUFHYzuCNbqiHpA2l5irgogtl4AmYowN5lTop4oiG48aEGcdFwny7P66rrAMjITYaS9PX0iNE0g8ncsMysTSufcAD8kL68v+GT06haqNTfbbQSoEl+xjOLmONOt9Iu3oxVOAqOM4zCF0j4oh+laLHllFQW6ekzjG/QXNNrkYok2QIDAQABAoIBABo14xs9x7Qw41SrbLHxpNigPdLtM1hG5HgpZFZ07aMfFjhrxoCJRuLsUEpGU5oP8gFuy+M+uWsDBJOD87aGEkg0tJasC4eUPJIpn3Jl1MFh4mfh2XQaTxvAp8dK6gD83WwrMH/igqZfmWp9QmlXEO6qq+wVVNnEwGqJtIf6O4oLzSBsgAe4L6kysx446omgm+JdmWphOiLCrso2iYlsjQB4yS1VvdPutX4pZwAtl89LWKbgvSd+fei7rnebINmASQxFaHu03FlqO5TLxZD7JNlOutbvf3tgC+OKYKK4As4mrzjYWCkabuRtq9pgNKFsGSJ0e3jAKs9KCyirJzCIXUECgYEA+MuTLRqsqcljU8wGsnsCMvJD92BCJnOGusoU2RfEdPe9yeTaYCjPjO3N+IEOPMg47hlyJ/06z5rPrYgRFtf3pRZlLbVb+9ihTBZx40FVyp72kI9Hh9GZ6VajLBLX3InNG4A/ikZlma5mhu5+uyKt7AGj9xq+iwcHj3uFuCss0gUCgYEA0LolKJavbYNRf4rxz1QPRxmMUZgGH+3oPxuyk8/VZHlpDTRb8ReMzsFu/ZOK5o0RbK4/HLZI9Y+Qo9/EJ+FGw+qlN6p9t+LR2EZcun7QScJrJzNmn3BxRE/YvQiZ5cuScP7KI3G67eqAAIN/kdPW6caEK4B0b7PRJJ9UPC0uG8UCgYAtsRWCfxeexwGa7il8teKdgKjC0cbUUPs5asuRYzANW0Jbxc/lQRl9BF+DeBApUYxDDiFM/tDCN+hUMl0RGPC+PPKwBlKyWgKleqnH1sPuxmr5+ZZldzURCXxGJ9/E/PnSRydkObGHG+Rwe3SC4ceXRGXch+jel5fn3gOc4zEEEQKBgCtaPlWL1qv8VnUTOt5BK1stJ0PbO4puM4rICfNBe6T+wp8HfQE6TviynIb2micArdnQ5zLjeYvnYbdnxqox1CzlE6PYOXx0E/nDw4bIyCJJy0+9EVeUVzJFE1NP8gcUnTny8bEi8hxwVF3G8jwrr7umi9UtPQHma04b+bb1+1mxAoGAJmd2MXeNPOFW5Ul8QcuQPwHeecp55PbfNIoxUTaxoIwJsl7Zb/CxYdZ2T/7t3bQ2rTWk+q8QHGPWWW015N/P4NE1Fvm/rKjaSObuSTq0jCu94zyCyk8qy1G5OrPqoZ3LJ2RQ7JhEdmH9akHgTezdjYw3TQxyDuAGBrcH3TcUZCQ=" //开发者私钥，由开发者自己生成
)

const (
	PartnerId = "2088721343882880"
	ReturnUrl = "http://cloud.360baige.com/callback"
	NotifyUrl = "http://cloud.360baige.com/callback"
)
