package alipay

// 统一收单下单并支付页面接口
type AlipayTradePagePayResponse struct {
	AlipayTradePagePay    `json:"code,omitempty"`
	sign string    `json:"sign,omitempty"`
}

type AlipayTradePagePay struct {
	/* 公共响应参数 */
	code     string    `json:"code,omitempty"`
	msg      string    `json:"msg,omitempty"`
	sub_code string    `json:"sub_code,omitempty"`
	sub_msg  string    `json:"sub_msg,omitempty"`
	/* 响应参数 */
}

// 统一收单交易退款接口
type AlipayTradeRefundResponse struct {
	AlipayTradeRefund    `json:"alipay_trade_refund_response,omitempty"`
	sign string    `json:"sign,omitempty"`
}

type AlipayTradeRefund struct {
	/* 公共响应参数 */
	code     string    `json:"code,omitempty"`
	msg      string    `json:"msg,omitempty"`
	sub_code string    `json:"sub_code,omitempty"`
	sub_msg  string    `json:"sub_msg,omitempty"`
	/* 响应参数 */
	trade_no                string    `json:"app_id,omitempty"`
	out_trade_no            string    `json:"app_id,omitempty"`
	buyer_logon_id          string    `json:"app_id,omitempty"`
	fund_change             string    `json:"app_id,omitempty"`
	refund_fee              string    `json:"app_id,omitempty"`
	gmt_refund_pay          string    `json:"app_id,omitempty"`
	refund_detail_item_list string    `json:"app_id,omitempty"`
	store_name              string    `json:"app_id,omitempty"`
	buyer_user_id           string    `json:"app_id,omitempty"`
}

// 统一收单交易退款查询接口
type AlipayTradeFastpayRefundQueryResponse struct {
	AlipayTradeFastpayRefundQuery    `json:"alipay_trade_fastpay_refund_query_response,omitempty"`
	sign string    `json:"sign,omitempty"`
}

type AlipayTradeFastpayRefundQuery struct {
	/* 公共响应参数 */
	code     string    `json:"code,omitempty"`
	msg      string    `json:"msg,omitempty"`
	sub_code string    `json:"sub_code,omitempty"`
	sub_msg  string    `json:"sub_msg,omitempty"`
	/* 响应参数 */
	trade_no       string    `json:"trade_no,omitempty"`
	out_trade_no   string    `json:"out_trade_no,omitempty"`
	out_request_no string    `json:"out_request_no,omitempty"`
	refund_reason  string    `json:"refund_reason,omitempty"`
	total_amount   string    `json:"total_amount,omitempty"`
	refund_amount  string    `json:"refund_amount,omitempty"`
}

// 统一收单线下交易查询接口
type AlipayTradeQueryResponse struct {
	AlipayTradeQuery    `json:"alipay_trade_query_response,omitempty"`
	sign string    `json:"sign,omitempty"`
}

type AlipayTradeQuery struct {
	/* 公共响应参数 */
	code     string    `json:"code,omitempty"`
	msg      string    `json:"msg,omitempty"`
	sub_code string    `json:"sub_code,omitempty"`
	sub_msg  string    `json:"sub_msg,omitempty"`
	/* 响应参数 */
	trade_no         string    `json:"trade_no,omitempty"`
	out_trade_no     string    `json:"out_trade_no,omitempty"`
	buyer_logon_id   string    `json:"buyer_logon_id,omitempty"`
	trade_status     string    `json:"trade_status,omitempty"`
	total_amount     string    `json:"total_amount,omitempty"`
	receipt_amount   string    `json:"receipt_amount,omitempty"`
	buyer_pay_amount string    `json:"buyer_pay_amount,omitempty"`
	point_amount     string    `json:"point_amount,omitempty"`
	invoice_amount   string    `json:"invoice_amount,omitempty"`
	send_pay_date    string    `json:"send_pay_date,omitempty"`
	store_id         string    `json:"store_id,omitempty"`
	terminal_id      string    `json:"terminal_id,omitempty"`
	fund_bill_list   []TradeFundBill    `json:"fund_bill_list,omitempty"`
	store_name       string    `json:"store_name,omitempty"`
	buyer_user_id    string    `json:"buyer_user_id,omitempty"`
}

type TradeFundBill struct {
	fund_channel string    `json:"fund_channel,omitempty"`
	amount       string    `json:"amount,omitempty"`
	real_amount  string    `json:"real_amount,omitempty"`
	fund_type    string    `json:"fund_type,omitempty"`
}

// 统一收单交易关闭接口
type AlipayTradeCloseResponse struct {
	AlipayTradeClose    `json:"alipay_trade_close_response,omitempty"`
	sign string    `json:"sign,omitempty"`
}

type AlipayTradeClose struct {
	/* 公共响应参数 */
	code     string    `json:"code,omitempty"`
	msg      string    `json:"msg,omitempty"`
	sub_code string    `json:"sub_code,omitempty"`
	sub_msg  string    `json:"sub_msg,omitempty"`
	/* 响应参数 */
	trade_no     string    `json:"trade_no,omitempty"`
	out_trade_no string    `json:"out_trade_no,omitempty"`
}

// 查询对账单下载地址
type AlipayDataDataserviceBillDownloadurlQueryResponse struct {
	AlipayDataDataserviceBillDownloadurlQuery    `json:"alipay_data_dataservice_bill_downloadurl_query_response,omitempty"`
	sign string    `json:"sign,omitempty"`
}

type AlipayDataDataserviceBillDownloadurlQuery struct {
	/* 公共响应参数 */
	code     string    `json:"code,omitempty"`
	msg      string    `json:"msg,omitempty"`
	sub_code string    `json:"sub_code,omitempty"`
	sub_msg  string    `json:"sub_msg,omitempty"`
	/* 响应参数 */
	bill_download_url string    `json:"bill_download_url,omitempty"`
}
