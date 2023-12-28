package model

// data contains type and one of the following structs

// type:attack
type AttackBoss struct {
	Type     int // 0-damage 1-defeat
	BossID   int
	Value    int64  // damage value
	FromName string // who caused
}

// type:revise
// ReviseBoss revise boss value
type ReviseBoss struct {
	BossID int
	Value  int64
	Round  int
}

// type:undo
// Undo the attack
type Undo struct {
	FromName string // who's attack
	BossID   int    // attack to who
}

// type:imin
// ImIn who is attacking boss
type ImIn struct {
	FromName string
	BossID   int
}

// type:imout
// ImOut undo ImIn
type ImOut struct {
	FromName string
	BossID   int
}
