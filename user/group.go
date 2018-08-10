package user

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/pkg/errors"
	"time"
	"fotomxq/gobase/filter"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//用户组
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

var(
	//用户组状态
	USER_GROUP_STATUS_PUBLIC string = "public"
	USER_GROUP_STATUS_TRASH string = "trash"
	USER_GROUP_STATUS_STOP string = "stop"
	USER_GROUP_STATUS_PRIVATE string = "private"

	//默认管理员用户组类别
	// 具备穿透一切功能的权限
	// 内置权限组，禁止修改、删除该权限
	USER_GROUP_MARK_ADMIN string = "admin"
	//默认用户的基础权限
	USER_GROUP_MARK_USER string = "user"
)

//获取用户组信息
//param groupID string 用户组ID
//return FieldsUserGroup 信息
//return error 错误
func GetGroup(groupID string) (FieldsUserGroup,error){
	res := FieldsUserGroup{}
	if bson.IsObjectIdHex(groupID) == false{
		return res,errors.New("Bad user group ID.")
	}
	err := GroupMgoDBC.FindId(bson.ObjectIdHex(groupID)).One(&res)
	return res,err
}

//获取用户组信息
//param groupID string 用户组ID
//return FieldsUserGroup 信息
//return error 错误
func GetGroupByMark(mark string) (FieldsUserGroup,error){
	res := FieldsUserGroup{}
	if filter.CheckUsername(mark) == false{
		return res,errors.New("get group by mark ,but mark is error, group mark : " + mark)
	}
	err := GroupMgoDBC.Find(bson.M{"Mark" : mark}).One(&res)
	return res,err
}

//获取用户组列表信息
//return []FieldsUserGroup 信息组
//return error 错误
func GetGroupList() ([]FieldsUserGroup,error){
	res := []FieldsUserGroup{}
	err := GroupMgoDBC.Find(nil).Sort("Level").All(&res)
	return res,err
}

//创建新的用户组
//param mark string 标识码
//param authority []string 权限marks组
//param name string 用户组名称
//param des string 描述
//param indexURL string 默认URL
//param level int 优先级
//return error 错误
func CreateGroup(mark string,status string,authority []string,name string,des string,indexURL string,level int) error {
	//过滤数据
	if filter.CheckNicename(name) == false || filter.CheckUsername(mark) == false{
		return errors.New("Submit information illegal.")
	}
	//标识码不能冲突
	res,err := GetGroupByMark(mark)
	if err == nil{
		if res.Mark == mark{
			return errors.New("User group ID cannot be duplicated.")
		}
	}
	switch status{
	case USER_GROUP_STATUS_PUBLIC:
	case USER_GROUP_STATUS_TRASH:
	case USER_GROUP_STATUS_STOP:
	case USER_GROUP_STATUS_PRIVATE:
	default:
		return errors.New("The wrong status value.")
	}
	//增加用户组
	newRes := FieldsUserGroup{
		bson.NewObjectId(),
		mark,
		status,
		authority,
		name,
		des,
		indexURL,
		level,
	}
	return GroupMgoDBC.Insert(newRes)
}

//更新用户组信息
// 注意检查权限组数据
// 不能修改用户组标识码
//param groupID string 用户组ID
//param mark string 标识码
//param status string 状态
//param authority []string 权限marks组
//param name string 名称
//param des string 描述
//param indexURL string 默认URL
//param level int 优先级
//return error 错误
func UpdateGroup(groupID string,mark string,status string,authority []string,name string,des string,indexURL string,level int) error {
	//过滤参数
	if bson.IsObjectIdHex(groupID) == false{
		return errors.New("Bad user group ID.")
	}
	if filter.CheckNicename(name) == false || filter.CheckUsername(mark) == false{
		return errors.New("Submit information illegal.")
	}
	switch status{
	case USER_GROUP_STATUS_PUBLIC:
	case USER_GROUP_STATUS_TRASH:
	case USER_GROUP_STATUS_STOP:
	case USER_GROUP_STATUS_PRIVATE:
	default:
		return errors.New("The wrong status value.")
	}
	//禁止修改管理员
	if mark == USER_GROUP_MARK_ADMIN {
		return errors.New("You cannot update the administrator's user group.")
	}
	//标识码不能重叠
	res,err := GetGroupByMark(mark)
	if err == nil{
		//不是管理员，但修改为管理员？
		if res.ID.Hex() != groupID{
			return errors.New("User group identifiers cannot be duplicated.")
		}
	}
	//获取
	resInfo,err := GetGroup(groupID)
	if err != nil{
		return err
	}
	//禁止修改管理员
	if resInfo.Mark == USER_GROUP_MARK_ADMIN {
		return errors.New("You cannot update the administrator's user group.")
	}
	//更新信息
	resInfo.Name = name
	resInfo.Mark = mark
	resInfo.Status = status
	resInfo.Des = des
	resInfo.Authority = authority
	resInfo.IndexURL = indexURL
	resInfo.Level = level
	return GroupMgoDBC.UpdateId(bson.ObjectIdHex(groupID),&resInfo)
}

//更新用户组信息
// 该方法不过滤任何数据
//param params FieldsUserGroup
//return error 错误信息
func UpdateGroupInParams(params FieldsUserGroup) error{
	return GroupMgoDBC.UpdateId(params.ID,&params)
}

//删除用户组
//param groupID string 用户组ID
//return error 错误
func DeleteGroup(groupID string) error {
	//检查参数
	if bson.IsObjectIdHex(groupID) == false{
		return errors.New("Bad user group ID.")
	}
	//不能删除admin
	res,err := GetGroup(groupID)
	if err != nil{
		return err
	}
	if res.Mark == USER_GROUP_MARK_ADMIN {
		return errors.New("You cannot delete the administrator's user group.")
	}
	//删除ID
	return GroupMgoDBC.RemoveId(bson.ObjectIdHex(groupID))
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//用户和用户组关系
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

//续订用户组
// 也可用于续订服务领域
// 如果期限为0，则默认取消该服务
//param userID string 用户ID
//param groupMark string 用户组标识码
//param timeType int64 时间类型 year 年 / month 月 / day 日 / hour 小时 / trash 取消 / unlimited 无限
//param timeUnit int64 单位长度，1则无限长度
//return error 错误信息
func SetGroupExpireTime(userID string,groupMark string,timeType string,timeUnit int64) error{
	//验证ID是否正确
	if bson.IsObjectIdHex(userID) == false{
		return errors.New("user id error.")
	}
	//获取用户信息
	userInfo,err := GetUserByID(userID)
	if err != nil{
		return errors.New("User info not exist,error : " + err.Error())
	}
	//获取组信息
	groupInfo,err := GetGroupByMark(groupMark)
	if err != nil{
		return errors.New("User group info not exist,error : " + err.Error())
	}
	//初始化服务
	groupBind := FieldsUserGroupBind{}
	//找到该服务
	isFind := false
	for _,v := range userInfo.Groups{
		if v.Mark == groupInfo.Mark{
			isFind = true
			//更新创建时间
			groupBind.CreateTime = v.CreateTime
			break
		}
	}
	//如果没有找到服务
	if isFind == false{
		//推入新的服务项目
		groupBind = FieldsUserGroupBind{
			groupInfo.Mark,
			time.Now().Unix(),
			0,
		}
		userInfo.Groups = append(userInfo.Groups,groupBind)
	}
	//重新查找该服务，进行叠加时间
	for k,v := range userInfo.Groups {
		if v.Mark == groupInfo.Mark {
			//获取当前用户组的时间
			groupTime := time.Unix(v.ExpireTime,0)
			//时间如果低于当前unix，则按照当前时间计算
			if groupTime.Unix() < time.Now().Unix(){
				groupTime = time.Now()
			}
			//计算叠加时间
			switch timeType{
			case "year":
				v.ExpireTime = groupTime.AddDate(int(timeUnit),0,0).Unix()
			case "month":
				v.ExpireTime = groupTime.AddDate(0,int(timeUnit),0).Unix()
			case "day":
				v.ExpireTime = groupTime.AddDate(0,0,int(timeUnit)).Unix()
			case "hour":
				v.ExpireTime += timeUnit * 24 * 60 * 60
			case "trash":
				v.ExpireTime = time.Now().Unix() - 1
			case "unlimited":
				v.ExpireTime = 1
			}
			userInfo.Groups[k] = v
		}
	}
	//更新数据
	return MgoDBC.UpdateId(userInfo.ID,&userInfo)
}