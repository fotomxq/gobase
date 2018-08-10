package source

import (
	"encoding/json"
	"encoding/base64"
	"fotomxq/gobase/file"
	"fotomxq/gobase/log"
)

//将资源打包为单一文件source.dat
// 只支持三个级别的文件结构
var(
	//资源集合
	Data map[string]map[string]map[string]string
)

//初始化
func Run(){
	Data = map[string]map[string]map[string]string{}
}

//将资源进行打包整理
//param dir string 需要打包的目录
//param target string 目标路径
//return error 错误
func SaveDat(dir string,target string) error {
	//所有html打包到html/html内
	//文件夹必须是三级目录分类，按照三级目录对照存放
	//尝试获取一级目录、所有html文件
	dataA,err := file.GetFileList(dir,[]string{},true)
	if err != nil{
		return err
	}
	//初始化
	Data["html"] = map[string]map[string]string{}
	Data["html"]["html"] = map[string]string{}
	//遍历一级目录
	for _,valA := range dataA{
		//获取基础名称
		valAInfo,err := file.GetFileInfo(valA)
		if err != nil{
			return err
		}
		valABaseName := valAInfo.Name()
		//所有文件存放到html/html下
		if file.IsFile(valA) == true{
			by,err := file.LoadFile(valA)
			if err != nil{
				return err
			}
			Data["html"]["html"][valABaseName] = base64.StdEncoding.EncodeToString(by)
			log.SendText("0.0.0.0","-1","Source.SaveDat",log.MessageTypeMessage,"载入数据：" + valA)
			continue
		}
		//建立数据
		_,ok := Data[valABaseName]
		if ok == false{
			Data[valABaseName] = map[string]map[string]string{}
		}
		//继续遍历二级目录
		dataB,err := file.GetFileList(valA,[]string{},true)
		if err != nil{
			return err
		}
		for _,valB := range dataB{
			//跳过所有文件，只看目录
			if file.IsFile(valB) == true{
				continue
			}
			//获取基础名称
			valBInfo,err := file.GetFileInfo(valB)
			if err != nil{
				return err
			}
			valBBaseName := valBInfo.Name()
			//建立数据
			_,ok := Data[valABaseName][valBBaseName]
			if ok == false{
				Data[valABaseName][valBBaseName] = map[string]string{}
			}
			//遍历三级目录
			dataC,err := file.GetFileList(valB,[]string{},true)
			if err != nil{
				return err
			}
			for _,valC := range dataC{
				//必须是文件，跳过所有目录
				if file.IsFolder(valC) == true{
					continue
				}
				//获取基础名称
				valCInfo,err := file.GetFileInfo(valC)
				if err != nil{
					return err
				}
				valCBaseName := valCInfo.Name()
				//存储文件
				by,err := file.LoadFile(valC)
				if err != nil{
					return err
				}
				Data[valABaseName][valBBaseName][valCBaseName] = base64.StdEncoding.EncodeToString(by)
				log.SendText("0.0.0.0","-1","Source.SaveDat",log.MessageTypeMessage,"载入数据：" + valC)
			}
		}
	}
	//将全体数据编译为json后，以二进制形式存储
	dataJson,err := json.Marshal(Data)
	if err != nil{
		return err
	}
	return file.WriteFile(target,dataJson)
}

//获取source.dat数据
//param src string 资源路径
//return error 错误
func GetDat(src string) error {
	dataByte,err := file.LoadFile(src)
	if err != nil{
		return err
	}
	err = json.Unmarshal(dataByte,&Data)
	if err != nil{
		return err
	}
	for k,v := range Data{
		for k2,v2 := range v{
			for k3,v3 := range v2{
				by,err := base64.StdEncoding.DecodeString(v3)
				if err != nil{
					return err
				}
				Data[k][k2][k3] = string(by)
			}
		}
	}
	return nil
}