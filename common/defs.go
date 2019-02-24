package common

import (
	"time"

	"github.com/gorilla/websocket"
)

//WsClient 客户端
type WsClient struct {
	ClientKey      string          //连接标识
	UserKey        string          //用户标识
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
	Method string         `json:"method"` //指令
	Body   *WsMessageBody `json:"body"`   //消息体
}

type WsMessageBody struct {
	ErrorCode string      `json:"err_code"`
	ErrorMsg  string      `json:"error_msg"`
	Data      interface{} `json:"data"`
}

//微信转发相关

type ProxyRecord struct {
	RequestTime  time.Time
	ResponseChan chan *LocalResponse
}

type WechatRequest struct {
	ID         string `json:"id"`
	AppID      string `json:"app_id"`
	GetParams  string `json:"get_params"`
	PostParams string `json:"post_params"`
	XmlData    string `json:"xml_data"`
}

type LocalResponse struct {
	ID       string `json:"id"`
	AppID    string `json:"app_id"`
	Response string `json:"response"`
}
