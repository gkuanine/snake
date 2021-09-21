package snake

type DeleteSnake struct {
	table      string
	where      map[string]interface{}
	lastSql    string
	lastParams []interface{}
}

func NewDeleteSnake() *DeleteSnake {
	m := DeleteSnake{}
	return &m
}

func (p *DeleteSnake) Table(table string) *DeleteSnake {
	p.table = table
	return p
}
func (p *DeleteSnake) Where(where map[string]interface{}) *DeleteSnake {
	p.where = where
	return p
}

func (p *DeleteSnake) BuildSql() *DeleteSnake {
	if p.table == "" {
		panic("table error ")
	}

	var sql = "delete from  " +
		p.table
		//p.field +

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
	} else {
		panic("where cannot null")
	}

	p.lastSql = sql

	return p
}
