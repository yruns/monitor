package dto

type Statistics struct {
	AttackName []string `json:"attack_name" form:"attack_name"`
	AttackNum  []uint   `json:"attack_num" form:"attack_num"`
}
