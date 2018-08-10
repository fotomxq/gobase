package server

import (
	"encoding/json"
	"html/template"
	"fotomxq/gobase/log"
	"fotomxq/gobase/user"
	"fotomxq/gobase/reg"
	"fotomxq/gobase/file"
)

//该文件为server分支文件
// 仅用于反馈头部分

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
//report 反馈头
//////////////////////////////////////////////////////////////////////////////////////////////////////////////

//反馈信息头的通用方法
//param w http.ResponseWriter
//param AppMark string 标识码
//param res interface{} 反馈结构体
func ReportJSONData(headerParams ServerHeaderParams,res interface{}) {
	headerParams.W.Header().Set("Content-Type","application/json")
	resJSON,err := json.Marshal(&res)
	if err != nil{
		log.SendError(headerParams.IP,headerParams.UserInfo.Username,"ServerType.ReportJSONData()",err)
		return
	}
	_,err = headerParams.W.Write(resJSON)
	if err != nil{
		log.SendError(headerParams.IP,headerParams.UserInfo.Username,"ServerType.ReportJSONData()",err)
		return
	}
	return
}

//通用输出页面
//param headerParams HeaderParams
//param templateName string 模版文件名称 必须和src对应
//param templateSrc string 模版文件路径
//param pageData interface{} 附带变量组
func ReportPage(headerParams ServerHeaderParams,templateName string, templateSrc string,pageData map[string]interface{}) {
	err := ReportPageTakeTemplates(headerParams,[]string{},templateName,templateSrc,pageData)
	if err != nil{
		log.SendError(headerParams.IP,headerParams.UserInfo.Username,"ServerType.ReportPage()",err)
	}
}

//通用输出页面，附加自定义模版文件
// 支持附加新的templates模版文件
//param headerParams HeaderParams
//param templateName string 模版文件名称 必须和src对应
//param templateSrc string 模版文件路径
//param pageData interface{} 附带变量组
//return error
func ReportPageTakeTemplates(headerParams ServerHeaderParams,otherTemplateFiles []string,templateName string, templateSrc string,pageData map[string]interface{}) error {
	//构建模版对象
	t,err := template.New(templateName).Funcs(template.FuncMap{
		"PageFuncUnescaped" : PageFuncUnescaped,
		"CheckLoginAuthority" : CheckLoginAuthority,
		"GetUnitToDatetime" : PageFuncGetUnitToDatetime,
		"GetConfigValue" : PageFuncGetConfigValue,
	}).ParseFiles(TemplatesDir + templateSrc)
	if err != nil{
		log.SendError(headerParams.IP,headerParams.UserInfo.Username,"server.ReportPageTakeTemplates()",err)
		return err
	}

	//附加全局缓冲文件
	for _,v := range TemplatesParseFiles{
		t,err = t.ParseFiles(TemplatesDir + file.Sep + "base" + file.Sep + v)
		if err != nil{
			log.SendError(headerParams.IP,headerParams.UserInfo.Username,"server.ReportPageTakeTemplates()",err)
			return err
		}
	}

	//如果存在附加文件
	for _,v := range otherTemplateFiles{
		t,err = t.ParseFiles(TemplatesDir + file.Sep + v)
		if err != nil{
			log.SendError(headerParams.IP,headerParams.UserInfo.Username,"server.ReportPageTakeTemplates()",err)
			return err
		}
	}

	//构建附带变量组
	for k,v := range PageGlobData{
		pageData[k] = v
	}

	//附带用户个人信息
	pageData["UserInfo"] = headerParams.UserInfo

	//输出页面
	err = t.Execute(headerParams.W,pageData)
	if err != nil{
		log.SendError(headerParams.IP,headerParams.UserInfo.Username,"server.ReportPageTakeTemplates()",err)
		return err
	}
	return nil
}

//输出通用页面下的模版文件
// 该数据将不带任何变量输出
// 主要用于错误输出
//param headerParams HeaderParams
//param name string 模版文件名
func ReportPageE(headerParams ServerHeaderParams,name string){
	pageData := map[string]interface{}{}
	ReportPage(headerParams,name,name,pageData)
	return
}

//输出错误页面
//param headerParams HeaderParams
func ReportErrorPage(headerParams ServerHeaderParams){
	pageData := map[string]interface{}{}
	ReportPage(headerParams,"page-error-404.tmpl","page-error-404.tmpl",pageData)
}

//输出更新KEY值页面
//param headerParams HeaderParams
//return bool 是否通过验证
func ReportKeyPage(headerParams ServerHeaderParams) bool {
	//检查密匙，是否可以启动服务？
	if reg.IsOK(PageGlobData["SystemKey"].(string)) == false{

		log.SendText(headerParams.IP,headerParams.UserInfo.Username,"GoBase.Server.ReportKeyPage",log.MessageTypeSafety,"system key error.")

		pageData := map[string]interface{}{
			"RegKey" : reg.GetKey(),
		}
		ReportPage(headerParams,"page-error-key.tmpl","page-error-key.tmpl",pageData)
		return false
	}
	return true
}

//获取一系列用户ID对应用户名数据组
//param headerParams HeaderParams
func ReportUsersByID(headerParams ServerHeaderParams) {
	//初始化
	res := ReportUsersByIDType{
		false,
		"",
		true,
		false,
		"",
		0,
		[]ReportUsersByIDDataType{},
	}
	//获取参数
	ids := GetPostArray(headerParams,"ids")
	//遍历参数，获取用户信息
	for _,v := range ids{
		userInfo,err := user.GetUserByID(v)
		if err != nil{
			log.SendError(headerParams.IP,headerParams.UserInfo.Username,"ServerType.ReportUsersByID",err)
			ReportJSONData(headerParams,res)
			return
		}
		res.Data = append(res.Data,ReportUsersByIDDataType{
			userInfo.ID.Hex(),
			userInfo.Username,
			userInfo.Name,
		})
	}
	res.Count = len(res.Data)
	res.Status = true
	//反馈数据
	ReportJSONData(headerParams,res)
}