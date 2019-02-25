package server

import (
	"encoding/json"
	"github.com/wowoniu/go_wechat_proxy/common"
)

/**
消息响应管理
*/

type MsgMgr struct {
	MessageChan chan []byte
}

var GMsgMgr *MsgMgr

//LoadMsgMgr 创建消息管理器单例
func LoadMsgMgr() *MsgMgr {
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
			err       error
		)
		//解析请求数据
		if wsMessage, err = c.ParseMsg(message); err != nil {
			common.Log(common.LogLevelInfo, "无效的消息:", err)
			response, _ := c.BuildMsg(common.WsMethodMessage, common.ErrorInvalidMessage)
			GClientMgr.Send(clientID, response)
			return
		}
		//业务处理
		switch wsMessage.Method {
		case common.WsMethodLocalResponse:
			//本地机器响应
			var localResponseJsonStr []byte
			if localResponseJsonStr, err = json.Marshal(wsMessage.Body.Data); err != nil {
				//无效的响应数据体
				response, _ := c.BuildMsg(common.WsMethodMessage, common.ErrorInvalidMessage)
				GClientMgr.Send(clientID, response)
				return
			}
			//格式解析
			localResponse := &common.LocalResponse{}
			if err = json.Unmarshal(localResponseJsonStr, localResponse); err != nil {
				response, _ := c.BuildMsg(common.WsMethodMessage, common.ErrorInvalidMessage)
				GClientMgr.Send(clientID, response)
				return
			}
			GWechatProxy.ToWechat(localResponse)
		case common.WsMethodInit:
			//初始化连接
			AppID := wsMessage.Body.Data.(string)
			common.Log(common.LogLevelInfo, "用户上线:", clientID, "-", AppID)
			GClientMgr.SetUserKey(clientID, AppID)
		}
	}()
}

//ParseMsg 解析收到的ws消息
func (c *MsgMgr) ParseMsg(message []byte) (wsMsg *common.WsMessage, err error) {
	wsMsg = &common.WsMessage{}
	err = json.Unmarshal(message, wsMsg)
	return
}

func (c *MsgMgr) BuildMsg(method string, data interface{}) (msg []byte, err error) {
	wsMsg := &common.WsMessage{
		Method: method,
		Body: &common.WsMessageBody{
			ErrorCode: "",
			ErrorMsg:  "",
			Data:      data,
		},
	}
	msg, err = json.Marshal(wsMsg)
	return
}
