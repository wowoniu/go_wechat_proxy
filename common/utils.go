package common

import (
	"net/http"
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
