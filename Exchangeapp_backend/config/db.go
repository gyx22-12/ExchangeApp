package config

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() {
	dsn := AppConfig.Database.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to initialize database, got error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection, got error: %v", err)
	}

	sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	global.Db = db

	// 启动时一次性迁移表结构，避免在请求处理里重复 AutoMigrate
	if err := db.AutoMigrate(&models.Article{}, &models.User{}, &models.ExchangeRate{}); err != nil {
		log.Fatalf("Failed to migrate database, got error: %v", err)
	}
}
