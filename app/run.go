package app

import (
	"time"
	"strings"
	"encoding/json"
	"github.com/pkg/errors"
	"fotomxq/gobase/file"
	"fotomxq/gobase/mgotool"
	"fotomxq/gobase/backup"
	"fotomxq/gobase/distribution"
	"fotomxq/gobase/config"
	"fotomxq/gobase/log"
)

//数据库连接配置
type SystemConfigType struct {
	//数据库URL
	DatabaseMgoURL string `json:"DatabaseMgoURL"`
	//数据库名称
	DatabaseMgoDatabaseName string `json:"DatabaseMgoDatabaseName"`

	//路由广播地址
	RouterHost string `json:"RouterHost"`
	//路由URL
	RouterURL string `json:"RouterURL"`

	//分布式配置信息
	DistributionServiceIP string
	DistributionServicePort string
	DistributionServiceAppMark string
	DistributionServiceSequence string
}

var(
	//系统配置文件信息
	SystemConfig SystemConfigType
	//应用标识码
	AppMark string
)

//初始化构建系统
//return error 错误代码
func Run() error{
	//初始化配置
	config.Run()

	//设置日志
	LogConsole,err := config.Get("LogConsole")
	if err != nil{
		log.LogConsole = true
	}else{
		log.LogConsole = LogConsole == "1"
	}
	LogSaveON,err := config.Get("LogSaveON")
	if err != nil{
		log.LogSaveON = true
	}else{
		log.LogSaveON = LogSaveON == "1"
	}
	log.LogSaveMode = "mgo"
	log.LogDir = "." + file.Sep + "log"
	//启动日志记录服务
	go log.Run()

	//初始化分布式应用体系
	distribution.Run()
	//注册分布式服务
	// 先尝试注销该服务
	_ = distribution.DeleteByMark(SystemConfig.DistributionServiceAppMark,SystemConfig.DistributionServiceSequence)
	err = distribution.Create(SystemConfig.DistributionServiceIP,SystemConfig.DistributionServicePort,SystemConfig.DistributionServiceAppMark,SystemConfig.DistributionServiceSequence,[]distribution.FieldsDistributionServiceSubTasksType{})
	if err != nil{
		return err
	}

	//完成
	return nil
}

//加载配置文件
//return error 错误代码
func RunLoadConfig() error{
	//读取数据库配置
	by,err := file.LoadFile("." + file.Sep + "config.json")
	if err != nil{
		return errors.New("cannot find config file , error : " + err.Error())
	}
	err = json.Unmarshal(by,&SystemConfig)
	if err != nil{
		return errors.New("cannot read config , error : " + err.Error())
	}
	//设置appmark
	AppMark = SystemConfig.DistributionServiceAppMark
	return nil
}

//连接数据库
//return error 错误代码
func RunMgoDB() error{
	//创建数据库连接
	err := mgotool.CreateDB(SystemConfig.DatabaseMgoURL,SystemConfig.DatabaseMgoDatabaseName)
	if err != nil{
		return err
	}
	return nil
}

//初始化备份系统
//return error 错误代码
func RunBackup() error {
	//启动备份
	backup.Run()
	//初始化
	var err error
	//设置参数
	backup.BackupDir,err = config.Get("BackupDir")
	if err != nil{
		return err
	}
	folderList,err := config.Get("BackupFolderList")
	if err != nil{
		return err
	}
	backup.BackupFolders = strings.Split(folderList,"|")
	dbList,err := config.Get("BackupDBList")
	if err != nil{
		return err
	}
	backup.BackupDbs = strings.Split(dbList,"|")
	//完成
	return nil
}

//定时更新一些基础配置信息
func RefConfig() {
	for {
		//更新debug
		debug, err := config.Get("Debug")
		if err != nil {
			debug = "0"
		}
		log.LogConsole = debug == "1"
		//2分钟更新一次
		time.Sleep(time.Minute * 2)
	}
}