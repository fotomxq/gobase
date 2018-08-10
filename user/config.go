package user

import (
	"strconv"
)

//刷新配置
//param userLoginTimeout string 超时时间配置信息
//param bindCookieAndIP string 是否绑定会话IP 1/0
//param pageAuthorityURLON string URL与权限是否绑定
//param encryptPasswdEncrypt string
//return error
func RefConfig(userLoginTimeout string,bindCookieAndIP string,pageAuthorityURLON string,encryptPasswdEncrypt string) error {
	var err error
	//userLoginTimeout
	UserLoginTimeout,err = strconv.ParseInt(userLoginTimeout,10,64)
	if err != nil{
		UserLoginTimeout = 86400
		return err
	}
	//bindCookieAndIP
	BindCookieAndIP = bindCookieAndIP == "1"
	//pageAuthorityURLON
	PageAuthorityURLON = pageAuthorityURLON == "1"
	//EncryptPasswdEncrypt
	PasswdEncrypt = encryptPasswdEncrypt
	//缓冲用户组
	CacheUserGroups,err = GetGroupList()
	if err != nil{
		return err
	}
	return nil
}