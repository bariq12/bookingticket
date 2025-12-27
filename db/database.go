package db

import (
	"fmt"
	"log"

	"github.com/bariq12/bookingticket/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(config *config.EnvConfig, DBMigrator func(db *gorm.DB) error) *gorm.DB {
	uri := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s port=5432",
		config.DBHost, config.DBUSer, config.DBPassword, config.DBName, config.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	log.Println("connected to database!")

	if err := DBMigrator(db); err != nil {
		log.Fatalf("unable to migrate tables: %v", err)
	}

	return db
}
