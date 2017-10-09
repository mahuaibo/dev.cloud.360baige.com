package alipay

import (
	"net/url"
)

type AlipayRequest interface {
	URL() string
	Action() string
	Method() string
	Params() map[string]string
}

/* 1 */
type AlipayTradePagePayRequest struct {
	Out_trade_no            string    `json:"out_trade_no,omitempty"`
	Product_code            string    `json:"product_code,omitempty"`
	Total_amount            string    `json:"total_amount,omitempty"`
	Subject                 string    `json:"subject,omitempty"`
	body                    string    `json:"body,omitempty"`
	goods_detail            string    `json:"goods_detail,omitempty"`
	passback_params         string    `json:"passback_params,omitempty"`
	extend_params           string    `json:"extend_params,omitempty"`
	goods_type              string    `json:"goods_type,omitempty"`
	timeout_express         string    `json:"timeout_express,omitempty"`
	enable_pay_channels     string    `json:"enable_pay_channels,omitempty"`
	disable_pay_channels    string    `json:"disable_pay_channels,omitempty"`
	auth_token              string    `json:"auth_token,omitempty"`
	qr_pay_mode             string    `json:"qr_pay_mode,omitempty"`
	qrcode_width            string    `json:"qrcode_width,omitempty"`
	sys_service_provider_id string    `json:"sys_service_provider_id,omitempty"`
	hb_fq_num               string    `json:"hb_fq_num,omitempty"`
	hb_fq_seller_percent    string    `json:"hb_fq_seller_percent,omitempty"`
}

func (this *AlipayTradePagePayRequest) URL() string {
	return URL
}

func (this *AlipayTradePagePayRequest) Action() string {
	return AlipayTradePagePayAction
}

func (this *AlipayTradePagePayRequest) Method() string {
	return POST_METHOD
}

func (this *AlipayTradePagePayRequest) Params() url.Values {
	return URLValues(this.Action(), this)
}

/* 2 */
type AlipayTradeRefundRequest struct {
	out_trade_no   string    `json:"out_trade_no,omitempty"`
	trade_no       string    `json:"trade_no,omitempty"`
	refund_amount  string    `json:"refund_amount,omitempty"`
	refund_reason  string    `json:"refund_reason,omitempty"`
	out_request_no string    `json:"out_request_no,omitempty"`
	operator_id    string    `json:"operator_id,omitempty"`
	store_id       string    `json:"store_id,omitempty"`
	terminal_id    string    `json:"terminal_id,omitempty"`
}

func (this *AlipayTradeRefundRequest) URL() string {
	return URL
}

func (this *AlipayTradeRefundRequest) Action() string {
	return AlipayTradeRefundAction
}

func (this *AlipayTradeRefundRequest) Method() string {
	return POST_METHOD
}

func (this *AlipayTradeRefundRequest) Params() url.Values {
	return URLValues(this.Action(), this)
}

/* 3 */
type AlipayTradeFastpayRefundQueryRequest struct {
	trade_no       string    `json:"trade_no,omitempty"`
	out_trade_no   string    `json:"out_trade_no,omitempty"`
	out_request_no string    `json:"out_request_no,omitempty"`
}

func (this *AlipayTradeFastpayRefundQueryRequest) URL() string {
	return URL
}

func (this *AlipayTradeFastpayRefundQueryRequest) Action() string {
	return AlipayTradeFastpayRefundQueryAction
}

func (this *AlipayTradeFastpayRefundQueryRequest) Method() string {
	return POST_METHOD
}

func (this *AlipayTradeFastpayRefundQueryRequest) Params() url.Values {
	return URLValues(this.Action(), this)
}

/* 4 */
type AlipayTradeQueryRequest struct {
	trade_no     string    `json:"trade_no,omitempty"`
	out_trade_no string    `json:"out_trade_no,omitempty"`
}

func (this *AlipayTradeQueryRequest) URL() string {
	return URL
}

func (this *AlipayTradeQueryRequest) Action() string {
	return AlipayTradeQueryAction
}

func (this *AlipayTradeQueryRequest) Method() string {
	return POST_METHOD
}

func (this *AlipayTradeQueryRequest) Params() url.Values {
	return URLValues(this.Action(), this)
}

/* 5 */
type AlipayTradeCloseRequest struct {
	trade_no     string    `json:"trade_no,omitempty"`
	out_trade_no string    `json:"out_trade_no,omitempty"`
	operator_id  string    `json:"operator_id,omitempty"`
}

func (this *AlipayTradeCloseRequest) URL() string {
	return URL
}

func (this *AlipayTradeCloseRequest) Action() string {
	return AlipayTradeCloseAction
}

func (this *AlipayTradeCloseRequest) Method() string {
	return POST_METHOD
}

func (this *AlipayTradeCloseRequest) Params() url.Values {
	return URLValues(this.Action(), this)
}

/* 6 */
type AlipayDataDataserviceBillDownloadurlQueryRequest struct {
	bill_type string    `json:"bill_type,omitempty"`
	bill_date string    `json:"bill_date,omitempty"`
}

func (this *AlipayDataDataserviceBillDownloadurlQueryRequest) URL() string {
	return URL
}

func (this *AlipayDataDataserviceBillDownloadurlQueryRequest) Action() string {
	return AlipayDataDataserviceBillDownloadurlQueryAction
}

func (this *AlipayDataDataserviceBillDownloadurlQueryRequest) Method() string {
	return POST_METHOD
}

func (this *AlipayDataDataserviceBillDownloadurlQueryRequest) Params() url.Values {
	return URLValues(this.Action(), this)
}
