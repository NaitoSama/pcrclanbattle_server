package model

// data contains type and one of the following structs

// type:attack
type AttackBoss struct {
	Type     int    `json:"type"` // 0-damage 1-defeat
	BossID   int    `json:"boss_id"`
	Value    int64  `json:"value"`     // damage value
	FromName string `json:"from_name"` // who caused
}
type AttackPayload struct {
	AttackBoss `json:"attack_boss"`
}

// type:revise
// ReviseBoss revise boss value
type ReviseBoss struct {
	BossID int   `json:"boss_id"`
	Value  int64 `json:"value"`
	Round  int   `json:"round"`
}
type RevisePayload struct {
	ReviseBoss `json:"revise_boss"`
}

// type:undo
// Undo the attack
type Undo struct {
	FromName string `json:"from_name"` // who's attack
	BossID   int    `json:"boss_id"`   // attack to who
}
type UndoPayload struct {
	Undo `json:"undo"`
}

// type:imin
// ImIn who is attacking boss
type ImIn struct {
	FromName string `json:"from_name"`
	BossID   int    `json:"boss_id"`
}
type ImInPayload struct {
	ImIn `json:"im_in"`
}

// type:imout
// ImOut undo ImIn
type ImOut struct {
	FromName string `json:"from_name"`
	BossID   int    `json:"boss_id"`
}
type ImOutPayload struct {
	ImOut `json:"im_out"`
}
