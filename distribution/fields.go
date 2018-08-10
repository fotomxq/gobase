package distribution

import "gopkg.in/mgo.v2/bson"

//服务注册表
type FieldsDistributionServiceType struct {
	//ID
	ID bson.ObjectId `bson:"_id"`
	//应用标识码
	AppMark string `bson:"AppMark"`
	//分布式序列标识码
	Sequence string `bson:"Sequence"`
	//IP
	IP string `bson:"IP"`
	//端口
	Port string `bson:"Port"`
	//注册时间
	CreateTime int64 `bson:"CreateTime"`
	//更新状态时间
	UpdateTime int64 `bson:"UpdateTime"`
	//子任务序列
	SubTasks []FieldsDistributionServiceSubTasksType `bson:"SubTasks"`
}

//服务注册，子任务注册序列
type FieldsDistributionServiceSubTasksType struct {
	//任务标识码
	Mark string `bson:"Mark"`
	//注册时间
	CreateTime int64 `bson:"CreateTime"`
	//更新状态时间
	UpdateTime int64 `bson:"UpdateTime"`
}
