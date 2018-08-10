package user

import (
	"gopkg.in/mgo.v2"
	"github.com/pkg/errors"
		"strings"
	"fotomxq/gobase/config"
	"fotomxq/gobase/encrypt"
)

//用户操作模块

var(
	//数据库集合操作句柄
	MgoDBC *mgo.Collection
	//用户组数据库集合
	GroupMgoDBC *mgo.Collection

	//用户自动退出时限，单位：秒
	//1600
	UserLoginTimeout int64

	//密码加密用字符串
	PasswdEncrypt string = ""

	//用户状态
	USER_STATUS_PUBLIC string = "public"
	USER_STATUS_TRASH string = "trash"
	USER_STATUS_STOP string = "stop"

	//配置信息组
	// session是否与ip绑定
	BindCookieAndIP bool = false
	// 是否 检查页面URL与用户组权限是否匹配
	PageAuthorityURLON bool = false
	// 缓冲用户组
	CacheUserGroups []FieldsUserGroup = []FieldsUserGroup{}
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//其他通用模块组
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

//计算密码SHA1值
//param passwd string 密码
//return string SHA1密匙，空则失败
//return error
func GetPassword(passwd string) (string,error){
	if passwd == ""{
		return "",errors.New("passwd is empty.")
	}
	sha1,err := encrypt.GetSha1([]byte(passwd + PasswdEncrypt))
	if err != nil{
		return "",err
	}
	return string(sha1),nil
}

//检查信息组
// 信息组如果存在数据，则必须存在对应的mark
//param info []FieldsUserInfo 信息组
//return error
func CheckInfo(info []FieldsUserInfo) error {
	//获取配置信息
	infoConfig,err := GetInfoConfig()
	if err != nil{
		return err
	}
	infoForceConfig,err := config.Get("UserInfoStructForce")
	if err != nil{
		return err
	}
	//信息组如果存在数据，则必须存在对应的mark
	for _,v := range info{
		if v.Value != ""{
			//是否存在标识码？
			if v.Mark == ""{
				return errors.New("info have value but not set mark.")
			}
			//如果未开启强制结构体，则跳过后续检查方案
			if infoForceConfig == "0"{
				continue
			}
			//该标识码是否存在于配置？
			isFind := false
			for _,vConfig := range infoConfig{
				if v.Mark == vConfig.Mark{
					isFind = true
					break
				}
			}
			if isFind == false{
				return errors.New("info mark not config.")
			}
		}
	}
	return nil
}

//获取信息组配置
//return map[string]string 信息组 标识码=>标题
func GetInfoConfig() ([]FieldsUserInfo,error){
	res := []FieldsUserInfo{}
	infoConfig,err := config.Get("UserInfoStruct")
	if err != nil{
		return res,err
	}
	//进行第一次拆分
	// a,b|c,d... => {{a,b},{c,d},...}
	infoConfigA := strings.Split(infoConfig,"|")
	for _,v := range infoConfigA{
		if v == ""{
			continue
		}
		infoConfigB := strings.Split(v,",")
		if len(infoConfigB) != 2{
			continue
		}
		res = append(res,FieldsUserInfo{
			infoConfigB[0],
			infoConfigB[1],
			"",
		})
	}
	return res,nil
}