package ipaddr

import (
	"strings"
	"time"
)

//IP处理器
var(
	//IP黑名单是否启动
	// = true
	IPBanON bool

	//IP白名单是否启动
	// = true
	IPWhiteON bool

	//IP黑名单地址池
	// = []IPAddrData{}
	IPBanContent []IPAddrData

	//IP白名单地址池
	// = []IPAddrData{}
	IPWhiteContent []IPAddrData

	//默认过期时间
	// = 1800
	ExpireTime int64

	//锁定机制
	// = make(chan int,1)
	Lock chan int
)

//IPAddrData结构体
type IPAddrData struct {
	//IP地址
	IPAddr string
	//过期时间，如果不设置则为-1
	ExpireTime int64
}

////////////////////////////////////////////////////////////////////////////////////
//Mgo存储模式
////////////////////////////////////////////////////////////////////////////////////

//初始化
func Run(){
	IPBanON = true
	IPWhiteON = true
	IPBanContent = []IPAddrData{}
	IPWhiteContent = []IPAddrData{}
	ExpireTime = 1800
	Lock = make(chan int,1)
}

//更新配置信息
//param ipBanOn string 是否启动黑名单
//param ipWhiteON string 是否启动白名单
//return error
func RefConfig(ipBanOn string,ipWhiteON string) error {
	IPBanON = ipBanOn == "1"
	IPWhiteON = ipWhiteON == "1"
	return nil
}

//检查IP是否可以通行
//param ipaddr string IP地址
//return bool
func CheckOK(ipaddr string) bool {
	Lock <- 1
	defer func(){
		<- Lock
	}()
	//如果关闭所有模式，则一律放行
	if IPBanON == false && IPWhiteON == false{
		return true
	}
	//如果启动黑名单，则检查
	if IPBanON == true{
		for _,v := range IPBanContent{
			if v.IPAddr == ipaddr{
				return false
			}
		}
	}
	//如果启动白名单，则检查
	if IPWhiteON == true{
		for _,v := range IPWhiteContent{
			if v.IPAddr == ipaddr{
				return false
			}
		}
	}
	//全部通过返回true
	return true
}

//根据配置文件给定的字符串，自动设置长期黑名单、白名单操作
//param IPBanContent string
//param IPWhiteContent string
func GetConfigs(IPBanContent string,IPWhiteContent string){
	if IPBanContent != ""{
		IPBanContents := strings.Split(IPBanContent,"|")
		for _,v := range IPBanContents{
			SetIPBan(v,true,false)
		}
	}
	if IPWhiteContent != ""{
		IPWhiteContents := strings.Split(IPWhiteContent,"|")
		for _,v := range IPWhiteContents{
			SetIPWhite(v,true,false)
		}
	}
}

//自动清空过期的IP地址
// 并发执行该操作即可，定时自动检测过期的IP地址
func AutoClear() {
	for{
		Lock <- 1
		//必须启动名单才会开启检查工具
		//如果启动黑名单，则检查
		if IPBanON == true{
			nowTime := time.Now().Unix()
			newData := []IPAddrData{}
			for _,v := range IPBanContent{
				if v.ExpireTime == -1 || v.ExpireTime < nowTime{
					newData = append(newData,v)
				}
			}
			IPBanContent = []IPAddrData{}
			for _,v := range newData{
				IPBanContent = append(IPBanContent,v)
			}
		}
		//如果启动白名单，则检查
		if IPWhiteON == true{
			nowTime := time.Now().Unix()
			newData := []IPAddrData{}
			for _,v := range IPWhiteContent{
				if v.ExpireTime == -1 || v.ExpireTime > nowTime{
					newData = append(newData,v)
				}
			}
			IPWhiteContent = []IPAddrData{}
			for _,v := range newData{
				IPWhiteContent = append(IPWhiteContent,v)
			}
		}
		<- Lock
		//5分钟检查一次
		time.Sleep(time.Minute * 5)
	}
}

//设置IP为黑名单
//param ipaddr string
//param b bool 是否列入黑名单
//param isExpireTime bool 是否过期
//return error
func SetIPBan(ipaddr string,b bool,isExpireTime bool) error{
	Lock <- 1
	newData := []IPAddrData{}
	for _,v := range IPBanContent{
		if v.IPAddr != ipaddr{
			newData = append(newData,v)
		}
	}
	if b == true{
		var expireTime int64 = -1
		if isExpireTime == true{
			expireTime = time.Now().Unix() + ExpireTime
		}
		newData = append(newData,IPAddrData{
			ipaddr,
			expireTime,
		})
	}
	IPBanContent = []IPAddrData{}
	for _,v := range newData{
		IPBanContent = append(IPBanContent,v)
	}
	<- Lock
	return nil
}

//设置IP为白名单
//param ipaddr string
//param b bool 是否列入白名单
//param isExpireTime int64 过期时间
//return error
func SetIPWhite(ipaddr string,b bool,isExpireTime bool) error{
	Lock <- 1
	newData := []IPAddrData{}
	for _,v := range IPWhiteContent{
		if v.IPAddr != ipaddr{
			newData = append(newData,v)
		}
	}
	if b == true{
		var expireTime int64 = -1
		if isExpireTime == true{
			expireTime = time.Now().Unix() + ExpireTime
		}
		newData = append(newData,IPAddrData{
			ipaddr,
			expireTime,
		})
	}
	IPWhiteContent = []IPAddrData{}
	for _,v := range newData{
		IPWhiteContent = append(IPWhiteContent,v)
	}
	<- Lock
	return nil
}