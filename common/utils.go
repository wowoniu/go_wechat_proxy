package common

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

//常用方法集

//获取所有的GET参数 仅支持一维
func HttpGetParamsString(r *http.Request) string {
	query := r.URL.Query()
	getParams := ""
	isFirst := true
	for key, list := range query {
		if isFirst {
			getParams = key + "=" + list[0]
			isFirst = false
		} else {
			getParams = (getParams + "&" + key + "=" + list[0])
		}
	}
	return getParams
}

//自定义LOG函数
func Log(level string, items ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	prefix := buildLogPrefix(file, line, level)
	//根据日志分级进行控制 TODO
	//switch level {
	//case LogLevelDebug:
	//
	//case logLevelInfo:
	//
	//case LogLevelWarn:
	//
	//case LogLevelError:
	//	fmt.Println()
	//}
	var logItems = make([]interface{}, len(items)+1)
	logItems[0] = prefix
	for k, v := range items {
		logItems[k+1] = v
	}
	//logItems=append(logItems,items)
	fmt.Println(logItems...)
}

func buildLogPrefix(file string, line int, level string) string {
	var now = time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006")
	var prefix = fmt.Sprintf("[%v]%v:%d\r\n[%v]", now, file, line, level)
	return prefix
}
