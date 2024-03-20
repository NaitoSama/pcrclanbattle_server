package db

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"pcrclanbattle_server/common"
	"pcrclanbattle_server/config"
)

var DB *gorm.DB

type cache struct {
	Bosses  []Boss
	Records []Record
	Users   map[string]*User
}

var Cache = &cache{
	Users: make(map[string]*User),
}

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
	dbCacheInit()
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
			Name:       username,
			Password:   password,
			Permission: 2,
		}
		DB.Create(&user)
	}

	// boss init
	boss := Boss{}
	result = DB.Model(boss).Where("id = ?", 3).First(&boss)
	if result.RowsAffected == 0 {
		for i := 0; i < 5; i++ {
			picETag, err := common.CalculateETag(fmt.Sprintf("./pic/%d.jpg", i+1))
			if err != nil {
				common.Logln(2, err)
			}
			boss = Boss{
				ID:      i + 1,
				Stage:   1,
				Round:   1,
				Value:   config.Config.Boss.StageOne[i],
				WhoIsIn: " ",
				Tree:    " ",
				ValueD:  config.Config.Boss.StageOne[i],
				PicETag: picETag,
			}
			DB.Create(&boss)
		}
	}
}

func dbCacheInit() {
	// boss cache init
	DB.Model(Boss{}).Find(&Cache.Bosses)
	// record cache init
	DB.Model(Record{}).Find(&Cache.Records)
	// user cache init
	var users []User
	DB.Model(User{}).Find(&users)
	for i := 0; i < len(users); i++ {
		Cache.Users[users[i].Name] = &users[i]
	}
}
