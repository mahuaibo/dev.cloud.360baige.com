package utils

import (
	"math/rand"
	"time"
	qrcode "github.com/skip2/go-qrcode"
)

func RandomName(prefix string, suffix string) string {
	return prefix + Datetime(CurrentTimestamp(), "20060102") + CreateAccessValue(RandomString(10)) + suffix
}

func RandomNum(length int) string {
	bytes := []byte("0123456789")
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func RandomString(length int) string {
	bytes := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func Qr(url string, size int) []byte {
	qr, _ := qrcode.Encode(url, qrcode.Medium, size)
	return qr
}
