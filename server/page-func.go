package server

import (
	"html/template"
	"strings"
	"time"
	"fotomxq/gobase/user"
)

//该文件为server分支文件
// 仅用于页面内置函数

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
//页面内函数组
//////////////////////////////////////////////////////////////////////////////////////////////////////////////

//页面输出非转义字符串
//param str string
//return interface{}
func PageFuncUnescaped (str string) interface{} {
	return template.HTML(str)
}


//页面输出 转换Unix时间戳为指定的时间格式
//param str int64
//param format string eg : YY / MM / DD H : M : S
//return string 格式化时间
func PageFuncGetUnitToDatetime (str int64,format string) string {
	//如果为空，则返回
	if str < 1 || format == ""{
		return ""
	}
	//替换结构中的字符
	format = strings.Replace(format,"YY","2006",-1)
	format = strings.Replace(format,"MM","01",-1)
	format = strings.Replace(format,"DD","02",-1)
	format = strings.Replace(format,"H","15",-1)
	format = strings.Replace(format,"M","04",-1)
	format = strings.Replace(format,"S","05",-1)
	//获取时间
	t := time.Unix(str,0)
	//返回标准时间
	return t.Format(format)
}

//页面中获取某个配置信息的数据组
//param name string
//return string
func PageFuncGetConfigValue(name string) string{
	return PageGlobData[name].(string)
}

//检查用户是否具备对应权限
// 只要满足其中一个权限，则通过
//param userID string 用户ID
//param authorityList string 权限标识码，以|分割，可选选项，只要有一个通过全部通过
//return bool 是否成功
func CheckLoginAuthority(userID string,authorityList string) bool {
	return user.CheckLoginAuthority(userID,authorityList)
}

//通过标记获取服务信息组
//param mark string 标识
//return FieldsUserGroup 信息组
//return bool 是否存在
func GetServiceInfoByMark(mark string) (user.FieldsUserGroup,bool){
	res,err := user.GetGroupByMark(mark)
	return res,err == nil
}

//通过ID获取服务组信息
//param id string 标识
//return FieldsUserGroup 信息组
//return bool 是否存在
func GetServiceInfoByID(id string) (user.FieldsUserGroup,bool){
	res,err := user.GetGroup(id)
	return res,err == nil
}
