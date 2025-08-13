package database

import (
	"fmt"
	"time"

	"fastdfs-migration-system/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 数据库实例
var DB *gorm.DB

// Config 数据库配置
type Config struct {
	Type     string `mapstructure:"type"`
	DSN      string `mapstructure:"dsn"`
	LogLevel string `mapstructure:"log_level"`
}

// Initialize 初始化数据库连接
func Initialize(config Config) error {
	var err error
	var dialector gorm.Dialector
	
	// 根据配置选择数据库驱动
	switch config.Type {
	case "sqlite":
		dialector = sqlite.Open(config.DSN)
	default:
		return fmt.Errorf("unsupported database type: %s", config.Type)
	}
	
	// 配置GORM日志级别
	var logLevel logger.LogLevel
	switch config.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Warn
	}
	
	// 创建数据库连接
	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	
	return nil
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	
	// 自动迁移所有模型
	err := DB.AutoMigrate(
		&models.Migration{},
		&models.Cluster{},
		&models.TaskLog{},
		&models.ScheduledTask{},
		&models.TransferState{},
	)
	
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}
	
	return nil
}

// Close 关闭数据库连接
func Close() error {
	if DB == nil {
		return nil
	}
	
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	
	return sqlDB.Close()
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}