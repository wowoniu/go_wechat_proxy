package client

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wowoniu/go_wechat_proxy/common"
)

type MsgMgr struct {
}

var GMsgMgr *MsgMgr

func LoadMsgMgr() *MsgMgr {
	if GMsgMgr == nil {
		GMsgMgr = &MsgMgr{}
	}
	return GMsgMgr
}

func (c *MsgMgr) Listen(wsConn *websocket.Conn, isClosedChan chan bool, proxyResponseChan chan *common.LocalResponse) {
	var (
		msgData []byte
		err     error
	)
	//连接初始化
	initResponse, _ := c.BuildWsResponse(common.WsMethodInit, "123456789")
	if err = wsConn.WriteMessage(websocket.TextMessage, initResponse); err != nil {
		wsConn.Close()
		isClosedChan <- true
		return
	}
	for {
		//监听消息
		if _, msgData, err = wsConn.ReadMessage(); err != nil {
			isClosedChan <- true
			return
		}
		//解析消息
		wsMsg := &common.WsMessage{}
		if err = json.Unmarshal(msgData, wsMsg); err != nil {
			fmt.Println("ws消息解析错误:", err)
			return
		}
		//开启协程 处理本地消息及转发
		go func() {
			c.HandleMsg(wsMsg, proxyResponseChan)
		}()
	}
}

func (c *MsgMgr) HandleMsg(wsMsg *common.WsMessage, resultChan chan *common.LocalResponse) {
	switch wsMsg.Method {
	case common.WsMethodHeartbeat:
		//心跳包
		//fmt.Println("心跳包")
	case common.WsMethodMessage:
		//普通消息 TODO
	case common.WsMethodWechatRequest:
		//微信请求 解析微信的请求
		var jsonStr []byte
		var err error
		var wechatRequest = &common.WechatRequest{}
		if jsonStr, err = json.Marshal(wsMsg.Body.Data); err != nil {
			fmt.Println("错误的微信转发消息体格式:", err)
			return
		}
		if err = json.Unmarshal(jsonStr, wechatRequest); err != nil {
			fmt.Println("错误的微信转发消息体格式:", err)
			return
		}
		fmt.Println("转发请求:", wechatRequest.XmlData)
		//本地转发
		GWechatProxy.ToLocal(wechatRequest, resultChan)
	}
}

func (c *MsgMgr) BuildLocalResponse(requestID string, appID string, response []byte) *common.LocalResponse {
	return &common.LocalResponse{
		ID:       requestID,
		AppID:    appID,
		Response: string(response),
	}
}

func (c *MsgMgr) BuildWsResponse(method string, data interface{}) (res []byte, err error) {
	var (
		wsRes *common.WsMessage
	)
	wsRes = &common.WsMessage{
		Method: method,
		Body: &common.WsMessageBody{
			ErrorCode: "",
			ErrorMsg:  "",
			Data:      data,
		},
	}
	//序列化
	res, err = json.Marshal(wsRes)
	return
}
