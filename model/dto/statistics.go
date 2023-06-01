package dto

type Statistics struct {
	AttackName []string `json:"attack_name" form:"attack_name"`
	AttackNum  []int64  `json:"attack_num" form:"attack_num"`
}
