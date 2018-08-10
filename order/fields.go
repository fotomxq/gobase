package order

import "gopkg.in/mgo.v2/bson"

//订单类
type FieldsOrder struct {
	//ID
	ID bson.ObjectId `bson:"_id"`
	//创建时间
	CreateTime int64 `bson:"CreateTime"`
	//创建时间
	CreateIP string `bson:"CreateIP"`
	//创建Token
	CreateToken string `bson:"CreateToken"`
	//状态
	Status string `bson:"Status"`
	//服务ID
	ServiceID string `bson:"ServiceID"`
	//用户ID
	ServiceUserID string `bson:"ServiceUserID"`
	//服务单位
	ServiceUnit string `bson:"ServiceUnit"`
	//订单金额
	OrderCost string `bson:"OrderCost"`
	//支付代码
	OrderCode string `bson:"OrderCode"`
}
