package Temple

import "time"

//页面输出通用模版、管道库处理

//将两个字符串合并
//param a sting
//param b string
//return string
func MargeString(a string,b string) string{
	return a + b
}

//将Unix时间戳转换为年月日时间
// 管道函数
//param unix int64 Unix时间戳
//return string 年月日时间 eg : 2017-9-6 10:15:06
func GetUnixDataTime(unix int64) string{
	return time.Unix(unix,0).Format("2006-01-02 15:04:05")
}
