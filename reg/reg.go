package reg

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"runtime"
	"strconv"
)

//算号器工具
// 自动根据当前计算机，计算出匹配的Key值
// 但需要注意的是，该系统只能由开发人员亲自操作，根据计算机核心数据清单计算出可行的方案

//计算当前Key是否可用于该计算机
func IsOK(key string) bool {
	thisKey := GetKey()
	thisKeyReg := GetKeyReg(thisKey)
	if key == thisKeyReg{
		return true
	}
	return false
}

//获取当前计算机Key值
func GetKey() string {
	//系统信息
	key,_ := os.Hostname()
	//GO环境信息
	key = key + runtime.GOARCH + runtime.GOOS + strconv.Itoa(runtime.NumCPU()) + runtime.Version()
	sysSha1 := string(GetSha1([]byte(key)))
	key2 := string(sysSha1[2]) + string(sysSha1[6]) + string(sysSha1[31]) + string(sysSha1[22]) + string(sysSha1[17]) + string(sysSha1[9]) + string(sysSha1[11]) + string(sysSha1[13]) + string(sysSha1[16]) + string(sysSha1[14]) + string(sysSha1[20]) + string(sysSha1[22])
	key2 += string(sysSha1[33]) + string(sysSha1[24]) + string(sysSha1[21]) + string(sysSha1[36]) + string(sysSha1[14]) + string(sysSha1[34]) + string(sysSha1[12]) + string(sysSha1[1])
	return key2
}

//根据Key获取注册码
func GetKeyReg(sysName string) string{
	sysSha1 := GetSha1([]byte(sysName))
	var key string = string(sysSha1[31]) + string(sysSha1[1]) + string(sysSha1[6]) + string(sysSha1[17]) + "-" + string(sysSha1[7]) + string(sysSha1[9]) + string(sysSha1[11]) + string(sysSha1[13]) + "-" + string(sysSha1[16]) + string(sysSha1[12]) + string(sysSha1[20]) + string(sysSha1[22])
	key += "-" + string(sysSha1[27]) + string(sysSha1[25]) + string(sysSha1[30]) + string(sysSha1[36]) + "-" + string(sysSha1[37]) + string(sysSha1[33]) + string(sysSha1[29]) + string(sysSha1[14])
	return key
}

//获取字符串SHA1摘要
//param str []byte 要加密的字符串
//return []byte SHA1值，加密失败则返回空字符串
func GetSha1(str []byte) []byte{
	hasher := sha1.New()
	_, err := hasher.Write(str)
	if err != nil {
		return nil
	}
	sha := hasher.Sum(nil)
	var dest []byte = make([]byte, hex.EncodedLen(len(str)))
	_ = hex.Encode(dest,sha)
	return dest
}