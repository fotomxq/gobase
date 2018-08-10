package server

import (
	"strings"
	"encoding/json"
	"net/http"
	"strconv"
	"fotomxq/gobase/log"
	"fotomxq/gobase/ipaddr"
	"fotomxq/gobase/user"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
//页面通用方法
//////////////////////////////////////////////////////////////////////////////////////////////////////////////

//跳转到其他URL的通用方法
//param headerParams HeaderParams
//param url string URL地址
func GoURL(headerParams ServerHeaderParams,url string){
	http.Redirect(headerParams.W,headerParams.R,url,http.StatusFound)
}

//访客处理
//param headerParams HeaderParams
func LogVisit(headerParams ServerHeaderParams){
	//记录来访
	log.SendText(headerParams.IP,headerParams.UserInfo.Username,AppMark+".ActionURL()",log.MessageTypeVisit,headerParams.R.URL.Path)
	//检查IP是否可以通过
	if ipaddr.CheckOK(headerParams.IP) == false {
		return
	}
}

//前置通用处理
// 获取头headerParams
//param w http.ResponseWriter
//param r *http.Request
//return ServerHeaderParams
func GetHeaderParams(w http.ResponseWriter, r *http.Request) ServerHeaderParams{
	//基础参数
	headerParams := ServerHeaderParams{
		w,
		r,
		strings.Split(r.URL.Path, "/"),
		ipaddr.IPGet(r,"all",true),
		"",
		user.FieldsUser{},
	}

	//记录访客信息
	// 不能在这里记录，因为会造成无法记录到用户名信息
	//LogVisit(headerParams)

	//返回头
	return headerParams
}

//系统设置数据格式
type PostDataType struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

//解析POST的JSON结构体
// 仅限于处理[]string{}类型
//param headerParams HeaderParams
//param name string 名称
//return []string 数据组
//return error 错误信息
func GetPostArray(headerParams ServerHeaderParams,name string) ([]string){
	res := []string{}
	val := headerParams.R.FormValue(name)
	if val == ""{
		return res
	}
	err := json.Unmarshal([]byte(val),&res)
	if err != nil{
		res = []string{
			val,
		}
	}
	return res
}

//解析POST的JSON结构体为Int格式
// 仅限于处理[]int{}类型
//param headerParams HeaderParams
//param name string 名称
//return []int 数据组
//return error 错误信息
func GetPostArrayToInt(headerParams ServerHeaderParams,name string) ([]int){
	res := []int{}
	val := headerParams.R.FormValue(name)
	if val == ""{
		return res
	}
	err := json.Unmarshal([]byte(val),&res)
	if err != nil{
		valInt,err := strconv.Atoi(val)
		if err != nil{
			res = []int{
				0,
			}
		}else{
			res = []int{
				valInt,
			}
		}
	}
	return res
}

//解析POST的JSON结构体为Int64格式
// 仅限于处理[]int64{}类型
//param headerParams HeaderParams
//param name string 名称
//return []int64 数据组
//return error 错误信息
func GetPostArrayToInt64(headerParams ServerHeaderParams,name string) ([]int64){
	res := []int64{}
	val := headerParams.R.FormValue(name)
	if val == ""{
		return res
	}
	err := json.Unmarshal([]byte(val),&res)
	if err != nil{
		valInt64,err := strconv.ParseInt(val,10,64)
		if err != nil{
			res = []int64{
				0,
			}
		}else{
			res = []int64{
				valInt64,
			}
		}
	}
	return res
}

//解析POST的JSON结构体为bool格式
// 仅限于处理[]bool{}类型
//param headerParams HeaderParams
//param name string 名称
//return []bool 数据组
//return error 错误信息
func GetPostArrayToBool(headerParams ServerHeaderParams,name string) ([]bool){
	res := []bool{}
	val := headerParams.R.FormValue(name)
	if val == ""{
		return res
	}
	err := json.Unmarshal([]byte(val),&res)
	if err != nil{
		if val == "true"{
			res = []bool{
				true,
			}
		}else{
			res = []bool{
				false,
			}
		}
	}
	return res
}