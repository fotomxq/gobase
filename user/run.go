package user

import "fotomxq/gobase/mgotool"

//初始化
func Run(){
	//用户组
	GroupMgoDBC = mgotool.MgoDB.C("group")
	//用户
	MgoDBC = mgotool.MgoDB.C("user")
}
