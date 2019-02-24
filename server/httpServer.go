package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/wowoniu/go_wechat_proxy/common"

	"github.com/gorilla/websocket"
)

func Start() {
	//微信请求监听
	http.HandleFunc("/wechat/proxy", handleWechat)
	http.HandleFunc("/wechat/proxymock", handleWechatLocalMock)
	//websocket
	http.HandleFunc("/ws", handleWs)
	http.Handle("/", http.FileServer(http.Dir("./webroot")))
	http.ListenAndServe("0.0.0.0:8082", nil)
}

func handleWs(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		wsMsgData []byte
		wsConn    *websocket.Conn
	)
	clientID := fmt.Sprintf("%s", time.Now())
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
	fmt.Println("客户端上线:", clientID)
	GClientMgr.PushClientEvent(GClientMgr.BuildClientEvent(wsClient, common.ClientEventOnlineType))
	for {
		//监听消息
		if _, wsMsgData, err = wsConn.ReadMessage(); err != nil {
			//下线
			fmt.Println("读取错误:连接关闭")
			GClientMgr.PushClientEvent(GClientMgr.BuildClientEvent(wsClient, common.ClientEventOfflineType))
			return
		}
		GMsgMgr.HandleMsg(clientID, wsMsgData)
	}

}

func handleWechat(w http.ResponseWriter, r *http.Request) {
	//接收到微信的请求后 将数据放入管道 向本地转发
	r.ParseForm()
	appID := r.Form["APPID"][0]
	//判断是否是初次服务器验证 暂不做验证 TODO
	if _, isExisted := r.Form["echostr"]; isExisted {
		echostr := r.Form["echostr"][0]
		w.Write([]byte(echostr))
		return
	}
	//其他响应解析并转发
	requestBody, _ := ioutil.ReadAll(r.Body)
	xmlData := string(requestBody)
	request := &common.WechatRequest{
		ID:        fmt.Sprintf("%s%d", appID, time.Now().Unix()),
		AppID:     appID,
		GetParams: common.HttpGetParamsString(r),
		XmlData:   xmlData,
	}
	responseChan := make(chan *common.LocalResponse)
	//转发至本地
	GWechatProxy.ToLocal(request, responseChan)
	//监听结果
	select {
	case response := <-responseChan:
		w.Write([]byte(response.Response))
	case <-time.Tick(3 * time.Second):
		w.Write([]byte("success"))
	}
}

func handleWechatLocalMock(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	requestBody, _ := ioutil.ReadAll(r.Body)
	getParams := common.HttpGetParamsString(r)
	w.Write([]byte("GET:" + getParams + "<br/>XML:" + string(requestBody)))
}
