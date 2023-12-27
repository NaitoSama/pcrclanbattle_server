package db

import "gorm.io/gorm"

type Record struct {
	gorm.Model
	AttackFrom string
	AttackTo   string
	Damage     int64
}

type Boss struct {
	ID    int
	Stage int
	Round int
	Value int64
}

type User struct {
	gorm.Model
	UserID     int
	Name       string
	Password   string
	Permission int
}
