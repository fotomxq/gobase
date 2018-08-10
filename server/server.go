package server

import (
	"net/http"
	"fotomxq/gobase/user"
)

//server/service通用的基础设置模块
// 所有广播或非广播类服务都必须先行调用该函数，以激活程序
// 其他工具类服务可根据业务情况调用

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
//标准化网络头套件
//////////////////////////////////////////////////////////////////////////////////////////////////////////////

//标准化网络头套件
// 将常用套件组合为单一的变量体进行传递，以强化代码的高复可用性
type ServerHeaderParams struct {
	//写入网络头
	W http.ResponseWriter
	//读取网络头
	R *http.Request
	//URL拆分结构体
	URLS []string
	//IP地址
	IP string
	//Token值
	Token string
	//用户基本信息
	UserInfo user.FieldsUser
}

var(
	//当前用用的标识码
	AppMark string

	//页面通用变量组
	PageGlobData map[string]interface{}

	//全局页面需要加载的文件列
	// 获取所有/base路径下的文件路径
	TemplatesParseFiles []string

	//服务地址
	RouterHost string

	//模板路径
	TemplatesDir string

	//是否检查权限URL是否匹配？
	PageAuthorityURLON bool
)