package dto

type Analysis struct {
	// 近七天攻击数量
	AttackNum [7]int64 `json:"attack_num" form:"attack_num"`
}
