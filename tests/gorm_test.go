package tests

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"monitor/database"
	"os"
	"testing"
	"time"
)

// 连接数据库
func TestConnect(t *testing.T) {
	db := database.Mysql
	fmt.Println(db)
}

func TestCRUD(t *testing.T) {
	customLog := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	connectURL := "root" + ":" + "root" + "@tcp(" + "localhost" + ")/" + "gin_mall"
	Mysql, _ := gorm.Open(mysql.Open(connectURL+"?charset=utf8mb4"), &gorm.Config{Logger: customLog})

	//var hashedPasswordFromDB string
	//Mysql.Raw("SELECT password_digest from `user` where username = ? AND `user`.`deleted_at` IS NULL", "294056734").
	//	Scan(&hashedPasswordFromDB)
	//fmt.Println(hashedPasswordFromDB)

	Mysql.Table("user").Where("id = ?", 1).Update("avatar", "333")
}
