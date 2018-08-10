package user

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"fotomxq/gobase/filter"
	"fotomxq/gobase/log"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//登录模块
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

//登陆用户操作
//param username string 用户名
//param password string 密码
//param token string token
//param ip string IP地址
//return bool 是否成功
func Login(username string,password string,token string,ip string) bool {
	//检查参数，必须确保token是一个sha1密匙
	if filter.CheckUsername(username) == false || filter.CheckPassword(password) == false || filter.CheckHexSha1(token) == false{
		return false
	}
	//普通用户登陆
	//计算密码值
	passwordSha1,err := GetPassword(password)
	if err != nil{
		log.SendError(ip,username,"user.login",err)
		return false
	}
	if passwordSha1 == ""{
		return false
	}
	//搜索数据库是否存在该数据，用户名、密码、状态
	loginUserInfo := FieldsUser{}
	err = MgoDBC.Find(bson.M{"Username" : username, "Password" : passwordSha1, "Status" : USER_STATUS_PUBLIC}).One(&loginUserInfo)
	//如果不存在数据，则返回
	if err != nil{
		return false
	}
	//如果不存在数据，则返回
	if loginUserInfo.ID.Hex() == ""{
		return false
	}
	//确定可以登陆后，将更新数据库中所有token和该token相同的ID为空，避免撞库风险
	// 其他维护工具，可根据此逻辑计算出出现撞库的情况并做风险性记录
	allNeedUpdateRes := []FieldsUser{}
	err = MgoDBC.Find(bson.M{"LoginToken" : token}).All(&allNeedUpdateRes)
	if err == nil{
		for _,v := range allNeedUpdateRes{
			if v.Username == loginUserInfo.Username{
				//出现这种情况是因为当前登陆的客户端重复登陆了，并不属于撞库问题
				continue
			}
			v.LoginToken = ""
			err = MgoDBC.UpdateId(v.ID,&v)
			if err != nil{
				return false
			}
			//记录日志
			log.SendText(ip,loginUserInfo.Username,"user.login()",log.MessageTypeSafety,"Found that there is a problem with the token crash, hit the user:" + v.Username + "，token：" + token)
		}
	}
	//更新token、ip、登陆记录等内容
	loginUserInfo.LoginTime = time.Now().Unix()
	loginUserInfo.LoginToken = token
	loginUserInfo.LoginIP = ip
	//更新数据
	err = MgoDBC.UpdateId(loginUserInfo.ID,&loginUserInfo)
	return err == nil
}

//退出登陆
// 如果IP为空，则跳过该条件
//param token string token
//param ip string IP地址
//return bool 是否成功
func LogoutByToken(token string,ip string) bool {
	//检查是否已经登陆？
	if CheckLogin(token,ip) == false{
		return true
	}
	//获取数据
	p := bson.M{}
	if ip == ""{
		p = bson.M{"LoginToken" : token}
	}else{
		p = bson.M{"LoginToken" : token, "LoginIP" : ip}
	}
	res := FieldsUser{}
	err := MgoDBC.Find(p).One(&res)
	if err != nil{
		return false
	}
	//修改值，确保退出
	res.LoginToken = ""
	return MgoDBC.UpdateId(res.ID, &res) == nil
}


//检查用户是否已经登陆
// 如果IP为空，则跳过IP审查
//param token string TOKEN
//param ip string IP地址
//return bool 是否登陆
func CheckLogin(token string,ip string) bool {
	//如果为空，则返回
	if token == ""{
		return false
	}
	//检查参数
	if filter.CheckHexSha1(token) == false {
		return false
	}
	//组织查询条件
	p := bson.M{}
	if ip == ""{
		p = bson.M{"LoginToken" : token}
	}else{
		if BindCookieAndIP == true{
			p = bson.M{"LoginToken" : token, "LoginIP" : ip}
		}
	}
	//获取用户信息
	loginUserInfo := FieldsUser{}
	err := MgoDBC.Find(p).One(&loginUserInfo)
	if err != nil{
		return false
	}
	//更新最后登陆时间
	loginUserInfo.LoginTime = time.Now().Unix()
	loginUserInfo.LoginIP = ip
	err = MgoDBC.UpdateId(loginUserInfo.ID,&loginUserInfo)
	//返回
	return err == nil
}

//清理所有超时用户
func ClearLogin(){
	expTime := time.Now().Unix() - UserLoginTimeout
	_,_ = MgoDBC.UpdateAll(bson.M{"LoginTime" : bson.M{"$lt" : expTime}},bson.M{"$set" : bson.M{"LoginToken" : ""}})
}