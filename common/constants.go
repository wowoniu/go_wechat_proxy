package common

const (
	//WsMethodHeartbeat 消息指令:心跳包
	WsMethodHeartbeat = "heartbeat"

	//WsMethodMessage 消息指令:普通消息
	WsMethodMessage = "message"

	//WsMethodResponse 消息指令:本地响应
	WsMethodLocalResponse = "localResponse"

	//WsMethodResponse 消息指令:微信请求
	WsMethodWechatRequest = "wechatRequest"

	//WsMethodMessage 消息指令:初始化
	WsMethodInit = "init"

	//ClientEventOnlineType 客户端事件类型:上线
	ClientEventOnlineType = "1"

	//ClientEventOfflineType 客户端事件类型:下线
	ClientEventOfflineType = "2"

	//LogLevelDebug 日志级别
	LogLevelDebug = "DEBUG"

	LogLevelInfo = "INFO"

	LogLevelWarn = "WARN"

	LogLevelError = "ERROR"
)
