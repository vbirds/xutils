// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package orm

import (
	"reflect"

	"gorm.io/gorm"
)

var gOrmDb *gorm.DB

// H 多列处理
type H map[string]interface{}

// SetDB gorm对象
func SetDB(db *gorm.DB) {
	gOrmDb = db
}

// SetDB gorm对象
func DB() *gorm.DB {
	return gOrmDb
}

// DbCount 数目
func DbCount(model, where interface{}) int64 {
	var count int64
	db := gOrmDb.Model(model)
	if where != nil {
		db = db.Where(where)
	}
	db.Count(&count)
	return count
}

// DbCreate 创建
func DbCreate(model interface{}) error {
	return gOrmDb.Create(model).Error
}

// DbSave 保存
func DbSave(value interface{}) error {
	return gOrmDb.Save(value).Error
}

// DbUpdateModel 更新
func DbUpdateModel(model interface{}) error {
	return gOrmDb.Model(model).Updates(model).Error
}

// DbUpdateModelBy 条件更新
func DbUpdateModelBy(model interface{}, where string, args ...interface{}) error {
	return gOrmDb.Where(where, args...).Updates(model).Error
}

// DbUpdatesById 更新
func DbUpdateById(model, id interface{}) error {
	return gOrmDb.Where("id = ?", id).Updates(model).Error
}

// DbUpdateColById 单列更新
func DbUpdateColById(model, id interface{}, column string, value interface{}) error {
	return gOrmDb.Model(model).Where("id = ?", id).Update(column, value).Error
}

func DbUpdateColBy(model interface{}, column string, value interface{}, where string, args ...interface{}) error {
	return gOrmDb.Model(model).Where(where, args...).Update(column, value).Error
}

// DbUpdateColsBy 更新多列
// 用于0不更新
func DbUpdateColsBy(model interface{}, value map[string]interface{}, where string, args ...interface{}) error {
	return gOrmDb.Model(model).Where(where, args...).Updates(value).Error
}

// 自定义字段，用于0不更新
func DbUpdates(model interface{}, cols ...string) error {
	return gOrmDb.Model(model).Select(cols).Updates(model).Error
}

// 自定义字段，用于0不更新
func DbUpdatesBy(model interface{}, cols []string, where string, args ...interface{}) error {
	return gOrmDb.Model(model).Select(cols).Where(where, args...).Updates(model).Error
}

// DbUpdateByIds 批量更新
// ids id数组
func DbUpdateByIds(model interface{}, ids []int, value map[string]interface{}) error {
	return gOrmDb.Model(model).Where("id in (?)", ids).Updates(value).Error
}

// DbDeletes 批量删除
func DbDeletes(value interface{}) error {
	return gOrmDb.Delete(value).Error
}

// DbDeleteByIds 批量删除
// ids id数组 []
func DbDeleteByIds(model, ids interface{}) error {
	return gOrmDb.Delete(model, ids).Error
}

// DbDeleteBy 删除
func DbDeleteBy(model interface{}, where string, args ...interface{}) (count int64, err error) {
	db := gOrmDb.Where(where, args...).Delete(model)
	err = db.Error
	if err != nil {
		return
	}
	count = db.RowsAffected
	return
}

// DbFirstBy 指定条件查找
func DbFirstBy(out interface{}, where string, args ...interface{}) (err error) {
	err = gOrmDb.Where(where, args...).First(out).Error
	return
}

// DbFirstById 查找
func DbFirstById(out interface{}, id uint) error {
	return gOrmDb.First(out, id).Error
}

// DbFirstWhere 查找
func DbFirstWhere(out, where interface{}) error {
	return gOrmDb.Where(where).First(out).Error
}

// DbFind 多个查找
func DbFind(out interface{}, orders ...string) error {
	db := gOrmDb
	if len(orders) > 0 {
		for _, order := range orders {
			db = db.Order(order)
		}
	}
	return db.Find(out).Error
}

// DbFindBy 多个条件查找
func DbFindBy(out interface{}, where string, args ...interface{}) (int64, error) {
	db := gOrmDb.Where(where, args...).Find(out)
	return db.RowsAffected, db.Error
}

type dbByWhere struct {
	db    *gorm.DB
	total int64
}

func (o *dbByWhere) Find(out interface{}, conds ...interface{}) (int64, error) {
	if o.total < 1 {
		return 0, nil
	}
	return o.total, o.db.Find(out, conds...).Error
}

func (o *dbByWhere) Scan(out interface{}) (int64, error) {
	if o.total < 1 {
		return 0, nil
	}
	return o.total, o.db.Scan(out).Error
}

// Preload 关联加载
func (o *dbByWhere) Preload(preloads ...string) *dbByWhere {
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
func (o *dbByWhere) Joins(query string, args ...interface{}) *dbByWhere {
	o.db = o.db.Joins(query, args...)
	return o
}

// DbByWhere
func DbByWhere(m interface{}, w *DbWhere) *dbByWhere {
	db := gOrmDb.Model(m)
	if w != nil {
		for _, wo := range w.Wheres {
			if wo.Where != "" {
				db = db.Where(wo.Where, wo.Value...)
			}
		}
		if len(w.Orders) > 0 {
			for _, order := range w.Orders {
				db = db.Order(order)
			}
		}
	}
	o := &dbByWhere{db: db}
	if db.Count(&o.total).Error == nil {
		// dbByWhere 分页
		if w.Page != nil && w.Page.Num > 0 {
			o.db = db.Offset((w.Page.Num - 1) * w.Page.Size).Limit(w.Page.Size)
		}
	}
	return o
}

func DbTableByWhere(table string, w *DbWhere) *dbByWhere {
	db := gOrmDb.Table(table)
	if w != nil {
		for _, wo := range w.Wheres {
			if wo.Where != "" {
				db = db.Where(wo.Where, wo.Value...)
			}
		}
		if len(w.Orders) > 0 {
			for _, order := range w.Orders {
				db = db.Order(order)
			}
		}
	}
	o := &dbByWhere{db: db}
	if db.Count(&o.total).Error == nil {
		// dbByWhere 分页
		if w.Page != nil && w.Page.Num > 0 {
			o.db = db.Offset((w.Page.Num - 1) * w.Page.Size).Limit(w.Page.Size)
		}
	}
	return o
}

// DbPageRawScan obj必须是数组类型
func DbPageRawScan(query string, obj interface{}, page, size int) (int64, error) {
	s := reflect.ValueOf(obj)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	switch s.Kind() {
	case reflect.Slice, reflect.Array:
	default:
		return 0, nil
	}
	if err := gOrmDb.Raw(query).Scan(obj).Error; err != nil {
		return 0, err
	}
	total := s.Len()
	start := size * (page - 1)
	end := start + size
	if end >= total {
		end = total
	}
	s.Set(s.Slice(start, end))
	return int64(total), nil
}
