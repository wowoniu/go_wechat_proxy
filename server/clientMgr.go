package server

import (
	"WECHAT_PROXY/common"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

/**
客户端相关管理
*/

//ClientMgr 客户端管理器
type ClientMgr struct {
	ClientEventChan chan *common.ClientEvent
	ActivityTable   map[string]*common.WsClient //在线的客户端列表
}

//GClientMgr 全局客户端管理器单利
var GClientMgr *ClientMgr

//NewClientMgr 返回客户端管理器单例
func NewClientMgr() *ClientMgr {
	if GClientMgr != nil {
		return GClientMgr
	}
	GClientMgr = &ClientMgr{
		ClientEventChan: make(chan *common.ClientEvent, 1000),
	}
	//心跳监测
	GClientMgr.Heartbeat()
	return GClientMgr
}

//WatchEvent 监听客户端连接等事件
func (c *ClientMgr) WatchEvent() {
	go func() {
		for {
			select {
			case clientEvent := <-c.ClientEventChan:
				c.HandleEvent(clientEvent)
			}
		}
	}()
}

//HandleEvent 处理客户端事件
func (c *ClientMgr) HandleEvent(clientEvent *common.ClientEvent) {
	switch clientEvent.Type {
	case common.ClientEventOnlineType:
		if client, isExisted := c.ActivityTable[clientEvent.ClientKey]; isExisted {
			client.Conn.Close()
		}
		c.ActivityTable[clientEvent.ClientKey] = clientEvent.Client
	}
}

//PushClientEvent 推送客户端事件
func (c *ClientMgr) PushClientEvent(clientEvent *common.ClientEvent) {
	c.ClientEventChan <- clientEvent
}

//Heartbeat 心跳检测 发送心跳包
func (c *ClientMgr) Heartbeat() {
	go func() {
		var (
			err error
			msg []byte
		)
		for {
			select {
			case <-time.Tick(3 * time.Second):
				for _, client := range c.ActivityTable {
					if msg, err = c.BuildHeartbeatPackage(); err == nil {
						client.Conn.WriteMessage(websocket.TextMessage, msg)
					}
				}
			}
		}

	}()
}

//BuildHeartbeatPackage 构建心跳包
func (c *ClientMgr) BuildHeartbeatPackage() (msg []byte, err error) {
	data := common.WsMessage{
		Method: common.WsMethodHeartbeat,
	}
	//JSON序列化
	if msg, err = json.Marshal(data); err != nil {
		return
	}
	return
}

//GetClient 获取连接客户端
func (c *ClientMgr) GetClient(clientID string) (client *common.WsClient, err error) {
	var isExisted bool
	if client, isExisted = c.ActivityTable[clientID]; isExisted {
		return client, nil
	}
	return nil, common.ErrorClientOffline
}

//BuildClient 构造客户端结构体
func (c *ClientMgr) BuildClient(clientKey string, conn *websocket.Conn) *common.WsClient {
	return &common.WsClient{
		ClientKey:      clientKey,
		Conn:           conn,
		ConnectTime:    time.Now(),
		LastActiveTime: time.Now(),
	}
}

//BuildClientEvent 构造客户端事件结构体
func (c *ClientMgr) BuildClientEvent(client *common.WsClient, eventType string) *common.ClientEvent {
	return &common.ClientEvent{
		Type:      eventType,
		ClientKey: client.ClientKey,
		Client:    client,
	}
}
