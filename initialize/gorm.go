package initialize

import (
	"fmt"
	"github.com/Demonx24/Vblog-backend/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

// InitGorm 初始化并返回一个使用 MySQL 配置的 GORM 数据库连接
func InitGorm() *gorm.DB {
	mysqlCfg := global.Config.Mysql

	// 使用给定的 DSN（数据源名称）和日志级别打开 MySQL 数据库连接
	db, err := gorm.Open(mysql.Open(mysqlCfg.Dsn()), &gorm.Config{
		Logger: logger.Default.LogMode(mysqlCfg.LogLevel()), // 设置日志级别
	})
	if err != nil {
		fmt.Println("Failed to connect to MySQL", err)
		os.Exit(1)
	}

	// 获取底层的 SQL 数据库连接对象
	sqlDB, _ := db.DB()
	// 设置数据库连接池中的最大空闲连接数
	sqlDB.SetMaxIdleConns(mysqlCfg.MaxIdleConns)
	// 设置数据库的最大打开连接数
	sqlDB.SetMaxOpenConns(mysqlCfg.MaxOpenConns)

	return db
}
