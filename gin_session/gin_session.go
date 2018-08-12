package gin_session

import (
	"strconv"
	"fotomxq/gobase/filter"
	"github.com/gin-gonic/gin"
)

//session处理器
// 存取cookie值并和数据库匹配，自动删除过期cookie
// 如果启动ip绑定，将确保ip的唯一性
var (
	//超时时间 秒
	// = 3600
	Timeout int64

	//AppName
	AppName string
)

//刷新配置
//param timeout string 超时时间
//return error
func RefConfig(timeout string) error {
	var err error
	Timeout,err = strconv.ParseInt(timeout,10,64)
	if err != nil{
		Timeout = 86400
		return err
	}
	return nil
}

//获取该用户标识码，不存在则创建
//param c *gin.Context
//return string 标识值，如果失败返回空值
//return error
func GetMark(c *gin.Context) (string,error){
	token,err := GetCookie(c)
	if err != nil || token == ""{
		token = c.PostForm("cookie")
		if token == ""{
			tokenStr := AppName + c.ClientIP() + filter.GetRandStr(95348671259875)
			token,err = filter.GetSha1(tokenStr)
			if err != nil{
				return "",err
			}
		}
	}
	SetCookie(c,token)
	return token,nil
}

//获取cookie
//param c *gin.Context
//return string cookie
//return error
func GetCookie(c *gin.Context) (string,error){
	return c.Cookie(AppName)
}

//设定cookie
//param c *gin.Context
//param token string token
func SetCookie(c *gin.Context,token string) {
	c.SetCookie(AppName,token,int(Timeout),"/","",false,true)
}