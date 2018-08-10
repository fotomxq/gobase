package user

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/pkg/errors"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//删除模块
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

//删除用户
//param userID string 用户ID
//return error 错误
func DeleteUserID(userID string) error {
	if bson.IsObjectIdHex(userID) == false{
		return errors.New("user id error.")
	}
	userInfo,err := GetUserByID(userID)
	if err != nil{
		return err
	}
	userInfo.Status = USER_STATUS_TRASH
	return MgoDBC.UpdateId(userInfo.ID,&userInfo)
}

//还原用户
//param userID string 用户ID
//return error 错误
func ReturnUserID(userID string) error {
	if bson.IsObjectIdHex(userID) == false{
		return errors.New("user id error.")
	}
	userInfo,err := GetUserByID(userID)
	if err != nil{
		return err
	}
	userInfo.Status = USER_STATUS_PUBLIC
	return MgoDBC.UpdateId(userInfo.ID,&userInfo)
}

