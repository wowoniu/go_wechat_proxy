package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

var (
	wsConn *websocket.Conn
)

func init() {

}

func main() {
	var err error
	if wsConn, _, err = websocket.DefaultDialer.Dial("ws://127.0.0.1:8082/ws", nil); err != nil {
		fmt.Println("连接失败")
		return
	}

	defer wsConn.Close()
	//监听消息
	go func() {
		for {
			if _, msgData, err := wsConn.ReadMessage(); err != nil {
				fmt.Println("消息读取错误")
				return
			} else {
				fmt.Println("收到消息:", string(msgData))
			}

		}
	}()
	for {
		select {
		case <-time.Tick(time.Second):
			wsConn.WriteMessage(websocket.TextMessage, []byte("你好啊"))
		}
	}
}
