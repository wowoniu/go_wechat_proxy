package server

import (
	"errors"
	"flag"
)

type Config struct {
	Port string
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
		port string
	)
	flag.StringVar(&port, "port", "80", "服务监听端口 默认80")
	flag.Parse()
	if port == "" {
		return errors.New("缺少启动参数: -port 80")
	}
	c.Port = port
	return
}
