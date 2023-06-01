package model

import (
	"time"
)

type Record struct {
	Id       int64     `json:"id" form:"id"`
	Date     time.Time `json:"date" form:"date"`           // 数据包时间
	Duration int64     `json:"duration" form:"duration"`   // 持续时间
	Protocol string    `json:"protocol" form:"protocol"`   // 协议
	Service  string    `json:"service" form:"service"`     // 服务类型
	Flag     string    `json:"flag" form:"flag"`           // flag
	DstPort  int64     `json:"dst_port" form:"dst_port"`   // 目标端口
	SrcHost  string    `json:"src_host" form:"src_host"`   // 源主机
	SrcBytes int64     `json:"src_bytes" form:"src_bytes"` // 原始字节数
	DstBytes int64     `json:"dst_bytes" form:"dst_bytes"` // 目标字节数
	Label    string    `json:"label" form:"label"`         // 攻击类型
}
