package distribution

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"github.com/pkg/errors"
	"fotomxq/gobase/mgotool"
)

var(
	//数据库句柄
	MgoDBC *mgo.Collection
	//本服务的ID
	ServerID string
)

//初始化
func Run(){
	MgoDBC = mgotool.MgoDB.C("distribution_service")
}

//获取服务列表

//注册一个服务
//param ip string
//param port string
//param appMark string
//param sequence string
//param subTask []FieldsDistributionServiceSubTasksType
//return error
func Create(ip string,port string,appMark string,sequence string,subTask []FieldsDistributionServiceSubTasksType) error {
	//查询服务是否存在？
	query := bson.M{"AppMark" : appMark, "Sequence" : sequence}
	res := FieldsDistributionServiceType{}
	err := MgoDBC.Find(query).One(&res)
	if err == nil{
		return errors.New("Server is exist.")
	}
	//不存在则创建
	nowTime := time.Now().Unix()
	res = FieldsDistributionServiceType{
		bson.NewObjectId(),
		appMark,
		sequence,
		ip,
		port,
		nowTime,
		nowTime,
		subTask,
	}
	ServerID = res.ID.Hex()
	return MgoDBC.Insert(&res)
}

//注销一个服务
//param id string ID
//return error
func Delete(id string) error {
	return MgoDBC.RemoveId(bson.ObjectIdHex(id))
}

//根据基本设定注销一个服务
//param appMark string
//param sequence string
//return error 错误代码
func DeleteByMark(appMark string,sequence string) error{
	query := bson.M{"AppMark" : appMark, "Sequence" : sequence}
	res := FieldsDistributionServiceType{}
	err := MgoDBC.Find(query).One(&res)
	if err != nil{
		return nil
	}
	return MgoDBC.RemoveId(res.ID)
}

//更新当前服务状态
//return error 错误代码
func UpdateServerBySelf() error {
	return UpdateServer(ServerID)
}

//更新当前服务的某个子任务
//param mark string 子任务标识码
//return error 错误代码
func UpdateSubTaskBySelf(mark string) error{
	return UpdateSubTask(ServerID,mark)
}

//更新服务状态
//param id string ID
//return error 错误代码
func UpdateServer(id string) error{
	res,err := GetID(id)
	if err != nil{
		return err
	}
	res.UpdateTime = time.Now().Unix()
	err = MgoDBC.UpdateId(res.ID,&res)
	return err
}

//更新子任务状态
//param id string
//param mark string 子任务标识码
//return error 错误代码
func UpdateSubTask(id string,mark string) error{
	res,err := GetID(id)
	if err != nil{
		return err
	}
	nowTime := time.Now().Unix()
	isOK := false
	for k,v := range res.SubTasks{
		if v.Mark == mark{
			res.SubTasks[k].UpdateTime = nowTime
			isOK = true
			break
		}
	}
	if isOK == false{
		res.SubTasks = append(res.SubTasks,FieldsDistributionServiceSubTasksType{
			mark,
			nowTime,
			nowTime,
		})
	}
	res.UpdateTime = nowTime
	err = MgoDBC.UpdateId(res.ID,&res)
	return err
}

//获取服务信息
//param id string
//return FieldsDistributionServiceType 数据集合
//return error 错误代码
func GetID(id string) (FieldsDistributionServiceType,error){
	res := FieldsDistributionServiceType{}
	if bson.IsObjectIdHex(id) == false{
		return res,errors.New("ID is error.")
	}
	err := MgoDBC.FindId(bson.ObjectIdHex(id)).One(&res)
	return res,err
}

//获取服务列表
//return []FieldsDistributionServiceType 数据集合
//return int 数据量
//return error 错误代码
func GetList() ([]FieldsDistributionServiceType,int,error){
	res := []FieldsDistributionServiceType{}
	err := MgoDBC.Find(nil).All(&res)
	if err != nil{
		return res,0,err
	}
	count,err := MgoDBC.Find(nil).Count()
	if err != nil{
		return res,0,err
	}
	return res,count,nil
}