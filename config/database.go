package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

// CONNECT TO DATABASE USING GORM
func ConnectDatabase(dsn string) *gorm.DB {
	tablePrefix := os.Getenv("DB_TABLE_PREFIX")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Silent in production, change to logger.Info for development
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		PrepareStmt: true, // Enable prepared statement for better performance
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// SET CONNECTION POOL SETTINGS
	sqlDB, err := db.DB()

	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}
