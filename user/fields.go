package user

import "gopkg.in/mgo.v2/bson"

//用户主体结构
type FieldsUser struct {
	//ID
	ID bson.ObjectId `bson:"_id"`
	//联系电话 Username
	Username string `bson:"Username"`
	//密码
	Password string `bson:"Password"`
	//姓名
	Name string `bson:"Name"`
	//上一级ID
	ParentID string `bson:"ParentID"`
	//用户状态 public 正常 / trash 删除 / stop 异常禁止访问
	Status string `bson:"Status"`
	//创建时间
	CreateTime int64 `bson:"CreateTime"`
	//创建IP
	CreateIP string `bson:"CreateIP"`
	//用户登陆token
	LoginToken string `bson:"LoginToken"`
	//用户登陆IP
	LoginIP string `bson:"LoginIP"`
	//登陆时间
	LoginTime int64 `bson:"LoginTime"`
	//用户组关联
	Groups []FieldsUserGroupBind `bson:"Groups"`
	//扩展信息组
	Info []FieldsUserInfo `bson:"Info"`
}

//用户组关系结构
type FieldsUserGroupBind struct{
	//用户组Mark
	Mark string `bson:"Mark"`
	//最初订阅时间
	CreateTime int64 `bson:"CreateTime"`
	//订阅到期时间
	ExpireTime int64 `bson:"ExpireTime"`
}

//用户信息组
type FieldsUserInfo struct {
	//名称
	Name string `bson:"Name"`
	//标识码
	Mark string `bson:"Mark"`
	//值
	Value string `bson:"Value"`
}

//用户组结构
type FieldsUserGroup struct{
	//ID
	ID bson.ObjectId `bson:"_id"`
	//标识码
	Mark string `bson:"Mark"`
	//状态 public 正常 / trash 删除 / stop 异常禁止访问 / private 只有管理员可以授权
	Status string `bson:"Status"`
	//权限Marks组
	Authority []string `bson:"Authority"`
	//名称
	Name string `bson:"Name"`
	//描述
	Des string `bson:"Des"`
	//登陆后进入的URL
	IndexURL string `bson:"IndexURL"`
	//优先级 index url所依据的排序
	Level int `bson:"Level"`
}
