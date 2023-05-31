package dto

type Overview struct {
	// 攻击概览
	AttackIncrement uint `json:"attack_increment" form:"attack_increment"`
	AttackTotal     uint `json:"attack_total" form:"attack_total"`

	// 正常类型
	NormalIncrement uint `json:"normal_increment" form:"normal_increment"`
	NormalTotal     uint `json:"normal_total" form:"normal_total"`

	// 总数
	Variation uint `json:"variation" form:"variation"`
	Total     uint `json:"total" form:"total"`
}
