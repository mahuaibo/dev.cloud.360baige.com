package utils

import (
	"regexp"
	"time"
	"strconv"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"path"
	"strings"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mime/multipart"
	"github.com/aliyun/aliyun-oss-go-sdk/sample"
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
	t := time.Unix(tm2.UnixNano() / 1e9, 0)
	s := time.Date(t.Year(), t.Month() + 1, t.Day(), 0, 0, 0, 0, t.Location())
	es := s.Format("2006-01-02")
	estm2, _ := time.ParseInLocation("2006-01-02", es, time.Local)
	etime := estm2.UnixNano() / 1e6
	return etime
}

//第二天
func GetNextDayUnix(current string) int64 {
	tm2, _ := time.ParseInLocation("2006-01-02", current, time.Local)
	t := time.Unix(tm2.UnixNano() / 1e9, 0)
	s := time.Date(t.Year(), t.Month(), t.Day() + 1, 0, 0, 0, 0, t.Location())
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

// 获取差集
func Minus(slice1 []int64, slice2 []int64) (diffSlice []int64) {
	for _, v := range slice1 {
		if !HasSlice(v, slice2) {
			diffSlice = append(diffSlice, v)
		}
	}
	return
}

func HasSlice(val int64, slice []int64) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func CreateAccessValue(value string) string {
	createTime := strconv.FormatInt(time.Now().UnixNano(), 10)
	mac := hmac.New(sha1.New, []byte(createTime))
	mac.Write([]byte(value))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

func UploadImage(handle *multipart.FileHeader, filePath string) (string, error) {
	objectKey := strconv.FormatInt(time.Now().Unix(), 10) + path.Base(handle.Filename)
	suffix := path.Ext(objectKey) //获取文件后缀
	objectKey = filePath + Reverse(strings.TrimSuffix(objectKey, suffix)) + suffix
	client, err := oss.New("oss-cn-beijing.aliyuncs.com", "ZCPOvqJlByc96mZb", "8MV5VOzaClwYAlJh0eQuI8M1norVAK")
	fmt.Println("client", client)
	fmt.Println("err", err)
	if err != nil {
		return "", err
	}
	bucket, err := client.Bucket("sdk-baige")
	fmt.Println("bucket", bucket)
	fmt.Println("err", err)
	if err != nil {
		return "", err
	}

	fd, err := handle.Open()
	defer fd.Close()

	err = bucket.PutObject(objectKey, fd)
	if err != nil {
		return "", err
	}

	return objectKey, err

}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r) - 1; i < j; i, j = i + 1, j - 1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func SignURLSample(objectKey string, validityTime int64) string {
	client, err := oss.New("oss-cn-beijing.aliyuncs.com", "ZCPOvqJlByc96mZb", "8MV5VOzaClwYAlJh0eQuI8M1norVAK")
	if err != nil {
		sample.HandleError(err)
	}
	bucket, err := client.Bucket("sdk-baige")
	if err != nil {
		sample.HandleError(err)
	}

	// put object
	signedURL, err := bucket.SignURL(objectKey, oss.HTTPGet, validityTime)
	if err != nil {
		sample.HandleError(err)
	}
	return signedURL
}