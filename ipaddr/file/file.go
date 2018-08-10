package file

import (
	"strings"
	"io/ioutil"
	"time"
	"fotomxq/gobase/file"
)

//录入文件操作方法

var(
	//是否启动黑名单
	// = true
	IPBanOn bool

	//是否启动白名单
	//= false
	IPWhiteOn bool

	//黑名单列表
	banList []string

	//白名单列表
	whiteList []string

	//配置文件路径
	// = "." + sep + "configs" + sep
	ConfigDir string

	//拉黑的IP地址
	// = configDir + "ip-ban.json"
	IPConfigBanSrc string

	//白名单IP地址
	// = configDir + "ip-white.json"
	IPConfigWhiteSrc string

	//是否编辑过数据
	// = false
	isEdit bool
)

//启动IP服务
func RunFile(){
	IPBanOn = true
	IPWhiteOn = false
	ConfigDir = "." + file.Sep + "configs" + file.Sep
	IPConfigBanSrc = ConfigDir + "ip-ban.json"
	IPConfigWhiteSrc = ConfigDir + "ip-white.json"
	isEdit = false

	loadIPConfig()
	go autoSaveIPConfig()
}

//获取IP原始数据
//param isBan bool 是否为黑名单数据
//return string 数据字符串，用|连接多个IP
func GetAllStr(isBan bool) string {
	var content string
	if isBan == true{
		content = strings.Join(banList,"|")
	}else{
		content = strings.Join(whiteList,"|")
	}
	return content
}

//是否禁止IP访问
//param ip string IP地址
//return bool 是否允许访问
func IsIP(ip string) bool{
	if IPBanOn == true{
		if IsIPBan(ip) == true{
			return false
		}
	}
	if IPWhiteOn == true{
		if IsIPWhite(ip) == false{
			return false
		}
	}
	return true
}

//判断IP是否在黑名单内
//param ip string IP地址
//return bool 是否在列表内
func IsIPBan(ip string) bool{
	return searchList(ip,banList)
}

//判断IP是否在白名单内
//param ip string IP地址
//return bool 是否在列表内
func IsIPWhite(ip string) bool{
	return searchList(ip,whiteList)
}

//IP是否列入黑名单和白名单
// 移除直接设置为false即可实现
//param ip string IP地址
//param isBan bool 是否列入黑名单
//param isWhite bool 是否列入白名单
func SetIPBan(ip string,isBan bool,isWhite bool){
	defer setIsEdit()
	banList = listAddDelete(banList,ip,isBan)
	whiteList = listAddDelete(whiteList,ip,isWhite)
}

//直接修改名单列队
// 直接替换列队数据
//param ipContent string IP数据列，用|符号分割
//param isBan bool 是否拉黑，如果不是则为白名单
//return bool 修改是否成功
func SetIPList(ipContent string,isBan bool) bool{
	defer setIsEdit()
	if isBan == true{
		banList = []string{}
	}else{
		whiteList = []string{}
	}
	if ipContent == ""{
		return true
	}
	ipList := strings.Split(ipContent,"|")
	if isBan == true{
		banList = ipList
	}else{
		whiteList = ipList
	}
	return true
}

//修改编辑状态
// 用于内部返回前处理
func setIsEdit(){
	isEdit = true
}

//读取IP配置数据
// 读取后的数据将直接存储到变量内调用
// 如果不存在文件则忽略
func loadIPConfig(){
	bb,err := ioutil.ReadFile(IPConfigBanSrc)
	if err == nil && len(bb) > 0{
		banList = strings.Split(string(bb),"|")
	}
	wb,err := ioutil.ReadFile(IPConfigWhiteSrc)
	if err == nil && len(wb) > 0{
		whiteList = strings.Split(string(wb),"|")
	}
}

//自动保存IP
func autoSaveIPConfig(){
	ticker := time.NewTicker(time.Millisecond  * 400)
	for _ = range ticker.C {
		saveIPConfig()
	}
}

//保存IP数据配置
// 只有在修改IP的时候触发
//return error
func saveIPConfig() error {
	if isEdit == false{
		return nil
	}
	banContent := strings.Join(banList,"|")
	err := ioutil.WriteFile(IPConfigBanSrc,[]byte(banContent),0666)
	if err != nil{
		return err
	}
	whiteContent := strings.Join(banList,"|")
	err = ioutil.WriteFile(IPConfigWhiteSrc,[]byte(whiteContent),0666)
	if err != nil{
		return err
	}
	return nil
}

//[]string列队写入和删除
//param list []string 列队
//param s string 值
//param b bool 是否写入列队
//return []string 新的列队
func listAddDelete(list []string,s string,b bool) []string{
	var newList []string
	ok := false
	for _,v := range list{
		if v == s{
			if b == true{
				newList = append(newList,s)
				ok = true
			}else{
				ok = true
			}
		}else{
			newList = append(newList,s)
		}
	}
	if ok == false{
		if b == true{
			newList = append(newList,s)
		}
	}
	return newList
}

//搜索IP是否在列表内
//param ip string IP地址
//param list []string 列表
//return bool 是否存在于列表内
func searchList(ip string,list []string) bool {
	for _,value := range list{
		if ip == value{
			return true
		}
	}
	return false
}