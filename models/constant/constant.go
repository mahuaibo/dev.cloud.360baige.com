package constant

// Response 状态码
const (
	ResponseSystemErr = "500"
	ResponseLogicErr  = "400"
	ResponseNormal    = "200"
)

// User 用户状态
const (
	UserStatusValid   = 1
	UserStatusInValid = 0
	UserStatusCancel  = -1
)

// Logger 日志信息
const (
	LoggerTypeAdd    = 1
	LoggerTypeUpdate = 2
	LoggerTypeDelete = 3
	LoggerTypeFind   = 4
	LoggerTypeOther  = 5
)
