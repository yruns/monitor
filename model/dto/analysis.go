package dto

type Analysis struct {
	// 近七天攻击数量
	AttackNum [7]uint `json:"attack_num" form:"attack_num"`
}
