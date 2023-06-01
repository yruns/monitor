package dto

type Overview struct {
	// 攻击概览
	AttackIncrement int64 `json:"attack_increment" form:"attack_increment"`
	AttackTotal     int64 `json:"attack_total" form:"attack_total"`

	// 正常类型
	NormalIncrement int64 `json:"normal_increment" form:"normal_increment"`
	NormalTotal     int64 `json:"normal_total" form:"normal_total"`

	// 总数
	Variation int64 `json:"variation" form:"variation"`
	Total     int64 `json:"total" form:"total"`
}
