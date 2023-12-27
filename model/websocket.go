package model

type AttackBoss struct {
	Type     int // 0-damage 1-defeat
	BossID   int
	Value    int64  // damage value
	FromName string // who caused
}

// ReviseBoss revise boss value
type ReviseBoss struct {
	BossID int
	Value  int64
	Round  int
}

// Undo the attack
type Undo struct {
	FromName string // who's attack
	BossID   int    // attack to who
}

// ImIn who is attacking boss
type ImIn struct {
	FromName string
	BossID   int
}

// ImOut undo ImIn
type ImOut struct {
	FromName string
	BossID   int
}
