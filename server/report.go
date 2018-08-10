package server

import (
	"fotomxq/gobase/user"
	"fotomxq/gobase/log"
	"fotomxq/gobase/backup"
	"fotomxq/gobase/system"
	"fotomxq/gobase/distribution"
)

//反馈头信息组

//通用标准头

//标准动作反馈头
// 任意字符串
type ReportActionType struct {
	//完成状态
	Status bool `json:"Status"`
	//错误信息
	Error string `json:"Error"`
	//是否已经登陆
	Login bool `json:"Login"`
	//反馈数据是否为缓冲数据
	Cache bool `json:"Cache"`
	//URL地址
	// 针对应用内指向请求URL地址，针对服务器是请求数据的URL地址，用于给cache缓冲数据
	URL string `json:"URL"`
	//附带数据内容
	Data string `json:"Data"`
}

//标准动作反馈头
// interface
type ReportActionInterfaceType struct {
	//完成状态
	Status bool `json:"Status"`
	//错误信息
	Error string `json:"Error"`
	//是否已经登陆
	Login bool `json:"Login"`
	//反馈数据是否为缓冲数据
	Cache bool `json:"Cache"`
	//URL地址
	// 针对应用内指向请求URL地址，针对服务器是请求数据的URL地址，用于给cache缓冲数据
	URL string `json:"URL"`
	//附带数据内容
	Data interface{} `json:"Data"`
}

//获取用户列队信息组
type ReportUserListType struct {
	//完成状态
	Status bool `json:"Status"`
	//错误信息
	Error string `json:"Error"`
	//是否已经登陆
	Login bool `json:"Login"`
	//反馈数据是否为缓冲数据
	Cache bool `json:"Cache"`
	//URL地址
	// 针对应用内指向请求URL地址，针对服务器是请求数据的URL地址，用于给cache缓冲数据
	URL string `json:"URL"`
	//数据总量
	Count int
	//数据列表
	Data []user.FieldsUser `json:"Data"`
}

//获取单一用户信息组
type ReportUserType struct {
	//完成状态
	Status bool `json:"Status"`
	//错误信息
	Error string `json:"Error"`
	//是否已经登陆
	Login bool `json:"Login"`
	//反馈数据是否为缓冲数据
	Cache bool `json:"Cache"`
	//URL地址
	// 针对应用内指向请求URL地址，针对服务器是请求数据的URL地址，用于给cache缓冲数据
	URL string `json:"URL"`
	//数据列表
	Data user.FieldsUser `json:"Data"`
}

//获取一列用户名称组
type ReportUsersByIDType struct {
	//完成状态
	Status bool `json:"Status"`
	//错误信息
	Error string `json:"Error"`
	//是否已经登陆
	Login bool `json:"Login"`
	//反馈数据是否为缓冲数据
	Cache bool `json:"Cache"`
	//URL地址
	// 针对应用内指向请求URL地址，针对服务器是请求数据的URL地址，用于给cache缓冲数据
	URL string `json:"URL"`
	//数据总量
	Count int
	//数据列表
	Data []ReportUsersByIDDataType `json:"Data"`
}
//获取一列用户名称组 数据组
type ReportUsersByIDDataType struct{
	//用户ID
	UserID string `json:"UserID"`
	//用户名
	Username string `json:"Username"`
	//昵称
	Name string `json:"Name"`
}

//获取日志列队
type ReportLogListType struct {
	//完成状态
	Status bool `json:"Status"`
	//错误信息
	Error string `json:"Error"`
	//是否已经登陆
	Login bool `json:"Login"`
	//反馈数据是否为缓冲数据
	Cache bool `json:"Cache"`
	//URL地址
	// 针对应用内指向请求URL地址，针对服务器是请求数据的URL地址，用于给cache缓冲数据
	URL string `json:"URL"`
	//数据总量
	Count int
	//数据列表
	Data []log.FieldsLogMgoData `json:"Data"`
}

//获取备份文件列队
type ReportBackupListType struct {
	//完成状态
	Status bool `json:"Status"`
	//错误信息
	Error string `json:"Error"`
	//是否已经登陆
	Login bool `json:"Login"`
	//反馈数据是否为缓冲数据
	Cache bool `json:"Cache"`
	//URL地址
	// 针对应用内指向请求URL地址，针对服务器是请求数据的URL地址，用于给cache缓冲数据
	URL string `json:"URL"`
	//数据总量
	Count int
	//数据列表
	Data []backup.BackupListType `json:"Data"`
}

//获取系统运行状态数据
type ReportSystemRunType struct {
	//完成状态
	Status bool `json:"Status"`
	//错误信息
	Error string `json:"Error"`
	//是否已经登陆
	Login bool `json:"Login"`
	//反馈数据是否为缓冲数据
	Cache bool `json:"Cache"`
	//URL地址
	// 针对应用内指向请求URL地址，针对服务器是请求数据的URL地址，用于给cache缓冲数据
	URL string `json:"URL"`
	//数据总量
	Count int
	//数据列表
	Data []system.FieldsSystemDataType `json:"Data"`
}

//系统运行状态
type ReportSystemRunStatType struct {
	//完成状态
	Status bool `json:"Status"`
	//错误信息
	Error string `json:"Error"`
	//是否已经登陆
	Login bool `json:"Login"`
	//反馈数据是否为缓冲数据
	Cache bool `json:"Cache"`
	//URL地址
	// 针对应用内指向请求URL地址，针对服务器是请求数据的URL地址，用于给cache缓冲数据
	URL string `json:"URL"`
	//数据总量
	Count int `json:"Count"`
	//数据列表
	Data []system.FieldsSystemDataType `json:"Data"`
	//应用名称
	AppName string `json:"AppName"`
	//分布式标识码
	DistributedMark string `json:"DistributedMark"`
	//进程启动总时间
	SystemStartTime int64 `json:"SystemStartTime"`
	//发生错误次数
	LogErrorCount int64 `json:"LogErrorCount"`
	//发生HTTP访问次数
	LogHttpCount int64 `json:"LogHttpCount"`
	//发生访问次数
	LogVisitCount int64 `json:"LogVisitCount"`
	//系统消息次数
	LogSystemCount int64 `json:"LogSystemCount"`
	//系统安全事件次数
	LogSafetyCount int64 `json:"LogSafetyCount"`
	//总的系统消息
	LogCount int64 `json:"LogCount"`
}

//获取服务列表
type ReportDistributionServiceListType struct {
	//完成状态
	Status bool `json:"Status"`
	//错误信息
	Error string `json:"Error"`
	//是否已经登陆
	Login bool `json:"Login"`
	//反馈数据是否为缓冲数据
	Cache bool `json:"Cache"`
	//URL地址
	// 针对应用内指向请求URL地址，针对服务器是请求数据的URL地址，用于给cache缓冲数据
	URL string `json:"URL"`
	//数据总量
	Count int
	//数据列表
	Data []distribution.FieldsDistributionServiceType `json:"Data"`
}