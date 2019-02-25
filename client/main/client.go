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
	if err = client.GConfig.ParseParams(); err != nil {
		fmt.Println(err.Error())
		return
	}
	//建立websocket连接
	if wsConn, _, err = websocket.DefaultDialer.Dial(client.GConfig.RemoteWsUrl, nil); err != nil {
		common.Log(common.LogLevelError, "连接失败", err)
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
		common.Log(common.LogLevelError, "连接断开")
		wsConn.Close()
	}
}

func init() {
	client.LoadConfig()
	client.LoadWechatProxy()
	client.LoadMsgMgr()
}
