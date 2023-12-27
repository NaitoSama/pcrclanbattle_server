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
	dbDataInit()
}

// dbDataInit Initialize data from config
func dbDataInit() {
	// user init
	user := User{
		Name:     config.Config.DB.AdminName,
		Password: common.PasswordEncryption(config.Config.DB.AdminPasswd),
	}
	result := DB.Take(&user)
	if result.RowsAffected == 0 {
		DB.Create(&user)
	}

	// boss init
	boss := Boss{
		ID: 1,
	}
	result = DB.Take(&boss)
	if result.RowsAffected == 0 {
		for i := 0; i < 5; i++ {
			boss = Boss{
				ID:    i + 1,
				Stage: 1,
				Round: 1,
				Value: config.Config.Boss.StageOne[i],
			}
			DB.Create(&boss)
		}
	}

}
