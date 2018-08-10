package user

import (
	"github.com/pkg/errors"
	"fotomxq/gobase/filter"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//更新模块组
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

//更新数据服务
// 使用该服务前，请务必确保该数据源是从数据库获取到的一手数据源，不要自行组件数据！
//param params FieldsUser 用户信息
//return bool 是否成功
func UpdateData(params FieldsUser) error {
	//检查参数
	// 仅检查用户名、电话、身份证、姓名信息，其他信息请根据情况提前核对
	// 一般其他信息的修改，仅管理员级别能使用，所以可选择信任来源
	if filter.CheckPhone(params.Username) == false || filter.CheckNicename(params.Name) == false{
		return errors.New("Incorrect username, nickname.")
	}
	//用户名不能重叠存在
	// 系统最初不允许用户名重叠，所以未来也不可能出现重叠的用户体系
	// 如果出现重叠的数据，则说明数据库存在异常问题，或中间某些环节直接修改了数据库内容
	// 用户所有修改方案只能通过该方法执行
	userInfo,err := GetUserByUsername(params.Username)
	if err == nil && userInfo.ID.Hex() != params.ID.Hex(){
		if userInfo.Username == params.Username{
			return errors.New("The user name already exists and cannot be modified.")
		}
	}
	//如果修改了密码？
	if params.Password != "" && params.Password != userInfo.Password{
		newPassword,err := GetPassword(params.Password)
		if err != nil{
			return err
		}
		if newPassword == userInfo.Password {
		}else{
			if filter.CheckPassword(params.Password) == false{
				return errors.New("update user info, password is error.")
			}
			//计算密码
			params.Password = newPassword
		}
	}
	//检查信息组
	//检查信息组
	err = CheckInfo(params.Info)
	if err != nil{
		return err
	}
	//检查状态
	switch params.Status{
	case USER_STATUS_PUBLIC:
	case USER_STATUS_TRASH:
	case USER_STATUS_STOP:
	default:
		params.Status = USER_STATUS_TRASH
	}
	//检查上级是否存在
	if params.ParentID != ""{
		parentInfo,err := GetUserByID(params.ParentID)
		//上级不存在
		if err != nil{
			return errors.New("The superior does not exist.")
		}
		//不能指定自己为上级
		if parentInfo.ID.Hex() == params.ID.Hex() {
			return errors.New("You cannot specify yourself as the upper level.")
		}
		//不能指定自己的上级为自己
		if parentInfo.ParentID == params.ID.Hex(){
			return errors.New("You cannot specify the superior's superior as yourself.")
		}
	}
	//用户组必须存在
	for _,v := range params.Groups{
		_,err := GetGroupByMark(v.Mark)
		if err != nil{
			return err
		}
	}
	//更新用户数据
	return MgoDBC.UpdateId(params.ID,&params)
}

//更新用户信息
// 指定参数的方式
//param userID string 用户ID
//param username string 用户名
//param password string 密码
//param name string 姓名
//param parent string 上级ID
//param status string 状态
//param info []FieldsUserInfo 信息组
//return error 错误代码
func UpdateDataInParams(userID string,username string,password string,name string,parent string,status string,info []FieldsUserInfo) error {
	//获取用户信息组
	userInfo,err := GetUserByID(userID)
	if err != nil{
		return err
	}
	//检查参数
	// 仅检查用户名、电话、身份证、姓名信息，其他信息请根据情况提前核对
	// 一般其他信息的修改，仅管理员级别能使用，所以可选择信任来源
	if filter.CheckPhone(username) == false || filter.CheckNicename(name) == false{
		return errors.New("Incorrect username, nickname.")
	}
	//用户名不能重叠存在
	// 系统最初不允许用户名重叠，所以未来也不可能出现重叠的用户体系
	// 如果出现重叠的数据，则说明数据库存在异常问题，或中间某些环节直接修改了数据库内容
	// 用户所有修改方案只能通过该方法执行
	searchUsernameInfo,err := GetUserByUsername(username)
	if err == nil && userInfo.ID.Hex() != searchUsernameInfo.ID.Hex(){
		if userInfo.Username == searchUsernameInfo.Username{
			return errors.New("The user name already exists and cannot be modified.")
		}
	}
	userInfo.Username = username
	//如果修改了密码？
	if searchUsernameInfo.Password == password{
	}else{
		if filter.CheckPassword(password) == false{
			return errors.New("update user info, password is error.")
		}
		//计算密码
		userInfo.Password,err = GetPassword(password)
		if err != nil{
			return err
		}
	}
	//给与昵称
	userInfo.Name = name
	//检查上级是否存在
	if parent != ""{
		parentInfo,err := GetUserByID(parent)
		//上级不存在
		if err != nil{
			return errors.New("The superior does not exist.")
		}
		//不能指定自己为上级
		if parentInfo.ID.Hex() == userInfo.ID.Hex() {
			return errors.New("You cannot specify yourself as the upper level.")
		}
		//不能指定自己的上级为自己
		if parentInfo.ParentID == userInfo.ID.Hex(){
			return errors.New("You cannot specify the superior's superior as yourself.")
		}
	}
	userInfo.ParentID = parent
	//检查状态
	switch status{
	case USER_STATUS_PUBLIC:
	case USER_STATUS_TRASH:
	case USER_STATUS_STOP:
	default:
		status = USER_STATUS_TRASH
	}
	userInfo.Status = status
	//检查信息组
	err = CheckInfo(info)
	if err != nil{
		return err
	}
	userInfo.Info = info
	//修改返回
	return MgoDBC.UpdateId(userInfo.ID,&userInfo)
}

