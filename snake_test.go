package snake

import (
	"fmt"
	"testing"
)

func TestPageAnd(t *testing.T) {
	var userId = 1
	var status = 2
	var keyword = "nice"
	var pageSize int64 = 10
	var pageNumber int64 = 2
	var m = make(map[string]interface{})
	m["t2.user_id"] = userId
	if status != 0 {
		m["t2.status"] = status
	}
	m["t1.title like"] = "%" + keyword + "%"
	var querySnake = NewQuerySnake().Field("count(1)").
		Table("table1 t1").
		LeftJoin("table2 "+" t2 on t2.id = t1.order_id").
		Where(AND, m).Order("t1.id desc").
		Limit(pageSize).
		Offset(pageSize * (pageNumber - 1)).
		BuildSql()
	var sql = querySnake.GetSql()
	var sqlParams = querySnake.GetSqlParams()
	fmt.Println(sql)
	for i := 0; i < len(sqlParams); i++ {
		fmt.Println(sqlParams[i])
	}
	/**
	select count(1) from table1 t1 left join table2  t2 on t2.id = t1.order_id where t1.title like ? and t2.user_id=? and t2.status=? order by t1.id desc limit 10 , 10
	%nice%
	1
	2
	*/

}

func TestPageOr(t *testing.T) {

	var status = 2
	var pageSize int64 = 10
	var pageNumber int64 = 2

	var m = make(map[string]interface{})
	m["t2.user_id"] = 12
	if status != 0 {
		m["t2.status"] = status
	}

	var querySnake = NewQuerySnake().Field("count(1)").
		Table("table1 t1").
		LeftJoin("table2 "+" t2 on t2.id = t1.order_id").
		Where(OR, m).Order("t1.id desc").
		Limit(pageSize).
		Offset(pageSize * (pageNumber - 1)).
		BuildSql()
	var sql = querySnake.GetSql()
	var sqlParams = querySnake.GetSqlParams()
	fmt.Println(sql)
	for i := 0; i < len(sqlParams); i++ {
		fmt.Println(sqlParams[i])
	}
	/**
	select count(1) from table1 t1 left join table2  t2 on t2.id = t1.order_id where t2.user_id=? or t2.status=? order by t1.id desc limit 10 , 10
	12
	2
	*/

}
func TestPageOrAndOr(t *testing.T) {

	var status = 2

	var pageSize int64 = 10
	var pageNumber int64 = 2
	var ms []map[string]interface{}
	for i := 0; i < 3; i++ {
		var m = make(map[string]interface{})
		m["t2.user_id"] = i
		if status != 0 {
			m["t2.status"] = status
		}
		ms = append(ms, m)
	}
	var querySnake = NewQuerySnake().Field("count(1)").
		Table("table1 t1").
		LeftJoin("table2 "+" t2 on t2.id = t1.order_id").
		Where(OR_AND_OR, ms...).Order("t1.id desc").
		Limit(pageSize).
		Offset(pageSize * (pageNumber - 1)).
		BuildSql()
	var sql = querySnake.GetSql()
	var sqlParams = querySnake.GetSqlParams()
	fmt.Println(sql)
	for i := 0; i < len(sqlParams); i++ {
		fmt.Println(sqlParams[i])
	}
	/**
	SELECT
		count( 1 )
	FROM
		table1 t1
		LEFT JOIN table2 t2 ON t2.id = t1.order_id
	WHERE
		( t2.user_id =? OR t2.STATUS =? )
		AND ( t2.user_id =? OR t2.STATUS =? )
		AND ( t2.user_id =? OR t2.STATUS =? )
	ORDER BY
		t1.id DESC
		LIMIT 10,
		10
	0
	2
	1
	2
	2
	2
	*/

}

func TestPageAndOrOr(t *testing.T) {

	var status = 2

	var pageSize int64 = 10
	var pageNumber int64 = 2
	var ms []map[string]interface{}
	for i := 0; i < 2; i++ {
		var m = make(map[string]interface{})
		m["t2.user_id"] = i
		if status != 0 {
			m["t2.status"] = status
		}
		ms = append(ms, m)
	}
	var querySnake = NewQuerySnake().Field("count(1)").
		Table("table1 t1").
		LeftJoin("table2 "+" t2 on t2.id = t1.order_id").
		Where(AND_OR_OR, ms...).Order("t1.id desc").
		Limit(pageSize).
		Offset(pageSize * (pageNumber - 1)).
		BuildSql()
	var sql = querySnake.GetSql()
	var sqlParams = querySnake.GetSqlParams()
	fmt.Println(sql)
	for i := 0; i < len(sqlParams); i++ {
		fmt.Println(sqlParams[i])
	}
	/**
	SELECT
		count( 1 )
	FROM
		table1 t1
		LEFT JOIN table2 t2 ON t2.id = t1.order_id
	WHERE
		( t2.user_id =? OR t2.STATUS =? )
		AND ( t2.user_id =? OR t2.STATUS =? )
		AND ( t2.user_id =? OR t2.STATUS =? )
	ORDER BY
		t1.id DESC
		LIMIT 10,
		10
	0
	2
	1
	2
	2
	2
	*/

}
