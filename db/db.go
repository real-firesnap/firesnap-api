package db

import (
	"log"
	"os"

	"github.com/real-firesnap/firesnap-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %s", err.Error())
	}

	DB = db
}

func MigrateDatabase() {
	DB.AutoMigrate(
		&models.User{},
		&models.Post{},
	)
}
