package orm

import (
	"log"
	"testing"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ParentId uint
	Name     string
	Children []User `json:"children,omitempty" gorm:"foreignKey:ParentId;"` // 这里注意，如果设置ParentId为0，要禁用外键约束
}

func TestOrm(t *testing.T) {
	db, err := NewGormV2("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&User{})
	user := &User{
		Name:     "test",
		ParentId: 1,
	}
	DbCreate(&user)
	var o User
	gOrmDb.Model(&User{}).Where("id = ?", 1).Preload("Children").First(&o)
	log.Println(o)
}
