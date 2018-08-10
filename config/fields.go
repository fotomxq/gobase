package config

import "gopkg.in/mgo.v2/bson"

//配置类型
type FieldsConfigData struct {
	//ID
	ID bson.ObjectId `bson:"_id"`
	//名称
	Name string `bson:"Name"`
	//值
	Value string `bson:"Value"`
	//默认值
	DefaultValue string `bson:"DefaultValue"`
}
