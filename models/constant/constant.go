package constant

// 用户状态
const (
	UserStatusValid   = 1
	UserStatusInValid = 0
	UserStatusCancel  = -1
)

// 返回状态码
const (
	ResponseSystemErr = "500"
	ResponseLogicErr = "400"
	ResponseNormal = "200"
)
