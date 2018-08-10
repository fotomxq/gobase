package server

import (
	"strings"
	"fotomxq/gobase/log"
	"fotomxq/gobase/file"
)

//下载文件处理
// 给定一个文件序列组，该序列组是经过严格判定符合标准的，且不允许出现../的字符串结构
//param headerParams ServerHeaderParams
//param src string 文件路径
//param name string 文件名称
func DownloadFile(headerParams ServerHeaderParams,src string,name string) {
	//检查src内是否包含..
	if strings.Count(src,"..") > 0{
		ReportErrorPage(headerParams)
		return
	}
	//打开文件
	fd,err := file.LoadFile(src)
	if err != nil{
		log.SendError(headerParams.IP,headerParams.UserInfo.Username,"ServerType.DownloadFile",err)
		ReportErrorPage(headerParams)
		return
	}
	//添加头信息
	headerParams.W.Header().Add("Content-Type","application/octet-stream")
	headerParams.W.Header().Add("content-disposition","attachment; filename=\""+name+"\"")
	//写入文件
	_,err = headerParams.W.Write(fd)
	if err != nil{
		log.SendError(headerParams.IP,headerParams.UserInfo.Username,"ServerType.DownloadFile",err)
		ReportErrorPage(headerParams)
		return
	}
}