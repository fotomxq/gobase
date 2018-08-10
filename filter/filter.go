package filter

import (
	"strings"
	"strconv"
	"time"
	"regexp"
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
)

//过滤器

//分解URL获取名称和类型
//param sendURL URL地址
//return map[string]string 返回值集合
func GetURLNameType(sendURL string) map[string]string {
	res := map[string]string{
		"full-name": "",
		"only-name": "",
		"type":      "",
	}
	urls := strings.Split(sendURL, "/")
	if len(urls) < 1 {
		return res
	}
	res["full-name"] = urls[len(urls)-1]
	if res["full-name"] == "" {
		res["only-name"] = res["full-name"]
		return res
	}
	names := strings.Split(res["full-name"], ".")
	if len(names) < 2 {
		return res
	}
	res["type"] = names[len(names)-1]
	for i := 0; i <= len(names); i++ {
		if i == len(names)-1 {
			break
		}
		if res["only-name"] == "" {
			res["only-name"] = names[i]
		} else {
			res["only-name"] += "." + names[i]
		}
	}
	return res
}

//验证搜索类型的字符串
//param str string 字符串
//return bool 是否正确
func CheckSearch(str string) bool {
	return MatchStr(`^[\u4e00-\u9fa5_a-zA-Z0-9]+$`, str)
}

//过滤非法字符
//param str string 要过滤的字符串
//return string 过滤后的字符串
func FilterStr(str string) string{
	//str = strings.Replace(str,"\r","",-1)
	//str = strings.Replace(str,"\n","",-1)
	//str = strings.Replace(str,"\t","",-1)
	str = strings.Replace(str,"~","～",-1)
	str = strings.Replace(str,"<","〈",-1)
	str = strings.Replace(str,">","〉",-1)
	str = strings.Replace(str,"$","￥",-1)
	str = strings.Replace(str,"!","！",-1)
	str = strings.Replace(str,"[","【",-1)
	str = strings.Replace(str,"]","】",-1)
	str = strings.Replace(str,"{","｛",-1)
	str = strings.Replace(str,"}","｝",-1)
	str = strings.Replace(str,"/","／",-1)
	str = strings.Replace(str,"\\","﹨",-1)
	return str
}

//过滤非法字符后判断其长度是否符合标准
//param str string 要过滤的字符串
//param min int 最短，包括该长度
//param max int 最长，包括该长度
//return string 过滤后的字符串，失败返回空字符串
func CheckFilterStr(str string,min int,max int) string{
	var newStr string
	newStr = FilterStr(str)
	if newStr == ""{
		return ""
	}
	var strLen int
	strLen = len(newStr)
	if strLen >= min && strLen <= max{
		return newStr
	}
	return ""
}
//验证是否为IP地址
//param str string IP地址
//return bool 是否正确
func CheckIP(str string) bool {
	if str == "[::1]" {
		return true
	}
	if MatchStr(`((2[0-4]\d|25[0-5]|[01]?\d\d?)\.){3}(2[0-4]\d|25[0-5]|[01]?\d\d?)`, str) == true {
		return true
	}
	if MatchStr(`^$(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`, str) == true {
		return true
	}
	return false
}

//处理page
//param postPage string 用户提交的page
//return int 过滤后的页数
func FilterPage(postPage string) int {
	res, err := strconv.Atoi(postPage)
	if err != nil {
		res = 1
	}
	if res < 1 {
		res = 1
	}
	return res
}

//处理max
//限制最小值为1，最大值为999
//param postMax string 用户提交的max
//return int 过滤后的页数
func FilterMax(postMax string) int {
	res, err := strconv.Atoi(postMax)
	if err != nil {
		res = 1
	}
	if res < 1 {
		res = 1
	}
	if res > 999 {
		res = 999
	}
	return res
}
//获取随机字符串
//param n int 随机码
//return string 新随机字符串
func GetRandStr(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	re := r.Intn(n)
	return strconv.Itoa(re)
}

//param mStr string 验证
//param str string 要验证的字符串
//return bool 是否成功
func MatchStr(mStr string, str string) bool {
	res, err := regexp.MatchString(mStr, str)
	if err != nil {
		return false
	}
	return res
}
//验证是否为SHA1
//param str string 字符串
//return bool 是否正确
func CheckHexSha1(str string) bool {
	return MatchStr(`^[a-z0-9]{10,45}$`, str)
}

//获取字符串的SHA1值
//param content string 要计算的字符串
//return string 计算出的SHA1值
//return error
func GetSha1(content string) (string,error) {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(content))
	if err != nil {
		return "",err
	}
	sha := hasher.Sum(nil)
	return hex.EncodeToString(sha),nil
}

//截取字符串
//param str string 要截取的字符串
//param star int 开始位置
//param length int 长度
//return string 新字符串
func SubStr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0
	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length
	if start > end {
		start, end = end, start
	}
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

//检查用户名
//param str string 用户名
//return bool 是否正确
func CheckUsername(str string) bool {
	if str == ""{
		return false
	}
	return MatchStr(`^[a-zA-Z0-9_-]{4,20}$`, str)
}

//检查昵称
//param str string 昵称
//return bool 是否正确
func CheckNicename(str string) bool {
	if str == ""{
		return false
	}
	//return matchStr(`^[\u4e00-\u9fa5_a-zA-Z0-9]+$`, str)
	return MatchStr(`^[\p{Han}_a-zA-Z0-9]{2,50}$`, str)
}

//验证邮箱
//param str string 邮箱地址
//return bool 是否正确
func CheckEmail(str string) bool {
	if str == ""{
		return false
	}
	return MatchStr(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, str)
}

//验证密码
//param str string 密码
//return bool 是否正确
func CheckPassword(str string) bool {
	if str == ""{
		return false
	}
	return MatchStr(`^[a-zA-Z0-9_-]{6,30}$`, str)
}

//验证身份证
// 因为复杂性，仅考虑验证身份证位数有效性
// 未来可根据实际需求加入外部API对身份证进行二次验证
//param str string 身份证号码
//return bool 是否正确
func CheckIDCard(str string) bool{
	if str == ""{
		return false
	}
	if len(str) > 10 && len(str) < 20{
		return true
	}
	return false
}

//验证电话号码
// 必须是手机电话号码或带区号的固定电话号码
// eg 03513168322
// eg 13066889999
func CheckPhone(str string) bool {
	if str == ""{
		return false
	}
	return MatchStr(`^[0-9]{11}$`, str)
}