package utils

import (
	"strings"
	"regexp"
	"time"
	"strconv"
)

func DetermineStringType(str string) (string, bool) {
	var strType string
	var b bool
	if (strings.ContainsAny(str, "@")) {
		strType = "email" // 邮箱
		b = IsEmail(strings.ToLower(str))
	} else {
		isnum := IsInteger(str)
		if (isnum) {
			strType = "phone" // 手机号
			b = IsMobile(str)
		} else {
			strType = "username" // 账号组合
			b = true
		}
	}
	return strType, b
}

//邮箱 最高30位
func IsEmail(str ...string) bool {
	var b bool
	for _, s := range str {
		b, _ = regexp.MatchString("^([a-z0-9_\\.-]+)@([\\da-z\\.-]+)\\.([a-z\\.]{2,6})$", s)
		if false == b {
			return b
		}
	}
	return b
}

//纯整数
func IsInteger(str ...string) bool {
	var b bool
	for _, s := range str {
		b, _ = regexp.MatchString("^[0-9]+$", s)
		if false == b {
			return b
		}
	}
	return b
}

//手机号码
func IsMobile(str ...string) bool {
	var b bool
	for _, s := range str {
		b, _ = regexp.MatchString("^(0|86|17951)?(13[0-9]|15[0-9]|17[0-9]|18[0-9]|14[0-9])[0-9]{8}$", s)
		if false == b {
			return b
		}
	}
	return b
}

func GetMonthStartUnix(current string) int64 {
	tm2, _ := time.ParseInLocation("2006-01-02", current, time.Local)
	stime := tm2.UnixNano() / 1e6
	return stime
}
func GetNextMonthStartUnix(current string) int64 {
	tm2, _ := time.ParseInLocation("2006-01-02", current, time.Local)
	t := time.Unix(tm2.UnixNano()/1e9, 0)
	s := time.Date(t.Year(), t.Month()+1, t.Day(), 0, 0, 0, 0, t.Location())
	es := s.Format("2006-01-02")
	estm2, _ := time.ParseInLocation("2006-01-02", es, time.Local)
	etime := estm2.UnixNano() / 1e6
	return etime
}

//第二天
func GetNextDayUnix(current string) int64 {
	tm2, _ := time.ParseInLocation("2006-01-02", current, time.Local)
	t := time.Unix(tm2.UnixNano()/1e9, 0)
	s := time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, t.Location())
	es := s.Format("2006-01-02")
	estm2, _ := time.ParseInLocation("2006-01-02", es, time.Local)
	etime := estm2.UnixNano() / 1e6
	return etime
}

//
func StrArrToInt64Arr(strArr []string) []int64 {
	var int64Arr []int64 = make([]int64, len(strArr), len(strArr))
	for index, str := range strArr {
		int64Arr[index], _ = strconv.ParseInt(str, 10, 64)
	}
	return int64Arr
}
