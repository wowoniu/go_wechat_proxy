package client

import (
	"errors"
	"flag"
)

type Config struct {
	AppID       string //代理微信的APPID
	LocalUrl    string //本地服务地址
	RemoteWsUrl string //代理服务器websocket地址
}

var GConfig *Config

func LoadConfig() (config *Config, err error) {
	if GConfig == nil {
		GConfig = &Config{}
	}
	return GConfig, err
}

func (c *Config) ParseParams() (err error) {
	var (
		appId       string
		localUrl    string
		remoteWsUrl string
	)
	flag.StringVar(&appId, "appid", "", "代理转发的微信appid")
	flag.StringVar(&localUrl, "local_url", "", "本地处理URL")
	flag.StringVar(&remoteWsUrl, "remote_ws_url", "", "代理服务器websocket监听地址")
	flag.Parse()
	if appId == "" {
		return errors.New("缺少启动参数: -appid 12345677")
	}
	if localUrl == "" {
		return errors.New("缺少启动参数: -local_url http://loclahost/wechat")
	}
	if remoteWsUrl == "" {
		return errors.New("缺少启动参数: -remote_ws_url ws://PROXY服务器/ws")
	}
	c.AppID = appId
	c.LocalUrl = localUrl
	c.RemoteWsUrl = remoteWsUrl
	return
}
