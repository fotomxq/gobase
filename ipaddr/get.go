package ipaddr

import (
	"strings"
	"net/http"
)

//通过r获取客户端IP地址
// 自动剔除端口部分，只获取IP地址
//param r *http.Request Http读取对象
//param format string 返回格式约定 all 全部返回 | remote 仅本地IP | forwarded 仅代理IP | real 仅真实IP | filter 只要存在就返回
//param needPort bool 是否保留端口
//return string IP地址
func IPGet(r *http.Request,format string,needPort bool) string{
	//获取代理的方案
	ipGetRes := []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"Remote_addr",
	}
	ipAddrs := map[string]string{}
	for _,v := range ipGetRes{
		ipAddrs[v] = r.Header.Get(v)
		if v == "Remote_addr"{
			if ipAddrs["Remote_addr"] == ""{
				ipAddrs["Remote_addr"] = r.RemoteAddr
				continue
			}
		}
	}
	//优先使用代理IP，之后才是本地IP
	var ipAddr string
	switch format{
	case "all":
		// forwarded -> real -> remote
		ipAddr = ipAddrs["X-Forwarded-For"]
		ipAddr += "->" + ipAddrs["X-Real-IP"]
		ipAddr += "->" + ipAddrs["Remote_addr"]
	case "remote":
		if ipAddrs["Remote_addr"] != "" {
			ipAddr = ipAddrs["Remote_addr"]
		}
	case "forwared":
		if ipAddrs["X-Forwarded-For"] != ""{
			ipAddr = ipAddrs["X-Forwarded-For"]
		}
	case "real":
		if ipAddrs["X-Real-IP"] != ""{
			ipAddr = ipAddrs["X-Real-IP"]
		}
	case "filter":
		if ipAddrs["X-Forwarded-For"] != ""{
			ipAddr = ipAddrs["X-Forwarded-For"]
		}
		if ipAddrs["X-Real-IP"] != ""{
			ipAddr = ipAddrs["X-Real-IP"]
		}
		if ipAddrs["Remote_addr"] != "" {
			ipAddr = ipAddrs["Remote_addr"]
		}
	}
	//处理IP地址
	if needPort == true{
		if ipAddr != ""{
			ipAddrs := strings.Split(ipAddr,":")
			ipAddrLast := ":"+ipAddrs[len(ipAddrs)-1]
			ipAddr = strings.Replace(ipAddr,ipAddrLast,"",-1)
		}
	}
	//返回
	return ipAddr
}