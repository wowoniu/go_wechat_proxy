package server

import (
	"WECHAT_PROXY/common"
	"encoding/json"

	"github.com/gorilla/websocket"
)

/**
消息响应管理
*/

type MsgMgr struct {
	MessageChan chan []byte
}

var GMsgMgr *MsgMgr

//NewMsgMgr 创建消息管理器单例
func NewMsgMgr() *MsgMgr {
	if GMsgMgr == nil {
		GMsgMgr = &MsgMgr{
			MessageChan: make(chan []byte, 1000),
		}
	}
	return GMsgMgr
}

//PushMsg 收到消息
func (c *MsgMgr) PushMsg(clientID string, message []byte) {
	c.MessageChan <- message
}

//HandleMsg 消息处理逻辑
func (c *MsgMgr) HandleMsg(clientID string, message []byte) {
	go func() {
		var (
			wsMessage *common.WsMessage
			wsClient  *common.WsClient
			err       error
		)
		//处理响应
		if wsMessage, err = c.ParseMsg(message); err != nil {
			return
		}
		//获取客户端连接
		if wsClient, err = GClientMgr.GetClient(clientID); err != nil {
			//
		}

		//todo 业务处理
		wsMessage = wsMessage
		wsClient.Conn.WriteMessage(websocket.TextMessage, []byte("你好"))

	}()
}

//ParseMsg 解析收到的ws消息
func (c *MsgMgr) ParseMsg(message []byte) (wsMsg *common.WsMessage, err error) {
	wsMsg = &common.WsMessage{}
	err = json.Unmarshal(message, wsMsg)
	return
}
