package dto

type Table struct {
	AttackType string  `json:"attack_type"`
	Proportion float64 `json:"proportion"`

	HazardLevel string `json:"hazard_level"`
	Frequency   int64  `json:"frequency"`
	Severity    int64  `json:"severity"`
	Grade       string `json:"grade"`
	Measure     string `json:"measure"`
}
