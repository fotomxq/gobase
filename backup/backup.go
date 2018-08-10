package backup

import (
	"time"
	"gopkg.in/mgo.v2"
	"encoding/json"
	"github.com/pkg/errors"
	"os/exec"
	"fotomxq/gobase/file"
	"fotomxq/gobase/log"
	"fotomxq/gobase/config"
)

//注意，该模块仅适用于非大型分布式文件系统、或超大型密集数据库结构
// 如您的数据库数据量极大，请运维人员针对性备份和还原操作。

var(
	//备份路径 eg : ./dir/
	BackupDir string
	//要备份的文件夹
	BackupFolders []string
	//要备份的数据表列
	BackupDbs []string
	//每隔N小时备份一次
	BackupTimeDay int64
	//每天备份的小时、分钟
	BackupTimeH int
	BackupTimeM int
	//应用标识码
	AppMark string
	//DatabaseMgoDatabaseName
	DatabaseMgoDatabaseName string
)

//初始化
// 初始化完成后，可自由修改相关参数
func Run(){
	BackupDir = "." + file.Sep + "backups" + file.Sep
	BackupFolders = []string{}
	BackupDbs = []string{}
	BackupTimeDay = 24
	BackupTimeH = 3
	BackupTimeM = 30
}

//备份文件结构体
// 反馈头专用
type BackupListType struct{
	//文件名称
	Name string `json:"Name"`
	//文件路径
	Src string `json:"Src"`
	//创建时间
	CreateTime int64 `json:"CreateTime"`
	//文件大小
	Size int64 `json:"Size"`
}

//获取备份目录所有文件
//return []BackupListType 文件序列
//return int 文件个数
func GetList() ([]BackupListType,int){
	//初始化
	res := []BackupListType{}
	count := 0
	//查询文件列表
	fileList,err := file.GetFileList(BackupDir,[]string{},true)
	if err != nil{
		log.SendError("0.0.0.0","",AppMark+".BackupType.GetList()",err)
		return res,count
	}
	//建立信息组
	for _,v := range fileList{
		//获取文件信息
		info,err := file.GetFileInfo(v)
		if err != nil{
			return res,count
		}
		//如果文件为temp临时文件夹，则跳过
		if info.Name() == "temp"{
			continue
		}
		//重新构建
		res = append(res, BackupListType{
			info.Name(),
			v,
			info.ModTime().Unix(),
			info.Size(),
		})
	}
	count = len(res)
	//返回
	return res,count
}

//备份数据库数据
// 自动备份到路径下的 ./backups/[时间]/
// log日志数据表不会被备份，将会直接按照标准化结构输出到log目录下，请注意检索
//return bool 是否成功
func RunBackup() bool {
	//状态是否允许？
	if CheckStatus() == false{
		return false
	}
	//构建新的路径
	newBackupFolderSrc := BackupDir + "temp"
	err := file.CreateFolder(newBackupFolderSrc)
	if err != nil{
		log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
		return false
	}
	//构建db目录
	newBackupFolderDBSrc := newBackupFolderSrc + file.Sep + "db"
	err = file.CreateFolder(newBackupFolderDBSrc)
	if err != nil{
		log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
		return false
	}
	//获取参数
	BackupMgoUsername,err := config.Get("BackupMgoUsername")
	if err != nil{
		log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
		return false
	}
	BackupMgoPassword,err := config.Get("BackupMgoPassword")
	if err != nil{
		log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
		return false
	}
	BackupMgoPort,err := config.Get("BackupMgoPort")
	if err != nil{
		log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
		return false
	}
	BackupMgoIP,err := config.Get("BackupMgoIP")
	if err != nil{
		log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
		return false
	}
	//使用mgo内置命令备份数据库
	var cmdExec *exec.Cmd
	if BackupMgoUsername != "" && BackupMgoPassword != ""{
		cmdExec = exec.Command("mongodump","-h",BackupMgoIP,"--port",BackupMgoPort,"-o",newBackupFolderDBSrc,"-d",DatabaseMgoDatabaseName)
	}else{
		cmdExec = exec.Command("mongodump","-h",BackupMgoIP,"--port",BackupMgoPort,"-u",BackupMgoPassword,"-p",BackupMgoPassword,"-o",newBackupFolderDBSrc,"-d",DatabaseMgoDatabaseName)
	}
	err = cmdExec.Run()
	if err != nil{
		log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
	}
	//构建Folders目录
	newBackupFolderFSrc := newBackupFolderSrc + file.Sep + "folders" + file.Sep
	err = file.CreateFolder(newBackupFolderFSrc)
	if err != nil{
		log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
		return false
	}
	//遍历复制文件路径
	for _,vFolderSrc := range BackupFolders{
		//基本名称
		vFolderNames,err := file.GetFileNames(vFolderSrc)
		if err != nil{
			log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
			continue
		}
		vFolderName := vFolderNames["name"]
		//将vFolder复制到对应目录下
		if file.CopyFolder(vFolderSrc,newBackupFolderFSrc + vFolderName) == false{
			log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",errors.New("cannot copy folder " + vFolderSrc + " to " + newBackupFolderFSrc + vFolderName))
			continue
		}
	}
	//压缩文件夹
	newBackupZipSrc := BackupDir + time.Now().Format("20060102_150405_backups") + ".zip"
	err = file.ZipDir(newBackupFolderSrc,newBackupZipSrc)
	if err != nil{
		//备份失败，则删除该备份，并退出
		log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
		err = file.DeleteF(newBackupFolderSrc)
		if err != nil{
			log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
		}
		return false
	}
	//删除压缩目录
	err = file.DeleteF(newBackupFolderSrc)
	if err != nil{
		log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
		return false
	}
	//完成提示
	log.SendText("0.0.0.0","",AppMark+".BackupType.Backup()",log.MessageTypeSystem,"Complete a standardized backup operation, backup file path:" + newBackupZipSrc)
	return true
}

//还原数据库
//param name string 文件名称
//return bool 是否成功
func RunReturn(name string) bool{
	//状态是否允许？
	if CheckStatus() == false{
		return false
	}
	//检查文件是否存在
	if Check(name) == false{
		return false
	}
	zipSrc := BackupDir + file.Sep + name
	//解压目录到temp
	tempDir := BackupDir + file.Sep + "temp"
	err := file.UnZip(zipSrc,tempDir)
	if err != nil{
		log.SendError("0.0.0.0","",AppMark+".BackupType.RunReturn()",err)
		return false
	}
	//如果存在db目录
	tempDbDir := tempDir + file.Sep + "db"
	//如果DB目录存在，遍历所有子文件
	// 每个子文件代表对应的数据表名称
	if file.IsExist(tempDbDir) == true{
		//获取文件列表
		folderList,err := file.GetFileList(tempDbDir,[]string{},true)
		if err != nil{
			log.SendError("0.0.0.0","",AppMark+".BackupType.RunReturn()",err)
		}
		for _,vFolder := range folderList{
			//获取文件信息
			info,err := file.GetFileNames(vFolder)
			if err != nil{
				log.SendError("0.0.0.0","",AppMark+".BackupType.RunReturn()",err)
				continue
			}
			//获取参数
			BackupMgoUsername,err := config.Get("BackupMgoUsername")
			if err != nil{
				log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
				return false
			}
			BackupMgoPassword,err := config.Get("BackupMgoPassword")
			if err != nil{
				log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
				return false
			}
			BackupMgoPort,err := config.Get("BackupMgoPort")
			if err != nil{
				log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
				return false
			}
			BackupMgoIP,err := config.Get("BackupMgoIP")
			if err != nil{
				log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
				return false
			}
			//执行还原脚本
			var cmdExec *exec.Cmd
			if BackupMgoUsername != "" && BackupMgoPassword != ""{
				cmdExec = exec.Command("mongorestore","-h",BackupMgoIP,"--port",BackupMgoPort,"--dir",vFolder,"--drop")
			}else{
				cmdExec = exec.Command("mongorestore","-h",BackupMgoIP,"--port",BackupMgoPort,"-u",BackupMgoUsername,"-p",BackupMgoPassword,"-d",info["name"],"--dir",vFolder,"--drop")
			}
			err = cmdExec.Run()
			if err != nil{
				log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
				continue
			}
		}
	}
	//如果存在folders目录
	tempFolderDir := tempDir + file.Sep + "folders"
	if file.IsExist(tempFolderDir) == true{
		//获取文件列表
		folderList,err := file.GetFileList(tempFolderDir,[]string{},true)
		if err != nil{
			log.SendError("0.0.0.0","",AppMark+".BackupType.RunReturn()",err)
		}
		for _,vFolder := range folderList{
			//获取文件信息
			info,err := file.GetFileNames(vFolder)
			if err != nil{
				log.SendError("0.0.0.0","",AppMark+".BackupType.RunReturn()",err)
				continue
			}
			//删除根目录下对应目录
			destDir := "." + file.Sep + info["name"]
			if file.IsExist(destDir) == true{
				err = file.DeleteF(destDir)
				if err != nil{
					//记录错误但不中断
					log.SendError("0.0.0.0","",AppMark+".BackupType.RunReturn()",err)
				}
			}
			//将备份目录复制到目标位置
			if file.CopyFolder(vFolder,destDir) == false{
				log.SendError("0.0.0.0","",AppMark+".BackupType.RunReturn()",errors.New("copy error."))
				continue
			}
		}
	}
	//删除临时文件
	err = file.DeleteF(tempDir)
	if err != nil{
		log.SendError("0.0.0.0","",AppMark+".BackupType.Backup()",err)
		return false
	}
	//完成提示
	log.SendText("0.0.0.0","",AppMark+".BackupType.Backup()",log.MessageTypeSystem,"Complete system restore and restore the data package:" + name + ",Please check the system.")
	//返回成功
	return true
}

//还原数据表数据
//param tableName string 数据表名称
//param data []byte 数据
//param c *mgo.Collection 集合对象
//return error 错误代码
func RunReturnTable(tableName string, data []byte,c *mgo.Collection) error {
	//测试map interface{}
	var res []map[string]interface{}
	err := json.Unmarshal(data,&res)
	if err != nil{
		return err
	}
	//将数据直接插入数据表
	err = c.Insert(&res)
	if err != nil{
		return err
	}
	return nil
}

//检查备份和还原状态
//return bool 是否正在进行
func CheckStatus() bool {
	//检查是否存在temp目录，如果存在，则请用户等待
	// 因为可能维护工具正在进行维护，否则需用户自行中断服务并删除该目录才能进行
	tempDir := BackupDir + file.Sep + "temp"
	if file.IsExist(tempDir) == true{
		return false
	}
	return true
}

//删除备份文件
//param name string 文件名称
//return bool 是否成功
func Delete(name string) bool{
	//检查文件是否存在
	if Check(name) == false{
		return false
	}
	//删除文件
	err := file.DeleteF(BackupDir + file.Sep + name)
	if err != nil{
		log.SendError("0.0.0.0","",AppMark+".BackupType.Delete",err)
		return false
	}
	//返回成功
	return true
}

//检查备份文件是否存在
//param name string 文件名称
//return bool 是否成功
func Check(name string) bool{
	//获取列表
	backupList,_ := GetList()
	//寻找文件是否存在？
	for _,v := range backupList{
		if v.Name == name{
			return true
			break
		}
	}
	return false
}

//获取文件路径
//param name string 文件名称
//return string 文件路径
func GetFileSrc(name string) string{
	return BackupDir + file.Sep + name
}