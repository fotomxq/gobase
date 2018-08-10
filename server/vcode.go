package server

import (
	"fotomxq/gobase/log"
	"fotomxq/gobase/vcode"
)

//验证码模块

//输出一个验证码
//param headerParams HeaderParams
func ImageVerificationCode(headerParams ServerHeaderParams) {
	imageBase64,err := vcode.GetImage(headerParams.Token)
	if err != nil{
		log.SendError(headerParams.IP,headerParams.UserInfo.Username,"ServerType.ImageVCode",err)
		return
	}
	headerParams.W.Write([]byte(imageBase64))
}