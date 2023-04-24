package config

import (
	"api/app/models"
	"api/app/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func ConnectToDatabase() *gorm.DB {
	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN:                  os.Getenv("DATABASE_CONNECTION_URI"),
				PreferSimpleProtocol: true,
			},
		),
		&gorm.Config{},
	)
	if err != nil {
		utils.ErrorLogger.Fatal("failed to connect database")
	}
	return db
}

func RunDatabaseMigrations(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{}, &models.MarketDeck{}, &models.MarketDeckImage{},
		&models.MarketDeckCard{}, &models.MarketDeckReview{},
	)
	if err != nil {
		utils.ErrorLogger.Fatal("failed to do migrations")
	}
}
