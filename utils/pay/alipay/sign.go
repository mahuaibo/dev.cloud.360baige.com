package alipay

import (
	"strings"
	"crypto"
	"encoding/base64"
	"dev.cloud.360baige.com/utils/pay/alipay/encoding"
	"net/url"
)

func sign_rsa2(keys []string, param url.Values, privateKey []byte) (s string) {
	if param == nil {
		param = make(url.Values, 0)
	}

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(param.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var src = strings.Join(pList, "&")
	var sig, err = encoding.SignPKCS1v15([]byte(src), privateKey, crypto.SHA256)
	if err != nil {
		return ""
	}
	s = base64.StdEncoding.EncodeToString(sig)
	return s
}

func sign_rsa(keys []string, param url.Values, privateKey []byte) (s string) {
	if param == nil {
		param = make(url.Values, 0)
	}

	var pList = make([]string, 0, 0)
	for _, key := range keys {
		var value = strings.TrimSpace(param.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	var src = strings.Join(pList, "&")
	var sig, err = encoding.SignPKCS1v15([]byte(src), privateKey, crypto.SHA1)
	if err != nil {
		return ""
	}
	s = base64.StdEncoding.EncodeToString(sig)
	return s
}
