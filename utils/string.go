package utils

import (
	"math/rand"
	"time"
	"strconv"
	"reflect"
	"github.com/astaxie/beego/context"
	qrcode "github.com/skip2/go-qrcode"
	"strings"
	"errors"
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

func Unable(params map[string]string, input *context.BeegoInput) error {
	typeErr := ""
	valueErr := ""
	for key, param := range params {
		value := input.Query(key)
		s := strings.Split(param, ":")
		if len(s) == 2 {
			if s[1] == "true" && value == "" {
				valueErr += " " + key + " "
			}
			if !checkType(value, s[0]) {
				typeErr += " " + key + " "
			}
		}
	}
	if typeErr != "" {
		typeErr = "类型错误:" + typeErr
	}
	if valueErr != "" {
		valueErr = "必传参数:" + valueErr
	}
	if typeErr == "" && valueErr == "" {
		return nil
	} else {
		return errors.New(typeErr + valueErr)
	}
}

func checkType(param string, paramType string) bool {
	switch paramType {
	case "string":
		if reflect.ValueOf(param).Type().String() == paramType {
			return true
		}
	case "int":
		value, err := strconv.ParseInt(param, 10, 64)
		if err == nil && param != "" && reflect.ValueOf(value).Type().String() == "int64" {
			return true
		}
	case "float":
		value, err := strconv.ParseFloat(param, 64)
		if err == nil && param != "" && reflect.ValueOf(value).Type().String() == "float64" {
			return true
		}
	}
	return false
}

func Qr(url string, size int) []byte {
	qr, _ := qrcode.Encode(url, qrcode.Medium, size)
	return qr
}
