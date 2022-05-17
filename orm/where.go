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

// Append  添加条件
func (p *DbWhere) String(w, v string) {
	if v != "" {
		p.Append(w, v)
	}
}

func (p *DbWhere) Int(w string, v int) {
	if v > 0 {
		p.Append(w, v)
	}
}

// Append  添加条件
func (p *DbWhere) Date(w string, date, e string) {
	if date != "" {
		p.Append(w, date+e)
	}
}
