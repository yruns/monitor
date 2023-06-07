package model

type Auc struct {
	Id        uint    `json:"id" form:"id" gorm:"primary_key"`
	VersionId string  `json:"version_id" form:"version_id"`
	Label     string  `json:"label" form:"label"`
	P         float64 `json:"precision" form:"p"`
	R         float64 `json:"recall" form:"r"`
	F1        float64 `json:"f1" form:"f1"`
}
