package common

import (
	"time"

	"github.com/gorilla/websocket"
)

//WsClient 客户端
type WsClient struct {
	ClientKey      string          //客户端标识
	Conn           *websocket.Conn //ws连接
	ConnectTime    time.Time       //连接时间
	LastActiveTime time.Time       //最后一次活跃时间
}

//ClientEvent 客户端事件
type ClientEvent struct {
	Type      string
	ClientKey string
	Client    *WsClient
}

//WsMessage 通信的消息体
type WsMessage struct {
	Method string      `json:"method"` //指令
	Data   interface{} `json:"data"`
}
