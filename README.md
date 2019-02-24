##  微信消息穿透代理(go_wechat_proxy)


#### 使用场景
微信公众号、企业号等接口提供了微信消息、用户关注等事件的异步通知机制
在开发过程中，微信异步通知只能推送到公网服务器上，项目在开发初期，在公网部署开发环境
需要不断进行代码上线、部署等操作，对开发工作十分不便；使用该代理软件，可在公网部署一个代理服务器
然后将微信公众号的服务器地址配置到代理服务器上，再在本地运行代理客户端，即可实现微信异步消息利用websocket从公网直接
推送到本地开发机器上进行开发调试。

#### 使用说明

bin目录下有已经编译好的服务端和客户端可执行文件(其他操作系统 可自行编译)

+ 公网部署代理服务端程序 监听微信通知 (注意：服务端端口 微信只支持80和443端口)
  + linux(64)
     + 启动命令： ./server_linux_64_linux -port 80
     + 配置微信公众号服务器地址: http://IP:80/wechat/proxy?APPID=微信公众号APPID
  +  win(64)如上
  
+ 本地开发机器运行客户端软件接收公网的转发请求
   + linux(64)
     + 启动客户端： ./client_linux_64_linux -appid 微信公众号APPID -local_url http://本地微信服务应用URL -remote_ws_url ws://公网部署IP:80/ws
   +  win(64)如上
   
   
#### 代理服务端快速推荐部署方式
    没有公网服务器的情况下 建议使用daocloud(https://www.daocloud.io/)的免费胶囊Ubuntu主机(每次申请可免费使用两个小时) 
    使用SSH登陆到胶囊主机后 下载服务端发行版本 运行服务 
    
    wget https://github.com/wowoniu/go_wechat_proxy/releases/download/v1.0/server_linux_64_linux 
    chmod +x server_linux_64_linux
    sudo ./server_linux_64_linux
    
    
#### <a href="https://github.com/wowoniu/go_wechat_proxy/releases" target="_blank">下载地址</a>
    