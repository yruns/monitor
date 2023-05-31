package model

import (
	"gorm.io/gorm"
	"time"
)

type Record struct {
	gorm.Model
	AttackType uint      `json:"attack_type" form:"attack_type"`
	Date       time.Time `json:"date" form:"date"`
}
