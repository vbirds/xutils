package orm

import (
	"log"
	"testing"
	"time"
)

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

func TestPartition(t *testing.T) {
	db, err := NewGormV2("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&PUser{})
	NewPartiton(PUser{}.TableName()).AlterRange("created_at", 2)
	user := &PUser{
		Name: "test",
	}
	log.Println(DbCreate(&user))
}
