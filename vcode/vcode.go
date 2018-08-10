package vcode

import (
	"time"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/mojocn/base64Captcha"
	"github.com/pkg/errors"
	"fotomxq/gobase/mgotool"
)

var(
	//参数配置
	VerificationCodeConfig = base64Captcha.ConfigCharacter{
		Height:             60,
		Width:              240,
		//const CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode:               base64Captcha.CaptchaModeNumber,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         6,
	}
	
	//数据库句柄
	MgoDBC *mgo.Collection

	//过期时间
	// 默认 1800 半小时
	ExpireTime int64
)

//初始化
func Run() {
	//验证码数据表
	MgoDBC = mgotool.MgoDB.C("verification_code")
	//设置过期时间
	ExpireTime = 1800
}

//定时任务
func RunAuto(){
	//删除所有过期数据
	_,_ = MgoDBC.RemoveAll(bson.M{"ExpireTime" : bson.M{"$lt" : time.Now().Unix()}})
}

//获取一个验证码
//param token string Token
//return string 验证码图形Base64结构
//return error 错误
func GetImage(token string) (string,error){
	//删除该token下所有验证码数据
	_,err := MgoDBC.RemoveAll(bson.M{"Token" : token})
	if err != nil{
		return "",err
	}
	//创建验证码
	key, cap := base64Captcha.GenerateCaptcha("", VerificationCodeConfig)
	newD := FieldsVerificationCodeDataType{
		bson.NewObjectId(),
		token,
		key,
		time.Now().Unix() + ExpireTime,
	}
	err = MgoDBC.Insert(&newD)
	if err != nil{
		return "",err
	}
	//获取Base64结构
	base64stringC := base64Captcha.CaptchaWriteToBase64Encoding(cap)
	//返回数据流
	return base64stringC,nil
}

//验证码一个验证码
//param token string Token
//param value string 要验证的文本
//return error 错误
func Check(token string,value string) error{
	res := FieldsVerificationCodeDataType{}
	//获取数据
	err := MgoDBC.Find(bson.M{"Token" : token}).One(&res)
	if err != nil{
		return errors.New("Token not have verification code.")
	}
	//检查是否相等？
	if base64Captcha.VerifyCaptcha(res.Value,value) == false{
		return errors.New("Verification code is error.")
	}
	//删除验证码并返回成功
	_,err = MgoDBC.RemoveAll(bson.M{"Token" : token})
	if err != nil{
		return err
	}
	return nil
}
