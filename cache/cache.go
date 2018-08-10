package cache

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"fotomxq/gobase/log"
	"fotomxq/gobase/distribution"
)

//缓冲套件

var(
	//缓冲器数据过期时间
	// 单位 秒，默认一天
	LimitTime int64

	//存储类型
	// go - go语言内部存储
	// mgo - 数据库存储方式
	SaveMode string

	//缓冲数据
	GoDatas map[string]CacheGoData

	//数据表句柄
	MgoDBC *mgo.Collection
)

//缓冲器核心部分
// 同一个应用可选择多种存储模式，所有方法需要调用对应模式下的方法才能实现
// 注意，模式均可能需要提前给定参数后才能有效执行，否则报错

//该文件指向一些通用的模块，不实现具体的方法

//自动清理过期内容
// 30分钟执行一次
func Run(){
	LimitTime = 86400
	GoDatas = map[string]CacheGoData{}
	switch SaveMode{
	case "go":
		for{
			nowTime := time.Now().Unix()
			for k,v := range GoDatas{
				if v.ExpireTime > nowTime{
					delete(GoDatas,k)
				}
			}
			//更新子服务项目 CacheType.Run
			err := distribution.UpdateSubTaskBySelf("CacheType.Run")
			if err != nil{
				log.SendError("0.0.0.0","","CacheType.Run",err)
			}
			time.Sleep(time.Minute * 30)
		}
	case "mgo":
		//定时删除过期缓冲数据
		// 30分钟执行一次
		for {
			_,_ = MgoDBC.RemoveAll(bson.M{"ExpireTime" : bson.M{"$gt" : time.Now().Unix()}})
			//更新子服务项目 CacheType.Run
			err := distribution.UpdateSubTaskBySelf("CacheType.Run")
			if err != nil{
				log.SendError("0.0.0.0","","CacheType.Run",err)
			}
			time.Sleep(time.Minute * 30)
		}
	}
}

//获取缓冲
//param mark string 标识码
//return []byte 缓冲内容
//return bool 是否存在数据
func Get(mark string) ([]byte,bool){
	res,b := GoDatas[mark]
	if b == true{
		if res.ExpireTime > time.Now().Unix(){
			return nil,false
		}
	}
	return res.Content,b
}

//设置缓冲
//param mark string 标识
//param content []byte 内容
func Set(mark string,content []byte) {
	GoDatas[mark] = CacheGoData{
		mark,
		time.Now().Unix() + LimitTime,
		content,
	}
}


//获取缓冲数据
//param mark string 标识
//return []byte 数据内容
//return error 错误
func MgoGet(mark string) ([]byte,error){
	var res CacheMgoData
	err := MgoDBC.Find(bson.M{"Mark" : mark}).One(&res)
	if err != nil{
		return nil,err
	}
	if res.ExpireTime > time.Now().Unix(){
		err = MgoDBC.RemoveId(res.ID)
		return nil,err
	}
	return res.Content,nil
}

//设置缓冲内容
//param mark string
//param content []byte
//return error
func MgoSet(mark string,content []byte) error{
	var res CacheMgoData
	err := MgoDBC.Find(bson.M{"Mark" : mark}).One(&res)
	if err != nil{
		res = CacheMgoData{
			bson.NewObjectId(),
			mark,
			time.Now().Unix() + LimitTime,
			content,
		}
		err = MgoDBC.Insert(&res)
		return err
	}else{
		res.Content = content
		res.ExpireTime = time.Now().Unix() + LimitTime
		err = MgoDBC.UpdateId(res.ID,res)
		return err
	}
	return nil
}