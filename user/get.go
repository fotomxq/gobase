package user

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/pkg/errors"
	"net/http"
	"fotomxq/gobase/mgotool"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//查看模块组
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

//根据用户ID获取信息组
//param userID string 用户ID
//return FieldsUser
//return error
func GetUserByID(userID string) (FieldsUser,error) {
	var res FieldsUser = FieldsUser{}
	if bson.IsObjectIdHex(userID) == false{
		return res,errors.New("Bad user ID.")
	}
	err := MgoDBC.FindId(bson.ObjectIdHex(userID)).One(&res)
	return res,err
}

//根据用户登录名获取信息组
//param username string 用户名
//return FieldsUser
//return error
func GetUserByUsername(username string) (FieldsUser,error) {
	var res FieldsUser = FieldsUser{}
	err := MgoDBC.Find(bson.M{"Username" : username}).One(&res)
	return res,err
}

//通过token获取用户信息
//param token string TOKEN值
//return FieldsUser
//return error
func GetUserByToken(token string) (FieldsUser,error) {
	var res FieldsUser = FieldsUser{}
	err := MgoDBC.Find(bson.M{"LoginToken" : token}).One(&res)
	return res,err
}

//获取用户列表
//param r *http.Request
//param search string 搜索内容
//param status []string 状态
//param parent string 上级ID
//param groups []string 用户组
//return []FieldsUser 用户组信息
//return int 用户总个数
//return error 错误
func GetList(r *http.Request,search string,status []string,parent string,groups []string) ([]FieldsUser,int,error){
	q := bson.M{}
	if search != "" {
		q["$or"] = []bson.M{
			bson.M{"Username": bson.M{"$regex": search}},
			bson.M{"Name": bson.M{"$regex": search}},
			bson.M{"Info.Value": bson.M{"$regex": search}},
		}
	}
	for k,v := range status{
		switch v{
		case USER_STATUS_TRASH:
			status[k] = USER_STATUS_TRASH
		case USER_STATUS_STOP:
			status[k] = USER_STATUS_STOP
		default:
			status[k] = USER_STATUS_PUBLIC
		}
	}
	q["Status"] = bson.M{"$in" : status}

	if parent != ""{
		q["ParentID"] = parent
	}

	for _,v := range groups{
		q["Groups.Mark"] = v
	}

	//获取数据
	resList := []FieldsUser{}
	err := mgotool.GetList(r,MgoDBC,q).All(&resList)
	if err != nil{
		return resList,0,err
	}
	count,err := MgoDBC.Find(q).Count()
	//去掉密码字段内容后返回
	if err == nil && count > 0{
		for k,_ := range resList{
			resList[k].Password = ""
		}
	}
	return resList,count,err
}