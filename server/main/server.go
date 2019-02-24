package main

import (
	"fmt"
	"github.com/wowoniu/go_wechat_proxy/server"
)

func main() {
	if err := server.GConfig.ParseParams(); err != nil {
		fmt.Println(err)
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
