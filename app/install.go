package app

import (
	"strings"
	"strconv"
	"encoding/json"
	"github.com/pkg/errors"
	"fotomxq/gobase/file"
	"fotomxq/gobase/session"
	"fotomxq/gobase/config"
	"fotomxq/gobase/authority"
	"fotomxq/gobase/user"
)

//安装配置
type ServerInstallConfigType struct {
	Configs [][]string `json:"Configs"`
	URL [][]string `json:"URL"`
	Group [][]string `json:"Group"`
	User [][]string `json:"User"`
}

//初始化安装系统
//return error 错误信息
func Install() error{
	//读取安装文件，检查数据完整性并创建
	by,err := file.LoadFile("." + file.Sep + "install.json")
	if err != nil{
		return err
	}
	installData := ServerInstallConfigType{}
	err = json.Unmarshal(by,&installData)
	if err != nil{
		return err
	}
	//检查配置完成行，如果不存在则创建对应配置
	for _,v := range installData.Configs{
		_,err = config.Get(v[0])
		if err != nil{
			err = config.Set(v[0],v[1])
			if err != nil{
				return err
			}
		}
	}
	//检查并构建URL
	for _,v := range installData.URL{
		_,err = authority.GetByURL(v[0])
		if err != nil{
			err = authority.Set(v[0],v[1],v[2],v[3])
			if err != nil{
				return err
			}
		}
	}
	//检查并创建用户组
	for _,v := range installData.Group{
		_,err = user.GetGroupByMark(v[0])
		if err != nil{
			v6Int,err := strconv.Atoi(v[6])
			if err != nil{
				v6Int = 99
			}
			//创建用户组
			err = user.CreateGroup(v[0],v[1],[]string{},v[2],v[3],v[4],v6Int)
			if err != nil{
				return err
			}
			//为用户组添加权限
			// 跳过管理员级别
			if v[0] != user.USER_GROUP_MARK_ADMIN{
				groupInfo,err := user.GetGroupByMark(v[0])
				if err != nil{
					return err
				}
				//生成权限组
				authority := strings.Split(v[5],"|")
				//如果存在权限，则添加
				// 否则跳过
				if len(authority) > 0{
					err = user.UpdateGroup(groupInfo.ID.Hex(),groupInfo.Mark,groupInfo.Status,authority,groupInfo.Name,groupInfo.Des,groupInfo.IndexURL,v6Int)
					if err != nil{
						return err
					}
				}
			}
		}
	}
	//装载用户加密
	user.PasswdEncrypt,err = config.Get("EncryptPasswdEncrypt")
	if err != nil{
		return err
	}
	//设置session
	sessionMark,err := config.Get("SessionMark")
	if err != nil{
		return err
	}
	version,err := config.Get("Version")
	if err != nil{
		return err
	}
	session.AppName = sessionMark + version
	//找到管理员级别组，将所有权限添加进去
	groupInfo,err := user.GetGroupByMark(user.USER_GROUP_MARK_ADMIN)
	if err != nil{
		return errors.New("cannot find admin user group,error : " + err.Error())
	}
	for _,v := range installData.URL{
		isFind := false
		for _,v2 := range groupInfo.Authority{
			if v2 == v[1]{
				isFind = true
			}
		}
		if isFind == false{
			groupInfo.Authority = append(groupInfo.Authority,v[1])
		}
	}
	err = user.UpdateGroupInParams(groupInfo)
	if err != nil{
		return errors.New("cannot update group in params ,error : " + err.Error())
	}
	//创建用户并添加用户组
	userInfoConfig,err :=user.GetInfoConfig()
	if err != nil{
		return err
	}
	for _,v := range installData.User{
		_,err = user.GetUserByUsername(v[0])
		if err != nil{
			userInfo,err := user.Create(v[0],v[1],v[2],userInfoConfig,"127.0.0.1",true)
			if err != nil{
				return err
			}
			v5int64,err := strconv.ParseInt(v[5],10,64)
			if err != nil{
				return err
			}
			err = user.SetGroupExpireTime(userInfo.ID.Hex(),v[3],v[4],v5int64)
			if err != nil{
				return err
			}
		}
	}
	return nil
}

