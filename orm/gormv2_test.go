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

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ParentId uint
	Name     string
	DeptID   uint
	Children []User `json:"children,omitempty" gorm:"foreignKey:ParentId;"` // 这里注意，如果设置ParentId为0，要禁用外键约束
}

var gUserTabs = 6

func (o User) TableName() string {
	l := len(strconv.Itoa(gUserTabs - 1))
	return fmt.Sprintf("t_user_%0*d", l, o.DeptID%uint(gUserTabs))
}

func (User) TableNameOf(id uint) string {
	return User{DeptID: id}.TableName()
}

func (User) TableCount() uint {
	return uint(gUserTabs)
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
	gOrmDb.AutoMigrate(&User{})
	user := &User{
		Name:     "test",
		ParentId: 1,
	}
	DbCreate(&user)
	var o User
	gOrmDb.Model(&User{}).Where("id = ?", 1).Preload("Children").First(&o)
	log.Println(o)
}

func TestXTablers(t *testing.T) {
	CreateTables(&User{})
	user := []User{{
		Name:     "test",
		ParentId: 0,
		DeptID:   1,
	}, {
		Name:     "test",
		ParentId: 0,
		DeptID:   1,
	},
	}
	// 批量插入， 需要人为保证数据中的数据在同一张表中
	gOrmDb.Table(user[0].TableName()).Create(&gOrmDb)
	lUser := user[0]
	// 单个插入
	gOrmDb.Table(lUser.TableName()).Create(&lUser)
	// 查询
	var o User
	gOrmDb.Table(lUser.TableName()).First(&o)
	//
	var data []User
	gOrmDb.Table(o.TableName()).Find(&data)

	log.Println(o)
}

func TestPreload(t *testing.T) {
	db := gOrmDb.Debug()
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
