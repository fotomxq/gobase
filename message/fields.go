package message

import "gopkg.in/mgo.v2/bson"

//用户消息处理器
type FieldsMessage struct {
	//ID
	ID bson.ObjectId `bson:"_id"`
	//创建时间
	CreateTime int64 `bson:"CreateTime"`
	//发送用户ID
	SendUserID string `bson:"SendUserID"`
	//到用户ID
	ToUserID string `bson:"ToUserID"`
	//类型
	Type string `bson:"Type"`
	//状态值
	Status string `bson:"Status"`
	//消息SHA1
	ContentSha1 string `bson:"ContentSha1"`
	//消息内容
	Content string `bson:"Content"`
}
