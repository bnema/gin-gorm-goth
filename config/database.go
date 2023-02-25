package config

import (
	"fmt"
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
	// If the connection is successful, the following message will be displayed
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
