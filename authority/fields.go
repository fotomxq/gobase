package authority

import "gopkg.in/mgo.v2/bson"

//数据库结构
type FieldsPageAuthority struct {
	//ID
	ID bson.ObjectId `bson:"_id"`
	//URL
	URL string `bson:"URL"`
	//标记
	Mark string `bson:"Mark"`
	//名称
	Name string `bson:"Name"`
	//描述
	Des string `bson:"Des"`
}
