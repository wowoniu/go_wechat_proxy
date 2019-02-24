package server

import (
	"github.com/wowoniu/go_wechat_proxy/common"
	"time"
)

//微信转发逻辑处理
type WechatProxy struct {
	RequestChan  chan *common.WechatRequest //微信服务器的请求
	ResponseChan chan *common.LocalResponse //本地机器的响应
	ProxyTable   map[string]*common.ProxyRecord
}

var GWechatProxy *WechatProxy

func LoadWechatProxy() *WechatProxy {
	if GWechatProxy == nil {
		GWechatProxy = &WechatProxy{
			RequestChan:  make(chan *common.WechatRequest, 1000),
			ResponseChan: make(chan *common.LocalResponse, 1000),
			ProxyTable:   make(map[string]*common.ProxyRecord),
		}
	}
	return GWechatProxy
}

func (c *WechatProxy) ToLocal(r *common.WechatRequest, responseChan chan *common.LocalResponse) {
	//存储转发记录
	c.ProxyTable[r.ID] = &common.ProxyRecord{
		RequestTime:  time.Now(),
		ResponseChan: responseChan,
	}
	//调用websocket客户端进行转发
	msg, _ := GMsgMgr.BuildMsg(common.WsMethodWechatRequest, r)
	GClientMgr.Send(r.AppID, msg)
}

func (c *WechatProxy) ToWechat(r *common.LocalResponse) {
	if proxyRecord, isExisted := c.ProxyTable[r.ID]; isExisted {
		proxyRecord.ResponseChan <- r
	}
}
