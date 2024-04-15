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
	BeforeBossValueD  int64
	BeforeBossWhoIsIn string
	BeforeBossTree    string
	ArchiveID         string
}

type Boss struct {
	ID      int
	Stage   int
	Round   int
	Value   int64
	WhoIsIn string
	Tree    string
	ValueD  int64
	PicETag string
}

type User struct {
	gorm.Model
	UserID     int
	Name       string
	Password   string
	Permission int
	UserPic    string
	PicETag    string
	UserPic16  string
	Pic16ETag  string
}
