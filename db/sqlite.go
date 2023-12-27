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
	dbDataInit()
	common.Logln(0, "database started")
}

// dbDataInit Initialize data from config
func dbDataInit() {
	// user init
	username := config.Config.DB.AdminName
	password := common.PasswordEncryption(config.Config.DB.AdminPasswd)
	user := User{}
	result := DB.Model(user).Where("name = ? and password = ?", username, password).First(&user)
	if result.RowsAffected == 0 {
		user = User{
			Name:     username,
			Password: password,
		}
		DB.Create(&user)
	}

	// boss init
	boss := Boss{}
	result = DB.Model(boss).Where("id = ?", 3).First(&boss)
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
