package server

import (
	"time"
	"fotomxq/gobase/log"
	"fotomxq/gobase/distribution"
	"fotomxq/gobase/config"
	"fotomxq/gobase/user"
	"fotomxq/gobase/session"
	"fotomxq/gobase/ipaddr"
)


//定时刷新配置信息，每10分钟进行一次
// 不推荐使用该模块，因为会增加数据库负荷，建议采用被动式修改方案，后台修改后直接修改对应的位置值
// 该刷新的配置和配置内部的缓冲结构不同，是针对不同包的配置信息进行修改
func AutoRefConfig(){
	for{
		//将所有配置信息加载到内存
		RefConfig()

		//更新子服务项目 Server.AutoRefConfig
		err := distribution.UpdateSubTaskBySelf("server.AutoRefConfig")
		if err != nil{
			log.SendError("0.0.0.0","","server.AutoRefConfig",err)
		}

		//10分钟更新一次
		time.Sleep(time.Minute * 5)
	}
}

//加载所有配置项目
func RefConfig(){
	//基础参数
	headerParams := ServerHeaderParams{}
	headerParams.IP = "0.0.0.0"

	//构建页面通用变量组
	// 将所有配置数据加载到页面通用配置组
	configDatas,err := config.GetAll()
	if err == nil{
		for _,v := range configDatas{
			PageGlobData[v.Name] = v.Value
		}
	}else{
		log.SendError(headerParams.IP,headerParams.UserInfo.Username,"server.RefConfig()",err)
	}

	//刷新ipadd设置
	err = ipaddr.RefConfig(PageGlobData["IPBanON"].(string),PageGlobData["IPWhiteON"].(string))
	if err != nil{
		log.SendError(headerParams.IP,headerParams.UserInfo.Username,"server.RefConfig()",err)
	}

	//设置用户组件
	err = user.RefConfing(PageGlobData["UserLoginTimeout"].(string),PageGlobData["BindCookieAndIP"].(string))
	if err != nil{
		log.SendError(headerParams.IP,headerParams.UserInfo.Username,"server.RefConfig()",err)
	}
	err = session.RefConfig(PageGlobData["SessionTimeout"].(string))
	if err != nil{
		log.SendError(headerParams.IP,headerParams.UserInfo.Username,"server.RefConfig()",err)
	}

	//设置模板
	TemplatesDir,err = config.Get("TemplatesDir")
	if err != nil{
		TemplatesDir = "./public/"
		log.SendError(headerParams.IP,headerParams.UserInfo.Username,"server.RefConfig()",err)
	}

	//设置页面权限URL锁定
	pageAuthorityURLON,err := config.Get("PageAuthorityURLON")
	if err != nil{
		PageAuthorityURLON = true
		log.SendError(headerParams.IP,headerParams.UserInfo.Username,"server.RefConfig()",err)
	}else{
		PageAuthorityURLON = pageAuthorityURLON == "1"
	}
}