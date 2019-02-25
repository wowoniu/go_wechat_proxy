package client

import (
	"github.com/gorilla/websocket"
	"github.com/wowoniu/go_wechat_proxy/common"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type WechatProxy struct {
}

var GWechatProxy *WechatProxy

func LoadWechatProxy() *WechatProxy {
	if GWechatProxy == nil {
		GWechatProxy = &WechatProxy{}
	}
	return GWechatProxy
}

func (c *WechatProxy) WatchLocalResponse(wsConn *websocket.Conn, proxyResponseChan chan *common.LocalResponse, isClosedChan chan bool, isWorkerClosed chan bool) {
	//监听本地转发结果
	for {
		select {
		case proxyLocalResponse := <-proxyResponseChan:
			if resStr, err := GMsgMgr.BuildWsResponse(common.WsMethodLocalResponse, proxyLocalResponse); err != nil {
				common.Log(common.LogLevelError, "构造本地响应失败:", err)
			} else {
				common.Log(common.LogLevelInfo, "本地响应:", string(resStr))
				if err = wsConn.WriteMessage(websocket.TextMessage, resStr); err != nil {
					common.Log(common.LogLevelError, "代理服务器连接失败", err)
					isClosedChan <- true
				}
			}
		case <-isClosedChan:
			isWorkerClosed <- true
			return
		}
	}
}

//http post至本地
func (c *WechatProxy) ToLocal(wechatRequest *common.WechatRequest, resultChan chan *common.LocalResponse) {
	var (
		request       *http.Request
		response      *http.Response
		localResponse []byte
		proxyResponse *common.LocalResponse
		err           error
		url           string
	)
	url = GConfig.LocalUrl + "?" + wechatRequest.GetParams
	if request, err = http.NewRequest("post", url, strings.NewReader(wechatRequest.XmlData)); err != nil {
		common.Log(common.LogLevelError, "无效的本地转发请求:", err)
		return
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	if response, err = client.Do(request); err != nil {
		common.Log(common.LogLevelError, "本地请求失败:", err)
		return
	}

	if localResponse, err = ioutil.ReadAll(response.Body); err != nil {
		common.Log(common.LogLevelError, "本地响应解析失败:", err)
		return
	}

	//构造本地回复体
	proxyResponse = GMsgMgr.BuildLocalResponse(wechatRequest.ID, wechatRequest.AppID, localResponse)
	//发送回应
	resultChan <- proxyResponse
}
