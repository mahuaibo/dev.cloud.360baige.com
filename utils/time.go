package utils

import (
	"time"
)

func CurrentTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

func Datetime(timestamp int64, format string) string {
	return time.Unix(timestamp/1000, 0).Format(format)
}
