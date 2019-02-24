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
	flag.StringVar(&port, "port", "8082", "服务监听端口 默认8082")
	flag.Parse()
	if port == "" {
		return errors.New("缺少启动参数: -port 8082")
	}
	c.Port = port
	return
}
