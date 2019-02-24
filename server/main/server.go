package main

import (
	"github.com/wowoniu/go_wechat_proxy/server"
)

func main() {
	server.Start()
}

func init() {
	server.LoadClientMgr()
	server.LoadMsgMgr()
	server.LoadWechatProxy()
}
