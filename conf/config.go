package conf

import (
	"fmt"
	"gopkg.in/ini.v1"
	"monitor/database"
)

var (
	AppMode  string
	HttpPort string

	MysqlHost     string
	MysqlPort     string
	MysqlUser     string
	MysqlPasswd   string
	MysqlDatabase string

	RedisHost     string
	RedisPort     int
	RedisPasswd   string
	RedisDatabase int

	GeoFilePath string

	MinioHost            string
	MinioPort            string
	MinioAccessKeyID     string
	MinioSecretAccessKey string

	ValidEmail string
	SmtpHost   string
	SmtpUser   string
	SmtpPass   string

	ProjectHost string
	ProductPath string
	AvatarPath  string
)

func Init() {
	file, err := ini.Load("./conf/conf.ini")
	if err != nil {
		panic(err)
	}

	loadServer(file)
	loadMysql(file)
	loadRedis(file)
	loadGeoIp(file)
	loadMinio(file)
	loadEmail(file)

	connRead := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MysqlUser, MysqlPasswd, MysqlHost, MysqlPort, MysqlDatabase)
	connWrite := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MysqlUser, MysqlPasswd, MysqlHost, MysqlPort, MysqlDatabase)
	database.InitMysql(connRead, connWrite)
	database.InitRedis(RedisHost, RedisPasswd, RedisPort, RedisDatabase)
	database.InitGeoIp(GeoFilePath)

}

func loadGeoIp(file *ini.File) {
	GeoFilePath = file.Section("geoip").Key("GeoFilePath").String()
}

func loadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}

func loadMysql(file *ini.File) {
	MysqlHost = file.Section("mysql").Key("MysqlHost").String()
	MysqlPort = file.Section("mysql").Key("MysqlPort").String()
	MysqlUser = file.Section("mysql").Key("MysqlUser").String()
	MysqlPasswd = file.Section("mysql").Key("MysqlPasswd").String()
	MysqlDatabase = file.Section("mysql").Key("MysqlDatabase").String()
}

func loadRedis(file *ini.File) {
	RedisHost = file.Section("redis").Key("RedisHost").String()
	RedisPort, _ = file.Section("redis").Key("RedisPort").Int()
	RedisPasswd = file.Section("redis").Key("RedisPasswd").String()
	RedisDatabase, _ = file.Section("redis").Key("RedisDatabase").Int()
}

func loadMinio(file *ini.File) {
	MinioHost = file.Section("minio").Key("MinioHost").String()
	MinioPort = file.Section("minio").Key("MinioPort").String()
	MinioAccessKeyID = file.Section("minio").Key("MinioAccessKeyID").String()
	MinioSecretAccessKey = file.Section("minio").Key("MinioSecretAccessKey").String()
}

func loadEmail(file *ini.File) {
	ValidEmail = file.Section("email").Key("ValidEmail").String()
	SmtpHost = file.Section("email").Key("SmtpHost").String()
	SmtpUser = file.Section("email").Key("SmtpUser").String()
	SmtpPass = file.Section("email").Key("SmtpPass").String()
}
