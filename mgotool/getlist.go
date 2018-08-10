package mgotool

import (
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fotomxq/gobase/filter"
)

//mgo数据库操作对象

//通用获取mgo dbc数据集合的列表一体式方案
// 自动获取page\max\sort\desc数据，并根据q条件查询返回query数据，用于后续工作
//param r *http.Request
//param q []bson.M 条件
//return *mgo.Query mgo查询句柄
func GetList(r *http.Request,dbc *mgo.Collection,q bson.M) (*mgo.Query){
	page := filter.FilterPage(r.FormValue("page"))
	max := filter.FilterMax(r.FormValue("max"))
	sort := r.FormValue("sort")
	if sort == ""{
		sort = "_id"
	}
	desc := r.FormValue("desc") == "true"
	if desc == true{
		sort = "-" + sort
	}
	return dbc.Find(q).Skip((page-1) * max).Limit(max).Sort(sort)
}

