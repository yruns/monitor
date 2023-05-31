package database

import (
	"fmt"
	"monitor/model"
)

func migration() {
	err := Mysql.Set("gorm.table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.User{},
			&model.Admin{},
			&model.Notice{},
			&model.Record{},
		)
	if err != nil {
		fmt.Println("自动迁移失败")
	}

	return
}
