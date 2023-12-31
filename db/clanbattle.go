package db

import "gorm.io/gorm"

type Record struct {
	gorm.Model
	AttackFrom        string
	AttackTo          int
	Damage            int64
	CanUndo           int
	BeforeBossStage   int
	BeforeBossRound   int
	BeforeBossValue   int64
	BeforeBossWhoIsIn string
	BeforeBossTree    string
}

type Boss struct {
	ID      int
	Stage   int
	Round   int
	Value   int64
	WhoIsIn string
	Tree    string
}

type User struct {
	gorm.Model
	UserID     int
	Name       string
	Password   string
	Permission int
}
