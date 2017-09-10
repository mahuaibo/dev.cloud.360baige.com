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

// replyApplicationTpl.PayCycle
// PayType:0:限免 1:永久免费 2:1次性收费 3:周期收费
// pay_cycle:0无1月2季3半年4年
func ServiceTime(PayType, PayCycle int, Num, currentTimestamp int64) int64 {
	var endTime int64
	if PayType == 0 {
		endTime = currentTimestamp + 3600000*24*30*Num
	} else if PayType == 1 {
		endTime = currentTimestamp + 3600000*24*30*12*100*Num
	} else if PayType == 2 {
		endTime = currentTimestamp + 3600000*24*30*12*100*Num
	} else if PayType == 3 {
		if PayCycle == 1 {
			endTime = currentTimestamp + 3600000*24*30*Num
		} else if PayCycle == 2 {
			endTime = currentTimestamp + 3600000*24*30*3*Num
		} else if PayCycle == 3 {
			endTime = currentTimestamp + 3600000*24*30*6*Num
		} else if PayCycle == 4 {
			endTime = currentTimestamp + 3600000*24*30*12*Num
		}
	} else {
		endTime = currentTimestamp
	}
	return endTime
}
