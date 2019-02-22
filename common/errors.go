package common

import "errors"

var (
	//ErrorClientOffline 错误：客户端已下线
	ErrorClientOffline error = errors.New("客户端已下线")
)
