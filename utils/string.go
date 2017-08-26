package utils

import (
	"math/rand"
	"time"
	"strconv"
	//"reflect"
)

func RandomName(prefix string, suffix string) string {
	return prefix + Datetime(CurrentTimestamp(), "20060102") + CreateAccessValue(RandomString(10)) + suffix
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

func Amount(price float64) string {
	return strconv.FormatFloat(price, 'f', 2, 64)
}

func Unable(params map[string]interface{}) string {
	str := ""
	for key, value := range params {
		//reflect.
	}
	return str
}
