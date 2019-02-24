package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wowoniu/go_wechat_proxy/client"
	"github.com/wowoniu/go_wechat_proxy/common"
)

var (
	wsConn *websocket.Conn
)

func main() {
	var (
		err               error
		clientClosedChan  = make(chan bool)
		workerClosedChan  = make(chan bool)
		proxyResponseChan = make(chan *common.LocalResponse, 1000)
	)
	//建立websocket连接
	if wsConn, _, err = websocket.DefaultDialer.Dial("ws://127.0.0.1:8082/ws", nil); err != nil {
		fmt.Println("连接失败")
		return
	}
	defer wsConn.Close()

	//协程 监听消息
	go client.GMsgMgr.Listen(wsConn, clientClosedChan, proxyResponseChan)

	//协程 监听本地结果
	go client.GWechatProxy.WatchLocalResponse(wsConn, proxyResponseChan, clientClosedChan, workerClosedChan)

	select {
	case <-workerClosedChan:
		//断开后 重连
		fmt.Println("连接断开")
		wsConn.Close()
	}
}

func init() {
	client.LoadWechatProxy()
	client.LoadMsgMgr()
}
