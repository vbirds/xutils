package orm

import (
	"encoding/json"
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
