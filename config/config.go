package config

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fotomxq/gobase/mgotool"
)

//全局配置处理器

var(
	//数据库句柄
	DBC *mgo.Collection
)

//初始化
func Run(){
	//配置表
	DBC = mgotool.MgoDB.C("configs")
}

//获取数据表所有值
//return []Config 配置列表
//return error 错误信息
func GetAll() (map[string]FieldsConfigData,error){
	//获取全部配置
	var res []FieldsConfigData
	result := map[string]FieldsConfigData{}
	err := DBC.Find(nil).All(&res)
	if err != nil{
		return result,err
	}
	//重组所有配置信息
	for _,v := range res{
		result[v.Name] = v
	}
	return result,err
}

//获取配置
//param name string 名称
//return string 值
//return error
func Get(name string) (string,error) {
	var res FieldsConfigData
	err := DBC.Find(bson.M{"Name" : name}).One(&res)
	if err != nil{
		return "",err
	}
	return res.Value,nil
}

//设定配置
// 如果不存在则创建新的配置
//param name string 名称
//param value string 值
//return error 成功则返回空
func Set(name string,value string) error{
	var res FieldsConfigData
	err := DBC.Find(bson.M{"Name" : name}).One(&res)
	if err != nil{
		res = FieldsConfigData{
			bson.NewObjectId(),
			name,
			value,
			value,
		}
		return DBC.Insert(res)
	}else{
		res.Value = value
		return DBC.UpdateId(res.ID,&res)
	}
}
