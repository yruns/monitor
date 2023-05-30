package database

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var Mysql *gorm.DB

func InitMysql(connRead, connWrite string) {
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       connWrite,
		DefaultStringSize:         256,  // string类型默认长度
		DisableDatetimePrecision:  true, // 禁止datetime精度
		DontSupportRenameIndex:    true, // 重命名索引
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(20)  // 最大连接池
	sqlDB.SetMaxOpenConns(100) // 打开连接数

	Mysql = db

	// 主从配置
	_ = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(connWrite)},
		Replicas: []gorm.Dialector{mysql.Open(connRead)},
		Policy:   dbresolver.RandomPolicy{}, // 负载均衡策略
	}))

	// 暂时关闭自动迁移
	//migration()
}

//func NewDBClient(ctx context.Context) *gorm.DB {
//	db := _db
//	return db.WithContext(ctx)
//}
