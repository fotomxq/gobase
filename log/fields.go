package log

import "gopkg.in/mgo.v2/bson"

//日志数据结构
type FieldsLogMgoData struct {
	//ID
	ID bson.ObjectId `bson:"_id"`
	//创建时间戳
	CreateTime int64 `bson:"CreateTime"`
	//创建IP
	CreateIP string `bson:"CreateIP"`
	//创建用户名
	CreateUsername string `bson:"CreateUsername"`
	//来源
	FromCode string `bson:"FromCode"`
	//类型
	MessageType string `bson:"MessageType"`
	//内容
	Message string `bson:"Message"`
}