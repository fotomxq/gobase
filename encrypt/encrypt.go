package encrypt

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
)

//该模块用于定义所有加密摘要、可反向加密方法的封装
//无需任何配置，直接调用对应方法即可使用

//SHA1加密

//获取字符串SHA1摘要
// 该模块返回string类型
//param str string 要加密的字符串
//return string SHA1值，加密失败则返回空字符串
func GetSha1Str(str string) string{
	hasher := sha1.New()
	_, err := hasher.Write([]byte(str))
	if err != nil {
		return ""
	}
	sha := hasher.Sum(nil)
	return hex.EncodeToString(sha)
}

//获取字符串SHA1摘要
// 该模块返回[]byte类型
//param str []byte 要加密的字符串
//return []byte SHA1值，加密失败则返回空字符串
func GetSha1(str []byte) ([]byte,error){
	if len(str) < 1{
		return nil,errors.New("encrypt get sha1 , str is empty.")
	}
	hasher := sha1.New()
	_, err := hasher.Write(str)
	if err != nil {
		return nil,err
	}
	sha := hasher.Sum(nil)
	var dest []byte = make([]byte, hex.EncodedLen(len(str)))
	if dest == nil || sha == nil{
		return nil,errors.New("encrypt get sha1 dest is nil.")
	}
	_ = hex.Encode(dest,sha)
	return dest,nil
}