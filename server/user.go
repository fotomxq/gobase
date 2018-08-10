package server

import (
	"time"
	"github.com/pkg/errors"
	"fotomxq/gobase/log"
	"fotomxq/gobase/session"
	"fotomxq/gobase/filter"
	"fotomxq/gobase/user"
	"fotomxq/gobase/authority"
)

//登陆前置处理
//param headerParams HeaderParams
//return error 错误
func LoginBefore(headerParams ServerHeaderParams) (ServerHeaderParams,error) {
	//构建token值
	headerParams.Token = session.GetMark(headerParams.W,headerParams.R,headerParams.IP)
	if headerParams.Token == "" || filter.CheckHexSha1(headerParams.Token) == false{
		return headerParams,errors.New("cannot create token.")
	}
	return headerParams,nil
}

//登陆后处理
//param headerParams HeaderParams
//return error 错误
func LoginAfter(headerParams ServerHeaderParams) (ServerHeaderParams,error) {
	//用户信息
	userInfo := user.FieldsUser{}
	userInfo,err := user.GetUserByToken(headerParams.Token)
	if err != nil{
		return headerParams,errors.New("Token is not exist , but user try to get user info.")
	}

	//确定用户的IP、Token
	if headerParams.IP != userInfo.LoginIP {
		return headerParams,errors.New("User ip not eq token ip.")
	}
	//更新用户时间
	userInfo.LoginTime = time.Now().Unix()
	err = user.UpdateData(userInfo)
	if err != nil{
		return headerParams,err
	}
	//设置用户信息
	headerParams.UserInfo = userInfo

	//检查用户是否具备URL权限？
	if PageAuthorityURLON == true{
		groupMark := authority.PageAuthorityDataURL[headerParams.R.URL.Path]
		if CheckLoginAuthority(headerParams.UserInfo.ID.Hex(),groupMark) == false{
			return headerParams,errors.New("This user cannnot visit url : " + headerParams.R.URL.Path)
		}
	}

	return headerParams,nil
}

//进入用户所属用户组的Index URL
//param headerParams HeaderParams
func GoUserIndex(headerParams ServerHeaderParams) {
	// 注意，改片段取消，因为会造成无法进入系统页面。将改为，默认进入第一个所属用户组的页面。
	//如果用户存在多个用户组，则进入默认的/user/center URL
	//if len(headerParams.UserInfo.Groups) > 1{
	//	//如果已经是center，则不动
	//	GoURL(headerParams,"/user/center")
	//	return
	//}
	//如果没有用户组，则进入退出动作
	if len(headerParams.UserInfo.Groups) < 1{
		log.SendError(headerParams.IP,headerParams.UserInfo.ID.Hex(),"ServerType.GoUserIndex",errors.New("user no have groups."))
		GoURL(headerParams,"/logout")
		return
	}
	//获取用户组列表
	groupsList,err := user.GetGroupList()
	if err != nil{
		log.SendError(headerParams.IP,headerParams.UserInfo.ID.Hex(),"ServerType.GoUserIndex",err)
		GoURL(headerParams,"/logout")
		return
	}
	//最优先的排前面，找到任意一个符合条件的用户组，即可跳出
	groupFirstInfo := user.FieldsUserGroup{}
	isFind := false
	for _,v := range groupsList{
		if v.Status == user.USER_GROUP_STATUS_TRASH || v.Status == user.USER_GROUP_STATUS_STOP{
			continue
		}
		for _,v2 := range headerParams.UserInfo.Groups{
			if v2.ExpireTime != 1 && v2.ExpireTime > time.Now().Unix(){
				continue
			}
			if v.Mark == v2.Mark{
				isFind = true
				groupFirstInfo = v
				break
			}
		}
		if isFind == true{
			break
		}
	}
	//如果没有找到？
	if isFind == false{
		log.SendError(headerParams.IP,headerParams.UserInfo.ID.Hex(),"ServerType.GoUserIndex",errors.New("cannot find user group."))
		GoURL(headerParams,"/logout")
		return
	}
	//如果用户组不存在URL，则进入默认地址
	if groupFirstInfo.IndexURL == ""{
		GoURL(headerParams,"/user/center?return=1")
		return
	}
	//进入用户组设定的URL
	if groupFirstInfo.IndexURL == "/user/center"{
		GoURL(headerParams,"/user/center?return=1")
		return
	}else{
		GoURL(headerParams,groupFirstInfo.IndexURL)
		return
	}
}