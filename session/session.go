package session

import (
	"net/http"
	"strconv"
	"github.com/pkg/errors"
	"fotomxq/gobase/filter"
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
//param w http.ResponseWriter
//param r *http.Request
//param ip string IP地址
//return string 标识值，如果失败返回空值
//return error
func GetMark(w http.ResponseWriter,r *http.Request,ip string) (string,error){
	token,err := GetCookie(w,r)
	if err != nil || token == ""{
		tokenStr := ip + filter.GetRandStr(95348671259875)
		token,err = filter.GetSha1(tokenStr)
		if err != nil{
			return "",err
		}
	}
	SetCookie(w,r,token)
	return token,nil
}

//获取cookie
//param w http.ResponseWriter
//param r *http.Request
//return string cookie
//return error
func GetCookie(w http.ResponseWriter,r *http.Request) (string,error){
	cookie,err := r.Cookie(AppName)
	if err != nil{
		return "",err
	}
	if cookie.Value == ""{
		return "",errors.New("cannot create cookie.")
	}
	return cookie.Value,nil
}

//设定cookie
//param w http.ResponseWriter
//param r *http.Request
//param token string token
func SetCookie(w http.ResponseWriter,r *http.Request,token string) {
	cookie := http.Cookie{
		Name:AppName,
		Value: token,
		Path:"/",
		HttpOnly: true,
		MaxAge:int(Timeout),
	}
	http.SetCookie(w,&cookie)
}