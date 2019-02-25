package main

import (
	"github.com/wowoniu/go_wechat_proxy/common"
	"github.com/wowoniu/go_wechat_proxy/server"
)

func main() {
	if err := server.GConfig.ParseParams(); err != nil {
		common.Log(common.LogLevelError, err)
		return
	}
	server.Start()
}

func init() {
	server.LoadConfig()
	server.LoadClientMgr()
	server.LoadMsgMgr()
	server.LoadWechatProxy()
}
