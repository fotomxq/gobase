package mgotool

import (
	"gopkg.in/mgo.v2"
	"github.com/pkg/errors"
)

//填装构建mgo句柄

var(
	MgoDB *mgo.Database
)

//创建数据库链接
//param url string 链接URL
//param dbName string 数据库name
//return error 错误代码
func CreateDB(url string,dbName string) error{
	//连接到数据库
	mgoSession,err := mgo.Dial(url)
	if err != nil{
		return errors.New("cannot connection mongodb , error : " + err.Error())
	}
	//defer mgoSession.Close()
	MgoDB = mgoSession.DB(dbName)
	return nil
}