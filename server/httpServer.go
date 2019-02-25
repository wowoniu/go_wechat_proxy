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
	//企业微信请求监听
	http.HandleFunc("/qywechat/proxy", handleQywechat)
	http.HandleFunc("/wechat/proxymock", handleWechatLocalMock)
	//websocket
	http.HandleFunc("/ws", handleWs)
	http.Handle("/", http.FileServer(http.Dir("./webroot")))
	common.Log(common.LogLevelInfo, "服务启动", "监听端口:"+GConfig.Port)
	if err := http.ListenAndServe("0.0.0.0:"+GConfig.Port, nil); err != nil {
		common.Log(common.LogLevelError, "服务器启动异常:", err)
	}
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
	common.Log(common.LogLevelInfo, "客户端上线:", clientID)
	GClientMgr.PushClientEvent(GClientMgr.BuildClientEvent(wsClient, common.ClientEventOnlineType))
	for {
		//监听消息
		if _, wsMsgData, err = wsConn.ReadMessage(); err != nil {
			//下线
			common.Log(common.LogLevelInfo, "读取错误:连接关闭")
			GClientMgr.PushClientEvent(GClientMgr.BuildClientEvent(wsClient, common.ClientEventOfflineType))
			return
		}
		GMsgMgr.HandleMsg(clientID, wsMsgData)
	}

}

//处理微信公众号的异步消息
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
	common.Log(common.LogLevelDebug, "转发请求:(get)", request.GetParams, " (xml)", request.XmlData)
	responseChan := make(chan *common.LocalResponse)
	//转发至本地
	GWechatProxy.ToLocal(request, responseChan)
	//监听结果
	select {
	case response := <-responseChan:
		w.Write([]byte(response.Response))
	case <-time.Tick(3 * time.Second):
		common.Log(common.LogLevelInfo, "客户端转发响应超时")
		w.Write([]byte("success"))
	}
}

//处理企业微信的异步消息
func handleQywechat(w http.ResponseWriter, r *http.Request) {
	//接收到微信的请求后 将数据放入管道 向本地转发
	r.ParseForm()
	appID := r.Form["APPID"][0]

	token := r.Form["_token"][0]
	corp_id := r.Form["_corpid"][0]
	encoding_aeskey := r.Form["_aeskey"][0]
	wxcpt := common.NewWXBizMsgCrypt(token, encoding_aeskey, corp_id, common.XmlType)

	//判断是否是初次服务器验证
	if _, isExisted := r.Form["echostr"]; isExisted {
		//第一次握手 加解密
		//* GET /cgi-bin/wxpush?msg_signature=5c45ff5e21c57e6ad56bac8758b79b1d9ac89fd3&timestamp=1409659589&nonce=263014780&echostr=P9nAzCzyDtyTWESHep1vC5X9xho%2FqYX3Zpb4yKa9SKld1DsH3Iyt3tP3zNdtp%2B4RPcs8TgAE7OaBO%2BFZXvnaqQ%3D%3D
		echoStr := r.Form["echostr"][0]
		msgSig := r.Form["msg_signature"][0]
		timeStamp := r.Form["timestamp"][0]
		nonce := r.Form["nonce"][0]
		echostr, crypt_err := wxcpt.VerifyURL(msgSig, timeStamp, nonce, echoStr)
		if nil != crypt_err {
			fmt.Println("verifyurl fail", crypt_err)
		}
		w.Write(echostr)
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
		common.Log(common.LogLevelInfo, "客户端转发响应超时")
		w.Write([]byte("success"))
	}
}

func handleWechatLocalMock(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	requestBody, _ := ioutil.ReadAll(r.Body)
	getParams := common.HttpGetParamsString(r)
	w.Write([]byte("GET:" + getParams + "<br/>XML:" + string(requestBody)))
}
