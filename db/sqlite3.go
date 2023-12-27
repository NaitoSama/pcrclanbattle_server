package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"pcrclanbattle_server/common"
	"pcrclanbattle_server/config"
)

var DB *gorm.DB

func DBInit() {
	db, err := gorm.Open(sqlite.Open("./db/clanbattle.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
		return
	}
	err = db.AutoMigrate(&Record{}, &Boss{}, &User{})
	if err != nil {
		panic("failed to migrate database")
		return
	}
	DB = db
	common.Logln(0, "database started")

}

// dbDataInit Initialize data from config
func dbDataInit() {
	user := User{
		Name:     config.Config.DB.AdminName,
		Password: common.PasswordEncryption(config.Config.DB.AdminPasswd),
	}
	result := DB.Take(&user)
	if result.RowsAffected == 0 {
		DB.Create(&user)
	}

}
