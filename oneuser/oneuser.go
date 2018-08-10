package oneuser

//单一用户模块处理器
// 用户名和密码不可更改，需要管理员手动修改配置文件，之后将其赋予给该内部变量
// 使用前务必给Username和Password变量赋值

type OneUserType struct {
	//用户名
	// = "admin"
	Username string

	//密码
	// = "adminadmin"
	Password string

	//session mark
	SessionMark string
}

var(
	//单一用户组对象
	OneUser OneUserType
)

//登陆用户
//param user string 用户名
//param passwd string 密码
//param mark string 会话标记
//return bool 是否登陆成功
func (this *OneUserType) Login(user string,passwd string,mark string) bool{
	if user == this.Username && this.Password == passwd{
		this.SessionMark = mark
		return true
	}
	return false
}

//退出登陆
//param mark string 会话标记
func (this *OneUserType) Logout(mark string) bool{
	if mark == this.SessionMark{
		this.SessionMark = ""
		return true
	}
	return false
}

//检查登陆状态
//param mark string 会话标记
func (this *OneUserType) Status(mark string) bool{
	if this.SessionMark == ""{
		return false
	}
	if mark == this.SessionMark{
		return true
	}
	return false
}