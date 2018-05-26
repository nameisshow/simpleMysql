package main

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"errors"
)

type SimpleMysql struct {
	dataSourceName string
	prefix         string
	Err            error
	RowsAffected   int64
	LastInsertId   int64
	db             *sql.DB
	// stm            *sql.Stmt

	// sql预备
	table  string
	field  string
	join   string
	group  string
	order  string
	having string
	index  int
	num    int

	or bool

	where interface{}

	data map[string]interface{}

	sql string

	isFind   bool
	isGet    bool
	isInsert bool
	isUpdate bool
}

type sqlString struct{
	// sql预备
	table  string
	field  string
	join   string
	group  string
	order  string
	having string
	index  int
	num    int

	or bool

	where interface{}

	data map[string]interface{}

	sql string
}

func (simple *SimpleMysql) setError(err error) {
	simple.Err = err
}

func (simple *SimpleMysql) getError() {
	log.Fatal(simple.Err)
}

func (simple *SimpleMysql) open(dataSourceName string) {
	db, err := sql.Open("mysql", dataSourceName)
	simple.setError(err)
	simple.db = db
}

func (simple *SimpleMysql) buildSql() {
	if simple.table == "" {
		simple.setError(errors.New("table is nil"))
		return
	}


}

/*
	返回操作类型
 */
func (simple *SimpleMysql) selectType() string {
	t := true
	switch t {
	case simple.isFind:
		return "find"
	case simple.isGet:
		return "get"
	case simple.isInsert:
		return "insert"
	case simple.isUpdate:
		return "update"
	default:
		return ""
	}
}

/*********操作方法部分*********/

/*
	设置表
 */
func (simple *SimpleMysql) Table(table string) *SimpleMysql {
	simple.table = table
	return simple
}

/*
	设置查询字段
 */
func (simple *SimpleMysql) Field(field string) *SimpleMysql {
	simple.field = field
	return simple
}

/*
	设置查询条件
 */
func (simple *SimpleMysql) Where(where interface{}) *SimpleMysql {
	simple.where = where
	return simple
}

/*
	设置条件or
 */
func (simple *SimpleMysql) Or() *SimpleMysql {
	simple.or = true
	return simple
}

/*
	设置join
 */
func (simple *SimpleMysql) Join(join string) *SimpleMysql {
	simple.join = join
	return simple
}

/*
	设置order
 */
func (simple *SimpleMysql) Order(order string) *SimpleMysql {
	simple.order = order
	return simple
}

/*
	设置group
 */
func (simple *SimpleMysql) Group(group string) *SimpleMysql {
	simple.group = group
	return simple
}

/*
	设置having
 */
func (simple *SimpleMysql) Having(having string) *SimpleMysql {
	simple.having = having
	return simple
}

/*
	设置limit
 */
func (simple *SimpleMysql) Limit(index int, num int) *SimpleMysql {
	simple.index = index
	simple.num = num
	return simple
}

/*
	find返回单条数据   map
 */
func (simple *SimpleMysql) Find() map[string]interface{} {
	simple.isFind = true
	simple.isGet = false
	simple.isInsert = false
	simple.isUpdate = false
	simple.buildSql()
}

/*
	get返回多条数据   切片-map
 */
func (simple *SimpleMysql) Get() []map[string]interface{} {

}

/*
	insert单条
 */
func (simple *SimpleMysql) Insert(data map[string]interface{}) int64 {

}

/*
	update单条
 */
func (simple *SimpleMysql) Update(data map[string]interface{}) int64 {

}

func (simple *SimpleMysql) Insertold(table string, data map[string]interface{}) int64 {

	// 字段字字符串
	keyString := ""
	// 字段个数
	valCount := len(data)
	// prepare传递的参数切片
	valSlice := make([]interface{}, 0, valCount)

	for k, v := range data {
		keyString += "`" + k + "`" + ","
		valSlice = append(valSlice, v)

	}

	// 构建sql语句
	sql := "INSERT INTO `" + table + "`(" + substr(keyString, 0, -1) + ")" + " VALUES("
	for i := 0; i < valCount; i++ {
		sql += "?, "
	}
	sql = substr(sql, 0, -2) + ")"

	stem, _ := simple.db.Prepare(sql)
	// 多参数传递的方法，再切片后面添加...
	res, err := stem.Exec(valSlice...)
	stem.Close()
	simple.setError(err)
	simple.RowsAffected, _ = res.RowsAffected()
	simple.LastInsertId, _ = res.LastInsertId()

	return simple.LastInsertId
}

func main() {
	dataSourceName := "root:root@tcp(localhost:3306)/dkb?charset=utf8"
	simple := SimpleMysql{}
	simple.open(dataSourceName)
	// fmt.Println(simple)

	data := make(map[string]interface{})
	data["configFlag"] = "new"
	data["configName"] = "新的"
	data["configValue"] = "newssss"
	data["configDesc"] = "asd"
	id := simple.Insert("p2_config", data)
	fmt.Println(id)
}

/**
	截取字符串
 */
func substr(str string, strat int, end int) string {
	str_s := []rune(str)
	length := len(str_s)

	// 如果end大于0，从头开始计算，如果end小于0，从结尾开始计算
	if (strat > length-1) || (strat < 0) {
		strat = 0
	}
	if end > length {
		end = length
	}
	if end < -length {
		end = -length
	}
	if end < 0 {
		end = length + end
	}

	// fmt.Println(len(str_s))
	return string(str_s[strat:end])
}
