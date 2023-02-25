package config

import (
	"fmt"
	"go-gorm-gauth/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dns := os.Getenv("DATABASE_URL")

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		panic(err)

	}

	DB = db

	// Migrate the schema
	migrateModels := db.AutoMigrate(&models.User{}, &models.Account{}, &models.Session{})
	if migrateModels != nil {
		panic(migrateModels)
	} else {
		fmt.Println("Database migration successful")
	}

	fmt.Println("Database connection successful")
}

func CloseDB() {
	db, err := DB.DB()
	if err != nil {
		panic(err)
	}
	db.Close()
	fmt.Println("Database connection closed")
}
