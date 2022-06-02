package orm

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"gorm.io/gorm"
)

type XUser struct {
	gorm.Model
	ParentID uint
	Name     string
	DeptID   uint
}

var gUserTabs = 6

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

func TestXTablers(t *testing.T) {
	db, err := NewGormV2("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln(err)
	}
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
	// 批量插入， 需要人为保证数据中的数据在同一张表中
	db.Table(user[0].TableName()).Create(&user)
	lUser := XUser{
		Name:     "test",
		ParentID: 0,
		DeptID:   2,
	}
	// 单个插入
	db.Table(lUser.TableName()).Create(&lUser)
	// 查询
	var o XUser
	db.Table(lUser.TableName()).First(&o)
	//
	var data []XUser
	db.Table(o.TableName()).Find(&data)
	log.Println(data)
}
