package conf

import (
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	gornmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"time"
)

type Config struct {
	Rabbitmq struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
		Port     string `json:"port"`
	} `json:"rabbitmq"`
	Mysql struct {
		Host        string `json:"host"`
		Port        string `json:"port"`
		User        string `json:"user"`
		Password    string `json:"password"`
		Name        string `json:"name"`
		MaxOpenConn int    `json:"max_open_conn"`
	} `json:"mysql"`
}

var instanceDb *gorm.DB
var DefaultConfig *Config

func GetConfig() *Config {
	if DefaultConfig == nil {
		conf, err := ioutil.ReadFile("./conf/config.json")
		if err != nil {
			fmt.Println(err)
			panic("读取配置失败")
		}
		err = json.Unmarshal(conf, &DefaultConfig)
		if err != nil {
			panic(err)
		}
		return DefaultConfig
	}
	return DefaultConfig
}
func NewDb() (*gorm.DB, error) {
	if instanceDb == nil {
		config := GetConfig().getConfig()
		db, err := gorm.Open(gornmysql.Open(config), &gorm.Config{})
		// SetMaxIdleConns 设置空闲连接池中连接的最大数量
		sqlDB, err := db.DB()
		// SetMaxIdleConns 设置空闲连接池中连接的最大数量
		sqlDB.SetMaxIdleConns(10)
		// SetMaxOpenConns 设置打开数据库连接的最大数量。
		sqlDB.SetMaxOpenConns(100)
		// SetConnMaxLifetime 设置了连接可复用的最大时间。
		sqlDB.SetConnMaxLifetime(time.Hour)
		instanceDb = db
		return instanceDb, err
	}
	return instanceDb, nil
}
func (config *Config) getConfig() string {
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.User = config.Mysql.User
	mysqlConfig.DBName = config.Mysql.Name
	mysqlConfig.Passwd = config.Mysql.Password
	mysqlConfig.ParseTime = true
	mysqlConfig.Net = "tcp"
	mysqlConfig.Addr = config.Mysql.Host + ":" + config.Mysql.Port
	return mysqlConfig.FormatDSN()
}
