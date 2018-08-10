package temporarydata

import (
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fotomxq/gobase/mgotool"
	"fotomxq/gobase/distribution"
	"fotomxq/gobase/log"
)

//临时数据读写模块
// 该数据提交后自动存储24小时，如果在此时间内不读取将销毁
var(
	//数据库句柄
	MgoDBC *mgo.Collection

	//过期时间
	// 86400
	ExpireTime int64
)

//初始化
func Run() {
	//临时存储工具
	MgoDBC = mgotool.MgoDB.C("temporary_data")
}

//维护程序
func RunAuto(){
	//删除所有过期数据
	_,_ = MgoDBC.RemoveAll(bson.M{"ExpireTime" : bson.M{"$lt" : time.Now().Unix()}})
	//更新子服务项目 TemporaryDataType.Auto
	err := distribution.UpdateSubTaskBySelf("TemporaryDataType.Auto")
	if err != nil{
		log.SendError("0.0.0.0","","TemporaryDataType.Auto",err)
	}
}

//创建一个数据节点
//param data string 数据内容
//return string 读取ID
func Create(data string) string{
	newD := FieldsTemporaryDataType{
		bson.NewObjectId(),
		data,
		time.Now().Unix() + ExpireTime,
	}
	err := MgoDBC.Insert(&newD)
	if err != nil{
		return ""
	}
	return newD.ID.Hex()
}

//读取数据内容
//param id string ID
//return string 数据值
func Get(id string) string{
	d := FieldsTemporaryDataType{}
	err := MgoDBC.FindId(bson.ObjectIdHex(id)).One(&d)
	if err != nil{
		return ""
	}
	err = MgoDBC.RemoveId(d.ID)
	if err != nil{
	}
	return d.Data
}
