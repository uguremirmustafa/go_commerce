package database

import (
	"log"
	"os"

	"github.com/uguremirmustafa/go_commerce/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	var err error
	dbstring := os.Getenv("DATABASE_URL")
	DB, err = gorm.Open(postgres.Open(dbstring), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database \n", err.Error())
		os.Exit(2)
	}

	log.Println("Connected to the database successfully")
	DB.Logger = logger.Default.LogMode(logger.Info)
	log.Println("running migrations")

	// TODO: Add migrations
	DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})
}
