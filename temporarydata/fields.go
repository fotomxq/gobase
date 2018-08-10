package temporarydata

import "gopkg.in/mgo.v2/bson"

//临时数据读写模块结构体
type FieldsTemporaryDataType struct {
	//ID
	ID bson.ObjectId `bson:"_id"`
	//变量数据
	Data string `bson:"Data"`
	//过期时间
	ExpireTime int64 `bson:"ExpireTime"`
}

