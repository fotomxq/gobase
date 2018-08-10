package pedometer

import (
	"time"
	"fotomxq/gobase/log"
	"fotomxq/gobase/distribution"
)

//计步器模块
// 指定一个唯一的标识符，记录其次数
// 用于IP超过N次后拉黑的功能实现
// 自动清理内存，避免内存占用超出系统承受范围
var(
	//存储器
	// eg : data['192.168.1.1'] = {1,8217832}
	Data map[string]PedometerData

	//上限值
	// 超过10次封顶，不再自动增加
	LimitMax int

	//全局超时时间 秒
	// 默认2天
	LimitTime int64
)

//计步器结构体
type PedometerData struct {
	//计步器
	T int
	//创建时间
	UpdateTime int64
}

//自动清理服务
func Run() {
	//初始化
	Data = map[string]PedometerData{}
	LimitMax = 10
	LimitTime = 172800
	//每10分钟执行一次清理服务
	for{
		nowTime := time.Now().Unix()
		for k,v := range Data{
			if v.UpdateTime + LimitTime > nowTime{
				delete(Data,k)
			}
		}
		//更新子服务项目 PedometerType.Run
		err := distribution.UpdateSubTaskBySelf("PedometerType.Run")
		if err != nil{
			log.SendError("0.0.0.0","","PedometerType.Run",err)
		}
		time.Sleep(time.Minute * 10)
	}
}

//获取标记次数
//param mark string 标记值
//return int 计步器次数
func Get(mark string) int {
	value,ok := Data[mark]
	if ok == true{
		return value.T
	}
	return 0
}

//增加计步器
//param mark string 标记值
func Add(mark string){
	_,ok := Data[mark]
	var newT int = 1
	if ok == true{
		if newT > LimitMax{
			newT = LimitMax
		}else{
			newT = Data[mark].T + 1
		}
	}
	Data[mark] = PedometerData{
		newT,
		time.Now().Unix(),
	}
}

//减少计步器
//param mark string 标记值
func Reduce(mark string){
	_,ok := Data[mark]
	if ok == true{
		Data[mark] = PedometerData{
			Data[mark].T - 1,
			time.Now().Unix(),
		}
		return
	}
}

//清空该标记的记录
//param mark string 标记值
func Clear(mark string){
	_,ok := Data[mark]
	if ok == false{
		return
	}
	delete(Data,mark)
}