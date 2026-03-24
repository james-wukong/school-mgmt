package models

import (
	"github.com/GoAdminGroup/go-admin/modules/db"
	"gorm.io/driver/postgres" // GORM v2 Postgres driver
	"gorm.io/gorm"
)

func Init(c db.Connection) (orm *gorm.DB) {
	orm, err := gorm.Open(postgres.New(postgres.Config{
		Conn: c.GetDB("default"),
	}), &gorm.Config{})

	if err != nil {
		panic("initialize orm failed")
	}

	return
}

// func GetORM() *gorm.DB {
//     return orm
// }
