package common

import "errors"

var (
	//ErrorClientOffline 错误：客户端已下线
	ErrorClientOffline = errors.New("客户端已下线")

	ErrorInvalidMessage = errors.New("无效的请求数据")
)
