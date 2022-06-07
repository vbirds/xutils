// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package orm

import (
	"errors"
	"reflect"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var _db *gorm.DB

// H 多列处理
type H map[string]interface{}

// SetDB gorm对象
func SetDB(db *gorm.DB) {
	_db = db
}

// SetDB gorm对象
func DB() *gorm.DB {
	return _db
}

func Model(v interface{}) *gorm.DB {
	if m, ok := v.(XTablers); ok {
		return _db.Table(m.TableName())
	}
	return _db.Model(v)
}

func Table(v interface{}) *gorm.DB {
	if m, ok := v.(XTablers); ok {
		return _db.Table(m.TableName())
	}
	return _db
}

// DbCount 数目
func DbCount(model, where interface{}) int64 {
	var count int64
	db := _db.Model(model)
	if where != nil {
		db = db.Where(where)
	}
	db.Count(&count)
	return count
}

// DbCreate 创建
func DbCreate(model interface{}) error {
	return _db.Create(model).Error
}

// DbSave 保存
func DbSave(value interface{}) error {
	return _db.Save(value).Error
}

// DbUpdateModel 更新
func DbUpdateModel(model interface{}) error {
	return _db.Model(model).Updates(model).Error
}

// DbUpdateModelBy 条件更新
func DbUpdateModelBy(model interface{}, where string, args ...interface{}) error {
	return _db.Where(where, args...).Updates(model).Error
}

// DbUpdatesById 更新
func DbUpdateById(model, id interface{}) error {
	return _db.Where("id = ?", id).Updates(model).Error
}

// DbUpdateColById 单列更新
func DbUpdateColById(model, id interface{}, column string, value interface{}) error {
	return _db.Model(model).Where("id = ?", id).Update(column, value).Error
}

func DbUpdateColBy(model interface{}, column string, value interface{}, where string, args ...interface{}) error {
	return _db.Model(model).Where(where, args...).Update(column, value).Error
}

// DbUpdateColsBy 更新多列
// 用于0不更新
func DbUpdateColsBy(model interface{}, value map[string]interface{}, where string, args ...interface{}) error {
	return _db.Model(model).Where(where, args...).Updates(value).Error
}

// 自定义字段，用于0不更新
func DbUpdates(model interface{}, cols ...string) error {
	return _db.Model(model).Select(cols).Updates(model).Error
}

// 自定义字段，用于0不更新
func DbUpdatesBy(model interface{}, cols []string, where string, args ...interface{}) error {
	return _db.Model(model).Select(cols).Where(where, args...).Updates(model).Error
}

// DbUpdateByIds 批量更新
// ids id数组
func DbUpdateByIds(model interface{}, ids []int, value map[string]interface{}) error {
	return _db.Model(model).Where("id in (?)", ids).Updates(value).Error
}

// DbDeletes 批量删除
func DbDeletes(value interface{}) error {
	return _db.Delete(value).Error
}

// DbDeleteByIds 批量删除
// ids id数组 []
func DbDeleteByIds(model, ids interface{}) error {
	return _db.Delete(model, ids).Error
}

// DbDeleteBy 删除
func DbDeleteBy(model interface{}, where string, args ...interface{}) (count int64, err error) {
	db := _db.Where(where, args...).Delete(model)
	err = db.Error
	if err != nil {
		return
	}
	count = db.RowsAffected
	return
}

// DbFirstBy 指定条件查找
func DbFirstBy(out interface{}, where string, args ...interface{}) (err error) {
	err = _db.Where(where, args...).First(out).Error
	return
}

// DbFirstById 查找
func DbFirstById(out interface{}, id uint) error {
	return _db.First(out, id).Error
}

// DbFirstWhere 查找
func DbFirstWhere(out, where interface{}) error {
	return _db.Where(where).First(out).Error
}

// DbFind 多个查找
func DbFind(out interface{}, orders ...string) error {
	db := _db
	if len(orders) > 0 {
		for _, order := range orders {
			db = db.Order(order)
		}
	}
	return db.Find(out).Error
}

// DbFindBy 多个条件查找
func DbFindBy(out interface{}, where string, args ...interface{}) (int64, error) {
	db := _db.Where(where, args...).Find(out)
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
	db := Model(m)
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
	if err := _db.Raw(query).Scan(obj).Error; err != nil {
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

var gconf = gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	},
	DisableForeignKeyConstraintWhenMigrating: true, // 禁用外键约束
}

func NewGormV2(name, address string) (db *gorm.DB, err error) {
	switch name {
	case "mysql":
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN: address,
			// DefaultStringSize:         64,    // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据版本自动配置
		}), &gconf)
	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(address), &gconf)
	case "postgresql":
		db, err = gorm.Open(postgres.Open(address), &gconf)
	default:
	}
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, errors.New("db invalid")
	}
	sqldb, err := db.DB()
	if err != nil {
		return nil, err
	}
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqldb.SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqldb.SetMaxOpenConns(100)
	_db = db
	return db, nil
}
