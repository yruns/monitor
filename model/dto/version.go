package dto

import "monitor/model"

type Version struct {
	model.Version
	AucList []model.Auc `json:"auc_list"`
}

func (v Version) name() {
}
