package user

import (
	"time"
	"strings"
	"gopkg.in/mgo.v2/bson"
	"github.com/pkg/errors"
	"fotomxq/gobase/filter"
)

//注意，请勿使用该模块，app api将内置全能处理器，方便处理

//检查用户是否具备某一组权限
// 只要满足其中一个权限，则通过
//param userID string 用户ID
//param authorityList string 权限标识码，以|分割，可选选项，只要有一个通过全部通过
//return bool 是否成功
func CheckLoginAuthority(userID string,authorityList string) bool {
	//获取用户信息
	userInfo,err := GetUserByID(userID)
	if err != nil{
		return false
	}
	//拆分查询的权限列表
	authority := strings.Split(authorityList,"|")
	//根据用户的用户组，遍历并查询
	for _,vUserGroup := range userInfo.Groups{
		//获取用户组信息
		groupInfo,err := GetGroupByMark(vUserGroup.Mark)
		if err != nil{
			return false
		}
		//用户组是否为可用状态？
		if groupInfo.Status != USER_GROUP_STATUS_PUBLIC && groupInfo.Status != USER_GROUP_STATUS_PRIVATE{
			return false
		}
		//检查用户组，是否具备该权限？
		for _,vAuthority := range groupInfo.Authority{
			for _,v := range authority{
				if vAuthority == v{
					//如果发现用户组具备该权限，则直接判断该用户的用户组是否可用
					_,err = CheckUserGroup(userInfo.ID.Hex(),groupInfo.Mark,true)
					if err == nil{
						return true
					}
				}
			}
		}
	}
	return false
}

//检查用户是否存在某用户组
// 同时将返回用户组数据
//param userID string 用户ID
//param groupMark string 用户组标识码
//param expireBool bool 是否判断过期
//return string 用户组mark
//return error 错误信息
func CheckUserGroup(userID string,groupMark string,expireBool bool) (string,error) {
	//检查参数
	if bson.IsObjectIdHex(userID) == false || filter.CheckUsername(groupMark) == false{
		return "",errors.New("user id or group mark error.")
	}
	//获取用户信息
	userInfo := FieldsUser{}
	err := MgoDBC.Find(bson.M{"_id" : bson.ObjectIdHex(userID), "Groups.Mark" : groupMark}).One(&userInfo)
	if err != nil{
		return "",errors.New("cannot find user id , or user mark not exist.")
	}
	//检查是否过期？
	if expireBool == true {
		for _, v := range userInfo.Groups {
			if v.Mark == groupMark {
				//1则无限授权时间
				if v.ExpireTime == 1 {
					break
				}
				//超出过期
				if v.ExpireTime < time.Now().Unix() {
					return "", errors.New("user group expired.")
				}
				//允许授权
				break
			}
		}
	}
	//返回成功
	return groupMark,err
}