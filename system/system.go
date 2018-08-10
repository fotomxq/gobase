package system

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"fotomxq/gobase/log"
	"fotomxq/gobase/mgotool"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/disk"
	"fotomxq/gobase/distribution"
)

//系统检测模块
//该模块可自动检测整个系统的运行情况

var(
	//Mgo数据库集合句柄
	MgoDBC *mgo.Collection
	//系统运行记录块
	SystemData []FieldsSystemDataType
	//进程启动总时间
	SystemStartTime int64
	//发生错误次数
	LogErrorCount int64
	//发生HTTP访问次数
	LogHttpCount int64
	//发生访问次数
	LogVisitCount int64
	//系统消息次数
	LogSystemCount int64
	//系统安全事件次数
	LogSafetyCount int64
	//总的系统消息
	LogCount int64
	//分布式标记
	DistributedMark string
	//应用标识码
	AppMark string
)

//启动系统记录器
//注意，请并发启动该方法
func Run(){
	//系统运行表
	MgoDBC = mgotool.MgoDB.C("system_run")
	//记录启动时间
	SystemStartTime = time.Now().Unix()
}

//监控程序
func RunAuto(){
	//启动自动循环，每1分钟记录一次
	// 考虑到系统效率和内存占用量，如果每分钟1次，每天将产生1440个记录，30天4.3万个。
	// 数据将自动清理超过30天的
	for{
		//初始化数据集合
		info := FieldsSystemDataType{}
		info.ID = bson.NewObjectId()
		info.AppName = AppMark
		info.DistributedMark = DistributedMark
		info.CreateTime = time.Now().Unix()
		//获取内存数据
		v, err := mem.VirtualMemory()
		if err != nil{
			log.SendError("0.0.0.0","","SystemType.Run",err)
		}else{
			info.MenmoryTotal = v.Total
			info.MenmoryUsed = v.Used
			info.MenmoryUsedPercent = v.UsedPercent
			info.MenmoryFree = v.Total - v.Used
		}
		//获取硬盘数据
		d,err := disk.Usage("/")
		if err != nil{
			log.SendError("0.0.0.0","","SystemType.Run",err)
		}else{
			info.DiskTotal = d.Total
			info.DiskUsed = d.Used
			info.DiskUsedPercent = d.UsedPercent
			info.DiskFree = d.Free
		}
		//写入数据
		err = MgoDBC.Insert(info)
		if err != nil {
			log.SendError("0.0.0.0", "", "SystemType.Run", err)
		}
		//清理超出30天数据
		beforeTime := time.Now().Unix() - 2592000
		_,err = MgoDBC.RemoveAll(bson.M{"CreateTime" : bson.M{"$lt" : beforeTime}})
		if err != nil {
			log.SendError("0.0.0.0", "", "SystemType.Run", err)
		}
		//更新子服务项目 SystemType.Run
		err = distribution.UpdateSubTaskBySelf("SystemType.Run")
		if err != nil{
			log.SendError("0.0.0.0","","SystemType.Run",err)
		}
		//1分钟后继续
		time.Sleep(time.Minute * 1)
	}
}

//监控程序清理工具
func RunAutoClear(){
	for{
		//清理超出30天数据
		beforeTime := time.Now().Unix() - 2592000
		_,err := MgoDBC.RemoveAll(bson.M{"CreateTime" : bson.M{"$lt" : beforeTime}})
		if err != nil {
			log.SendError("0.0.0.0", "", "SystemType.Run", err)
		}
		//更新子服务项目 SystemType.Run
		err = distribution.UpdateSubTaskBySelf("SystemType.Run")
		if err != nil{
			log.SendError("0.0.0.0","","SystemType.Run",err)
		}
		//12小时清理一次
		time.Sleep(time.Hour * 12)
	}
}

//根据日志标识码，递增某部分
//param logType string 日志标识码
func Add(logType string){
	switch logType{
	case log.MessageTypeError:
		LogErrorCount += 1
	case log.MessageTypeHTTP:
		LogHttpCount += 1
	case log.MessageTypeSafety:
		LogSafetyCount += 1
	case log.MessageTypeVisit:
		LogVisitCount += 1
	case log.MessageTypeSystem:
		LogSystemCount += 1
	}
	LogCount += 1
}

//获取数据
//param search string 搜索内容
//param lastTime int64 获取某个时间之后数据，如果小于1则判断为无效
//param page int 页数
//param max int 页长
//return []FieldsSystemDataType 数据集合
//return int 数据总量
//return error 错误代码
func GetList(search string,lastTime int64,page int,max int) ([]FieldsSystemDataType,int,error){
	//如果低于1，则判定为当前时间戳
	if lastTime < 1{
		lastTime = 1
	}
	//修正页数和页码
	if page < 1{
		page = 1
	}
	if max < 1 {
		max = 1
	}
	//组合条件
	q := bson.M{"CreateTime" : bson.M{"$gt" : lastTime}}
	//搜索选项
	if search != ""{
		q["$or"] = []bson.M{
			bson.M{"AppName": bson.M{"$regex": search}},
			bson.M{"DistributedMark": bson.M{"$regex": search}},
		}
	}
	//数据集合
	res := []FieldsSystemDataType{}
	//获取数据
	err := MgoDBC.Find(q).Sort("-CreateTime").Skip((page - 1) * max).Limit(max).All(&res)
	if err != nil{
		return res,0,err
	}
	//获取数据量
	count,err := MgoDBC.Find(q).Count()
	if err != nil{
		return res,0,err
	}
	//返回数据
	return res,count,nil
}