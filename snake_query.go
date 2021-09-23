package snake

/**
  构建查询 sql
*/
import (
	"fmt"
	"reflect"
	"strings"
)

type NullType byte
type WhereType byte

const (
	_ NullType = iota
	// IsNull the same as `is null`
	IsNull
	// IsNotNull the same as `is not null`
	IsNotNull
	Linker_AND WhereType = iota
	Linker_OR
	Linker_or_AND_or
	Linker_and_AND_or
)

type QuerySnake struct {
	field      string
	table      string
	leftJoin   []string
	where      []map[string]interface{}
	whereType  WhereType
	groupBy    []string
	having     []string
	order      []string
	limit      int64
	offset     int64
	lastSql    string
	lastParams []interface{}
}

func NewQuerySnake() *QuerySnake {
	m := QuerySnake{}
	return &m
}
func (p *QuerySnake) Select(fields ...string) *QuerySnake {
	p.field = strings.Join(fields, ",")
	return p
}

func (p *QuerySnake) Table(table string) *QuerySnake {
	p.table = table
	return p
}

func (p *QuerySnake) LeftJoin(leftJoin ...string) *QuerySnake {
	p.leftJoin = leftJoin
	return p
}

func (p *QuerySnake) Where(whereType WhereType, where ...map[string]interface{}) *QuerySnake {
	p.where = where
	p.whereType = whereType
	return p
}

func (p *QuerySnake) GroupBy(groupBy ...string) *QuerySnake {
	p.groupBy = groupBy
	return p
}
func (p *QuerySnake) Having(having ...string) *QuerySnake {
	p.having = having
	return p
}
func (p *QuerySnake) Order(order ...string) *QuerySnake {
	p.order = order
	return p
}
func (p *QuerySnake) Limit(limit int64) *QuerySnake {
	p.limit = limit
	return p
}
func (p *QuerySnake) Offset(offset int64) *QuerySnake {
	p.offset = offset
	return p
}
func (p *QuerySnake) GetSql() string {
	return p.lastSql
}
func (p *QuerySnake) GetSqlParams() []interface{} {
	return p.lastParams
}

//select join where、group by、having、order by
func (p *QuerySnake) BuildSql() *QuerySnake {
	var sql = "select " +
		p.field + " from " +
		p.table
	if len(p.leftJoin) > 0 {
		joinSql := ""
		for i := 0; i < len(p.leftJoin); i++ {
			joinSql = joinSql + " left join " + p.leftJoin[i]
		}
		sql = sql + joinSql
	}
	if len(p.where) > 0 {
		var whereSql string
		var vals []interface{}
		var err error
		switch p.whereType {
		case Linker_AND:
			whereSql, vals, err = p.buildAndWhere(p.where...)
		case Linker_OR:
			whereSql, vals, err = p.buildOrdWhere(p.where...)
		case Linker_or_AND_or:
			whereSql, vals, err = p.buildComplexWhere(p.where...)
		case Linker_and_AND_or:
			whereSql, vals, err = p.buildAndOrOrWhere(p.where...)

		}
		if err != nil {
			panic("whereSql error ")
		}
		p.lastParams = vals
		sql = sql + " where " + whereSql
	}
	if len(p.groupBy) > 0 {
		groupSql := ""
		for i := 0; i < len(p.groupBy); i++ {
			groupSql = groupSql + p.groupBy[i]
		}
		sql = sql + " group by " + groupSql
	}
	if len(p.having) > 0 {
		havingSql := ""
		for i := 0; i < len(p.having); i++ {
			havingSql = havingSql + p.having[i]
		}
		sql = sql + " having " + havingSql
	}
	if len(p.order) > 0 {
		orderSql := ""
		for i := 0; i < len(p.order); i++ {
			orderSql = orderSql + p.order[i]
		}
		sql = sql + " order by " + orderSql
	}
	if p.limit > 0 {
		sql = sql + " limit " + fmt.Sprint(p.limit)
		if p.offset > 0 {
			sql = sql + " , " + fmt.Sprint(p.offset)
		}
	}
	p.lastSql = sql

	return p
}

//if len(andWhere)>1
//  (or or or ) and (or or or ) ....
func (p *QuerySnake) buildComplexWhere(complexWhere ...map[string]interface{}) (whereSQL string, vals []interface{}, err error) {
	if len(complexWhere) == 1 {
		panic(" too simple to build whereSql")
	}
	var valss []interface{}
	var whereSQLs []string
	for i := 0; i < len(complexWhere); i++ {
		var orWhere = complexWhere[i]
		var whereSql = ""
		whereSql, vals, err = p.buildOrdWhere(orWhere)
		if err != nil {
			panic("buildComplexWhere error ")
		}
		whereSQLs = append(whereSQLs, "("+whereSql+")")
		valss = append(valss, vals...)
	}
	whereSQL = strings.Join(whereSQLs, " and ")
	return whereSQL, valss, err

}

//if len(andWhere)==1
//and and and
func (p *QuerySnake) buildAndWhere(andWhere ...map[string]interface{}) (whereSQL string, vals []interface{}, err error) {
	if len(andWhere) > 1 {
		panic("whereSql error ")
	}
	return p.buildAndOrWhere(andWhere[0], "and")

}

//if len(andWhere)==1
// or or or
func (p *QuerySnake) buildOrdWhere(andWhere ...map[string]interface{}) (whereSQL string, vals []interface{}, err error) {
	if len(andWhere) > 1 {
		panic("whereSql error ")
	}
	return p.buildAndOrWhere(andWhere[0], "or")
}

//where (t2.user_id=? and t2.status=?) AND (t2.status=? or t2.user_id=?) o
func (p *QuerySnake) buildAndOrOrWhere(andWhere ...map[string]interface{}) (whereSQL string, vals []interface{}, err error) {
	if len(andWhere) != 2 {
		panic("buildAndOrOrWhere error ")
	}
	//var andWhereSql,andVals,err= p.buildAndOrWhere(andWhere[0],"or")
	whereSQL, vals, err = p.buildAndWhere(andWhere[0])
	if err != nil {
		panic("buildAndOrOrWhere error")
	}
	whereSQL2, vals2, err2 := p.buildOrdWhere(andWhere[1])
	if err2 != nil {
		panic("buildAndOrOrWhere error")
	}
	whereSQL = whereSQL + "  AND " + whereSQL2
	vals = append(vals, vals2...)
	return whereSQL, vals, nil

}

// sql build where
func (p *QuerySnake) buildAndOrWhere(where map[string]interface{}, andOr string) (whereSQL string, vals []interface{}, err error) {
	for k, v := range where {
		k = strings.Trim(k, " ")
		ks := strings.Split(k, " ")
		if len(ks) > 2 {
			return "", nil, fmt.Errorf("Error in query condition: %s. ", k)
		}
		if whereSQL != "" {
			whereSQL += " " + andOr + " "
		}
		strings.Join(ks, ",")
		switch len(ks) {
		case 1:
			fmt.Println(reflect.TypeOf(v))
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					whereSQL += fmt.Sprint(k, " IS NOT NULL")
				} else {
					whereSQL += fmt.Sprint(k, " IS NULL")
				}
			default:
				t := reflect.TypeOf(v)
				reflectType := fmt.Sprint(t)
				switch reflectType {
				case "string":
					whereSQL += fmt.Sprint(k, "=?")
					vals = append(vals, v)
				case "uint", "uint8", "uint16", "uint32", "uint64", "int", "int8", "int16", "int32", "int64", "float32", "float64", "complex64", "complex128":
					whereSQL += fmt.Sprint(k, "=?")
					vals = append(vals, v)
				case "decimal.Decimal":
					whereSQL += fmt.Sprint(k, "=?")
					vals = append(vals, v)
				}
			}
			break
		case 2:
			k = ks[0]
			switch ks[1] {
			case "=":
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
				break
			case ">":
				whereSQL += fmt.Sprint(k, ">?")
				vals = append(vals, v)
				break
			case ">=":
				whereSQL += fmt.Sprint(k, ">=?")
				vals = append(vals, v)
				break
			case "<":
				whereSQL += fmt.Sprint(k, "<?")
				vals = append(vals, v)
				break
			case "<=":
				whereSQL += fmt.Sprint(k, "<=?")
				vals = append(vals, v)
				break
			case "!=":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "<>":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "in":
				whereSQL += fmt.Sprint(k, " in (?)")
				vals = append(vals, v)
				break
			case "like":
				whereSQL += fmt.Sprint(k, " like ?")
				vals = append(vals, v)
			}
			break
		}
	}
	whereSQL = "(" + whereSQL + ")"
	return
}
