package sqltool

import (
	"database/sql"
	"strconv"
	"github.com/pkg/errors"
)

//sql工具盒
// 集成通用生成工具，帮助快速生成和执行复杂的sql

//删除模块

//删除ID
//param db *sql.DB 数据库句柄
//param table string 表名称
//param id int64 ID
//return bool 是否成功
func DeleteID(db *sql.DB,table string,id int64) error {
	sql := "delete from `" + table + "` where `id` = ?"
	stmt,err := db.Prepare(sql)
	if err != nil{
		return err
	}
	defer stmt.Close()
	result,err := stmt.Exec(id)
	if err != nil{
		return err
	}
	row,err := result.RowsAffected()
	if err != nil{
		return err
	}
	if row > 0{
		return nil
	}
	return errors.New("cannot delete id, row < 1.")
}

//删除指定的值
//param db *sql.DB 数据库句柄
//param table string 表名称
//param field string 字段名称
//param value string 值
//return bool 是否成功
func DeleteFieldValue(db *sql.DB,table string,field string,value string) error {
	sql := "delete from `" + table + "` where `" + field + "` = ?"
	stmt,err := db.Prepare(sql)
	if err != nil{
		return err
	}
	defer stmt.Close()
	result,err := stmt.Exec(value)
	if err != nil{
		return err
	}
	row,err := result.RowsAffected()
	if err != nil{
		return err
	}
	if row > 0{
		return nil
	}
	return errors.New("cannot delete id, row < 1.")
}

//字段组模块

//检查字段是否存在
// 如果不存在则返回默认值
//param fields []string 字段组
//param field string 需要检查的字段名称
//param defaultKey int 默认值，字段组键值
//return string 该字段值，如果失败则返回默认值
func CheckField(fields []string,field string,defaultKey int) string {
	for _,v := range fields{
		if v == field{
			return v
		}
	}
	return fields[defaultKey]
}

//建立数据模块

//插入新的数据
//param db *sql.DB
//param table string
//param sqlFields string
//param sqlValues string
//param args ...interface{}
//return int64 插入的ID
func Insert(db *sql.DB,table string,sqlFields string,sqlValues string,args ...interface{}) (int64,error){
	sql := "insert into `" + table + "`(" + sqlFields + ") values(" + sqlValues + ")"
	stmt,err := db.Exec(sql,args)
	if err != nil{
		return -1,err
	}
	id,err := stmt.LastInsertId()
	if err != nil{
		return -1,err
	}
	return id,nil
}

//获取页数部分
//param page int 页数
//param max int 页长
//return string sql页数部分
func PageLimit(page int,max int) string{
	return "limit " + strconv.Itoa((page-1)*max) + "," + strconv.Itoa(max)
}

//获取排序部分
//param sort string 要排序的字段
//param desc bool 是否倒叙
//return string sql排序部分
func PageSort(sort string,desc bool) string{
	var descStr string
	if desc == true{
		descStr = "desc"
	}else{
		descStr = "asc"
	}
	return "order by `" + sort + "` " + descStr
}

//查询模块

//查询ID
//param db *sql.DB 数据库
//param table string 数据表
//param fields string 字段组
//param id int64 ID
//return *sql.Row 结果集
func SelectID(db *sql.DB,table string,fields string,id int64) (*sql.Row){
	sql := "select " + fields + " from `" + table + "` where `id` = ?"
	row := db.QueryRow(sql,id)
	return row
}

//查询特定值
//param db *sql.DB 数据库
//param table string 数据表
//param fields string 字段组
//param searchField string 查询字段名称
//param value string 值
//return *sql.Row 结果集
func SelectValue(db *sql.DB,table string,fields string,searchField string,value string) (*sql.Row){
	sql := "select " + fields + " from `" + table + "` where `" + searchField + "` = ?"
	row := db.QueryRow(sql,value)
	return row
}


//查询列表
//param db *sql.DB 数据库句柄
//param table string 数据表
//param field string 字段
//param fields string 字段组
//param value string 查询值
//param sort string 要排序的字段
//param desc bool 是否倒叙
//return *sql.Rows,bool 数据集,是否成功
func SelectList(db *sql.DB,table string,fields string,where string,page int,max int,sort string,desc bool,args ...interface{}) (*sql.Rows,error){
	sql := "select " + fields + " from `" + table + "` where " + where + " " + PageSort(sort,desc) + " " + PageLimit(page,max)
	rows,err := db.Query(sql,args)
	if err != nil{
		return rows,err
	}
	return rows,nil
}

//更新数据表组
//param db *sql.DB 数据库句柄
//param table string 表名称
//param id int64 ID
//param field string 字段
//param value string 修改值
//return bool 是否成功
func UpdateID(db *sql.DB,table string,id int64,field string,value string) error{
	sql := "update `" + table + "` set `" + field + "` = ? where `id` = ?"
	result,err := db.Exec(sql,value,id)
	if err != nil{
		return err
	}
	row,err := result.RowsAffected()
	if err != nil{
		return err
	}
	if row > 0{
		return nil
	}
	return errors.New("cannot update id, row < 1.")
}

//更新数据表组
//param db *sql.DB 数据库句柄
//param table string 表名称
//param sets string 修改组
//param where string 条件
//param args ...interface{} 参数组
//return bool 是否成功
func UpdateValues(db *sql.DB,table string,sets string,where string,args ...interface{}) error{
	sql := "update `" + table + "` set " + sets + " where " + where
	result,err := db.Exec(sql,args)
	if err != nil{
		return err
	}
	row,err := result.RowsAffected()
	if err != nil{
		return err
	}
	if row > 0{
		return nil
	}
	return errors.New("cannot update value, row < 1.")
}

//标签关联模块
// 管理标签、将标签和指定ID关联、统计标签存在关联的个数
//对应的列必须如下：
//	标签表
//		id \ name
//	标签关联表
//		id \ target_id \ tag_id

var(
//标签表
//tagsTableName string = "tags"

//标签关联表
//tagsBindTableName string = "tags_bind"
)

//获取标签列表
//param db *sql.DB 数据库句柄
//param tagsTableName string 标签表
//param page int 页数
//param max int 长度
//return []string 组
func TagsList(db *sql.DB,tagsTableName string,page int,max int) (map[int]string){
	sql := "select `id`,`name` from `" + tagsTableName + "` " + PageLimit(page,max) + " " + PageSort("name",false)
	rows,err := db.Query(sql)
	if err != nil{
		return map[int]string{}
	}
	var result map[int]string
	for rows.Next(){
		var id int
		var name string
		err = rows.Scan(&id,&name)
		if err != nil{
			return map[int]string{}
		}
		result[id] = name
	}
	return result
}

//添加一个新的标签
// 如果重复将自动反馈存在标签的ID
//param db *sql.DB 数据库句柄
//param tagsTableName string 标签表
//param name string 标签名称
//return int 新的ID，-1则表示失败
func TagsAdd(db *sql.DB,tagsTableName string,name string) int {
	sql := "select `id` from `" + tagsTableName + "` where `name` = ?"
	row := db.QueryRow(sql,name)
	var id int
	if row.Scan(&id) != nil{
		return -1
	}
	return id

}

//删除标签
// 将自动删除和该标签关联的所有记录
//param db *sql.DB 数据库句柄
//param tagsTableName string 标签表
//param tagsBindTableName string 标签关联表
//param id int 主键
//return bool 是否成功
func TagsDelete(db *sql.DB,tagsTableName string,tagsBindTableName string,id int) bool {
	sql := "delete from `" + tagsBindTableName + "` where `tags_id` = ?"
	_,err := db.Exec(sql,id)
	if err != nil{
		return false
	}
	sql = "delete from `" + tagsTableName + "` where `id` = ?"
	_,err = db.Exec(sql,id)
	if err != nil{
		return false
	}
	return true
}

//将所有标签组添加到指定ID
//param db *sql.DB 数据库句柄
//param tagsTableName string 标签表
//param tagsBindTableName string 标签关联表
//param targetID int 目标ID
//param tags []string 标签组
//return bool 是否成功
func TagsBind(db *sql.DB,tagsTableName string,tagsBindTableName string,targetID int,tags []string) error {
	//调取现有的数据，做对比
	var nowTags map[int]string
	sql := "select `id`,`tag_id` from `" + tagsBindTableName + "` where `target_id` = ?"
	rows,err := db.Query(sql,targetID)
	rowsNum := 0
	if err == nil{
		for rows.Next(){
			var id int
			var tagID int
			err = rows.Scan(&id,&tagID)
			if err != nil{
				return err
			}
			nowTags[id] = string(tagID)
			rowsNum += 1
		}
	}
	//如果存在数据，但编辑删除了所有标签，则删除所有标签
	if rowsNum < 1 && len(tags) < 1{
		sql = "delete from `" + tagsBindTableName + "` where `target_id` = ?"
		_,err = db.Exec(sql,targetID)
		return err
	}
	//遍历修改后的标签组，查看是否有不一致的地方
	// 添加新的标签
	sql = "insert into `" + tagsBindTableName + "`(`id`,`target_id`,`tag_id`) values(null,?,?)"
	sqlSelect := "select `id` from `" + tagsTableName + "` where `tag_name` = ?"
	for _,tagName := range tags{
		var exist bool = false
		//遍历现在存在的标签，跳过存在的内容
		for _,nowTagName := range nowTags{
			//如果发现相同内容，跳过
			if nowTagName == tagName{
				exist = true
				break
			}
		}
		//完成遍历后，剩下的都是需要添加的
		if exist == false{
			var tagID int
			row := db.QueryRow(sqlSelect,tagName)
			err = row.Scan(tagID)
			if err != nil{
				continue
			}
			_,err = db.Exec(sql,targetID,tagID)
			if err != nil{
				continue
			}
		}
	}
	//遍历当前标签组，删除不存在的
	sql = "delete from `" + tagsBindTableName + "` where `id` = ?"
	for nowID,nowName := range tags{
		var exist bool = false
		for _,tagName := range nowTags{
			if nowName == tagName{
				exist = true
			}
		}
		if exist == false{
			_,err = db.Exec(sql,nowID)
			if err != nil{
				continue
			}
		}
	}
	//返回成功
	return nil
}