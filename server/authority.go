package server

import (
	"github.com/pkg/errors"
	"fotomxq/gobase/log"
)

//检查用户是否具备权限
// 不具备权限，直接进入404页面
//param headerParams HeaderParams
//param authority string 权限authority
//return bool 是否成功
func CheckUserAuthorityOrError(headerParams ServerHeaderParams,authority string) bool {
	//没有找到，返回404页面
	if CheckLoginAuthority(headerParams.UserInfo.ID.Hex(),authority) == false{
		log.SendError(headerParams.IP,headerParams.UserInfo.Username,"ServerType.CheckUserAuthorityOrError",errors.New("User no have authority , error url : " + headerParams.R.URL.Path + " , authority : " + authority))
		ReportErrorPage(headerParams)
		return false
	}
	return true
}

//权限检查反馈JSON类型
//param headerParams HeaderParams
//param authority string 权限authority
//return bool 是否成功
func CheckUserAuthorityOrErrorJSON(headerParams ServerHeaderParams,authority string) bool {
	//初始化反馈头
	res := ReportActionType{}
	res.Login = true
	//没有找到，返回404页面
	if CheckLoginAuthority(headerParams.UserInfo.ID.Hex(),authority) == false{
		log.SendError(headerParams.IP,headerParams.UserInfo.Username,"ServerType.CheckUserAuthorityOrError",errors.New("User no have authority , error url : " + headerParams.R.URL.Path + " , authority : " + authority))
		res.Error = "不具备操作权限。"
		ReportJSONData(headerParams,res)
		return false
	}
	return true
}
