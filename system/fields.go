package system

import "gopkg.in/mgo.v2/bson"

//系统记录块
type FieldsSystemDataType struct {
	//ID
	ID bson.ObjectId `bson:"_id"`
	//所属应用
	AppName string `bson:"AppName"`
	//所属分布结构标识码
	DistributedMark string `bson:"DistributedMark"`
	//创建时间
	CreateTime int64 `bson:"CreateTime"`
	//系统内存总量
	MenmoryTotal uint64 `bson:"MenmoryTotal"`
	//系统内存占用
	MenmoryUsed uint64 `bson:"MenmoryUsed"`
	//系统内存占用百分比
	MenmoryUsedPercent float64 `bson:"MenmoryUsedPercent"`
	//系统内存剩余
	MenmoryFree uint64 `bson:"MenmoryFree"`
	//应用所在硬盘空间大小
	DiskTotal uint64 `bson:"DiskTotal"`
	//硬盘使用大小
	DiskUsed uint64 `bson:"DiskUsed"`
	//硬盘使用比例
	DiskUsedPercent float64 `bson:"DiskUsedPercent"`
	//硬盘所用空间百分比
	DiskFree uint64 `bson:"DiskFree"`
	//CPU占用率
}