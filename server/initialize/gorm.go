package initialize

import (
	"fmt"
	"server/global"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitGorm() {
	cfg := global.Config.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&collation=utf8mb4_unicode_ci",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	var logMode logger.LogLevel
	if global.Config.Server.Mode == "debug" {
		logMode = logger.Info
	} else {
		logMode = logger.Silent
	}

	baseLogger := logger.Default.LogMode(logMode)
	slowLogger := NewSlowLogger(baseLogger, 1000*time.Millisecond) // 超过200ms记录慢查询

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: slowLogger,
	})
	if err != nil {
		panic(fmt.Errorf("连接数据库失败: %w", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("获取数据库实例失败: %w", err))
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	global.DB = db
	global.Log.Info("数据库连接成功")
}
