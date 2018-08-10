package message

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"errors"
	"fotomxq/gobase/encrypt"
	"fotomxq/gobase/mgotool"
)

var(
	//消息类型
	// 系统消息
	MESSAGE_TYPE_SYSTEM string = "system"
	// 公开通知
	MESSAGE_TYPE_PUBLIC string = "public"
	// 私密消息
	MESSAGE_TYPE_PRIVATE string = "private"

	//消息状态
	// 正常
	MESSAGE_STATUS_PUBLIC string = "public"
	// 删除
	MESSAGE_STATUS_TRASH string = "trash"
	// 未读
	MESSAGE_STATUS_UNREAD string = "unread"

	//Mgo数据库集合句柄
	MgoDBC *mgo.Collection
)

//初始化
func Run(){
	//消息模块
	MgoDBC = mgotool.MgoDB.C("message")
}

//获取此人发送的消息
//param id string
//return FieldsMessage 消息组
//return error 错误信息
func Get(id string) (FieldsMessage,error){
	res := FieldsMessage{}
	if bson.IsObjectIdHex(id) == false{
		return res,errors.New("message id is error.")
	}
	err := MgoDBC.FindId(bson.ObjectIdHex(id)).One(&res)
	return res,err
}

//获取此人发送的消息列队
//param userID string 用户ID
//param messageType string 消息类型
//param status string 状态
//param page int 页数
//param max int 页长
//return []FieldsMessage 数据组
//return error 错误代码
func GetSendList(userID string,messageType string,status string,page int,max int) ([]FieldsMessage,error){
	res := []FieldsMessage{}
	if bson.IsObjectIdHex(userID) == false{
		return res,errors.New("user id is error.")
	}
	if messageType != MESSAGE_TYPE_SYSTEM && messageType != MESSAGE_TYPE_PUBLIC && messageType != MESSAGE_TYPE_PRIVATE{
		return res,errors.New("message type error.")
	}
	if status != MESSAGE_STATUS_PUBLIC && status != MESSAGE_STATUS_TRASH && status != MESSAGE_STATUS_UNREAD{
		return res,errors.New("message status error.")
	}
	err := MgoDBC.Find(bson.M{"SendUserID" : userID,"Type" : messageType,"Status" : status}).Sort("-_CreateTime").Skip((page-1) * max).Limit(max).All(&res)
	return res,err
}

//获取收件人收到的消息列队
//param userID string 用户ID
//param messageType string 消息类型
//param status string 状态
//param page int 页数
//param max int 页长
//return []FieldsMessage 数据组
//return error 错误代码
func GetToList(userID string,messageType string,status string,page int,max int) ([]FieldsMessage,error){
	res := []FieldsMessage{}
	if bson.IsObjectIdHex(userID) == false{
		return res,errors.New("user id is error.")
	}
	if messageType != MESSAGE_TYPE_SYSTEM && messageType != MESSAGE_TYPE_PUBLIC && messageType != MESSAGE_TYPE_PRIVATE{
		return res,errors.New("message type error.")
	}
	if status != MESSAGE_STATUS_PUBLIC && status != MESSAGE_STATUS_TRASH && status != MESSAGE_STATUS_UNREAD{
		return res,errors.New("message status error.")
	}
	err := MgoDBC.Find(bson.M{"ToUserID" : userID,"Type" : messageType,"Status" : status}).Sort("-_CreateTime").Skip((page-1) * max).Limit(max).All(&res)
	return res,err
}

//获取系统通知
//param status string 状态
//param page int 页数
//param max int 页长
//return []FieldsMessage 数据组
//return error 错误代码
func GetSystemList(status string,page int,max int) ([]FieldsMessage,error){
	res := []FieldsMessage{}
	if status != MESSAGE_STATUS_PUBLIC && status != MESSAGE_STATUS_TRASH && status != MESSAGE_STATUS_UNREAD{
		return res,errors.New("message status error.")
	}
	err := MgoDBC.Find(bson.M{"Type" : MESSAGE_TYPE_SYSTEM,"Status" : status}).Sort("-_CreateTime").Skip((page-1) * max).Limit(max).All(&res)
	return res,err
}

//发送一个消息
// 注意，请务必事先过滤用户信息
//param sendUserID string 发送用户ID
//param toUserID string 收件用户ID
//param messageType string 消息类型
//param content string 消息内容
//return error 错误信息
func Send(sendUserID string,toUserID string,messageType string,content string) error {
	//检查参数
	if messageType == MESSAGE_TYPE_SYSTEM{
		return errors.New("message type error,cannot set system.")
	}
	if messageType != MESSAGE_TYPE_PUBLIC && messageType != MESSAGE_TYPE_PRIVATE{
		return errors.New("message type error.")
	}
	//内容长度不能超过500汉字 / 1000个英文单词
	if len(content) > 1000{
		return errors.New("message too langer.")
	}
	//检查摘要，并检查最近时间10秒内的内容是否重复
	nowTime := time.Now().Unix()
	lastTime := nowTime + 10
	contentSha1 := encrypt.GetSha1Str(content)
	if contentSha1 == ""{
		return errors.New("content sha1 error.")
	}
	res := FieldsMessage{}
	err := MgoDBC.Find(bson.M{"CreateTime" : bson.M{"$lt" : lastTime} , "SendUserID" : sendUserID , "ContentSha1" : contentSha1}).One(&res)
	if err == nil{
		return err
	}
	//创建消息
	res = FieldsMessage{
		bson.NewObjectId(),
		nowTime,
		sendUserID,
		toUserID,
		messageType,
		MESSAGE_STATUS_UNREAD,
		contentSha1,
		content,
	}
	return MgoDBC.Insert(&res)
}

//修改消息状态
//param id string 消息ID
//param status string 状态
//return error 错误信息
func SetStatus(id string,status string) error {
	//过滤参数
	if status != MESSAGE_STATUS_PUBLIC && status != MESSAGE_STATUS_TRASH && status != MESSAGE_STATUS_UNREAD{
		return errors.New("message status error.")
	}
	//找到消息并修改
	if bson.IsObjectIdHex(id) == false{
		return errors.New("message id is error.")
	}
	//获取信息
	res,err := Get(id)
	if err != nil{
		return err
	}
	if res.Status == status{
		return nil
	}
	res.Status = status
	//返回
	return MgoDBC.UpdateId(bson.ObjectIdHex(id),&res)
}