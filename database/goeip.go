package database

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
)

var GeoIP *geoip2.Reader

func InitGeoIp(filePath string) {
	db, err := geoip2.Open(filePath)
	if err != nil {
		fmt.Println("GeoIp文件导入失败")
	}

	GeoIP = db
}
