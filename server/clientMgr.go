package server

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wowoniu/go_wechat_proxy/common"

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
		ActivityTable:   make(map[string]*common.WsClient),
	}
	//监听客户端事件
	GClientMgr.WatchEvent()
	//心跳包发送
	//GClientMgr.Heartbeat()
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
	//上线
	case common.ClientEventOnlineType:
		if client, isExisted := c.ActivityTable[clientEvent.ClientKey]; isExisted {
			client.Conn.Close()
		}
		c.ActivityTable[clientEvent.ClientKey] = clientEvent.Client
	case common.ClientEventOfflineType:
		if client, isExisted := c.ActivityTable[clientEvent.ClientKey]; isExisted {
			client.Conn.Close()
			fmt.Println("客户端下线:", client.ClientKey, "-", client.UserKey)
			delete(c.ActivityTable, clientEvent.ClientKey)
		}
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
			case <-time.Tick(5 * time.Second):
				fmt.Println("客户端数量:", len(c.ActivityTable))
				for _, client := range c.ActivityTable {
					if msg, err = c.BuildHeartbeatPackage(); err == nil {
						//fmt.Println("发送心跳:",string(msg))
						if err = client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
							//心跳发送失败 下线
							c.PushClientEvent(c.BuildClientEvent(client, common.ClientEventOfflineType))
						}
					} else {
						fmt.Println("心跳包构造失败:", err)
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
		Body:   &common.WsMessageData{},
	}
	//JSON序列化
	if msg, err = json.Marshal(data); err != nil {
		return
	}
	return
}

//GetClient 获取连接客户端
func (c *ClientMgr) GetClient(userKey string) (client *common.WsClient, err error) {
	for _, client := range c.ActivityTable {
		if client.UserKey == userKey {
			return client, nil
		}
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

func (c *ClientMgr) Send(userKey string, wsMsg []byte) (err error) {
	var (
		wsClient *common.WsClient
	)
	//获取客户端连接
	if wsClient, err = GClientMgr.GetClient(userKey); err != nil {
		//客户端已下线
		//fmt.Println("获取客户端失败")
		return
	}
	err = wsClient.Conn.WriteMessage(websocket.TextMessage, wsMsg)
	return
}

func (c *ClientMgr) SetUserKey(clientID string, userKey string) {
	if client, isExisted := c.ActivityTable[clientID]; isExisted {
		client.UserKey = userKey
	}
	return
}
