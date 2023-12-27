package db

import "gorm.io/gorm"

type Record struct {
	gorm.Model
	AttackFrom string
	AttackTo   string
	Damage     int64
}

type Boss struct {
	Stage int
	Round int
	Value int64
}

type User struct {
	gorm.Model
	Name     string
	Password string
}
