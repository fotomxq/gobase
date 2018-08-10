package server

import (
	"net/http"
	"github.com/pkg/errors"
	"fotomxq/gobase/log"
	"fotomxq/gobase/file"
	"fotomxq/gobase/user"
	"fotomxq/gobase/session"
	"fotomxq/gobase/config"
)

//初始化系统
//param appMark string 应用标识码
//return error 错误信息
func Run(appMark string) error {
	//当前用用的标识码
	AppMark = appMark
	//页面通用变量组
	PageGlobData = map[string]interface{}{}
	//全局页面需要加载的文件列
	// 获取所有/base路径下的文件路径
	TemplatesParseFiles = []string{}

	//将配置信息加载到内存中
	RefConfig()

	//自动自动刷新功能
	go AutoRefConfig()

	//返回成功
	return nil
}

//登陆系统初始化
//param appMark string 应用标识码
//return error 错误信息
func RunLogin(appMark string) error {
	//初始化
	err := Run(appMark)
	if err != nil{
		return err
	}

	//设置用户加密块
	user.PasswdEncrypt,err = config.Get("EncryptPasswdEncrypt")
	if err != nil{
		return err
	}
	//设置session
	sessionMark,err := config.Get("SessionMark")
	if err != nil{
		return err
	}
	version,err := config.Get("Version")
	if err != nil{
		return nil
	}
	session.AppName = sessionMark + version

	return err
}

//加载模板基础文件组
// 对外广播API等服务类应用不需要加载该设定
//return error 错误信息
func RunLoadTemplateBase() error {
	//更新模版路径
	// 获取模版路径
	var err error
	TemplatesParseFiles,err = file.GetFileList(TemplatesDir + file.Sep + "base",[]string{},false)
	if err != nil{
		return errors.New("cannot load template base direction, error : " + err.Error())
	}
	return nil
}

//自动化广播
//return error 错误信息
func RunRouter() error{
	//启动服务
	log.SendText("127.0.0.1","",AppMark+".ServerType.RunRouter",log.MessageTypeSystem,"启动服务，地址："+RouterHost)
	err := http.ListenAndServe(RouterHost, nil)
	if err != nil {
		return err
	}
	return nil
}