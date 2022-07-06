// Copyright 2021 xutils. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package orm

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ParentId uint
	Name     string
	DeptID   uint
	Children []User `json:"children,omitempty" gorm:"foreignKey:ParentId;"` // 这里注意，如果设置ParentId为0，要禁用外键约束
}

type Dog struct {
	gorm.Model
	Name   string
	GirlID uint
}

type Girl struct {
	gorm.Model
	Name string
	Dog  []Dog
}

func init() {
	_, err := NewGormV2("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln(err)
	}
}

func TestOrm(t *testing.T) {
	_db.AutoMigrate(&User{})
	user := &User{
		Name:     "test",
		ParentId: 1,
	}
	DbCreate(&user)
	var o *User
	_db.Model(&User{}).Where("id = ?", 1).First(o)
	log.Println(o)
}

type Param struct {
	DbPage
	Name   string `form:"name"`
	DeptID string `form:"deptID"`
}

// 条件组合查询
func TestWhere(t *testing.T) {
	p := Param{}
	w := p.DbWhere() // 获取分页
	w.Equal("name", p.Name)
	w.EqualNumber("dept_id", p.DeptID)
	var data []User
	w.Model(&User{}).Find(&data)
}


func TestPreload(t *testing.T) {
	db := _db.Debug()
	// db.AutoMigrate(&Girl{}, &Dog{})
	// g1 := &Girl{
	// 	Name: "Girl2",
	// }
	// db.Create(g1)
	// d1 := &Dog{
	// 	GirlID: 2,
	// 	Name:   "Girl2",
	// }
	// d2 := &Dog{
	// 	GirlID: 2,
	// 	Name:   "Girl2",
	// }
	// db.Create(d1)
	// db.Create(d2)
	var girls []Girl
	// db.Model(&Girl{}).Preload("Dogs").Find(&girls)
	db.Joins("LEFT JOIN dog on dog.girl_id = girl.id").Find(&girls)
	data, _ := json.Marshal(&girls)
	log.Printf("%s\n", data)
}

type XUser struct {
	gorm.Model
	ParentID uint
	Name     string
	DeptID   uint
}

var gUserTabs = 6

// 接口实现
func (o XUser) TableName() string {
	l := len(strconv.Itoa(gUserTabs - 1))
	return fmt.Sprintf("t_xuser_%0*d", l, o.DeptID%uint(gUserTabs))
}

func (XUser) TableNameOf(id uint) string {
	return XUser{DeptID: id}.TableName()
}

func (XUser) TableCount() uint {
	return uint(gUserTabs)
}

// 分表
func TestTables(t *testing.T) {
	CreateTables(&XUser{})
	user := []XUser{{
		Name:     "test",
		ParentID: 0,
		DeptID:   1,
	}, {
		Name:     "test",
		ParentID: 0,
		DeptID:   1,
	},
	}
	// 批量插入，需要人为保证数据中的数据在同一张表中
	_db.Table(user[0].TableName()).Create(&user)
	lUser := XUser{
		Name:     "test",
		ParentID: 0,
		DeptID:   2,
	}
	// 单个插入
	_db.Table(lUser.TableName()).Create(&lUser)
	// 查询
	var o XUser
	_db.Table(lUser.TableName()).First(&o)
	//
	var data []XUser
	_db.Table(o.TableName()).Find(&data)
	log.Println(data)
}

type PUser struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"primarykey"` // 定义为主键
	UpdatedAt time.Time
	Name      string
	DeptID    uint
}

func (PUser) TableName() string {
	return "t_puser"
}

// 分区
func TestPartition(t *testing.T) {
	_db.AutoMigrate(&PUser{})
	NewPartition(PUser{}.TableName()).AlterRange("created_at", 2)
	user := &PUser{
		Name: "test",
	}
	log.Println(DbCreate(&user))
}
