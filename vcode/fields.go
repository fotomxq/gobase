package vcode

import "gopkg.in/mgo.v2/bson"

//验证码记录
type FieldsVerificationCodeDataType struct {
	//ID
	ID bson.ObjectId `bson:"_id"`
	//Token
	Token string `bson:"Token"`
	//验证码Key/Value
	Value string `bson:"Value"`
	//过期时间
	ExpireTime int64 `bson:"ExpireTime"`
}
