// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package orm

// wValues 分页条件
type wValue struct {
	Where string
	Value []interface{}
}

type DbPage struct {
	Num  int `form:"pageNum"`  // 当前页码
	Size int `form:"pageSize"` // 每页数
}

// DbWhere 搜索条件条件
type DbWhere struct {
	Wheres []wValue
	Orders []string
	Page   *DbPage
}

func (o *DbPage) DbWhere() *DbWhere {
	return &DbWhere{Page: o}
}

// Append  添加条件
func (p *DbWhere) Append(w string, v ...interface{}) {
	if v == nil {
		return
	}
	p.Wheres = append(p.Wheres, wValue{Where: w, Value: v})
}

// Like
func (p *DbWhere) Like(field, v string) {
	if v == "" {
		return
	}
	p.Append(field+" like ?", v)
}

// Equal
func (p *DbWhere) Equal(field string, v interface{}) {
	switch v := v.(type) {
	case int:
		if v == 0 {
			return
		}
	case string:
		if v == "" {
			return
		}
	case *int:
		if v == nil {
			return
		}
	case *string:
		if v == nil {
			return
		}
	default:
		return
	}
	p.Append(field+" = ?", v)
}

// DateRange
func (p *DbWhere) DateRange(field string, st, et string) {
	if st == "" || et == "" {
		return
	}
	p.Append(field+" >= ? AND "+field+" <= ?", st+" 00:00:00", et+" 23:59:59")
}

// TimeRange
func (p *DbWhere) TimeRange(field string, st, et string) {
	if st == "" || et == "" {
		return
	}
	p.Append(field+" >= ? AND "+field+" <= ?", st, et)
}
