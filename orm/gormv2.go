package orm

import (
	"reflect"

	"gorm.io/gorm"
)

var _db *gorm.DB

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

// H 多列处理
type H map[string]interface{}

// Append  添加条件
func (p *DbWhere) Append(where string, value ...interface{}) {
	var w whereValues
	w.Where = where
	w.Value = value
	p.Wheres = append(p.Wheres, w)
}

// SetDB gorm对象
func SetDB(db *gorm.DB) {
	_db = db
}

// SetDB gorm对象
func DB() *gorm.DB {
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

// DbUpdateColsById 更新多列
// 用于0不更新 fix
// func DbUpdateColsById(model, id interface{}, value H) error {
// 	return _db.Model(model).Where("id = ?", id).Updates(value).Error
// }

// DbUpdateColsBy 更新多列
// 用于0不更新
func DbUpdateColsBy(model interface{}, value map[string]interface{}, where string, args ...interface{}) error {
	return _db.Model(model).Where(where, args...).Updates(value).Error
}

// 自定义字段，用于0不更新
func DbUpdates(model interface{}, args []string) error {
	return _db.Model(model).Select(args).Updates(model).Error
}

// 自定义字段，用于0不更新
func DbUpdateSelect(model interface{}, args ...string) error {
	return DbUpdates(model, args)
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
func DbFirstById(out interface{}, id uint64) error {
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

// DbPage 分页
type dbPage struct {
	db    *gorm.DB
	total int64
}

// Find 分页
func (o *dbPage) Find(page, size int, out interface{}, conds ...interface{}) (int64, error) {
	if o.total <= 0 {
		return 0, nil
	}
	if page > 0 {
		return o.total, o.db.Offset((page-1)*size).Limit(size).Find(out, conds...).Error
	}
	return o.total, o.db.Find(out, conds...).Error
}

// Find 分页
func (o *dbPage) Scan(page, size int, out interface{}) (int64, error) {
	if o.total <= 0 {
		return 0, nil
	}
	if page > 0 {
		return o.total, o.db.Offset((page - 1) * size).Limit(size).Scan(out).Error
	}
	return o.total, o.db.Scan(out).Error
}

// Preload 关联加载
func (o *dbPage) Preload(preloads ...string) *dbPage {
	if len(preloads) > 0 {
		for _, preload := range preloads {
			o.db = o.db.Preload(preload)
		}
	}
	return o
}

// Joins join
func (o *dbPage) Joins(query string, args ...interface{}) *dbPage {
	o.db = o.db.Joins(query, args...)
	return o
}

// DbPage
func DbPage(model interface{}, where *DbWhere) *dbPage {
	db := _db.Model(model)
	if where != nil {
		for _, wo := range where.Wheres {
			if wo.Where != "" {
				db = db.Where(wo.Where, wo.Value...)
			}
		}
		if len(where.Orders) > 0 {
			for _, order := range where.Orders {
				db = db.Order(order)
			}
		}
	}
	var total int64
	if db.Count(&total).Error != nil {
		total = 0
	}
	return &dbPage{db: db, total: total}
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
