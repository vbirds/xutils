package orm

// whereValues 分页条件
type whereValues struct {
	Where string
	Value []interface{}
}

// ModelWhere 分页条件
type DbWhere struct {
	Wheres []whereValues
	Orders []string
}

// Append  添加条件
func (p *DbWhere) Append(where string, value ...interface{}) {
	var w whereValues
	w.Where = where
	w.Value = value
	p.Wheres = append(p.Wheres, w)
}

// Append  添加条件
func (p *DbWhere) String(where string, value string) {
	if value != "" {
		p.Append(where, value)
	}
}

func (p *DbWhere) Int(where string, value int) {
	if value > 0 {
		p.Append(where, value)
	}
}

// Append  添加条件
func (p *DbWhere) Date(where string, date, e string) {
	if date != "" {
		p.Append(where, date+e)
	}
}
