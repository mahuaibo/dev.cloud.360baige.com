package utils

import (
	"github.com/astaxie/beego/context"
	"strings"
	"errors"
	"reflect"
	"strconv"
)

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
