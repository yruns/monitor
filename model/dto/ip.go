package dto

import "time"

type IpStatistics struct {
	IP         string    `json:"ip"`
	Count      int64     `json:"count"`
	Address    string    `json:"address"`
	AttackName []string  `json:"attack_name"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}
