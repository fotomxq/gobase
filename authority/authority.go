package authority

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"errors"
	"fotomxq/gobase/mgotool"
)

//页面权限组
// 该模块可识别URL对应的权限
// 1\加载URL配置到内存中
// 2\自动判断用户是否可访问URL，如果不能访问则退出

var(
	//数据库操作
	MgoDBC *mgo.Collection

	//内存数据结构
	// URL -> Mark
	PageAuthorityDataURL map[string]string
	//内存数据结构
	// Mark -> URL
	PageAuthorityDataMark map[string]string
)

//初始化装载数据
//return error 错误代码
func Run() {
	//页面权限组
	MgoDBC = mgotool.MgoDB.C("page_authority")
}

//检查权限组是否存在
//param authority []string 需要检查的权限组
//return error 发现权限不存在则返回错误信息
func Check(authority []string) error {
	for _,v2 := range authority{
		isFind := false
		for _,v := range PageAuthorityDataURL{
			if v2 == v{
				isFind = true
			}
		}
		if isFind == false{
			return errors.New("cannot find anathority.")
		}
	}
	return nil
}

//将数据装载到内存
//return error 错误代码
func GetAllIn() error{
	//从数据库装载数据到内存
	PageAuthorityDataURL = map[string]string{}
	PageAuthorityDataMark = map[string]string{}
	pageAuthorityData := []FieldsPageAuthority{}
	err := MgoDBC.Find(nil).All(&pageAuthorityData)
	if err != nil{
		return err
	}
	for _,v := range pageAuthorityData{
		PageAuthorityDataURL[v.URL] = v.Mark
		PageAuthorityDataMark[v.Mark] = v.URL
	}
	return nil
}

//获取所有权限
//return []FieldsPageAuthority 权限组
//return error 错误信息
func GetAll() ([]FieldsPageAuthority,error){
	res := []FieldsPageAuthority{}
	err := MgoDBC.Find(nil).All(&res)
	return res,err
}

//根据URL获取数据组
//param url string URL地址
//return FieldsPageAuthority 权限信息
//return error 错误信息
func GetByURL(url string) (FieldsPageAuthority,error){
	res := FieldsPageAuthority{}
	err := MgoDBC.Find(bson.M{"URL" : url}).One(&res)
	return res,err
}

//设置URL
//param url string URL
//param mark string 标识码
//param name string 名称
//param des string 描述
//return error 错误
func Set(url string,mark string,name string,des string) error {
	//查询是否存在
	res,err := GetByURL(url)
	if err != nil{
		//不存在，则创建
		res = FieldsPageAuthority{
			bson.NewObjectId(),
			url,
			mark,
			name,
			des,
		}
		return MgoDBC.Insert(&res)
	}
	//存在，直接修改
	res.URL = url
	res.Mark = mark
	res.Name = name
	res.Des = des
	err = MgoDBC.Update(bson.M{"URL" : url},&res)
	if err != nil{
		return err
	}
	//PageAuthorityDataMark[res.Mark]= res.URL
	//PageAuthorityDataURL[res.URL] = res.Mark
	return nil
}

//删除URL
//param url string URL
//return error 错误
func Delete(url string) error {
	err := MgoDBC.Remove(bson.M{"URL" : url})
	if err != nil{
		return err
	}
	return GetAllIn()
}