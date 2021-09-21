package snake

type UpdateSnake struct {
	table                 string
	field                 map[string]interface{}
	where                 map[string]interface{}
	lastSql               string
	lastUpdateFieldParams []interface{}
	lastParams            []interface{}
}

func NewUpdateSnake() *UpdateSnake {
	m := UpdateSnake{}
	return &m
}

func (p *UpdateSnake) Table(table string) *UpdateSnake {
	p.table = table
	return p
}
func (p *UpdateSnake) Update(field map[string]interface{}) *UpdateSnake {
	p.field = field
	return p
}
func (p *UpdateSnake) Where(where map[string]interface{}) *UpdateSnake {
	p.where = where
	return p
}

func (p *UpdateSnake) BuildSql() *UpdateSnake {
	if p.table == "" {
		panic("table error ")
	}

	var sql = "update " +
		p.table + " set "
	//p.field +
	if len(p.field) > 0 {
		fieldSQL, vals, err := BuildUpdateField(p.field)
		if err != nil {
			panic(" fieldSql error ")
		}
		p.lastUpdateFieldParams = vals
		sql = sql + fieldSQL
	}

	if len(p.where) > 0 {
		var whereSql string
		var vals []interface{}
		var err error
		whereSql, vals, err = BuildAndOrWhere(p.where, "and")

		if err != nil {
			panic("whereSql error ")
		}
		p.lastParams = vals
		sql = sql + " where " + whereSql
		p.lastParams = append(p.lastUpdateFieldParams, p.lastParams...)
	} else {
		panic("where cannot null")
	}

	p.lastSql = sql

	return p
}
