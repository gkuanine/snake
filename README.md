##偶尔需要 写sql拼接查询感觉麻烦就 花了一小时写了个工具 针对有的参数需要工具情况动态拼接 问题 结合 gorm 使用
###拼接 sql 查询  需要结合、gorm的  model.Raw(sql,sqlParams).Scan() 或者 Count()
#####golang  mybatis plus  
```golang
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






and or 连用

 
		var andMap = make(map[string]interface{})
	andMap["t1.shop_id"] = admin.ShopId
	if skuId != 0 {
		andMap["t1.sku_id"] = skuId
	}
	if goodsId != 0 {
		andMap["t1.goods_id"] = goodsId
	}
	
	var orMap = make(map[string]interface{})
	orMap["t2.goods_sn"] = code
	orMap["t2.barcode"] = code
	
	var query= snake.NewQuerySnake().Table(global.MallBatch().TableName()+" t1").Select("*").
		LeftJoin(global.MallSku().TableName()+" t2 on t2.id = t1.sku_id").
		Where(snake.Linker_and_AND_or, andMap, orMap).Limit(pageSize).Offset(pageSize * (pageNumber - 1)).
		Order("t1.stock desc,t1.sku_id desc,t1.goods_id desc, t1.updated_at desc,t1.id desc").BuildSql()
			
var dataList []model.MallBatchDTO
	model.DB.Raw(query.GetSql(), query.GetSqlParams()...).Scan(&dataList)
	var total int64
	query.Select("count(1)").BuildSql()
	model.DB.Raw(query.GetSql(), query.GetSqlParams()...).Count(&total)




多left join 

	var admin adminmodel.AdminDTO
	var orMap = make(map[string]interface{})
	orMap["phone"] = username
	orMap["email"] = username
	var query = snake.NewQuerySnake().Select("t1.*,t3.name as role_name").
		Table(global.Admin().TableName()+" t1").
		LeftJoin(global.AdminUserRole().TableName()+" t2 on t2.admin_id = t1.id ", global.AdminRole().TableName()+" t3 on t3.id = 			t2.role_id").Where(snake.Linker_OR, orMap).BuildSql()
	model.DB.Raw(query.GetSql(), query.GetSqlParams()...).Find(&admin)
	if admin.ID == 0 {
		global.Error("登录帐号不存在") 
	}
	if admin.Status == 1 {
		global.Error("帐号禁用") 
	}
	
	
	
	
and or like 
	var andM = make(map[string]interface{})
	andM["t2.user_id"] = userId
	if status != 0 {
		andM["t2.status"] = status
	}
	var orM = make(map[string]interface{})
	if !strs.IsEmpty(keyword) {
		orM["t1.title like"] = "%" + keyword + "%"
		orM["t2.order_no like"] = "%" + keyword + "%"
	}
	
		var  querySnake = snake.NewQuerySnake().Select("count(1)").
			Table(global.MallOrderGoods().TableName()+" t1").
			LeftJoin(global.MallOrder().TableName()+" t2 on t2.id = t1.order_id").
			Where(snake.Linker_and_AND_or, andM, orM).Order("t2.updated_at desc,t1.id desc").
			Limit(pageSize).
			Offset(pageSize * (pageNumber - 1)).
			BuildSql()
			
```
