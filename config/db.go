package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Frhnmj2004/LabMonitoring-server/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Set connection pool settings
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance: ", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	// First create the extension for UUID support
	err = DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		log.Fatal("Failed to create extension: ", err)
	}

	// Drop existing tables to start fresh
	err = DB.Exec("DROP TABLE IF EXISTS internet_usages, alerts, resource_logs, computers, users CASCADE;").Error
	if err != nil {
		log.Fatal("Failed to drop tables: ", err)
	}

	// Create tables with proper references
	createSchema()

	log.Println("Database connection established successfully")
}

func createSchema() {
	// Create User model
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to create User table: ", err)
	}

	// Create Computer model
	err = DB.AutoMigrate(&models.Computer{})
	if err != nil {
		log.Fatal("Failed to create Computer table: ", err)
	}

	// Create ResourceLog model
	err = DB.AutoMigrate(&models.ResourceLog{})
	if err != nil {
		log.Fatal("Failed to create ResourceLog table: ", err)
	}

	// Create Alert model
	err = DB.AutoMigrate(&models.Alert{})
	if err != nil {
		log.Fatal("Failed to create Alert table: ", err)
	}

	// Create InternetUsage model
	err = DB.AutoMigrate(&models.InternetUsage{})
	if err != nil {
		log.Fatal("Failed to create InternetUsage table: ", err)
	}
}
