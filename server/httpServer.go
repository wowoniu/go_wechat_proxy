package server

import (
	"WECHAT_PROXY/common"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func Listen() {
	http.HandleFunc("/ws", handleWs)
	http.Handle("/", http.FileServer(http.Dir("./webroot")))
	http.ListenAndServe("0.0.0.0:8082", nil)
}

func handleWs(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		wsMsg    []byte
		msgType  int
		wsConn   *websocket.Conn
		clientID string
	)
	clientID = fmt.Sprintf("%s", time.Unix)
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		w.Write([]byte("建立连接错误"))
		return
	}
	//建立连接
	defer wsConn.Close()
	wsClient := GClientMgr.BuildClient(clientID, wsConn)
	//客户端上线
	GClientMgr.PushClientEvent(GClientMgr.BuildClientEvent(wsClient, common.ClientEventOnlineType))
	for {
		//监听消息
		if msgType, wsMsg, err = wsConn.ReadMessage(); err != nil {
			//连接关闭 下线
			fmt.Println("读取错误:连接关闭")
			GClientMgr.PushClientEvent(GClientMgr.BuildClientEvent(wsClient, common.ClientEventOfflineType))
			return
		}
		//消息处理
		GMsgMgr.PushMsg(clientID, wsMsg)
	}

}
