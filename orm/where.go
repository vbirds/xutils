// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package orm

import (
	"strconv"

	"gorm.io/gorm"
)

// wValues 分页条件
type wValue struct {
	Where string
	Value []interface{}
}

type DbPage struct {
	Num  int `form:"pageNum"`  // 当前页码
	Size int `form:"pageSize"` // 每页数
}

// DbWhere 搜索条件
type DbWhere struct {
	db     *gorm.DB
	total  int64
	page   *DbPage
	wheres []wValue
	Orders []string
}

func (o *DbPage) DbWhere() *DbWhere {
	return &DbWhere{page: o}
}

// Where  添加条件
func (o *DbWhere) Where(w string, v ...interface{}) *DbWhere {
	if v != nil {
		o.wheres = append(o.wheres, wValue{Where: w, Value: v})
	}
	return o
}

// Equal
func (o *DbWhere) Equal(field string, v interface{}) *DbWhere {
	switch v := v.(type) {
	case string:
		if v == "" {
			return o
		}
	}
	return o.Where(field+" = ?", v)
}

// EqualNumber
func (o *DbWhere) EqualNumber(field, v string) *DbWhere {
	if v != "" {
		n, _ := strconv.Atoi(v)
		o.Equal(field, n)
	}
	return o
}

// Like
func (o *DbWhere) Like(field, v string) *DbWhere {
	if v != "" {
		// o.Where("INSTR("+field+", ?)>0", v)
		o.Where(field+" LIKE ?", "%"+v+"%")
	}
	return o
}

// DateRange
func (o *DbWhere) TimeRange(field string, st, et string) *DbWhere {
	if st != "" && et != "" {
		o.Where(field+" BETWEEN ? AND ?", st, et)
	}
	return o
}

// DateRange
func (o *DbWhere) DateRange(field string, r []string) *DbWhere {
	if r != nil {
		o.Where(field+" BETWEEN ? AND ?", r[0]+" 00:00:00", r[1]+" 23:59:59")
	}
	return o
}

func (o *DbWhere) Find(out interface{}, conds ...interface{}) (int64, error) {
	if o.total < 1 {
		return 0, nil
	}
	return o.total, o.db.Find(out, conds...).Error
}

func (o *DbWhere) Scan(out interface{}) (int64, error) {
	if o.total < 1 {
		return 0, nil
	}
	return o.total, o.db.Scan(out).Error
}

// Preload 关联加载
func (o *DbWhere) Preload(preloads ...string) *DbWhere {
	if o.total < 1 {
		return o
	}
	if len(preloads) > 0 {
		for _, preload := range preloads {
			o.db = o.db.Preload(preload)
		}
	}
	return o
}

// Joins join
func (o *DbWhere) Joins(query string, args ...interface{}) *DbWhere {
	o.db = o.db.Joins(query, args...)
	return o
}

// DbByWhere
func (o *DbWhere) Model(m interface{}) *DbWhere {
	db := Model(m)
	for _, v := range o.wheres {
		db = db.Where(v.Where, v.Value...)
	}
	for _, order := range o.Orders {
		db = db.Order(order)
	}
	if db.Count(&o.total).Error == nil {
		// dbByWhere 分页
		if o.page != nil && o.page.Num > 0 {
			db = db.Offset((o.page.Num - 1) * o.page.Size).Limit(o.page.Size)
		}
	}
	o.db = db
	return o
}
