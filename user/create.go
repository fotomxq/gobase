package user

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/pkg/errors"
	"strconv"
	"fotomxq/gobase/filter"
	"fotomxq/gobase/config"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//创建模块组
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

//构建一个新的用户信息组
// 因为信息非常庞大，所以仅返回部分必要性的内容，如创建日期等，其他内容需自行创建
//param username string 用户名即联系电话
//param password string 密码
//param name string 姓名
//param info []FieldsUserInfo 扩展信息组
//param ip string 创建IP
//param noDefaultGroup bool 是否不给定默认用户组
//return FieldsUser 数据集合
//return error 错误
func Create(username string,password string,name string,info []FieldsUserInfo,ip string,noDefaultGroup bool) (FieldsUser,error) {
	//初始化参数
	var res FieldsUser
	//检查参数
	if filter.CheckPhone(username) == false || filter.CheckPassword(password) == false || filter.CheckNicename(name) == false{
		return res,errors.New("The submitted user name, password, and other information do not meet the standards and cannot create new users.")
	}
	//根据用户名，先尝试查询该数据，如果找到则返回该数据集合
	err := MgoDBC.Find(bson.M{"Username" : username}).One(&res)
	if err == nil{
		return res,errors.New("The user already exists.")
	}
	passwordSha1,err := GetPassword(password)
	if err != nil{
		return res,err
	}
	if passwordSha1 == ""{
		return res,errors.New("Incorrect password.")
	}
	//检查信息组
	err = CheckInfo(info)
	if err != nil{
		return res,err
	}
	//如果找不到该数据集合，则创建
	res = FieldsUser{
		bson.NewObjectId(),
		username,
		passwordSha1,
		name,
		"",
		USER_STATUS_PUBLIC,
		time.Now().Unix(),
		ip,
		"",
		"",
		0,
		[]FieldsUserGroupBind{},
		info,
	}
	//如果需要给定默认用户组，则执行给定操作
	if noDefaultGroup == false{
		//默认给定user权限组
		res.Groups = append(res.Groups,FieldsUserGroupBind{
			USER_GROUP_MARK_USER,
			time.Now().Unix(),
			1,
		})
		//获取默认用户组
		defaultGroupMark,err := config.Get("CreateUserToGroupMark")
		if err != nil{
			return res,errors.New("The default user group configuration information does not exist.")
		}
		defaultGroupExpireHour,err := config.Get("CreateUserToGroupExpireHour")
		if err != nil{
			return res,errors.New("The default user group expiration time configuration does not exist.")
		}
		if defaultGroupMark != ""{
			//转化过期时间
			defaultGroupExpireHourInt64,err := strconv.ParseInt(defaultGroupExpireHour,10,64)
			if err != nil{
				return res,errors.New("The default user group expiration time is incorrectly configured.")
			}
			//获取用户组
			groupInfo,err := GetGroupByMark(defaultGroupMark)
			if err != nil{
				return res,errors.New("The default user group does not exist.")
			}
			var expireTime int64
			if defaultGroupExpireHourInt64 > 0{
				expireTime = time.Now().Unix() + (defaultGroupExpireHourInt64 * 60 * 60)
			}else{
				expireTime = 1
			}
			res.Groups = append(res.Groups,FieldsUserGroupBind{
				groupInfo.Mark,
				time.Now().Unix(),
				expireTime,
			})
		}
	}
	err = MgoDBC.Insert(&res)
	//返回
	return res,err
}
