package cache

import "gopkg.in/mgo.v2/bson"

//数据结构
type CacheMgoData struct {
	//ID
	ID bson.ObjectId `bson:"ID"`
	//标识
	Mark string `bson:"Mark"`
	//过期时间戳
	ExpireTime int64 `bson:"ExpireTime"`
	//数据内容
	Content []byte `bson:"Content"`
}