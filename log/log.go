package log

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fotomxq/gobase/file"
	"fotomxq/gobase/filter"
	"fotomxq/gobase/mgotool"
)


var(
	//信息类型
	MessageTypeError string = "ERROR" //错误
	MessageTypeSystem string = "SYSTEM" //系统消息
	MessageTypeHTTP string = "HTTP" //HTTP消息
	MessageTypeClient string = "CLIENT" //客户端消息
	MessageTypeMessage string = "MESSAGE" //普通消息
	MessageTypeSafety string = "SAFETY" //安全事件
	MessageTypeVisit string = "VISIT" //访问记录

	//列队进程锁定，避免同时读写
	LogListCH chan int = make(chan int,1)

	//Mgo数据库集合句柄
	MgoDBC *mgo.Collection

	//是否向控制台发送日志
	//true
	LogConsole bool

	//日志列队
	// = []LogData{}
	MessageList []LogData

	//是否启动日志存储
	//= false
	LogSaveON bool

	//存储模式
	//mgo Mgo数据库 default
	//file 文件存储
	LogSaveMode string

	//日志存储路径
	// = "." + file.Sep + "log"
	LogDir string

	//从mgo导出的日志目录
	// = "." + file.Sep + "log"
	LogMgoToFolder string
)

//日志数据结构
type LogData struct {
	//创建时间戳
	CreateTime int64
	//创建IP
	CreateIP string
	//创建用户名
	CreateUsername string
	//来源
	FromCode string
	//类型
	MessageType string
	//内容
	Message string
}

//发送一个标准化日志
//param ipaddr string IP地址
//param createUsername string 创建用户
//param from string 来源
//param t int 类型
//param m string 消息内容
func SendText(ipaddr string,createUsername,from string,t string,m string){
	//将超长日志截取前半部分，剔除后半部分
	// 截取1000个字符以上的数据
	if len(m) > 1000{
		m = filter.SubStr(m,0,1000)
	}
	//组合信息
	nowTime := time.Now()
	d := LogData{
		nowTime.Unix(),
		ipaddr,
		createUsername,
		from,
		t,
		m,
	}
	LogListCH <- 1
	MessageList = append(MessageList, d)
	<- LogListCH
	m = nowTime.Format("2006-01-02 15:04:05.999999999") + " [" + t + "] " + " [" + from + "] " + " [" + ipaddr + "] " + " [" + createUsername + "] " + m
	if LogConsole == true{
		fmt.Println(m)
	}
}

//发送一个错误
//param ipaddr string IP地址
//param createUsername string 创建用户名
//param from string 来源
//param e error 错误
func SendError(ipaddr string,createUsername string,from string,e error){
	SendText(ipaddr,createUsername,from,MessageTypeError,e.Error())
}

//发送一个标准化网络日志
//param r *http.Request
//param ipaddr string IP地址
//param createUsername string 创建用户名
//param from string 来源
//param t string 类型
//param m string 消息内容
func SendHttp(r *http.Request,ipaddr string,createUsername string,from string,t string,m string){
	m = "[URL:" + r.URL.Path + "] " + m
	SendText(ipaddr,createUsername,from,t,m)
}

//构建标准化日志
func GetMsg(v LogData) string{
	var m string
	m += strconv.FormatInt(v.CreateTime,10)
	if v.FromCode != ""{
		m += " [" + v.FromCode + "]"
	}
	m += " [" + v.MessageType + "] " + v.Message
	return m
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//自动存储模块
// 该模块用于快速增加消息列队，并以400毫秒速度写入日志文件
// 日志数据将写入指定的mgo数据表内
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//启动自动保存日志服务
// 以一个新的子进程，自动保存日志列队信息
func Run(){
	//日志服务
	MgoDBC = mgotool.MgoDB.C("log")
	//定时任务
	ticker := time.NewTicker(time.Millisecond  * 400)
	for _ = range ticker.C {
		if LogSaveON == true{
			switch LogSaveMode {
			case "mgo":
				saveLogMgo()
			case "file":
				saveLogFile()
			}
		}
	}
}

//日志存储独立函数
// 用于定时存储日志列队
func saveLogMgo(){
	LogListCH <- 1
	defer func(){
		<- LogListCH
	}()
	if len(MessageList) < 1{
		return
	}
	for _,v := range MessageList{
		vf := FieldsLogMgoData{
			bson.NewObjectId(),
			v.CreateTime,
			v.CreateIP,
			v.CreateUsername,
			v.FromCode,
			v.MessageType,
			v.Message,
		}
		err := MgoDBC.Insert(&vf)
		if err != nil{
			continue
		}
	}
	MessageList = []LogData{}
	return
}

//日志存储独立函数
// 用于定时存储日志列队
func saveLogFile(){
	LogListCH <- 1
	defer func(){
		<- LogListCH
	}()
	if len(MessageList) < 1{
		return
	}
	var m string = ""
	for _,v := range MessageList{
		m += "\n" + GetMsg(v)
	}
	dir := LogDir + file.Sep + time.Now().Format("200601"+file.Sep+"02")
	err := os.MkdirAll(dir, 666)
	if err != nil{
		return
	}
	src := dir + file.Sep + time.Now().Format("2006010215") + ".log"
	f, err := os.OpenFile(src, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer f.Close()
	mb := []byte(m)
	_, err = f.Write(mb)
	if err != nil{
		return
	}
	MessageList = []LogData{}
	return
}

//查询模块

//查找内容
// 参数可任意选择
// 返回日志标准格式
//param search string 查询内容
//param messageType []string 消息性质
//param fromCode string 来源限定
//param findMinTime int64
//param findMaxTime int64
//param page int
//param max int
//return []LogT
//return int 数量
//return bool
func Find(r *http.Request,search string,messageType []string,fromCode string,findMinTime int64,findMaxTime int64) ([]FieldsLogMgoData,int,error){
	var res []FieldsLogMgoData
	q := bson.M{}
	if search != ""{
		q["$or"] = []bson.M{
			bson.M{"CreateIP": bson.M{"$regex": search}},
			bson.M{"CreateUsername": bson.M{"$regex": search}},
			bson.M{"Message": bson.M{"$regex": search}},
		}
	}
	if fromCode != ""{
		q["FromCode"] = fromCode
	}
	if len(messageType) > 0{
		q["MessageType"] = bson.M{"$in" : messageType}
	}
	if findMinTime > 0{
		q["CreateTime"] = bson.M{"$gte":findMinTime}
	}
	if findMaxTime > 0{
		q["CreateTime"] = bson.M{"$lte":findMaxTime}
	}
	err := mgotool.GetList(r,MgoDBC,q).All(&res)
	if err != nil{
		return nil,0,err
	}
	count,err := MgoDBC.Find(q).Count()
	if err != nil{
		return nil,0,err
	}
	return res,count,nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//将日志从mgo数据库导出到指定文件夹内
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//保存日志数据
//return error
func SaveLogFile() error {
	//构建新的目录路径
	dir := LogDir + file.Sep + time.Now().Format("20060102")
	err := file.CreateFolder(dir)
	if err != nil{
		return err
	}

	//构建条件
	getTime := time.Now().Unix() - 10
	p := bson.M{"CreateTime" : bson.M{"$lt" : getTime}}

	//读取当前Unix时间戳10秒之前的数据
	res := []interface{}{}
	err = MgoDBC.Find(p).All(&res)
	if err != nil{
		return err
	}

	//输出到JSON文件
	fileSrc := dir + file.Sep + time.Now().Format("20060102_150405.log")
	resJSON,err := json.Marshal(res)
	if err != nil{
		return err
	}
	err = file.WriteFile(fileSrc,resJSON)
	if err != nil{
		return err
	}

	//删除旧的日志数据
	_,err = MgoDBC.RemoveAll(p)
	if err != nil{
		return err
	}

	return nil
}

//获取日志数据
//param page int
//param max int
//return []logcore.LogT
//return error
func GetLog(page int,max int) ([]FieldsLogMgoData,error){
	res := []FieldsLogMgoData{}
	err := MgoDBC.Find(nil).Sort("-CreateTime").Skip((page-1) * max).Limit(max).All(&res)
	return res,err
}